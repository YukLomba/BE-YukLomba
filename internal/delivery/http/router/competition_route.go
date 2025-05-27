package router

import (
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/middleware"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
)

func SetupCompetitionRoute(router fiber.Router, competitionController *controller.CompetitionController, authService service.AuthService) {
	competitions := router.Group("/competitions")

	// public routes
	// Get all competitions
	competitions.Get("/", competitionController.GetAllCompetitions)

	// Get competition by ID
	competitions.Get("/:id", competitionController.GetCompetition)

	// Get competitions by organizer ID
	competitions.Get("/organizer/:id", competitionController.GetCompetitionsByOrganizer)

	// Register user to competition
	competitions.Post("/:id/register", middleware.AuthMiddleware(authService), competitionController.RegisterToCompetition)

	// protected routes
	protected := competitions.Use(middleware.AuthMiddleware(authService), middleware.RoleMiddleware("admin", "organizer"))

	// Create new competition
	protected.Post("/", competitionController.CreateCompetition)

	// Update competition
	protected.Put("/:id", competitionController.UpdateCompetition)

	// Delete competition
	protected.Delete("/:id", competitionController.DeleteCompetition)
}
