package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
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
	filterQuery := new(dto.CompetitionFilter)
	if err := ctx.QueryParser(filterQuery); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	competitions, err := c.competitionService.GetAllCompetitions(filterQuery)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch competitions"})
	}
	return ctx.JSON(fiber.Map{"data": competitions})
}

// CreateCompetition creates a new competition
func (c *CompetitionController) CreateCompetition(ctx *fiber.Ctx) error {
	req := new(dto.CompetitionCreateRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if !validateDeadlineFuture(req.Deadline) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Deadline must be a future date"})
	}

	user, err := c.getUserFromCtx(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if user.Role != "organizer" || user.OrganizationID == nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is not authorized to create competition"})
	}

	*req.OrganizerID = *user.OrganizationID

	if err := c.competitionService.CreateCompetition(req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create competition"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Competition created successfully"})
}

func (c *CompetitionController) CreateManyCompetitition(ctx *fiber.Ctx) error {
	req := new(dto.MultiCompetitionCreateRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	filtered := new(dto.MultiCompetitionCreateRequest)
	var errors []error
	for _, c := range req.Competitions {
		if !validateDeadlineFuture(c.Deadline) {
			errors = append(errors, fmt.Errorf("competition with ID %s has an invalid deadline", c.Deadline))
			continue
		}
		if c.OrganizerID == nil {
			errors = append(errors, fmt.Errorf("competition with ID %s has an invalid organizer ID", c.OrganizerID))
			continue
		}
		filtered.Competitions = append(filtered.Competitions, c)
	}
	if err := c.competitionService.CreateManyCompetitition(filtered); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create competitions"})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Competitions created successfully", "errors": errors})
}

// UpdateCompetition updates an existing competition
func (c *CompetitionController) UpdateCompetition(ctx *fiber.Ctx) error {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req := new(dto.CompetitionUpdateRequest)
	if err := ctx.BodyParser(req); err != nil {
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

	if !c.isAuthorizedOrganizer(user, existingCompetition.Organizer.ID) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User not authorized to update this competition"})
	}

	if !validateDeadlineFuture(req.Deadline) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Deadline must be a future date"})
	}

	if err := c.competitionService.UpdateCompetition(id, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update competition"})
	}

	return ctx.JSON(fiber.Map{"message": "Competition updated successfully"})
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

	if !c.isAuthorizedOrganizer(user, competition.Organizer.ID) {
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

func (c *CompetitionController) RegisterToCompetition(ctx *fiber.Ctx) error {
	competitionID, err := parseUUIDParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := c.getUserFromCtx(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.competitionService.RegisterUserToCompetition(user.ID, competitionID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register for competition"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Successfully registered for competition"})
}
