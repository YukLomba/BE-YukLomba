package controller

import (
	"errors"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrganizationController struct {
	organizationService service.OrganizationService
	userService         service.UserService
}

func NewOrganizationController(
	organizationService service.OrganizationService,
	userService service.UserService,
) *OrganizationController {
	return &OrganizationController{
		organizationService: organizationService,
		userService:         userService,
	}
}

// getUserFromCtx fetches user from context locals ("userID")
func (c *OrganizationController) getUserFromCtx(ctx *fiber.Ctx) (*entity.User, error) {
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

// GetOrganization retrieves an organization by ID
func (c *OrganizationController) GetOrganization(ctx *fiber.Ctx) error {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	organization, err := c.organizationService.GetOrganization(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Organization not found"})
	}

	return ctx.JSON(fiber.Map{"data": organization})
}

// GetAllOrganizations retrieves all organizations
func (c *OrganizationController) GetAllOrganizations(ctx *fiber.Ctx) error {
	organizations, err := c.organizationService.GetAllOrganizations()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch organizations"})
	}
	return ctx.JSON(fiber.Map{"data": organizations})
}

// CreateOrganization creates a new organization
func (c *OrganizationController) CreateOrganization(ctx *fiber.Ctx) error {
	req := new(dto.OrganizationCreateRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := c.getUserFromCtx(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if user.Role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is not authorized to create organization"})
	}

	if err := c.organizationService.CreateOrganization(req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create organization"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Organization created successfully"})
}

// UpdateOrganization updates an existing organization
func (c *OrganizationController) UpdateOrganization(ctx *fiber.Ctx) error {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req := new(dto.OrganizationUpdateRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := c.getUserFromCtx(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if user.Role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User not authorized to update organization"})
	}

	if err := c.organizationService.UpdateOrganization(id, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update organization"})
	}

	return ctx.JSON(fiber.Map{"message": "Organization updated successfully"})
}

// DeleteOrganization deletes an organization by ID
func (c *OrganizationController) DeleteOrganization(ctx *fiber.Ctx) error {
	id, err := parseUUIDParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := c.getUserFromCtx(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if user.Role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User not authorized to delete organization"})
	}

	if err := c.organizationService.DeleteOrganization(id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete organization"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Organization deleted successfully"})
}
