package main

import (
	"fmt"
	"log"

	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/controller"
	"github.com/YukLomba/BE-YukLomba/internal/delivery/http/router"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/config"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/repository"
	"github.com/YukLomba/BE-YukLomba/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(
		&entity.User{},
		&entity.Competition{},
		&entity.Registration{},
	)

	err = db.SetupJoinTable(&entity.User{}, "JoinedCompetitions", &entity.Registration{})

	if err != nil {
		log.Fatalf("Failed to set up join table: %v", err)
	}

	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo)

	userController := controller.NewUserController(userService)

	app := fiber.New()

	api := app.Group("/api")

	router.SetupUserRoute(api, userController)

	app.Listen(fmt.Sprintf(":%s", cfg.AppPort))
}
