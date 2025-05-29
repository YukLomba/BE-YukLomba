package util

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

// OAuthStateClaims stores nonce and redirect URL in the JWT claims
type OAuthStateClaims struct {
	Nonce string `json:"nonce"`
	jwt.RegisteredClaims
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}

// GenerateOAuthStateJWT creates a JWT string encoding nonce and redirect URL
func GenerateOAuthStateJWT(jwtSecret string, expiry time.Duration) (string, error) {
	nonce, err := GenerateNonce(16)
	if err != nil {
		return "", err
	}
	claims := OAuthStateClaims{
		Nonce: nonce,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ParseOAuthStateJWT parses and validates the JWT state string
func ParseOAuthStateJWT(tokenString string, jwtSecret string) (*OAuthStateClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OAuthStateClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*OAuthStateClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func GetGoogleUserInfo(ctx context.Context, idToken *oauth2.Token) (*GoogleUserInfo, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(idToken))
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info")
	}
	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func GenerateNonce(size int) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Use base64 URL encoding to ensure it's safe for URLs
	return base64.URLEncoding.EncodeToString(b), nil
}
