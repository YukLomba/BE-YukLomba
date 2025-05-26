package database

import (
	"log"

	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if db == nil {
		log.Fatal("❌ Cannot migrate: DB connection is nil")
	}
	err := db.AutoMigrate(&entity.User{}, &entity.Competition{}, &entity.Registration{})
	if err != nil {
		log.Fatalf("❌ Failed to migrate: %v", err)
	}
}
