package seeder

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedUsers seeds the users table with sample data
func SeedUsers(db *gorm.DB) error {
	// Check if users already exist
	var count int64
	db.Model(&entity.User{}).Count(&count)
	if count > 0 {
		return nil // Skip seeding if users already exist
	}

	// Define sample users
	users := []entity.User{
		{
			Username:   "organizer1",
			Email:      "organizer1@example.com",
			Password:   "$2a$10$1qAz2wSx3eCc1/CuUTCTfez5YzaBPQuE31ASx8O0tYVvCmZmJUaLm", // "password"
			Role:       stringPtr("organizer"),
			University: "University of Technology",
			Interests:  "Technology, Innovation",
		},
		{
			Username:   "organizer2",
			Email:      "organizer2@example.com",
			Password:   "$2a$10$1qAz2wSx3eCc1/CuUTCTfez5YzaBPQuE31ASx8O0tYVvCmZmJUaLm", // "password"
			Role:       stringPtr("organizer"),
			University: "State University",
			Interests:  "Business, Entrepreneurship",
		},
		{
			Username:   "participant1",
			Email:      "participant1@example.com",
			Password:   "$2a$10$1qAz2wSx3eCc1/CuUTCTfez5YzaBPQuE31ASx8O0tYVvCmZmJUaLm", // "password"
			Role:       stringPtr("participant"),
			University: "Technical Institute",
			Interests:  "Programming, AI",
		},
		{
			Username:   "participant2",
			Email:      "participant2@example.com",
			Password:   "$2a$10$1qAz2wSx3eCc1/CuUTCTfez5YzaBPQuE31ASx8O0tYVvCmZmJUaLm", // "password"
			Role:       stringPtr("participant"),
			University: "Arts Academy",
			Interests:  "Design, UI/UX",
		},
		{
			Username:   "admin",
			Email:      "admin@example.com",
			Password:   "$2a$10$1qAz2wSx3eCc1/CuUTCTfez5YzaBPQuE31ASx8O0tYVvCmZmJUaLm", // "password"
			Role:       stringPtr("admin"),
			University: "Admin University",
			Interests:  "System Administration",
		},
	}

	// Insert users
	for i := range users {
		// Ensure each user has a UUID
		users[i].ID = uuid.New()

		// Create the user
		if err := db.Create(&users[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
