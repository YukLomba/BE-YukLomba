package controller

import (
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

	*req.OrganizerID = *authInfo.OrganizationID

	competition := mapper.ToCompetitionFromCreate(req)

	if err := c.competitionService.CreateCompetition(authInfo, competition); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create competition"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Competition created successfully"})
}

func (c *CompetitionController) CreateManyCompetitition(ctx *fiber.Ctx) error {
	req := new(dto.MultiCompetitionCreateRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := util.ValidateStruct(req); err != nil {
		errors := util.GenerateValidationErrorMessage(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors})
	}
	authInfo := util.GetAuthInfo(ctx)

	competitions := mapper.ToCompetitionsFromCreate(req.Competitions)

	notValidMessage, err := c.competitionService.CreateManyCompetitition(authInfo, competitions)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create competitions"})
	}
	response := fiber.Map{
		"message": "Competitions created successfully",
	}
	if notValidMessage != nil {
		response["not_valid"] = *notValidMessage
	}
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

// UpdateCompetition updates an existing competition
func (c *CompetitionController) UpdateCompetition(ctx *fiber.Ctx) error {
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
func (c *CompetitionController) DeleteCompetition(ctx *fiber.Ctx) error {
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
