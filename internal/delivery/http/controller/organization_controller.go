package controller

import (
	"errors"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	errs "github.com/YukLomba/BE-YukLomba/internal/domain/error"
	"github.com/YukLomba/BE-YukLomba/internal/domain/mapper"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/util"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
)

type OrganizationController struct {
	organizationService service.OrganizationService
}

func NewOrganizationController(
	organizationService service.OrganizationService,
) *OrganizationController {
	return &OrganizationController{
		organizationService: organizationService,
	}
}

// GetOrganization retrieves an organization by ID
func (c *OrganizationController) GetOrganization(ctx *fiber.Ctx) error {
	id, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	organization, err := c.organizationService.GetOrganization(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrganizationNotFound):
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Organization not found"})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch organization"})
		}
	}
	response := mapper.ToOrganizationResponse(organization)

	return ctx.JSON(response)
}

// GetAllOrganizations retrieves all organizations
func (c *OrganizationController) GetAllOrganizations(ctx *fiber.Ctx) error {
	organizations, err := c.organizationService.GetAllOrganizations()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch organizations"})
	}
	response := mapper.ToOrganizationsResponse(organizations)
	return ctx.JSON(response)
}

// CreateOrganization creates a new organization
func (c *OrganizationController) CreateOrganization(ctx *fiber.Ctx) error {
	req := new(dto.OrganizationCreateRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	authInfo := util.GetAuthInfo(ctx)

	Organization := mapper.ToOrganizationFromCreate(req)

	if err := c.organizationService.CreateOrganization(Organization, authInfo); err != nil {
		switch {
		case errors.Is(err, errs.ErrUnauthorized):
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		case errors.Is(err, errs.ErrInternalServer):
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create organization"})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create organization"})
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Organization created successfully"})
}

// UpdateOrganization updates an existing organization
func (c *OrganizationController) UpdateOrganization(ctx *fiber.Ctx) error {
	id, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req := new(dto.OrganizationUpdateRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	authInfo := util.GetAuthInfo(ctx)
	Organization := mapper.ToOrganizationFromUpdate(req, id)

	if err := c.organizationService.UpdateOrganization(Organization, authInfo); err != nil {
		switch {
		case errors.Is(err, errs.ErrUnauthorized):
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		case errors.Is(err, service.ErrOrganizationNotFound):
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Organization not found"})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update organization"})
		}
	}

	return ctx.JSON(fiber.Map{"message": "Organization updated successfully"})
}

// DeleteOrganization deletes an organization by ID
func (c *OrganizationController) DeleteOrganization(ctx *fiber.Ctx) error {
	id, err := util.ParseCtxParam(ctx, "id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	authInfo := util.GetAuthInfo(ctx)

	if err := c.organizationService.DeleteOrganization(id, authInfo); err != nil {
		switch {
		case errors.Is(err, errs.ErrUnauthorized):
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		case errors.Is(err, service.ErrOrganizationNotFound):
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Organization not found"})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete organization"})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Organization deleted successfully"})
}
