package sequencer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserInfo struct {
	Email         string `json:"email"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Identity      string `json:"identity"`
	IdentityCurve string `json:"identity_curve"`
}

type OAuthParameters struct {
	config     *oauth2.Config
	session    *sessions.CookieStore
	jwtSecret  string
	oauthState string
	verifier   string
	message    string
}

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

	// State string for CSRF protection
	oauthState := generateState()
	return &OAuthParameters{googleOauthConfig, sessionStore, jwtSecret, oauthState, verifier, message}

}
func generateIdToken(user UserInfo, identity string, identityCurve string, signedMessage string) (string, error) {
	claims := jwt.MapClaims{}
	claims["sub"] = user.ID
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["identity"] = identity
	claims["identityCurve"] = identityCurve
	claims["userSignedMessage"] = signedMessage

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	return token.SignedString([]byte(oauthInfo.jwtSecret))
}

// GenerateAccessToken creates a JWT access token
func generateAccessToken(userId string, identity string, identityCurve string) (string, error) {
	// Define token claims
	claims := jwt.MapClaims{}
	claims["sub"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Token expires after 1 hour
	claims["iat"] = time.Now().Unix()
	claims["iss"] = "StripChain"
	claims["identity"] = identity
	claims["identityCurve"] = identityCurve

	// Create the token using HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	return token.SignedString([]byte(oauthInfo.jwtSecret))
}

func login(w http.ResponseWriter, r *http.Request) {
	// adding nonce and hd
	url := oauthInfo.config.AuthCodeURL(oauthInfo.oauthState, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(oauthInfo.verifier))

	// Redirect user to Google's consent page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify the state parameter
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

	// generate the idToken here without identity and identityCurve
	session, err := oauthInfo.session.Get(r, "session-name")
	fmt.Println("session", session)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["authenticated"] = true
	session.Values["user"] = userInfo
	session.Save(r, w)

	// Generate a JWT id token
	idToken, err := generateIdToken(userInfo, "", "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a JWT access token
	accessToken, err := generateAccessToken(userInfo.ID, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("accessToken", accessToken)
	fmt.Println("idToken", idToken)
	SetTokenCookie(w, idToken, "id_token")
	SetTokenCookie(w, accessToken, "access_token")

	// Redirect user to the home page
	fmt.Println("w callback response 1", w.Header().Get("Set-Cookie"))
	fmt.Println("w callback response 2", w.Header().Values("access_token"))
	fmt.Println("w callback response 3", w.Header().Values("id_token"))
	fmt.Println("w callback response 4", w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func SetTokenCookie(w http.ResponseWriter, token string, tokenType string) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenType,
		Value:    token,
		HttpOnly: true,                          // Prevents JavaScript access
		Secure:   true,                          // Ensures cookie is sent over HTTPS
		SameSite: http.SameSiteStrictMode,       // CSRF protection
		Path:     "/",                           // Cookie path
		Expires:  time.Now().Add(time.Hour * 1), // Matches token expiration
	})
}

func verifyToken(tokenStr string, secretKey string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) { // interface to define
		return secretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func generateState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "strip_chain"
	}
	return base64.URLEncoding.EncodeToString(b)
}

func handleIdentityVerification(w http.ResponseWriter, r *http.Request) {
	// get and verify access token in authorization header
	// verifyToken(tokenStr string, secretKey string)
	prefix := "Bearer "
	authHeader := r.Header.Get("Authorization")
	accessToken := strings.TrimPrefix(authHeader, prefix)

	if authHeader == "" || accessToken == authHeader {
		http.Error(w, "Authentication header not present or malformed", http.StatusInternalServerError)
		return
	}

	err := verifyToken(accessToken, oauthInfo.jwtSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userId := r.URL.Query().Get("userId")
	userName := r.URL.Query().Get("Name")
	userEmail := r.URL.Query().Get("Email")
	identity := r.URL.Query().Get("identity")
	identityCurve := r.URL.Query().Get("identityCurve")

	signature := r.URL.Query().Get("signature")
	message := oauthInfo.message
	session, err := oauthInfo.session.Get(r, "session-name")
	fmt.Println("session", session.Values["authenticated"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userInfo := &UserInfo{
		ID:    userId,
		Name:  userName,
		Email: userEmail,
	}
	_, err = common.VerifySignature(
		identity,
		identityCurve,
		message,
		signature,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//get the old access and id token
	idToken, _ := generateIdToken(*userInfo, identity, identityCurve, signature)
	accessToken, _ = generateAccessToken(userInfo.ID, identity, identityCurve)
	SetTokenCookie(w, accessToken, "access_token")
	SetTokenCookie(w, idToken, "id_token")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// user sign a message
// message verified
// jwt stored and sent to the client
// use the jwt in the header authorization bearer
// adding nonce --> replay attack
// adding oauth protected to existing routes
// https://www.unicorn.studio/embed/SaCYz48FXaFwo5ifY36I?preview=true
