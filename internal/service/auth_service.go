package service

import (
	"context"
	"errors"
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	errs "github.com/YukLomba/BE-YukLomba/internal/domain/error"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/config"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidRole        = errors.New("invalid role")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid Email or Password")
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
	_, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			break
		default:
			return nil, errs.ErrInternalServer
		}
	}

	// Check if username already exists
	_, err = s.userRepo.FindByUsername(req.Username)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			break
		default:
			return nil, errs.ErrInternalServer
		}
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errs.ErrInternalServer
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
		return nil, errs.ErrInternalServer
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthServiceImpl) Login(req *dto.LoginRequest) (*dto.TokenResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrInvalidCredentials
		default:
			return nil, errs.ErrInternalServer
		}
	}

	// Check password
	if !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}
	expirationTime := time.Duration(24) * time.Hour * 7

	// Generate token
	token, expiresIn, err := util.GenerateToken(user, s.config.JWTSecret, expirationTime)
	if err != nil {
		return nil, errs.ErrInternalServer
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
		return "", errs.ErrInternalServer
	}
	url := s.config.AuthCodeURL(state)
	return url, nil
}

func (s *AuthServiceImpl) SignInWithGoogle(code string, state string) (*dto.TokenResponse, error) {
	_, err := util.ParseOAuthStateJWT(state, s.config.JWTSecret)
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	// Exchange code for token
	oauthToken, err := s.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	userInfo, err := util.GetGoogleUserInfo(context.Background(), oauthToken)
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	// Create user if not exists
	user, err := s.userRepo.FindByEmail(userInfo.Email)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			// Create user
		default:
			return nil, errs.ErrInternalServer
		}
	}
	if user == nil {
		user = &entity.User{
			Username:   userInfo.GivenName,
			Email:      userInfo.Email,
			Password:   "",
			Role:       "pending",
			University: "",
			Interests:  "",
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, errs.ErrInternalServer
		}
	}
	// Generate token
	expirationTime := time.Duration(24) * time.Hour * 7
	token, expiresIn, err := util.GenerateToken(user, s.config.JWTSecret, expirationTime)
	if err != nil {
		return nil, errs.ErrInternalServer
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
		return nil, ErrInvalidRole
	}

	// Find the user by ID
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrUserNotFound
		default:
			return nil, errs.ErrInternalServer
		}
	}

	// Update role
	data := &map[string]interface{}{
		"role": role,
	}

	// Save updated user
	if err := s.userRepo.Update(user.ID, data); err != nil {
		return nil, errs.ErrInternalServer
	}

	return user, nil
}

// ValidateToken validates a JWT token
func (s *AuthServiceImpl) ValidateToken(token string) (*util.JWTClaims, error) {
	return util.ValidateToken(token, s.config.JWTSecret)
}
