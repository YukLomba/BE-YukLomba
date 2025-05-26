package router

import (
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupCompetitionRoute(router fiber.Router, competitionController *controller.CompetitionController) {
	competitions := router.Group("/competitions")

	// Get all competitions
	competitions.Get("/", competitionController.GetAllCompetitions)

	// Get competition by ID
	competitions.Get("/:id", competitionController.GetCompetition)

	// Create new competition
	competitions.Post("/", competitionController.CreateCompetition)

	// Update competition
	competitions.Put("/:id", competitionController.UpdateCompetition)

	// Delete competition
	competitions.Delete("/:id", competitionController.DeleteCompetition)

	// Get competitions by organizer ID
	competitions.Get("/organizer/:id", competitionController.GetCompetitionsByOrganizer)
}
