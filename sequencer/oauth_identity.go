package sequencer

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuth struct {
	config         *oauth2.Config
	jwtSecret      string
	oauthState     string
	verifier       string
	walletSeedSalt string
}

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Identity      string `json:"identity"`
	IdentityCurve string `json:"identity_curve"`
}

type M map[string]interface{}

type contextKey string
type callbackKey string

const (
	// Key for the ID
	identityAccess          contextKey  = "identityAccess"
	tokensCallbackInfoKey   callbackKey = "tokensCallbackInfo"
	stripchainGoogleSession string      = "stripchain-google-session"
	tokenIssuer             string      = "StripChain"
)

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
	ErrTokenExpired                 = errors.New("token has expired")
	ErrInvalidToken                 = errors.New("invalid token")
	ErrAuthorization                = errors.New("authorization failed")
	ErrInvalidTokenIdentityRequired = errors.New("invalid token: identity and identity curve are required")
	ErrRefreshTokenExpired          = errors.New("refresh token has expired")
	ErrInvalidTokenId               = errors.New("invalid token id")
	ErrNotAuthenticated             = errors.New("not authenticated")
	ErrRefreshTokenNotFound         = errors.New("refresh token not found")
)

func NewGoogleAuth(redirectUrl string, clientId string, clientSecret string, sessionSecret string, jwtSecret string, walletSeedSalt string) *GoogleAuth {

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

	gob.Register(&UserInfo{})
	gob.Register(&M{})

	// State string for CSRF protection
	oauthState := generateState()
	return &GoogleAuth{googleOauthConfig, jwtSecret, oauthState, verifier, walletSeedSalt}

}

// DeriveIdentity derives a deterministic public key from a Google ID
// This becomes the identity for the MPC wallet
func (s *GoogleAuth) deriveIdentity(userId string) (address string, curve string, err error) {
	// Create deterministic seed from Google ID and server secret
	seed := crypto.Keccak256([]byte(userId + s.walletSeedSalt))
	// Generate private key deterministically
	privateKey, err := crypto.ToECDSA(seed)
	if err != nil {
		return "", "", fmt.Errorf("failed to derive private key: %v", err)
	}

	// Get public key
	pubKey := privateKey.Public()
	publicKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", fmt.Errorf("error casting public key to ECDSA")
	}

	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, "ecdsa", nil
}

func (s *GoogleAuth) sign(userId string, message string) (string, error) {
	// Derive private key (same seed as identity derivation)
	seed := crypto.Keccak256([]byte(userId + s.walletSeedSalt))
	privateKey, err := crypto.ToECDSA(seed)
	if err != nil {
		return "", fmt.Errorf("failed to derive private key: %v", err)
	}
	// Hash the message
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message)
	hash := crypto.Keccak256Hash(hashedMessage)

	// Sign the hash
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}
	if len(signature) != 65 {
		return "", fmt.Errorf("invalid signature length")
	}
	signature[64] += 27
	return hex.EncodeToString(signature), nil
}

func (s *GoogleAuth) generateIdToken(user UserInfo, identity string, identityCurve string) (string, error) {
	claims := ClaimsWithIdentity{
		Email:         user.Email,
		Name:          user.Name,
		Identity:      identity,
		IdentityCurve: identityCurve,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			Issuer:    tokenIssuer,
			Subject:   user.ID,
		},
	}

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	return token.SignedString([]byte(s.jwtSecret))
}

// GenerateAccessToken creates a JWT access token
func (s *GoogleAuth) generateAccessToken(userId string, identity string, identityCurve string) (string, error) {
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
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *GoogleAuth) generateRefreshToken(userId string, identity string, identityCurve string) (string, error) {
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

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *GoogleAuth) verifyToken(tokenStr string, tokenType string, verifyIdentity bool, secretKey string) (*ClaimsWithIdentity, error) {
	claims := &ClaimsWithIdentity{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
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
	fmt.Println("claims", claims)
	if verifyIdentity && (tokenType == "access_token" || tokenType == "refresh_token") && (claims.Identity == "" || claims.IdentityCurve == "") {
		return nil, ErrInvalidTokenIdentityRequired
	}
	if tokenType == "refresh_token" {
		token, _ := GetRefreshToken(tokenStr)
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

func GenerateRandomSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random salt: %v", err)
	}
	return salt, nil
}
