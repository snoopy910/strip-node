package sequencer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type SignInfo struct {
	UserId  string `json:"userId"`
	Message string `json:"message"`
}

type Signature struct {
	Signature string `json:"signature"`
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

		// Derive identity from Google ID
		identity, identityCurve, err := oauthInfo.deriveIdentity(userInfo.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("userInfo", userInfo)
		// generate the idToken
		idToken, err = oauthInfo.generateIdToken(userInfo, identity, identityCurve)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Generate a JWT access token
		accessToken, err = oauthInfo.generateAccessToken(userInfo.ID, identity, identityCurve)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Generate a JWT refresh
		refreshToken, err = oauthInfo.generateRefreshToken(userInfo.ID, "", "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tokens = &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IdToken:      idToken,
	}
	json.NewEncoder(w).Encode(*tokens)
	fmt.Println("tokens", tokens)
}

func requestAccess(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request access-1")
	tokens, err := getAccess(r)
	if err != nil {
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
		return nil, ErrRefreshTokenNotFound
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
			accessToken, err = oauthInfo.generateAccessToken(refreshClaims.Subject, refreshClaims.Identity, refreshClaims.IdentityCurve)
			if err != nil {
				return nil, err
			}
			AddRefreshToken(refreshToken, true)
			refreshToken, err = oauthInfo.generateRefreshToken(refreshClaims.Subject, refreshClaims.Identity, refreshClaims.IdentityCurve)
			if err != nil {
				return nil, err
			}
		}
	}
	fmt.Println("handle tokens", accessToken, refreshToken)
	return &Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func handleSigning(w http.ResponseWriter, r *http.Request) {
	var data *SignInfo
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signature, err := oauthInfo.sign(data.UserId, data.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Signature{
		Signature: signature,
	}

	json.NewEncoder(w).Encode(response)
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
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		claims, err := verifyToken(tokensData.AccessToken, "access_token", true, oauthInfo.jwtSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
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

// https://www.unicorn.studio/embed/SaCYz48FXaFwo5ifY36I?preview=true
// http://localhost/oauth/verifySignature?identity=0x76C09917EF1A6E885affCb8B14c0E09df271F393&identityCurve=ecdsa&signature=0x2b7fe067cf63bbfff8df636002eb6f71f4610b3958e75e759a2f6633c75c0a03147da8cbcbf788174b3c771e829ac5089a3cbc6511f9b76f2575ba2dd64dbfdd1b
// http://localhost/oauth/verifySignature?identity=0x2c8251052663244f37BAc7Bde1C6Cb02bBffff93&identityCurve=ecdsa&signature=0xbc490764bf20e3e55f100555d9e1fd84c41fa658850332c388cd8f40d554983b68db49713591539f24cf7d4bcb68f17c89930c314ef7581cf7805b247683b8561b
// http://localhost/createWallet?identity=0x2c8251052663244f37BAc7Bde1C6Cb02bBffff93&identityCurve=ecdsa&auth=oauth
// http://localhost/createWallet?identity=0x2c8251052663244f37BAc7Bde1C6Cb02bBffff93&identityCurve=ecdsa&auth=oauth
