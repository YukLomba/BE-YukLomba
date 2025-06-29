package router

import (
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoute sets up the authentication routes
func SetupAuthRoute(router fiber.Router, authController *controller.AuthController, authMiddleware *fiber.Handler) {
	auth := router.Group("/auth")

	// Public routes
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Get("/google", authController.GoogleAuth)
	auth.Get("/google/callback", authController.GoogleCallback)
	auth.Post("/complete-registration", authController.CompleteRegistration)

	// Protected routes
	auth.Get("/profile", *authMiddleware, authController.GetProfile)
}
