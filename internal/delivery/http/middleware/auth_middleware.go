package middleware

import (
	"strings"

	"slices"

	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/config"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/util"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware is a middleware for JWT authentication
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Check if the authorization header has the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		JwtSecret := config.JwtSecret

		// Validate the token
		claims, err := util.ValidateToken(token, JwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		if claims.Role == "pending" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User registration is pending",
			})
		}

		// Set user information in the context
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		// Continue to the next middleware or handler
		return c.Next()
	}
}

// RoleMiddleware is a middleware for role-based authorization
func RoleMiddleware(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Check if the user has the required role
		userRole := role.(string)
		if slices.Contains(roles, userRole) {
			return c.Next()
		}

		// If the user doesn't have the required role, return forbidden
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}
}
