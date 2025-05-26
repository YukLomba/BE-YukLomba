package controller

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CompetitionController struct {
	competitionService service.CompetitionService
}

func NewCompetitionController(competitionService service.CompetitionService) *CompetitionController {
	return &CompetitionController{
		competitionService: competitionService,
	}
}

// GetCompetition retrieves a competition by ID
func (c *CompetitionController) GetCompetition(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	competition, err := c.competitionService.GetCompetition(parsedID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Competition not found",
		})
	}

	return ctx.JSON(competition)
}

// GetAllCompetitions retrieves all competitions
func (c *CompetitionController) GetAllCompetitions(ctx *fiber.Ctx) error {
	competitions, err := c.competitionService.GetAllCompetitions()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch competitions",
		})
	}

	return ctx.JSON(competitions)
}

// CreateCompetition creates a new competition
func (c *CompetitionController) CreateCompetition(ctx *fiber.Ctx) error {
	competition := new(entity.Competition)

	if err := ctx.BodyParser(competition); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := c.competitionService.CreateCompetition(competition); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create competition",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(competition)
}

// UpdateCompetition updates an existing competition
func (c *CompetitionController) UpdateCompetition(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	competition := new(entity.Competition)
	if err := ctx.BodyParser(competition); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	competition.ID = parsedID

	if err := c.competitionService.UpdateCompetition(competition); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update competition",
		})
	}

	return ctx.JSON(competition)
}

// DeleteCompetition deletes a competition by ID
func (c *CompetitionController) DeleteCompetition(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	if err := c.competitionService.DeleteCompetition(parsedID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete competition",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Competition deleted successfully",
	})
}

// GetCompetitionsByOrganizer retrieves competitions by organizer ID
func (c *CompetitionController) GetCompetitionsByOrganizer(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	competitions, err := c.competitionService.GetCompetitionsByOrganizer(parsedID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch competitions",
		})
	}

	return ctx.JSON(competitions)
}
