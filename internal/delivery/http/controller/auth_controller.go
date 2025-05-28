package controller

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/util"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	if err := util.ValidateStruct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
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

// GoogleAuth handles Google OAuth2 authentication
func (c *AuthController) GoogleAuth(ctx *fiber.Ctx) error {
	// Get Google OAuth URL
	url, err := c.authService.GetGoogleOauthUrl()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get Google OAuth URL",
		})
	}
	// Return redirect URL
	return ctx.Redirect(url)
}
func (c *AuthController) GoogleCallback(ctx *fiber.Ctx) error {
	// Parse request body
	req := new(dto.GoogleAuthCallbackRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	token, err := c.authService.SignInWithGoogle(req.Code, req.State)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Return token
	return ctx.JSON(token)
}

func (c *AuthController) CompleteRegistration(ctx *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Role string `json:"role"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Validate role here or in service layer
	if req.Role != "student" && req.Role != "organizer" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid role"})
	}
	userId := ctx.Locals("user_id").(uuid.UUID)

	// Complete registration
	user, err := c.authService.CompleteRegistration(userId, req.Role)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Return success response
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registration completed successfully",
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
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
