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
	migrateFlag := flag.Bool("migrate", false, "Migrate the database schema")
	freshFlag := flag.Bool("fresh", false, "Drop and recreate the database")
	flag.Parse()

	cfg := config.LoadConfig()
	db := database.Init(cfg)

	if *freshFlag {
		database.DropAllTables(db)
	}
	if *migrateFlag {
		database.Migrate(db)
	}
	if *seedFlag {
		seeder.SeedAll(db)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	competitionRepo := repository.NewCompetitionRepository(db)
	organizationRepo := repository.NewOrganizationRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	competitionService := service.NewCompetitionService(competitionRepo)
	authService := service.NewAuthService(userRepo, cfg)
	organizationService := service.NewOrganizationService(organizationRepo)

	// Initialize controllers
	userController := controller.NewUserController(userService)
	competitionController := controller.NewCompetitionController(competitionService, userService)
	authController := controller.NewAuthController(authService)
	organizationController := controller.NewOrganizationController(organizationService, userService)

	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
	}))

	api := app.Group("/api")

	// Setup routes
	router.SetupUserRoute(api, userController, authService)
	router.SetupCompetitionRoute(api, competitionController, authService)
	router.SetupAuthRoute(api, authController, authService)
	router.SetupOrganizationRoute(api, organizationController, authService)

	app.Listen(fmt.Sprintf(":%s", cfg.AppPort))
}
