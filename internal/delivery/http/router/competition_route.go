package router

import (
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupCompetitionRoute(router fiber.Router, competitionController *controller.CompetitionController, authMiddleware *fiber.Handler) {
	// Student routes (public)
	competition := router.Group("/competitions")
	{
		// Get all approved competitions
		competition.Get("/", competitionController.GetAllCompetitions)

		// Get competition details
		competition.Get("/:id", competitionController.GetCompetition)

		// Register to competition
		competition.Post("/:id/register", *authMiddleware, competitionController.RegisterToCompetition)

		// Competition reviews
		competition.Get("/:id/reviews", competitionController.GetCompetitionReviews)
		competition.Post("/:id/reviews", *authMiddleware, middleware.RoleMiddleware("student"), competitionController.SubmitReview)
	}

	// Management routes (protected)
	manageRoutes := router.Group("/manage/competitions")
	manageRoutes.Use(*authMiddleware, middleware.RoleMiddleware("admin", "organizer"))
	{
		// Get all managed competitions
		manageRoutes.Get("/", competitionController.GetManagedCompetitions)

		// CRUD operations (status will be pending)
		manageRoutes.Post("/", competitionController.CreateCompetition)
		manageRoutes.Get("/:id", competitionController.GetManagedCompetition)
		manageRoutes.Put("/:id", competitionController.UpdateManagedCompetition)
		manageRoutes.Delete("/:id", competitionController.DeleteManagedCompetition)

		// Admin-only approval endpoint
		manageRoutes.Post("/:id/approve", middleware.RoleMiddleware("admin"), competitionController.ApproveCompetition)
	}
}
