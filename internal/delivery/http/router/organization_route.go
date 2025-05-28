package router

import (
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupOrganizationRoute(router fiber.Router, organizationController *controller.OrganizationController) {
	organizations := router.Group("/organizations")

	// public routes
	// Get all organizations
	organizations.Get("/", organizationController.GetAllOrganizations)

	// Get organization by ID
	organizations.Get("/:id", organizationController.GetOrganization)

	// protected routes
	protected := organizations.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))

	// Create new organization
	protected.Post("/", organizationController.CreateOrganization)

	// Update organization
	protected.Put("/:id", organizationController.UpdateOrganization)

	// Delete organization
	protected.Delete("/:id", organizationController.DeleteOrganization)
}
