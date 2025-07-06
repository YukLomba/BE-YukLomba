package router

import (
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupCompetitionRoute(router fiber.Router, competitionController *controller.CompetitionController, authMiddleware *fiber.Handler) {
	competitions := router.Group("/competitions")

	// public routes
	// Get all competitions
	competitions.Get("/", competitionController.GetAllCompetitions)

	// Get competition by ID
	competitions.Get("/:id", competitionController.GetCompetition)

	// Get competitions by organizer ID
	competitions.Get("/organizer/:id", competitionController.GetCompetitionsByOrganizer)

	competitions.Get("/:id/reviews", competitionController.GetCompetitionReviews)

	competitions.Post("/:id/reviews", *authMiddleware, middleware.RoleMiddleware("student"), competitionController.SubmitReview)

	// Register user to competition
	competitions.Post("/:id/register", *authMiddleware, competitionController.RegisterToCompetition)

	// protected routes
	protected := competitions.Use(*authMiddleware, middleware.RoleMiddleware("admin", "organizer"))

	// Create new competition
	protected.Post("/", competitionController.CreateCompetition)

	// Create multiple competitions
	protected.Post("/multi", competitionController.CreateManyCompetitition)

	// Update competition
	protected.Put("/:id", competitionController.UpdateCompetition)

	// Delete competition
	protected.Delete("/:id", competitionController.DeleteCompetition)
}
