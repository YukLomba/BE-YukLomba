package seeder

import (
	"log"

	"gorm.io/gorm"
)

// Seed runs all seeders
func SeedAll(db *gorm.DB) {
	log.Println("Running seeders...")

	// Seed users first since competitions depend on users
	if err := SeedUsers(db); err != nil {
		log.Printf("Error seeding users: %v\n", err)
	} else {
		log.Println("Users seeded successfully")
	}

	// Seed competitions
	if err := SeedCompetitions(db); err != nil {
		log.Printf("Error seeding competitions: %v\n", err)
	} else {
		log.Println("Competitions seeded successfully")
	}

	if err := SeedUserCompetition(db); err != nil {
		log.Printf("Error seeding Registration: %v\n", err)
	} else {
		log.Println("Regsitration seeded successfully")
	}

	log.Println("Seeding completed")
}
