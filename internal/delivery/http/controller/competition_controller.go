package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CompetitionController struct {
	competitionService service.CompetitionService
	userService        service.UserService
}

func NewCompetitionController(
	competitionService service.CompetitionService,
	userService service.UserService,
) *CompetitionController {
	return &CompetitionController{
		competitionService: competitionService,
		userService:        userService,
	}
}

// parseUUIDParam parses UUID from path param and returns error if invalid/missing
func parseUUIDParam(ctx *fiber.Ctx, param string) (uuid.UUID, error) {
	id := ctx.Params(param)
	if id == "" {
		return uuid.Nil, fmt.Errorf("missing '%s' parameter", param)
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format for '%s'", param)
	}
	return parsedID, nil
}

// getUserFromCtx fetches user from context locals ("userID")
func (c *CompetitionController) getUserFromCtx(ctx *fiber.Ctx) (*entity.User, error) {
	rawUserID := ctx.Locals("userID")
	userID, ok := rawUserID.(uuid.UUID)
	if !ok {
		return nil, errors.New("unauthorized: user ID missing or invalid in context")
	}
	user, err := c.userService.GetUser(userID)
	if err != nil {
		return nil, errors.New("unauthorized: user not found")
	}
	return user, nil
}

// isAuthorizedOrganizer checks if user is organizer and belongs to given organization ID
func (c *CompetitionController) isAuthorizedOrganizer(user *entity.User, organizerID uuid.UUID) bool {
	return user.Role == "organizer" && user.OrganizationID != nil && *user.OrganizationID == organizerID
}

// validateDeadlineFuture checks if deadline is a future date
func validateDeadlineFuture(deadline time.Time) bool {
	return deadline.After(time.Now())
}

// GetCompetition retrieves a competition by ID
func (c *CompetitionController) GetCompetition(ctx *fiber.Ctx) error {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	competition, err := c.competitionService.GetCompetition(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Competition not found"})
	}

	return ctx.JSON(fiber.Map{"data": competition})
}

// GetAllCompetitions retrieves all competitions
func (c *CompetitionController) GetAllCompetitions(ctx *fiber.Ctx) error {
	competitions, err := c.competitionService.GetAllCompetitions()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch competitions"})
	}
	return ctx.JSON(fiber.Map{"data": competitions})
}

// CreateCompetition creates a new competition
func (c *CompetitionController) CreateCompetition(ctx *fiber.Ctx) error {
	competition := new(entity.Competition)
	if err := ctx.BodyParser(competition); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if !validateDeadlineFuture(competition.Deadline) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Deadline must be a future date"})
	}

	user, err := c.getUserFromCtx(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if user.Role != "organizer" || user.OrganizationID == nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is not authorized to create competition"})
	}

	competition.OrganizerID = *user.OrganizationID

	if err := c.competitionService.CreateCompetition(competition); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create competition"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": competition})
}

// UpdateCompetition updates an existing competition
func (c *CompetitionController) UpdateCompetition(ctx *fiber.Ctx) error {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	competitionData := new(entity.Competition)
	if err := ctx.BodyParser(competitionData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	existingCompetition, err := c.competitionService.GetCompetition(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Competition not found"})
	}

	user, err := c.getUserFromCtx(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if !c.isAuthorizedOrganizer(user, existingCompetition.OrganizerID) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User not authorized to update this competition"})
	}

	if !validateDeadlineFuture(competitionData.Deadline) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Deadline must be a future date"})
	}

	// Update allowed fields only
	existingCompetition.Title = competitionData.Title
	existingCompetition.Type = competitionData.Type
	existingCompetition.Description = competitionData.Description
	existingCompetition.Deadline = competitionData.Deadline
	existingCompetition.Category = competitionData.Category
	existingCompetition.Rules = competitionData.Rules
	existingCompetition.EventLink = competitionData.EventLink
	existingCompetition.Results = competitionData.Results

	if err := c.competitionService.UpdateCompetition(existingCompetition); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update competition"})
	}

	return ctx.JSON(fiber.Map{"data": existingCompetition})
}

// DeleteCompetition deletes a competition by ID
func (c *CompetitionController) DeleteCompetition(ctx *fiber.Ctx) error {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	competition, err := c.competitionService.GetCompetition(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Competition not found"})
	}

	user, err := c.getUserFromCtx(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if !c.isAuthorizedOrganizer(user, competition.OrganizerID) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User not authorized to delete this competition"})
	}

	if err := c.competitionService.DeleteCompetition(id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete competition"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Competition deleted successfully"})
}

// GetCompetitionsByOrganizer retrieves competitions by organizer ID
func (c *CompetitionController) GetCompetitionsByOrganizer(ctx *fiber.Ctx) error {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	competitions, err := c.competitionService.GetCompetitionsByOrganizer(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch competitions"})
	}

	return ctx.JSON(fiber.Map{"data": competitions})
}
