package controller

import (
	"errors"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/mapper"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/util"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
)

type CompetitionController struct {
	competitionService service.CompetitionService
}

func NewCompetitionController(
	competitionService service.CompetitionService,
) *CompetitionController {
	return &CompetitionController{
		competitionService: competitionService,
	}
}

// GetCompetition retrieves a competition by ID
func (c *CompetitionController) GetCompetition(ctx *fiber.Ctx) error {
	id, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	competition, err := c.competitionService.GetCompetition(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Competition not found"})
	}

	response := mapper.ToCompetitionResponse(competition)

	return ctx.JSON(response)
}

func (c *CompetitionController) GetCompetitionEventLink(ctx *fiber.Ctx) error {
	id, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	authInfo := util.GetAuthInfo(ctx)

	eventLink, err := c.competitionService.GetRegisteredEventLink(authInfo, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Competition not found"})
	}

	if eventLink == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not registered for this competition"})
	}

	return ctx.JSON(fiber.Map{"event_link": eventLink})
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
	response := mapper.ToCompetitionsResponse(competitions)
	return ctx.JSON(response)
}

func (c *CompetitionController) GetManagedCompetition(ctx *fiber.Ctx) error {
	id, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	authInfo := util.GetAuthInfo(ctx)

	competition, err := c.competitionService.GetManagedCompetitionByID(authInfo, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Competition not found"})
	}

	response := mapper.ToCompetitionResponse(competition)

	return ctx.JSON(response)
}

// CreateCompetition creates a new competition
func (c *CompetitionController) CreateCompetition(ctx *fiber.Ctx) error {
	req := new(dto.CompetitionCreateRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := util.ValidateStruct(req); err != nil {
		errors := util.GenerateValidationErrorMessage(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors})
	}
	authInfo := util.GetAuthInfo(ctx)

	competition := mapper.ToCompetitionFromCreate(req)

	if err := c.competitionService.CreateCompetition(authInfo, competition); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create competition"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Competition created successfully"})
}

// UpdateCompetition updates an existing competition
func (c *CompetitionController) UpdateManagedCompetition(ctx *fiber.Ctx) error {
	id, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req := new(dto.CompetitionUpdateRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	authInfo := util.GetAuthInfo(ctx)
	data := mapper.ToCompetitionFromUpdate(req)

	if err := c.competitionService.UpdateCompetition(authInfo, id, data); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update competition"})
	}

	return ctx.JSON(fiber.Map{"message": "Competition updated successfully"})
}

// DeleteCompetition deletes a competition by ID
func (c *CompetitionController) DeleteManagedCompetition(ctx *fiber.Ctx) error {
	CompetitionId, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	authInfo := util.GetAuthInfo(ctx)

	if err := c.competitionService.DeleteCompetition(authInfo, CompetitionId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete competition"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Competition deleted successfully"})
}

// GetCompetitionsByOrganizer retrieves competitions by organizer ID
func (c *CompetitionController) GetCompetitionsByOrganizer(ctx *fiber.Ctx) error {
	id, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	competitions, err := c.competitionService.GetCompetitionsByOrganizer(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch competitions"})
	}
	response := mapper.ToCompetitionsResponse(competitions)

	return ctx.JSON(response)
}

func (c *CompetitionController) RegisterToCompetition(ctx *fiber.Ctx) error {
	competitionID, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	authInfo := util.GetAuthInfo(ctx)

	if err := c.competitionService.RegisterUserToCompetition(authInfo, competitionID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register for competition", "message": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Successfully registered for competition"})
}

func (c *CompetitionController) SubmitReview(ctx *fiber.Ctx) error {
	competitionID, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid competition ID"})
	}

	req := new(dto.CompetititionReview)
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid review data"})
	}

	review := mapper.ToCompetitionReview(req)

	authInfo := util.GetAuthInfo(ctx)

	if err := c.competitionService.SubmitReview(authInfo, competitionID, review); err != nil {
		switch {
		case errors.Is(err, service.ErrCompetitionNotRegistered):
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Review cannot be submitted because you are not registered for this competition"})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to submit review"})
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Review submitted successfully"})
}

func (c *CompetitionController) GetCompetitionReviews(ctx *fiber.Ctx) error {
	competitionID, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid competition ID"})
	}

	reviews, err := c.competitionService.GetCompetitionReviews(competitionID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get reviews"})
	}

	response := mapper.ToCompetitionReviewsResponse(reviews)

	return ctx.JSON(response)
}

// GetManagedCompetitions retrieves all competitions managed by the current user (organizer/admin)
func (c *CompetitionController) GetManagedCompetitions(ctx *fiber.Ctx) error {
	authInfo := util.GetAuthInfo(ctx)

	competitions, err := c.competitionService.GetManagedCompetitions(authInfo)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch managed competitions"})
	}

	response := mapper.ToCompetitionsResponse(competitions)
	return ctx.JSON(response)
}

// ApproveCompetition approves a competition (admin only)
func (c *CompetitionController) ApproveCompetition(ctx *fiber.Ctx) error {
	competitionID, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid competition ID"})
	}

	authInfo := util.GetAuthInfo(ctx)

	if err := c.competitionService.ApproveCompetition(authInfo, competitionID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to approve competition", "message": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Competition approved successfully"})
}
