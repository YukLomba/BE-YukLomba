package main

import (
	"flag"
	"fmt"

	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/router"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/config"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/database"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/database/seeder"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/repository"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Define command line flags
	seedFlag := flag.Bool("seed", false, "Seed the database with sample data")
	flag.Parse()

	cfg := config.LoadConfig()
	db := database.Init(cfg)
	database.Migrate(db)
	if *seedFlag {
		seeder.SeedAll(db)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	competitionRepo := repository.NewCompetitionRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	competitionService := service.NewCompetitionService(competitionRepo)
	authService := service.NewAuthService(userRepo, cfg)

	// Initialize controllers
	userController := controller.NewUserController(userService)
	competitionController := controller.NewCompetitionController(competitionService)
	authController := controller.NewAuthController(authService)

	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	api := app.Group("/api")

	// Setup routes
	router.SetupUserRoute(api, userController)
	router.SetupCompetitionRoute(api, competitionController)
	router.SetupAuthRoute(api, authController, authService)

	app.Listen(fmt.Sprintf(":%s", cfg.AppPort))
}
