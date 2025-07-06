package main

import (
	"flag"
	"fmt"

	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/middleware"
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

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db := database.Init(cfg.DB)

	if *freshFlag {
		database.DropAllTables(db)
	}
	if *migrateFlag {
		database.Migrate(db)
	}
	if *seedFlag {
		database.TruncateAllTables(db)
		seeder.SeedAll(db)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	competitionRepo := repository.NewCompetitionRepository(db)
	organizationRepo := repository.NewOrganizationRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	competitionService := service.NewCompetitionService(competitionRepo, reviewRepo)
	authService := service.NewAuthService(userRepo, cfg.Auth)
	organizationService := service.NewOrganizationService(organizationRepo, userRepo)
	analyticsService := service.NewAnalyticsService(userRepo, competitionRepo, reviewRepo)

	// Initialize controllers
	userController := controller.NewUserController(userService)
	competitionController := controller.NewCompetitionController(competitionService)
	authController := controller.NewAuthController(authService)
	organizationController := controller.NewOrganizationController(organizationService)
	analyticsController := controller.NewAnalyticsController(analyticsService)

	// Initialize AuthMiddleware
	authMiddlware := middleware.AuthMiddleware(userService)

	app := fiber.New(fiber.Config{
		AppName: "YukLomba API",
	})

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
	router.SetupUserRoute(api, userController, &authMiddlware)
	router.SetupCompetitionRoute(api, competitionController, &authMiddlware)
	router.SetupAuthRoute(api, authController, &authMiddlware)
	router.SetupOrganizationRoute(api, organizationController, &authMiddlware)
	router.SetupAnalyticsRoute(api, analyticsController, &authMiddlware)

	app.Listen(fmt.Sprintf(":%s", cfg.Server.Port))
}
