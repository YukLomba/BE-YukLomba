package dto

// LoginRequest represents the login request data
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents the registration request data
type RegisterRequest struct {
	Username   string  `json:"username"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	University string  `json:"university"`
	Interests  string  `json:"interests"`
	Role       *string `json:"role"`
}

type CompleteRegistrationRequest struct {
	Role string `json:"role" binding:"required"`
}

// TokenResponse represents the authentication token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// GoogleAuthRequest represents the Google OAuth2 authentication request
type GoogleAuthRequest struct {
	IdToken string `json:"id_token"`
}

type GoogleAuthCallbackRequest struct {
	Code  string `form:"code" json:"code" binding:"required"`
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
