package util

import (
	"context"
	"errors"

	"google.golang.org/api/idtoken"
)

// GoogleUserInfo represents the user information from Google OAuth2
type GoogleUserInfo struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	GoogleID string `json:"sub"`
}

// VerifyGoogleIDToken verifies the Google ID token and returns the user information
func VerifyGoogleIDToken(idToken string, clientID string) (*GoogleUserInfo, error) {
	if idToken == "" {
		return nil, errors.New("id token is required")
	}

	// Create a context
	ctx := context.Background()

	// Verify the ID token
	payload, err := idtoken.Validate(ctx, idToken, clientID)
	if err != nil {
		return nil, err
	}

	// Extract user information from the payload
	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, errors.New("email not found in token")
	}

	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	googleID, _ := payload.Claims["sub"].(string)

	// Create and return the user information
	userInfo := &GoogleUserInfo{
		Email:    email,
		Name:     name,
		Picture:  picture,
		GoogleID: googleID,
	}

	return userInfo, nil
}
