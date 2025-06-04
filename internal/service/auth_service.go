package service

import (
	"context"
	"errors"
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/config"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/util"
	"github.com/google/uuid"
)

var (
	ErrInvalidRole = errors.New("invalid role")
)

// AuthService defines the authentication service interface
type AuthService interface {
	Register(req *dto.RegisterRequest) (*entity.User, error)
	Login(req *dto.LoginRequest) (*dto.TokenResponse, error)
	GetGoogleOauthUrl() (string, error)
	SignInWithGoogle(code string, state string) (*dto.TokenResponse, error)
	ValidateToken(token string) (*util.JWTClaims, error)
	CompleteRegistration(userID uuid.UUID, role string) (*entity.User, error)
}

// AuthServiceImpl implements the AuthService interface
type AuthServiceImpl struct {
	userRepo repository.UserRepository
	config   config.Auth
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repository.UserRepository, cfg config.Auth) AuthService {
	return &AuthServiceImpl{
		userRepo: userRepo,
		config:   cfg,
	}
}

// Register registers a new user
func (s *AuthServiceImpl) Register(req *dto.RegisterRequest) (*entity.User, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Check if username already exists
	existingUser, err = s.userRepo.FindByUsername(req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Set default role
	role := "pending"
	if req.Role != nil {
		if *req.Role != "student" && *req.Role != "organizer" {
			return nil, ErrInvalidRole
		}
		role = *req.Role
	}

	// Create user
	user := &entity.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   hashedPassword,
		Role:       role,
		University: req.University,
		Interests:  req.Interests,
	}

	// Save user to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthServiceImpl) Login(req *dto.LoginRequest) (*dto.TokenResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}
	expirationTime := time.Duration(24) * time.Hour * 7

	// Generate token
	token, expiresIn, err := util.GenerateToken(user, s.config.JWTSecret, expirationTime)
	if err != nil {
		return nil, err
	}

	// Return token response
	return &dto.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}, nil
}

func (s *AuthServiceImpl) GetGoogleOauthUrl() (string, error) {
	expirationTime := time.Duration(5) * time.Minute
	state, err := util.GenerateOAuthStateJWT(s.config.JWTSecret, expirationTime)
	if err != nil {
		return "", err
	}
	url := s.config.AuthCodeURL(state)
	return url, nil
}

func (s *AuthServiceImpl) SignInWithGoogle(code string, state string) (*dto.TokenResponse, error) {
	_, err := util.ParseOAuthStateJWT(state, s.config.JWTSecret)
	if err != nil {
		return nil, err
	}
	// Exchange code for token
	oauthToken, err := s.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	userInfo, err := util.GetGoogleUserInfo(context.Background(), oauthToken)
	if err != nil {
		return nil, err
	}
	// Create user if not exists
	user, err := s.userRepo.FindByEmail(userInfo.Email)
	if err != nil {
		user = &entity.User{
			Username:   userInfo.GivenName,
			Email:      userInfo.Email,
			Password:   "",
			Role:       "pending",
			University: "",
			Interests:  "",
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, err
		}
	}
	// Generate token
	expirationTime := time.Duration(24) * time.Hour * 7
	token, expiresIn, err := util.GenerateToken(user, s.config.JWTSecret, expirationTime)
	if err != nil {
		return nil, err
	}
	// Return token response
	return &dto.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}, nil

}

func (s *AuthServiceImpl) CompleteRegistration(userID uuid.UUID, role string) (*entity.User, error) {
	// Validate role
	if role != "student" && role != "organizer" {
		return nil, errors.New("invalid role, must be 'student' or 'organizer'")
	}

	// Find the user by ID
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update role
	data := &map[string]interface{}{
		"role": role,
	}

	// Save updated user
	if err := s.userRepo.Update(user.ID, data); err != nil {
		return nil, err
	}

	return user, nil
}

// ValidateToken validates a JWT token
func (s *AuthServiceImpl) ValidateToken(token string) (*util.JWTClaims, error) {
	return util.ValidateToken(token, s.config.JWTSecret)
}
