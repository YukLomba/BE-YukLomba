package dto

// LoginRequest represents the login request data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest represents the registration request data
type RegisterRequest struct {
	Username   string  `json:"username" validate:"required"`
	Email      string  `json:"email" validate:"required,email"`
	Password   string  `json:"password" validate:"required,min=8"`
	University string  `json:"university" validate:"required"`
	Interests  string  `json:"interests" validate:"required"`
	Role       *string `json:"role" validate:"omitempty"`
}

type CompleteRegistrationRequest struct {
	Role string `json:"role"`
}

// TokenResponse represents the authentication token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type GoogleAuthCallbackRequest struct {
	Code  string `form:"code" json:"code"`
	State string `form:"state" json:"state"`
}

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}
