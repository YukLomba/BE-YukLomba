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
	err := db.AutoMigrate(&entity.User{}, &entity.Competition{}, &entity.Registration{}, &entity.Organization{})
	if err != nil {
		log.Fatalf("❌ Failed to migrate: %v", err)
	}
}

func DropAllTables(db *gorm.DB) {
	db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
}

func TruncateAllTables(db *gorm.DB) {
	if db == nil {
		log.Fatal("❌ Cannot truncate: DB connection is nil")
		return
	}

	err := db.Exec("TRUNCATE users, competitions, registrations, organizations RESTART IDENTITY CASCADE").Error
	if err != nil {
		log.Fatalf("❌ Failed to truncate tables: %v", err)
	}
}
