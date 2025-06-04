package router

import (
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoute(router fiber.Router, userController *controller.UserController, authMiddleware *fiber.Handler) {
	users := router.Group("/users", *authMiddleware)
	users.Get("/", userController.GetAllUsers, middleware.RoleMiddleware("admin"))
	users.Get("/:id", userController.GetUser)
	users.Get("/:id/registrations", userController.GetAllUserPastCompetition)
	users.Put("/:id", userController.UpdateUser)
}
