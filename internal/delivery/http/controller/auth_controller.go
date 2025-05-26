package controller

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	authService service.AuthService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register handles user registration
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	// Parse request body
	req := new(dto.RegisterRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, email, and password are required",
		})
	}

	// Register user
	user, err := c.authService.Register(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login handles user login
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	// Parse request body
	req := new(dto.LoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Login user
	token, err := c.authService.Login(req)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return token
	return ctx.JSON(token)
}

// GoogleLogin handles Google OAuth2 login
func (c *AuthController) GoogleLogin(ctx *fiber.Ctx) error {
	// Parse request body
	req := new(dto.GoogleAuthRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if req.IdToken == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID token is required",
		})
	}

	// Login with Google
	token, err := c.authService.GoogleLogin(req)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return token
	return ctx.JSON(token)
}

// GetProfile returns the authenticated user's profile
func (c *AuthController) GetProfile(ctx *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := ctx.Locals("user_id")
	email := ctx.Locals("email")
	role := ctx.Locals("role")

	// Return user profile
	return ctx.JSON(fiber.Map{
		"id":    userID,
		"email": email,
		"role":  role,
	})
}
