package seeder

import (
	"fmt"

	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"gorm.io/gorm"
)

func SeedUserCompetition(db *gorm.DB) error {
	var count int64
	db.Model(&entity.Registration{}).Count(&count)

	if count > 0 {
		return nil // Skip seeding if registrations already exist
	}

	var users []entity.User
	if err := db.Limit(2).
		Find(&users).Where("role = ?", "student").Error; err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	var competitions []entity.Competition
	if err := db.Limit(5).Find(&competitions).Error; err != nil {
		return fmt.Errorf("failed to get competitions: %w", err)
	}

	for _, user := range users {
		for _, competition := range competitions {
			registration := entity.Registration{
				UserID:        user.ID,
				CompetitionID: competition.ID,
			}

			if err := db.Create(&registration).Error; err != nil {
				return fmt.Errorf("failed to create registration: %w", err)
			}
		}
	}

	return nil
}
