package controller

import (
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/util"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AnalyticsController struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsController(analyticsService service.AnalyticsService) *AnalyticsController {
	return &AnalyticsController{
		analyticsService: analyticsService,
	}
}

func (c *AnalyticsController) GetDashboard(ctx *fiber.Ctx) error {
	authInfo := util.GetAuthInfo(ctx)

	dashboard, err := c.analyticsService.GetDashboard(authInfo)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(dashboard)
}

func (c *AnalyticsController) GetCompetitionAnalytics(ctx *fiber.Ctx) error {
	competitionID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid competition ID")
	}

	analytics, err := c.analyticsService.GetCompetitionAnalytics(competitionID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(analytics)
}
