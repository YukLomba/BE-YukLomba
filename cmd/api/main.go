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

	// Initialize controllers
	userController := controller.NewUserController(userService)
	competitionController := controller.NewCompetitionController(competitionService)

	app := fiber.New()

	api := app.Group("/api")

	// Setup routes
	router.SetupUserRoute(api, userController)
	router.SetupCompetitionRoute(api, competitionController)

	app.Listen(fmt.Sprintf(":%s", cfg.AppPort))
}
