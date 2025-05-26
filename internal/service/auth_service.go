package service

import (
	"errors"
	"strings"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/config"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/util"
	"github.com/google/uuid"
)

// AuthService defines the authentication service interface
type AuthService interface {
	Register(req *dto.RegisterRequest) (*entity.User, error)
	Login(req *dto.LoginRequest) (*dto.TokenResponse, error)
	GoogleLogin(req *dto.GoogleAuthRequest) (*dto.TokenResponse, error)
	ValidateToken(token string) (*util.JWTClaims, error)
	CompleteRegistration(userID uuid.UUID, role string) (*entity.User, error)
}

// AuthServiceImpl implements the AuthService interface
type AuthServiceImpl struct {
	userRepo repository.UserRepository
	config   *config.Config
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
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

	// Generate token
	token, expiresIn, err := util.GenerateToken(user, s.config.JWTSecret)
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

// GoogleLogin authenticates a user with Google OAuth2
func (s *AuthServiceImpl) GoogleLogin(req *dto.GoogleAuthRequest) (*dto.TokenResponse, error) {
	// Verify Google ID token
	googleUser, err := util.VerifyGoogleIDToken(req.IdToken, s.config.GoogleClientID)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	user, err := s.userRepo.FindByEmail(googleUser.Email)
	if err != nil {
		// Create new user if not exists
		username := strings.Split(googleUser.Email, "@")[0]
		role := "pending"

		// Generate random password for OAuth users
		randomPassword, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		hashedPassword, err := util.HashPassword(randomPassword.String())
		if err != nil {
			return nil, err
		}

		user = &entity.User{
			Username: username,
			Email:    googleUser.Email,
			Password: hashedPassword,
			Role:     role,
		}

		if err := s.userRepo.Create(user); err != nil {
			return nil, err
		}
	}

	// Generate token
	token, expiresIn, err := util.GenerateToken(user, s.config.JWTSecret)
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
	user.Role = role

	// Save updated user
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ValidateToken validates a JWT token
func (s *AuthServiceImpl) ValidateToken(token string) (*util.JWTClaims, error) {
	return util.ValidateToken(token, s.config.JWTSecret)
}
