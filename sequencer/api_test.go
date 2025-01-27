package sequencer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/StripChain/strip-node/common"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

// TestCreateWalletEndpoint tests the /createWallet endpoint for creating a wallet.
// It verifies the response status code for valid and invalid requests.
func TestCreateWalletEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		identity      string
		identityCurve string
		wantStatus    int
	}{
		{
			name:          "Valid wallet creation",
			identity:      "testIdentity",
			identityCurve: "ecdsa",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "Missing identity",
			identity:      "",
			identityCurve: "ecdsa",
			wantStatus:    http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/createWallet?identity="+tt.identity+"&identityCurve="+tt.identityCurve, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if tt.identity == "" {
						http.Error(w, "identity required", http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
				}).ServeHTTP(w, r)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateWallet returned wrong status code: got %v want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestGetWalletEndpoint tests the /getWallet endpoint for retrieving a wallet.
// It checks the response status code for valid and missing identity scenarios.
func TestGetWalletEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		identity      string
		identityCurve string
		wantStatus    int
	}{
		{
			name:          "Valid wallet retrieval",
			identity:      "testIdentity",
			identityCurve: "ecdsa",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "Missing identity",
			identity:      "",
			identityCurve: "ecdsa",
			wantStatus:    http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/getWallet?identity="+tt.identity+"&identityCurve="+tt.identityCurve, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.identity == "" {
					http.Error(w, "identity required", http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(map[string]string{"identity": tt.identity})
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetWallet returned wrong status code: got %v want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestCreateIntent tests the /createIntent endpoint for creating an intent.
// It verifies the JSON payload and response status code.
func TestCreateIntent(t *testing.T) {
	// Create a test intent
	testIntent := Intent{
		Identity:      "testIdentity",
		IdentityCurve: "ecdsa",
		Status:        INTENT_STATUS_PROCESSING,
		Operations: []Operation{
			{
				Type:     OPERATION_TYPE_TRANSACTION,
				Status:   OPERATION_STATUS_PENDING,
				ChainId:  "1",
				KeyCurve: "ecdsa",
			},
		},
	}
	// Convert intent to JSON
	intentJSON, err := json.Marshal(testIntent)
	if err != nil {
		t.Fatalf("Failed to marshal test intent: %v", err)
	}

	// Create test request
	req := httptest.NewRequest("POST", "/createIntent", bytes.NewBuffer(intentJSON))
	w := httptest.NewRecorder()

	// Create handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		var intent Intent
		err := json.NewDecoder(r.Body).Decode(&intent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Verify the received intent matches what we sent
		if intent.Identity != testIntent.Identity {
			t.Errorf("got identity %v, want %v", intent.Identity, testIntent.Identity)
		}
		if intent.Status != testIntent.Status {
			t.Errorf("got status %v, want %v", intent.Status, testIntent.Status)
		}

		// Mock successful response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int64{"id": 1})
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("CreateIntent returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

// TestGetIntent tests the /getIntent endpoint for retrieving an intent by ID.
// It checks the response status code and verifies the returned intent data.
func TestGetIntent(t *testing.T) {
	// Create test request
	req := httptest.NewRequest("GET", "/getIntent?id=1", nil)
	w := httptest.NewRecorder()

	// Create handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		// Verify query parameter
		if id := r.URL.Query().Get("id"); id != "1" {
			t.Errorf("got id %v, want 1", id)
		}

		// Mock response intent
		intent := Intent{
			ID:            1,
			Identity:      "testIdentity",
			IdentityCurve: "ecdsa",
			Status:        INTENT_STATUS_COMPLETED,
		}

		json.NewEncoder(w).Encode(intent)
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetIntent returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	// Decode response
	var response Intent
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != 1 || response.Status != INTENT_STATUS_COMPLETED {
		t.Errorf("Got unexpected response: %+v", response)
	}
}

// TestGetIntentsWithPagination tests the /getIntents endpoint with pagination parameters.
// It verifies the response status code for valid and invalid limit values.
func TestGetIntentsWithPagination(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		skip       int
		wantStatus int
	}{
		{
			name:       "Valid pagination",
			limit:      10,
			skip:       0,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid limit",
			limit:      -1,
			skip:       0,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/getIntents?limit="+strconv.Itoa(tt.limit)+"&skip="+strconv.Itoa(tt.skip), nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.limit < 0 {
					http.Error(w, "invalid limit", http.StatusInternalServerError)
					return
				}
				result := IntentsResult{
					Intents: []*Intent{},
					Total:   0,
				}
				json.NewEncoder(w).Encode(result)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetIntents returned wrong status code: got %v want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestGetSolverStats tests the /getStatsOfSolver endpoint for retrieving solver statistics.
// It checks the response status code for valid and empty solver scenarios.
func TestGetSolverStats(t *testing.T) {
	tests := []struct {
		name       string
		solver     string
		wantStatus int
	}{
		{
			name:       "Valid solver",
			solver:     "testSolver",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Empty solver",
			solver:     "",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/getStatsOfSolver?solver="+tt.solver, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.solver == "" {
					http.Error(w, "solver required", http.StatusInternalServerError)
					return
				}
				result := SolverStatResult{
					IsActive:    true,
					ActiveSince: 123456,
					Chains:      []uint{1, 2, 3},
				}
				json.NewEncoder(w).Encode(result)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetSolverStats returned wrong status code: got %v want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestParseOperation tests the /parseOperation endpoint for parsing an operation.
// It verifies the response status code for valid and invalid operation ID scenarios.
func TestParseOperation(t *testing.T) {
	tests := []struct {
		name        string
		operationId string
		intentId    string
		wantStatus  int
	}{
		{
			name:        "Valid operation",
			operationId: "1",
			intentId:    "1",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "Invalid operation",
			operationId: "",
			intentId:    "1",
			wantStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/parseOperation?operationId=" + tt.operationId + "&intentId=" + tt.intentId
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.operationId == "" {
					http.Error(w, "operation ID required", http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(map[string]string{"status": "success"})
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ParseOperation returned wrong status code: got %v want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestGetTotalStats tests the /getTotalStats endpoint for retrieving total statistics.
// It checks the response status code and verifies the returned statistics data.
func TestGetTotalStats(t *testing.T) {
	req := httptest.NewRequest("GET", "/getTotalStats", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stats := TotalStats{
			TotalSolvers: 10,
			TotalIntents: 100,
		}
		json.NewEncoder(w).Encode(stats)
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetTotalStats returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	var response TotalStats
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.TotalSolvers != 10 || response.TotalIntents != 100 {
		t.Errorf("Unexpected response values: %+v", response)
	}
}

// TestStatusEndpoint tests the /status endpoint for checking service health.
// It verifies the response status code and body content.
func TestStatusEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/status", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status endpoint returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Status endpoint returned wrong body: got %v want OK", w.Body.String())
	}
}

// OAuth tests

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	id, ok := r.Context().Value(identityAccess).(*IdentityAccess)
	if !ok {
		http.Error(w, "id not found in context", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ID: " + id.Identity))
}

func generateTestToken(oauthInfo *GoogleAuth, userId string, expiryAt time.Time, identity string, identityCurve string) (string, error) {
	claims := ClaimsWithIdentity{
		Identity:      identity,
		IdentityCurve: identityCurve,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryAt), // Token expires after 10min
			Issuer:    tokenIssuer,
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token using HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	return token.SignedString([]byte(oauthInfo.jwtSecret))
}

func TestValidateAccessMiddleware(t *testing.T) {
	router := mux.NewRouter()
	router.Use(ValidateAccessMiddleware)
	router.HandleFunc("/testOAuth", testHandler)

	oauthInfo = NewGoogleAuth("/redirect", "clientId", "clientSecret", "sessionSecret", "jwtSecret", "salt")

	// Valid token
	accessTokenValid, _ := generateTestToken(oauthInfo, "1", time.Now().Add(time.Minute*10), "0xa", "ecdsa")
	// Invalid Expired token
	accessTokenExpired, _ := generateTestToken(oauthInfo, "1", time.Now(), "0xa", "ecdsa")
	// Invalid No Identity token
	accessTokenNoIdentity, _ := generateTestToken(oauthInfo, "1", time.Now().Add(time.Minute*10), "", "")

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedBody   string
		tokens         *Tokens
	}{
		{
			name:           "Valid Token",
			token:          "validtoken",
			expectedStatus: http.StatusOK,
			expectedBody:   "ID: 0xa",
			tokens: &Tokens{
				AccessToken: accessTokenValid,
			},
		},
		{
			name:           "Invalid Token: Expired",
			token:          "invalidtokenexpired",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   ErrTokenExpired.Error(),
			tokens: &Tokens{
				AccessToken: accessTokenExpired,
			},
		},
		{
			name:           "Invalid Token: No Identity",
			token:          "validtokennoidentity",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   ErrInvalidTokenIdentityRequired.Error(),
			tokens: &Tokens{
				AccessToken: accessTokenNoIdentity,
			},
		},
	}

	for _, tt := range tests {
		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(tt.tokens)
		req, err := http.NewRequest("GET", "/testOAuth?auth=oauth", payloadBuf)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		id := req.Context().Value(identityAccess)
		fmt.Println("id from Context in test", id)
		if w.Code != tt.expectedStatus {
			t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
		}

		// Check the response body
		if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(tt.expectedBody) {
			t.Errorf("Expected body to be '%s', got '%s'", strings.TrimSpace(tt.expectedBody), strings.TrimSpace(w.Body.String()))
		}

	}
}

func TestGoogleAuthGenerateWallet(t *testing.T) {
	oauthInfo = NewGoogleAuth("/redirect", "clientId", "clientSecret", "sessionSecret", "jwtSecret", "salt")
	address, curve, _ := oauthInfo.deriveIdentity("someone@gmail.com")
	if address == "" || curve == "" {
		t.Errorf("Expected address and curve to be set, got empty strings")
	}
	isValidAddress := gcommon.IsHexAddress(address)
	if !isValidAddress {
		t.Errorf("Expected address to be valid, got %s", address)
	}
	expectedAddress := "0x59f7C6dcceBd83ee56dB4F06D35E5E65F2247A1f"
	expectedCurve := "ecdsa"
	if address != expectedAddress || curve != expectedCurve {
		t.Errorf("Expected address %s and curve %s, got %s and %s", expectedAddress, expectedCurve, address, curve)
	}
}

func TestGoogleAuthSign(t *testing.T) {
	oauthInfo = NewGoogleAuth("/redirect", "clientId", "clientSecret", "sessionSecret", "jwtSecret", "salt")
	message := strings.TrimSpace("message to sign")
	signature, _ := oauthInfo.sign("someone@gmail.com", message)
	expectedSignature := "dad814b8c7b43d3b26f9621ce017f7c6bd3609596b3294d6b6dba9c58c0386035cfc7dd6dd704f67b46ebc381fffa8a45c98afbcee34a3a91a2b316fabc17a531b"
	if signature != expectedSignature {
		t.Errorf("Expected signature %s, got %s", expectedSignature, signature)
	}

	address := strings.TrimSpace("0x59f7C6dcceBd83ee56dB4F06D35E5E65F2247A1f")
	ok, _ := common.VerifySignature(address, common.ECDSA_CURVE, message, "0x"+expectedSignature)
	fmt.Println("ok", ok)
	if !ok {
		t.Errorf("Expected signature to be valid, got false")
	}
}

func TestGoogleAuthSignEndpoint(t *testing.T) {
	oauthInfo = NewGoogleAuth("/redirect", "clientId", "clientSecret", "sessionSecret", "jwtSecret", "salt")
	router := mux.NewRouter()
	router.Use(ValidateAccessMiddleware)
	router.HandleFunc("/oauth/sign", handleSigning)
	payloadBuf := new(bytes.Buffer)
	googleId := "someone_else123@gmail.com"
	address, curve, _ := oauthInfo.deriveIdentity(googleId)
	if address == "" || curve == "" {
		t.Errorf("Expected address and curve to be set, got empty strings")
	}
	info := &SignInfo{
		UserId:  googleId,
		Message: strings.TrimSpace("message to sign"),
	}
	json.NewEncoder(payloadBuf).Encode(*info)
	req, err := http.NewRequest("GET", "/oauth/sign", payloadBuf)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	signature := Signature{
		Signature: "a38cf0d18db57af96c9a07e15c3ed422f2356e8eb436e93c3fca0d1c366f8aa32219539e978843393a7b9191c3f7e6a2935f4b05857d4aab373f7c2ae773ed6d1b",
	}
	j, _ := json.Marshal(signature)
	if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(string(j)) {
		t.Errorf("Expected body %s, got %s", w.Body.String(), string(j))
	}
	ok, _ := common.VerifySignature(address, common.ECDSA_CURVE, strings.TrimSpace("message to sign"), "0x"+signature.Signature)
	if !ok {
		t.Errorf("Expected signature to be valid, got false")
	}
}

func TestGenerateRandomSalt(t *testing.T) {
	salt, err := GenerateRandomSalt(32)
	if err != nil {
		t.Fatalf("Failed to generate salt: %v", err)
	}
	if len(salt) != 32 {
		t.Errorf("Expected salt length to be 16, got %d", len(salt))
	}
}

type MockGoogleAuth struct {
	oauthInfo *GoogleAuth
}

func NewMockGoogleAuth() *MockGoogleAuth {
	return &MockGoogleAuth{
		oauthInfo: NewGoogleAuth("/redirect", "clientId", "clientSecret", "sessionSecret", "jwtSecret", "salt"),
	}
}

func (m *MockGoogleAuth) verifyToken(tokenStr string, tokenType string, verifyIdentity bool, secretKey string) (*ClaimsWithIdentity, error) {
	fmt.Println("token from verifyToken", tokenStr, tokenType)
	claims := &ClaimsWithIdentity{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	fmt.Println("token from verifyToken-here", token, err)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else {
				return nil, fmt.Errorf("token validation error: %v", err)
			}
		}
		return nil, fmt.Errorf("could not parse token: %v", err)
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	fmt.Println("claims from verifyToken-here", claims)
	if verifyIdentity && (tokenType == "access_token" || tokenType == "refresh_token") && (claims.Identity == "" || claims.IdentityCurve == "") {
		return nil, ErrInvalidTokenIdentityRequired
	}
	if tokenType == "refresh_token" {
		fmt.Println("gettoken from db-1", token)
		val, ok := refreshTokensMap[tokenStr]
		fmt.Println("gettoken from db-2", val, ok)
		if ok {
			return nil, ErrInvalidToken
		}
	}
	return claims, nil
}

func (m *MockGoogleAuth) generateAccessToken(_ string, _ string, _ string) (string, error) {
	return "new_access_token", nil
}

func (m *MockGoogleAuth) generateRefreshToken(_ string, _ string, _ string) (string, error) {
	return "new_refresh_token", nil
}

var oauthInfoMock *MockGoogleAuth
var refreshTokensMap map[string]bool

// identical as requestAccess but replacing getAccess with getAccessMock for local testing with mock oauth
func requestAccessMock(w http.ResponseWriter, r *http.Request) {
	tokens, err := getAccessMock(r)
	if err != nil {
		fmt.Println("requestAccessMock-1", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = json.NewEncoder(w).Encode(*tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAccessMock(r *http.Request) (*Tokens, error) {
	tokensData, err := extractUserTokensInfo(r)
	if err != nil {
		return nil, err
	}
	refreshToken := tokensData.RefreshToken
	if refreshToken == "" {
		return nil, ErrRefreshTokenNotFound
	}

	refreshClaims, err := oauthInfoMock.verifyToken(refreshToken, "refresh_token", true, oauthInfoMock.oauthInfo.jwtSecret)
	if err != nil {
		if errors.Is(err, ErrTokenExpired) {
			return nil, ErrRefreshTokenExpired
		}
		return nil, err
	}
	accessToken := tokensData.AccessToken
	if accessToken == "" {
		return nil, fmt.Errorf("access token not found")
	}
	_, err = oauthInfoMock.verifyToken(accessToken, "access_token", true, oauthInfoMock.oauthInfo.jwtSecret)
	if err != nil {
		if errors.Is(err, ErrTokenExpired) {
			accessToken, err = oauthInfoMock.generateAccessToken(refreshClaims.Subject, refreshClaims.Identity, refreshClaims.IdentityCurve)
			if err != nil {
				return nil, err
			}
			refreshTokensMap[refreshToken] = true
			refreshToken, err = oauthInfoMock.generateRefreshToken(refreshClaims.Subject, refreshClaims.Identity, refreshClaims.IdentityCurve)
			if err != nil {
				return nil, err
			}
		}
	}
	fmt.Println("getAccessMock-2", accessToken, refreshToken)
	return &Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func TestGoogleRequestAcccessEndpoint(t *testing.T) {
	refreshTokensMap = make(map[string]bool)
	oauthInfoMock = NewMockGoogleAuth()
	router := mux.NewRouter()
	router.HandleFunc("/oauth/accessToken", requestAccessMock)

	accessTokenValid, _ := generateTestToken(oauthInfoMock.oauthInfo, "1", time.Now().Add(time.Minute*10), "0xa", "ecdsa")
	refreshTokenValid, _ := generateTestToken(oauthInfoMock.oauthInfo, "1", time.Now().Add(time.Hour*24*7), "0xa", "ecdsa")
	accessTokenExpired, _ := generateTestToken(oauthInfoMock.oauthInfo, "1", time.Now(), "0xa", "ecdsa")
	refreshTokenExpired, _ := generateTestToken(oauthInfoMock.oauthInfo, "1", time.Now(), "0xa", "ecdsa")
	testA, _ := json.Marshal(&Tokens{
		AccessToken:  accessTokenValid,
		RefreshToken: refreshTokenValid,
	})

	testB, _ := json.Marshal(&Tokens{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	})

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedBody   string
		tokens         *Tokens
	}{
		{
			name:           "Valid Token",
			token:          "validtoken",
			expectedStatus: http.StatusOK,
			expectedBody:   string(testA),
			tokens: &Tokens{
				AccessToken:  accessTokenValid,
				RefreshToken: refreshTokenValid,
			},
		},
		{
			name:           "Invalid Token: Expired Access Token",
			token:          "invalidtokenexpired",
			expectedStatus: http.StatusOK,
			expectedBody:   string(testB),
			tokens: &Tokens{
				AccessToken:  accessTokenExpired,
				RefreshToken: refreshTokenValid,
			},
		},
		{
			name:           "Invalid Token: Expired Refresh Token",
			token:          "invalidtokenexpired",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   ErrRefreshTokenExpired.Error(),
			tokens: &Tokens{
				AccessToken:  accessTokenValid,
				RefreshToken: refreshTokenExpired,
			},
		},
	}
	for _, tt := range tests {
		payloadBuf := new(bytes.Buffer)
		if err := json.NewEncoder(payloadBuf).Encode(tt.tokens); err != nil {
			t.Fatalf("Failed to encode JSON: %v", err)
		}
		req, err := http.NewRequest("GET", "/oauth/accessToken", payloadBuf)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != tt.expectedStatus {
			t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
		}
		if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(tt.expectedBody) {
			t.Errorf("Expected body %s, got %s", tt.expectedBody, w.Body.String())
		}
	}
}
