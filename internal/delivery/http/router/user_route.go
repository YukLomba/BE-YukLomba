package router

import (
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/middleware"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoute(router fiber.Router, userController *controller.UserController, authService service.AuthService) {
	users := router.Group("/users", middleware.AuthMiddleware(authService))
	users.Get("/", userController.GetAllUsers, middleware.RoleMiddleware("admin"))
	users.Get("/:id", userController.GetUser)
	users.Get("/:id/registration", userController.GetAllUserPastCompetition)
	users.Put("/:id", userController.UpdateUser)
	// users.Post("/", userController.CreateUser)
}
