package sequencer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

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

func handleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting handleGoogleAuth")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Code string `json:"code"`
		// State string `json:"state"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.Code == "" {
		fmt.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received code: %s\n", body.Code)

	// // Validate state parameter
	// if body.State != oauthInfo.oauthState {
	// 	http.Error(w, "Invalid state parameter", http.StatusBadRequest)
	// 	return
	// }

	// Exchange the authorization code for tokens
	token, err := oauthInfo.config.Exchange(r.Context(), body.Code)
	if err != nil {
		fmt.Printf("Error exchanging code for token: %v\n", err)
		http.Error(w, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Successfully exchanged code for token: %+v\n", token)

	// Get user information from Google
	client := oauthInfo.config.Client(r.Context(), token)
	userInfoResponse, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		fmt.Printf("Error getting user info: %v\n", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer userInfoResponse.Body.Close()

	var userInfo UserInfo
	err = json.NewDecoder(userInfoResponse.Body).Decode(&userInfo)
	if err != nil {
		fmt.Printf("Error decoding user info: %v\n", err)
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Derive identity from Google ID
	identity, identityCurve, err := oauthInfo.deriveIdentity(userInfo.ID)
	if err != nil {
		fmt.Printf("Error deriving identity: %v\n", err)
		http.Error(w, "Failed to derive identity", http.StatusInternalServerError)
		return
	}

	// Generate our custom tokens
	idToken := userInfo.ID
	// idToken, err := oauthInfo.generateIdToken(userInfo, identity, identityCurve)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Failed to generate ID token: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	accessToken, err := oauthInfo.generateAccessToken(userInfo.ID, identity, identityCurve)
	if err != nil {
		fmt.Printf("Error generating access token: %v\n", err)
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := oauthInfo.generateRefreshToken(userInfo.ID, identity, identityCurve)
	if err != nil {
		fmt.Printf("Error generating refresh token: %v\n", err)
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	response := struct {
		Tokens *Tokens         `json:"tokens"`
		Wallet *IdentityAccess `json:"wallet"`
		User   *UserInfo       `json:"user"`
	}{
		Tokens: &Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			IdToken:      idToken,
		},
		Wallet: &IdentityAccess{
			Identity:      identity,
			IdentityCurve: identityCurve,
		},
		User: &userInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleRefreshToken(w http.ResponseWriter, r *http.Request) {
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

	AddRefreshToken(body.RefreshToken, true)
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

	// Generate tokens
	idToken, err := oauthInfo.generateIdToken(userInfo, identity, identityCurve)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, err := oauthInfo.generateAccessToken(userInfo.ID, identity, identityCurve)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	info := &CallbackInfo{
		Tokens: tokens,
		Wallet: id,
	}

	ctx := context.WithValue(r.Context(), tokensCallbackInfoKey, info)
	r = r.WithContext(ctx)

	// Redirect to the wallet application with tokens in URL fragment
	//"http://localhost:5173/auth/callback#access_token=%s&refresh_token=%s&id_token=%s"
	redirectURL := fmt.Sprintf("%s/auth/callback#access_token=%s&refresh_token=%s&id_token=%s",
		oauthInfo.stripchainWalletUrl, accessToken, refreshToken, idToken)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle redirect")
	// Get the access token from the request
	http.Redirect(w, r, oauthInfo.stripchainWalletUrl, http.StatusMovedPermanently)
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
	signature := data.Signature
	if strings.TrimSpace(signature) == "" {
		http.Error(w, "empty signature", http.StatusBadRequest)
		return
	}
	signature = strings.TrimSpace(signature)
	if !strings.HasPrefix(signature, "0x") {
		signature = "0x" + signature
	}

	sig, err := hex.DecodeString(signature[2:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(sig) != 65 {
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	identity, identityCurve, err := oauthInfo.deriveIdentity(data.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	isValid, err := common.VerifySignature(identity, identityCurve, data.Message, signature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValid {
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	response := Signature{
		Signature: signature,
	}

	json.NewEncoder(w).Encode(response)
}

func ValidateAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/login" || r.URL.Path == "/oauth/callback" || r.URL.Path == "/oauth/sign" || r.URL.Path == "/oauth/accessToken" || r.URL.Path == "/oauth/verifySignature" || r.URL.Path == "/oauth/redirect" {
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
