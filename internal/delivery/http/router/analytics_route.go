package router

import (
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupAnalyticsRoute(router fiber.Router, controller *controller.AnalyticsController, authMiddleware *fiber.Handler) {
	group := router.Group("/analytics")
	group.Use(*authMiddleware)

	group.Get("/", controller.GetDashboard)
	group.Get("/competitions/:id", controller.GetCompetitionAnalytics)
}
