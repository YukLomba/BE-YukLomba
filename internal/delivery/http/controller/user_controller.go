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

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (h *UserController) GetUser(c *fiber.Ctx) error {
	id, err := util.ParseCtxParam(c, "id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Id",
		})
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch user",
			})
		}
	}

	response := mapper.ToUserResponse(user)

	return c.JSON(response)
}

func (h *UserController) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	response := mapper.ToUsersResponse(users)

	return c.JSON(response)
}

func (h *UserController) GetAllUserPastCompetition(c *fiber.Ctx) error {
	id, err := util.ParseCtxParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Id",
		})
	}

	pastCompetition, err := h.userService.GetAllUserRegistration(id)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch user",
			})
		}
	}
	data := mapper.ToCompetitionsResponse(pastCompetition)

	return c.JSON(data)
}

func (h *UserController) UpdateUser(c *fiber.Ctx) error {
	id, err := util.ParseCtxParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Id",
		})
	}

	userData := new(dto.UserProfileUpdate)

	if err := c.BodyParser(userData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	authInfo := util.GetAuthInfo(c)

	if err := util.ValidateStruct(userData); err != nil {
		errors := util.GenerateValidationErrorMessage(err)
		return c.Status(400).JSON(fiber.Map{
			"errors": errors, // Direct usage, no dereferencing needed
		})
	}
	data := mapper.MapUserUpdate(userData)

	if err := h.userService.UpdateUser(authInfo, id, data); err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update user",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}
