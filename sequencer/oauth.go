package sequencer

import (
	"context"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type contextKey string

const (
	// Key for the ID
	identityAccess          contextKey = "identityAccess"
	stripchainGoogleSession string     = "stripchain-google-session"
	tokenIssuer             string     = "StripChain"
)

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Identity      string `json:"identity"`
	IdentityCurve string `json:"identity_curve"`
}

type M map[string]interface{}

type OAuthParameters struct {
	config     *oauth2.Config
	session    *sessions.CookieStore
	jwtSecret  string
	oauthState string
	verifier   string
	message    string
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
}

type UserTokensInfo struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	Expiry       uint64 `json:"expiry"`
	Signature    string `json:"signature"`
}

type ClaimsWithIdentity struct {
	Email             string `json:"email"`
	Name              string `json:"name"`
	Identity          string `json:"identity"`
	IdentityCurve     string `json:"identity_curve"`
	UserSignedMessage string `json:"user_signed_message"`
	jwt.RegisteredClaims
}

type IdentityAccess struct {
	Identity      string `json:"identity"`
	IdentityCurve string `json:"identity_curve"`
}

var (
	ErrTokenExpired                 = errors.New("token is expired")
	ErrInvalidToken                 = errors.New("invalid token")
	ErrAuthorization                = errors.New("authorization failed")
	ErrInvalidTokenIdentityRequired = errors.New("invalid token: identity and identity curve are required")
	ErrRefreshTokenExpired          = errors.New("refresh token is expired")
	ErrInvalidTokenId               = errors.New("invalid token id")
	ErrNotAuthenticated             = errors.New("not authenticated")
)

func initializeGoogleOauth(redirectUrl string, clientId string, clientSecret string, sessionSecret string, jwtSecret string, message string) *OAuthParameters {

	googleOauthConfig := &oauth2.Config{
		RedirectURL:  redirectUrl,
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	// use PKCE to protect against CSRF attacks
	verifier := oauth2.GenerateVerifier()

	sessionStore := sessions.NewCookieStore([]byte(sessionSecret))

	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30, // 30 days
		HttpOnly: true,
		// Secure:   true,                   // Ensures cookie is sent over HTTPS
		Secure:   false,                   // for localhost testing
		SameSite: http.SameSiteStrictMode, // CSRF protection
	}

	gob.Register(&UserInfo{})
	gob.Register(&M{})

	// State string for CSRF protection
	oauthState := generateState()
	return &OAuthParameters{googleOauthConfig, sessionStore, jwtSecret, oauthState, verifier, message}

}
func generateIdToken(user UserInfo, identity string, identityCurve string, signedMessage string) (string, error) {
	claims := ClaimsWithIdentity{
		Email:             user.Email,
		Name:              user.Name,
		Identity:          identity,
		IdentityCurve:     identityCurve,
		UserSignedMessage: signedMessage,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			Issuer:    tokenIssuer,
			Subject:   user.ID,
		},
	}

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	return token.SignedString([]byte(oauthInfo.jwtSecret))
}

// GenerateAccessToken creates a JWT access token
func generateAccessToken(userId string, identity string, identityCurve string) (string, error) {
	claims := ClaimsWithIdentity{
		Identity:      identity,
		IdentityCurve: identityCurve,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)), // Token expires after 10min
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

func generateRefreshToken(userId string, identity string, identityCurve string) (string, error) {
	claims := ClaimsWithIdentity{
		Identity:      identity,
		IdentityCurve: identityCurve,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Token expires after 7 days
			Issuer:    tokenIssuer,
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(oauthInfo.jwtSecret))
}

func login(w http.ResponseWriter, r *http.Request) {
	// adding nonce and hd
	url := oauthInfo.config.AuthCodeURL(oauthInfo.oauthState, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(oauthInfo.verifier))

	// Redirect user to Google's consent page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {

	tokens, err := getAccess(r)
	fmt.Println("err handle access", err, tokens)
	var idToken string
	var accessToken string
	var refreshToken string

	if tokens != nil {
		idToken = tokens.IdToken
		accessToken = tokens.AccessToken
		refreshToken = tokens.RefreshToken
	} else {
		fmt.Println("acces token or refresh token empty")
		if r.FormValue("state") != oauthInfo.oauthState {
			http.Error(w, "invalid state parameter", http.StatusBadRequest)
			return
		}

		// Verify the code exists
		code := r.FormValue("code")
		if code == "" {
			http.Error(w, "code parameter not found", http.StatusBadRequest)
			return
		}

		// Exchange the authorization code for an access token
		token, err := oauthInfo.config.Exchange(context.Background(), code, oauth2.VerifierOption(oauthInfo.verifier))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Exchange token", token)
		// Get user information from Google
		client := oauthInfo.config.Client(context.Background(), token)
		userInfoResponse, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer userInfoResponse.Body.Close()

		var userInfo UserInfo
		err = json.NewDecoder(userInfoResponse.Body).Decode(&userInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("userInfo", userInfo)
		// generate the idToken here without identity and identityCurve
		session, err := oauthInfo.session.Get(r, stripchainGoogleSession)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["authenticated"] = true
		session.Values["user"] = &userInfo
		session.Save(r, w)
		fmt.Println("session", session.Values["user"])

		// Generate a JWT id token
		idToken, err = generateIdToken(userInfo, "", "", "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Generate a JWT access token
		accessToken, err = generateAccessToken(userInfo.ID, "", "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Generate a JWT access token
		refreshToken, err = generateRefreshToken(userInfo.ID, "", "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	SetTokenCookie(w, idToken, "id_token")
	SetTokenCookie(w, accessToken, "access_token")
	SetTokenCookie(w, refreshToken, "refresh_token")

	tokens = &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IdToken:      idToken,
	}
	json.NewEncoder(w).Encode(*tokens)
	fmt.Println("tokens", tokens)
}

func handleIdentityVerification(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleIdentityVerification-1")
	fmt.Println("r.Body", r.Body)
	tokensData, err := extractUserTokensInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("handleIdentityVerification-2")
	fmt.Println("tokensData handleIdentityVerification", tokensData)
	identity := r.URL.Query().Get("identity")
	identityCurve := r.URL.Query().Get("identityCurve")
	signature := tokensData.Signature
	message := oauthInfo.message
	fmt.Println("handleIdentityVerification-3")
	userInfo := &UserInfo{
		ID:            tokensData.ID,
		Name:          tokensData.Name,
		Email:         tokensData.Email,
		Identity:      identity,
		IdentityCurve: identityCurve,
	}

	verified, err := common.VerifySignature(
		identity,
		identityCurve,
		message,
		signature,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("handleIdentityVerification-4")
	//get the old access and id token
	if verified {
		idToken, _ := generateIdToken(*userInfo, identity, identityCurve, signature)
		accessToken, _ := generateAccessToken(userInfo.ID, identity, identityCurve)
		refreshToken, _ := generateRefreshToken(userInfo.ID, identity, identityCurve)

		tokens := Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			IdToken:      idToken,
		}
		err = json.NewEncoder(w).Encode(tokens)
		fmt.Println("tokens handleIdentityVerification-5", tokens)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
	}
}

func requestAccess(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request access-1")
	tokens, err := getAccess(r)
	if err != nil {
		if errors.Is(err, ErrNotAuthenticated) {
			http.Error(w, ErrNotAuthenticated.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(*tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("request access", w)
}

func getAccess(r *http.Request) (*Tokens, error) {
	fmt.Println("get access-1")
	tokensData, err := extractUserTokensInfo(r)
	if err != nil {
		return nil, err
	}
	refreshToken := tokensData.RefreshToken
	fmt.Println("get access-2", refreshToken)
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token not found")
	}

	refreshClaims, err := verifyToken(refreshToken, "refresh_token", true, oauthInfo.jwtSecret)
	if err != nil {
		if errors.Is(err, ErrTokenExpired) {
			return nil, ErrRefreshTokenExpired
		}
		return nil, err
	}
	fmt.Println("get access-3")
	accessToken := tokensData.AccessToken
	if accessToken == "" {
		return nil, fmt.Errorf("access token not found")
	}
	fmt.Println("get access-4", accessToken)
	_, err = verifyToken(accessToken, "access_token", true, oauthInfo.jwtSecret)
	if err != nil {
		if errors.Is(err, ErrTokenExpired) {
			fmt.Println("handle access-5")
			accessToken, err = generateAccessToken(refreshClaims.Subject, refreshClaims.Identity, refreshClaims.IdentityCurve)
			if err != nil {
				return nil, err
			}
			AddRefreshToken(refreshToken, true)
			refreshToken, err = generateRefreshToken(refreshClaims.Subject, refreshClaims.Identity, refreshClaims.IdentityCurve)
			if err != nil {
				return nil, err
			}
		}
	}
	fmt.Println("handle tokens", accessToken, refreshToken)
	return &Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := oauthInfo.session.Get(r, stripchainGoogleSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// check if authenticated
	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		http.Error(w, ErrNotAuthenticated.Error(), http.StatusUnauthorized)
		return
	}

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Values["user"] = nil
	session.Options.MaxAge = -1 // Delete session
	session.Save(r, w)
	SetTokenCookie(w, "", "access_token")
	SetTokenCookie(w, "", "id_token")
	SetTokenCookie(w, "", "refresh_token")
}

func ValidateAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/login" || r.URL.Path == "/oauth/callback" || r.URL.Path == "/oauth/verifySignature" || r.URL.Path == "/oauth/accessToken" || r.URL.Path == "/oauth/logout" {
			next.ServeHTTP(w, r)
			return
		}
		auth := r.URL.Query().Get("auth")
		if auth != "oauth" {
			next.ServeHTTP(w, r)
			return
		}
		log.Printf("Auth middleware triggered for: %s\n", r.URL.Path)
		tokensData, err := extractUserTokensInfo(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusBadRequest)
			return
		}
		claims, err := verifyToken(tokensData.AccessToken, "access_token", true, oauthInfo.jwtSecret)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("claims handle access", claims)
		scId := IdentityAccess{
			Identity:      claims.Identity,
			IdentityCurve: claims.IdentityCurve,
		}

		ctx := context.WithValue(r.Context(), identityAccess, &scId)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func SetTokenCookie(w http.ResponseWriter, token string, tokenType string) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenType,
		Value:    token,
		HttpOnly: true, // Prevents JavaScript access
		Secure:   true, // Ensures cookie is sent over HTTPS
		// Secure:   false,                         // just for localhost testing
		SameSite: http.SameSiteStrictMode,       // CSRF protection
		Path:     "/",                           // Cookie path
		Expires:  time.Now().Add(time.Hour * 1), // Matches token expiration
	})
}

func verifyToken(tokenStr string, tokenType string, verifyIdentity bool, secretKey string) (*ClaimsWithIdentity, error) {
	claims := &ClaimsWithIdentity{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) { //interface to define
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	fmt.Println("token from verifyToken", token)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else {
				return nil, fmt.Errorf("token validation error: %v", err)
			}
		}
		return nil, fmt.Errorf("could not parse token: %v", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	fmt.Println("claims from verifyToken", claims)
	if verifyIdentity && (tokenType == "access_token" || tokenType == "refresh_token") && (claims.Identity == "" || claims.IdentityCurve == "") {
		return nil, ErrInvalidTokenIdentityRequired
	}
	if tokenType == "refresh_token" {
		token, _ := GetRefreshToken(tokenStr)
		fmt.Println("gettoken from db", token)
		if token != nil {
			return nil, ErrInvalidToken
		}
	}
	return claims, nil
}

func verifyIdentity(r *http.Request) (bool, error) {
	id, _ := r.Context().Value(identityAccess).(*IdentityAccess)
	identity := r.URL.Query().Get("identity")
	identityCurve := r.URL.Query().Get("identityCurve")
	if id != nil && (identity != id.Identity || identityCurve != id.IdentityCurve) {
		return false, errors.New("mismatch identity between url and token")
	}
	return true, nil
}

func extractUserTokensInfo(r *http.Request) (*UserTokensInfo, error) {
	var tokensData *UserTokensInfo
	err := json.NewDecoder(r.Body).Decode(&tokensData)
	if err != nil {
		return nil, err
	}
	return tokensData, nil
}

func generateState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "strip_chain"
	}
	return base64.URLEncoding.EncodeToString(b)
}

// https://www.unicorn.studio/embed/SaCYz48FXaFwo5ifY36I?preview=true
// http://localhost/oauth/verifySignature?identity=0x76C09917EF1A6E885affCb8B14c0E09df271F393&identityCurve=ecdsa&signature=0x2b7fe067cf63bbfff8df636002eb6f71f4610b3958e75e759a2f6633c75c0a03147da8cbcbf788174b3c771e829ac5089a3cbc6511f9b76f2575ba2dd64dbfdd1b
// http://localhost/oauth/verifySignature?identity=0x2c8251052663244f37BAc7Bde1C6Cb02bBffff93&identityCurve=ecdsa&signature=0xbc490764bf20e3e55f100555d9e1fd84c41fa658850332c388cd8f40d554983b68db49713591539f24cf7d4bcb68f17c89930c314ef7581cf7805b247683b8561b
// http://localhost/createWallet?identity=0x2c8251052663244f37BAc7Bde1C6Cb02bBffff93&identityCurve=ecdsa&auth=oauth
// http://localhost/createWallet?identity=0x2c8251052663244f37BAc7Bde1C6Cb02bBffff93&identityCurve=ecdsa&auth=oauth
