package middleware

import (
	"strings"

	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware is a middleware for JWT authentication
func AuthMiddleware(authService service.AuthService) fiber.Handler {
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

		// Validate the token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
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
		for _, r := range roles {
			if r == userRole {
				return c.Next()
			}
		}

		// If the user doesn't have the required role, return forbidden
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}
}
