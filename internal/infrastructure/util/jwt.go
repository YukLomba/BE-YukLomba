package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	UserID         uuid.UUID  `json:"user_id"`
	Role           string     `json:"role"`
	OrganizationID *uuid.UUID `json:"organization,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for the given user
func GenerateToken(user *entity.User, jwtSecret string, expiry time.Duration) (string, int64, error) {

	// Set role to "user" if it's nil
	role := user.Role
	expirationTime := time.Now().Add(expiry)
	expiresIn := expirationTime.Unix() - time.Now().Unix()

	// Create claims
	claims := &JWTClaims{
		UserID: user.ID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}
	if user.Organization != nil {
		claims.OrganizationID = user.OrganizationID
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
		// Check if token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// helper to get ctx.locals
func GetAuthInfo(ctx *fiber.Ctx) *dto.AuthInfo {
	AuthInfo := new(dto.AuthInfo)
	AuthInfo.ID = ctx.Locals("user_id").(uuid.UUID)
	AuthInfo.Role = ctx.Locals("role").(string)
	AuthInfo.OrganizationID = ctx.Locals("organization_id").(*uuid.UUID)
	return AuthInfo
}
