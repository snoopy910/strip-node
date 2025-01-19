package sequencer

import (
	"context"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
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

	fmt.Println("userInfo", userInfo)
	// generate the idToken here without identity and identityCurve
	session, err := oauthInfo.session.Get(r, "stripchain-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["authenticated"] = true
	session.Values["user"] = &userInfo
	session.Save(r, w)
	fmt.Println("session", session.Values["user"])
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

	SetTokenCookie(w, idToken, "id_token")
	SetTokenCookie(w, accessToken, "access_token")

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: "",
		IdToken:      idToken,
	}
	json.NewEncoder(w).Encode(tokens)
	fmt.Println("tokens", tokens)
	// Redirect user to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleIdentityVerification(w http.ResponseWriter, r *http.Request) {
	// get and verify access token in authorization header
	// verifyToken(tokenStr string, secretKey string)

	accessCookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	accessToken := accessCookie.Value
	idCookie, err := r.Cookie("id_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	idToken := idCookie.Value
	fmt.Println("accessToken", accessToken)
	fmt.Println("idToken", idToken)

	err = verifyToken(accessToken, oauthInfo.jwtSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	identity := r.URL.Query().Get("identity")
	identityCurve := r.URL.Query().Get("identityCurve")
	signature := r.URL.Query().Get("signature")
	message := oauthInfo.message
	session, err := oauthInfo.session.Get(r, "stripchain-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// check if authenticated
	if !session.Values["authenticated"].(bool) {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}
	userInfo := &UserInfo{
		ID:    session.Values["user"].(*UserInfo).ID,
		Name:  session.Values["user"].(*UserInfo).Name,
		Email: session.Values["user"].(*UserInfo).Email,
	}
	fmt.Println("userInfo", userInfo)
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
	//get the old access and id token
	if verified {
		idToken, _ = generateIdToken(*userInfo, identity, identityCurve, signature)
		accessToken, _ = generateAccessToken(userInfo.ID, identity, identityCurve)
		SetTokenCookie(w, accessToken, "access_token")
		SetTokenCookie(w, idToken, "id_token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
	}
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

func verifyToken(tokenStr string, secretKey string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) { //interface to define
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
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

// user sign a message
// message verified
// jwt stored and sent to the client
// use the jwt in the header authorization bearer
// adding nonce --> replay attack
// adding oauth protected to existing routes
// https://www.unicorn.studio/embed/SaCYz48FXaFwo5ifY36I?preview=true
// http://localhost/oauth/verifySignature?identity=0x76C09917EF1A6E885affCb8B14c0E09df271F393&identityCurve=ecdsa&signature=0x2b7fe067cf63bbfff8df636002eb6f71f4610b3958e75e759a2f6633c75c0a03147da8cbcbf788174b3c771e829ac5089a3cbc6511f9b76f2575ba2dd64dbfdd1b
// http://localhost/oauth/verifySignature?identity=0x2c8251052663244f37BAc7Bde1C6Cb02bBffff93&identityCurve=ecdsa&signature=0xbc490764bf20e3e55f100555d9e1fd84c41fa658850332c388cd8f40d554983b68db49713591539f24cf7d4bcb68f17c89930c314ef7581cf7805b247683b8561b
