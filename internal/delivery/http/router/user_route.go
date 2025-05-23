package router

import (
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoute(router fiber.Router, userController *controller.UserController) {
	users := router.Group("/users")
	users.Get("/", userController.GetAllUsers)
	users.Get("/:id", userController.GetUser)
	users.Get("/:id/registration", userController.GetAllUserPastCompetition)
	users.Post("/", userController.CreateUser)
	users.Put("/:id", userController.UpdateUser)
}
