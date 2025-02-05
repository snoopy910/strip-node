package sequencer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

// roundTripFunc is a helper to create a custom RoundTripper from a function.
// It implements the RoundTripper interface.
type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// fakeOAuthConfig implements a fake oauth2 config for testing.
// It implements both Exchange, Client and AuthCodeURL methods.
type fakeOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Endpoint     oauth2.Endpoint
}

// Exchange simulates exchanging an authorization code for a token.
func (f *fakeOAuthConfig) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	// For testing, always return a valid token when code is "valid"
	if code == "valid" {
		return &oauth2.Token{
			AccessToken:  "fake_access_token",
			TokenType:    "Bearer",
			RefreshToken: "fake_refresh_token",
			Expiry:       time.Now().Add(time.Hour),
		}, nil
	}
	return nil, fmt.Errorf("invalid code")
}

// Client returns an HTTP client that simulates fetching user info from Google.
func (f *fakeOAuthConfig) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			// Simulate Google userinfo response
			userInfo := map[string]interface{}{
				"id":             "123456789",
				"email":          "test@example.com",
				"verified_email": true,
				"name":           "Test User",
				"given_name":     "Test",
				"family_name":    "User",
				"picture":        "https://example.com/photo.jpg",
				"locale":         "en",
			}

			jsonData, err := json.Marshal(userInfo)
			if err != nil {
				return nil, err
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(jsonData)),
				Header:     make(http.Header),
			}, nil
		}),
	}
}

// AuthCodeURL returns a URL to OAuth 2.0 provider's authorization endpoint
// that requests an authorization code.
func (f *fakeOAuthConfig) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	// For testing, just return a fixed URL
	return "https://test.example.com/auth?state=" + state
}

// Note: We assume that the oauthInfo global variable and its type are defined in oauth_handlers.go.
// For testing, we override its fields with fake implementations.

// fake functions for deriveIdentity, generateAccessToken, and generateRefreshToken

// In the test, we assume oauthInfo is mutable. Save the original and restore after the test.

func TestHandleGoogleAuth(t *testing.T) {
	// Save original oauthInfo to restore later (assuming it's a package level variable)
	oldOauthInfo := oauthInfo
	defer func() { oauthInfo = oldOauthInfo }()

	// Create a new test instance of GoogleAuth
	fakeConfig := &fakeOAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"test.scope"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://test.example.com/auth",
			TokenURL: "https://test.example.com/token",
		},
	}
	oauthInfo = &GoogleAuth{
		config:              fakeConfig,
		jwtSecret:           "test-jwt-secret",
		oauthState:          "test-state",
		verifier:            "test-verifier",
		walletSeedSalt:      "test-salt",
		stripchainWalletUrl: "strip-chain-wallet",
	}

	// Prepare a valid POST request with a JSON body containing a valid code
	requestBody := []byte(`{"code": "valid"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/google", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	// Call the handler
	handleGoogleAuth(w, req)

	fmt.Printf("Response Status: %d\n", w.Code)
	fmt.Printf("Response Body: %s\n", w.Body.String())

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Decode the response JSON
	var responseData struct {
		Tokens struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			IdToken      string `json:"id_token"`
		} `json:"tokens"`
		Wallet struct {
			Identity      string `json:"identity"`
			IdentityCurve string `json:"identity_curve"`
		} `json:"wallet"`
		User struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
	}

	bodyBytes, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		t.Fatalf("failed to unmarshal response JSON: %v", err)
	}

	// Validate the tokens and user info
	if responseData.Tokens.AccessToken == "" {
		t.Errorf("expected non-empty access token, got empty string")
	}
	if responseData.Tokens.RefreshToken == "" {
		t.Errorf("expected non-empty refresh token, got empty string")
	}

	claims := &ClaimsWithIdentity{}
	_, err = jwt.ParseWithClaims(responseData.Tokens.IdToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(oauthInfo.jwtSecret), nil
	})
	if err != nil {
		t.Errorf("failed to parse id token: %v", err)
	}

	if claims.RegisteredClaims.Subject != "123456789" {
		t.Errorf("expected id token '123456789', got '%s'", claims.RegisteredClaims.Subject)
	}
	if responseData.Wallet.Identity != "0x623e01B359e01549Ffd21E7b7aC7853afc227803" {
		t.Errorf("expected wallet identity '0x623e01B359e01549Ffd21E7b7aC7853afc227803', got '%s'", responseData.Wallet.Identity)
	}
	if responseData.Wallet.IdentityCurve != "ecdsa" {
		t.Errorf("expected wallet identity curve to be 'ecdsa', got '%s'", responseData.Wallet.IdentityCurve)
	}
	if responseData.User.ID != "123456789" {
		t.Errorf("expected user id '123456789', got '%s'", responseData.User.ID)
	}
	if responseData.User.Name != "Test User" {
		t.Errorf("expected user name 'Test User', got '%s'", responseData.User.Name)
	}
}

var mockRefreshTokens = make(map[string]*RefreshTokenSchema)

func mockGetRefreshToken(token string) (*RefreshTokenSchema, error) {
	if rt, ok := mockRefreshTokens[token]; ok {
		return rt, nil
	}
	return nil, nil
}

func TestHandleRefreshToken(t *testing.T) {
	// Save original functions to restore later
	oldOauthInfo := oauthInfo
	oldGetRefreshToken := GetRefreshToken
	defer func() {
		oauthInfo = oldOauthInfo
		GetRefreshToken = oldGetRefreshToken
	}()

	// Set up mock
	GetRefreshToken = mockGetRefreshToken

	// Create a fake OAuth config template
	fakeConfig := &fakeOAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"test.scope"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://test.example.com/auth",
			TokenURL: "https://test.example.com/token",
		},
	}

	// Create a test instance template
	testAuth := &GoogleAuth{
		config:              fakeConfig, // Pass pointer to fakeOAuthConfig
		jwtSecret:           "test-jwt-secret",
		oauthState:          "test-state",
		verifier:            "test-verifier",
		walletSeedSalt:      "test-salt",
		stripchainWalletUrl: "test-wallet-seed",
	}

	// Generate a valid refresh token for testing
	refreshToken, err := testAuth.generateRefreshToken("123456789", "0x623e01B359e01549Ffd21E7b7aC7853afc227803", "ecdsa")
	if err != nil {
		t.Fatalf("Failed to generate test refresh token: %v", err)
	}

	tests := []struct {
		name           string
		method         string
		refreshToken   string
		expectedStatus int
		expectNewToken bool
	}{
		{
			name:           "Valid refresh token",
			method:         http.MethodPost,
			refreshToken:   refreshToken,
			expectedStatus: http.StatusOK,
			expectNewToken: true,
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			refreshToken:   "valid-token",
			expectedStatus: http.StatusMethodNotAllowed,
			expectNewToken: false,
		},
		{
			name:           "Empty refresh token",
			method:         http.MethodPost,
			refreshToken:   "",
			expectedStatus: http.StatusBadRequest,
			expectNewToken: false,
		},
		{
			name:           "Invalid refresh token format",
			method:         http.MethodPost,
			refreshToken:   "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectNewToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear mock refresh tokens for each test
			mockRefreshTokens = make(map[string]*RefreshTokenSchema)

			// Create a fresh instance of GoogleAuth for each test
			oauthInfo = &GoogleAuth{
				config:              fakeConfig, // Pass pointer to fakeOAuthConfig
				jwtSecret:           "test-jwt-secret",
				oauthState:          "test-state",
				verifier:            "test-verifier",
				walletSeedSalt:      "test-salt",
				stripchainWalletUrl: "test-wallet-seed",
			}

			// Create request body
			body := struct {
				RefreshToken string `json:"refresh_token"`
			}{
				RefreshToken: tt.refreshToken,
			}
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			fmt.Printf("Test case '%s': oauthInfo pointer before request: %p\n", tt.name, oauthInfo)

			// Create request
			req := httptest.NewRequest(tt.method, "/auth/refresh", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			// Call the handler
			handleRefreshTokenMock(w, req)
			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			// For successful cases, verify the response
			if tt.expectNewToken {
				var response struct {
					AccessToken string `json:"access_token"`
					TokenType   string `json:"token_type"`
					ExpiresIn   int    `json:"expires_in"`
				}

				err = json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if response.AccessToken == "" {
					t.Error("expected non-empty access token")
				}
				if response.TokenType != "Bearer" {
					t.Errorf("expected token type 'Bearer', got '%s'", response.TokenType)
				}
				if response.ExpiresIn != 600 {
					t.Errorf("expected expires_in 600, got %d", response.ExpiresIn)
				}
			}
		})
	}
}

func handleRefreshTokenMock(w http.ResponseWriter, r *http.Request) {
	refreshTokensMap = make(map[string]bool)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.RefreshToken == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Verify the refresh token and extract claims
	claims, err := oauthInfo.verifyToken(body.RefreshToken, "refresh_token", true, oauthInfo.jwtSecret)
	if err != nil {
		if errors.Is(err, ErrTokenExpired) {
			http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		} else {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		}
		return
	}

	// Generate new access token using claims from refresh token
	accessToken, err := oauthInfo.generateAccessToken(claims.Subject, claims.Identity, claims.IdentityCurve)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate access token: %v", err), http.StatusInternalServerError)
		return
	}

	refreshTokensMap[body.RefreshToken] = true
	refreshToken, err := oauthInfo.generateRefreshToken(claims.Subject, claims.Identity, claims.IdentityCurve)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate refresh token: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    600, // 10 minutes in seconds
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("handleRefreshToken: failed to encode response: %v\n", err)
	}
}
