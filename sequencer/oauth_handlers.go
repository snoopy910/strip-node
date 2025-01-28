package sequencer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/StripChain/strip-node/common"
	"golang.org/x/oauth2"
)

var oauthInfo *GoogleAuth

type SignInfo struct {
	UserId  string `json:"userId"`
	Message string `json:"message"`
}

type Signature struct {
	Signature string `json:"signature"`
}

type SignatureInfo struct {
	UserId    string `json:"userId"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

type CallbackInfo struct {
	Tokens *Tokens         `json:"tokens"`
	Wallet *IdentityAccess `json:"wallet"`
}

func login(w http.ResponseWriter, r *http.Request) {
	// adding nonce and hd
	url := oauthInfo.config.AuthCodeURL(oauthInfo.oauthState, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(oauthInfo.verifier))

	// Redirect user to Google's consent page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle google callback")
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

	// Derive identity from Google ID
	identity, identityCurve, err := oauthInfo.deriveIdentity(userInfo.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("userInfo", userInfo)
	// generate the idToken
	idToken, err := oauthInfo.generateIdToken(userInfo, identity, identityCurve)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a JWT access token
	accessToken, err := oauthInfo.generateAccessToken(userInfo.ID, identity, identityCurve)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a JWT refresh
	refreshToken, err := oauthInfo.generateRefreshToken(userInfo.ID, identity, identityCurve)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokens := &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IdToken:      idToken,
	}

	id := &IdentityAccess{
		Identity:      identity,
		IdentityCurve: identityCurve,
	}

	response := &CallbackInfo{
		Tokens: tokens,
		Wallet: id,
	}

	json.NewEncoder(w).Encode(*response)
}

func requestAccess(w http.ResponseWriter, r *http.Request) {
	tokens, err := getAccess(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = json.NewEncoder(w).Encode(*tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAccess(r *http.Request) (*Tokens, error) {
	tokensData, err := extractUserTokensInfo(r)
	if err != nil {
		return nil, err
	}
	refreshToken := tokensData.RefreshToken
	if refreshToken == "" {
		return nil, ErrRefreshTokenNotFound
	}

	refreshClaims, err := oauthInfo.verifyToken(refreshToken, "refresh_token", true, oauthInfo.jwtSecret)
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

	_, err = oauthInfo.verifyToken(accessToken, "access_token", true, oauthInfo.jwtSecret)
	if err != nil {
		if errors.Is(err, ErrTokenExpired) {
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

func handleVerifySignature(w http.ResponseWriter, r *http.Request) {
	var data *SignatureInfo
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	identity, identityCurve, err := oauthInfo.deriveIdentity(data.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	isValid, err := common.VerifySignature(identity, identityCurve, data.Message, data.Signature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValid {
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	response := Signature{
		Signature: data.Signature,
	}

	json.NewEncoder(w).Encode(response)
}

func ValidateAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/login" || r.URL.Path == "/oauth/callback" || r.URL.Path == "/oauth/sign" || r.URL.Path == "/oauth/accessToken" || r.URL.Path == "/oauth/logout" {
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
		claims, err := oauthInfo.verifyToken(tokensData.AccessToken, "access_token", true, oauthInfo.jwtSecret)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		fmt.Println("claims middleware access", claims)
		scId := IdentityAccess{
			Identity:      claims.Identity,
			IdentityCurve: claims.IdentityCurve,
		}

		ctx := context.WithValue(r.Context(), identityAccess, &scId)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
