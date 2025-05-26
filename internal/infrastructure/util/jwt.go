package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for the given user
func GenerateToken(user *entity.User, jwtSecret string) (string, int64, error) {
	// Set token expiration time (e.g., 24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)
	expiresIn := expirationTime.Unix() - time.Now().Unix()

	// Set role to "user" if it's nil
	role := user.Role

	// Create claims
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresIn, nil
}

// ValidateToken validates the JWT token and returns the claims
func ValidateToken(tokenString string, jwtSecret string) (*JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
