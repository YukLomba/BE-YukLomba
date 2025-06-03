package seeder

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RawUser struct {
	Username     string
	Email        string
	Password     string
	Role         string
	University   string
	Interests    string
	Organization *entity.Organization
}

// SeedUsers seeds the users table with sample data
func SeedUsers(db *gorm.DB) error {
	// Check if users already exist
	var count int64
	db.Model(&entity.User{}).Count(&count)
	if count > 0 {
		return nil // Skip seeding if users already exist
	}

	db.Model(&entity.Organization{}).Count(&count)
	if count > 0 {
		return nil
	}

	organization := []entity.Organization{
		{
			Name:        "Google LLC",
			Description: "is an American multinational corporation and technology company focusing on online advertising, search engine technology, cloud computing, computer software, quantum computing, e-commerce, consumer electronics, and artificial intelligence (AI).[9] It has been referred to as the most powerful company in the world and is one of the world's most valuable brands due to its market dominance, data collection, and technological advantages in the field of AI.[11][12][13] Alongside Amazon, Apple, Meta, and Microsoft, Google's parent company, Alphabet Inc. is one of the five Big Tech companies",
			Logo:        "https://img.icons8.com/?size=100&id=17949&format=png&color=000000",
		},
		{
			Name:        "University of Indonesia",
			Description: "UI is one of the leading research universities or academic institutions in the world that continues to pursue the highest achievements in terms of discovery, development and knowledge diffusion regionally and globally.",
			Logo:        "https://air.eng.ui.ac.id/images/8/82/University_of_Indonesia_logo.svg.png",
		},
		{
			Name:        "Bandung Institute of Technology",
			Description: "Bandung Institute of Technology traces its origin to the Technische Hoogeschool te Bandoeng (THB) which was founded during the centuries-old Dutch colonialism of Indonesia. The project was founded by Karel Albert Rudolf Bosscha, a German-Dutch entrepreneur and philanthropist. His proposal was later approved by the colonial government to meet increasing demand of technical know-how in the colony.",
			Logo:        "Bandung Institute of Technology traces its origin to the Technische Hoogeschool te Bandoeng (THB) which was founded during the centuries-old Dutch colonialism of Indonesia. The project was founded by Karel Albert Rudolf Bosscha, a German-Dutch entrepreneur and philanthropist. His proposal was later approved by the colonial government to meet increasing demand of technical know-how in the colony.",
		},
	}

	// Define sample users
	users_data := []RawUser{
		{
			Username:     "organizer1",
			Email:        "organizer1@example.com",
			Password:     "abcd1234",
			Role:         "organizer",
			University:   "University of Technology",
			Interests:    "Technology, Innovation",
			Organization: &organization[0],
		},
		{
			Username:     "organizer2",
			Email:        "organizer2@example.com",
			Password:     "abcd1234",
			Role:         "organizer",
			University:   "State University",
			Interests:    "Business, Entrepreneurship",
			Organization: &organization[1],
		},
		{
			Username:     "organizer3",
			Email:        "organizer3@example.com",
			Password:     "abcd1234",
			Role:         "organizer",
			University:   "State University",
			Interests:    "Business, Entrepreneurship",
			Organization: &organization[2],
		},
		{
			Username:   "participant1",
			Email:      "participant1@example.com",
			Password:   "test1234",
			Role:       "student",
			University: "Technical Institute",
			Interests:  "Programming, AI",
		},
		{
			Username:   "participant2",
			Email:      "participant2@example.com",
			Password:   "test1234",
			Role:       "student",
			University: "Arts Academy",
			Interests:  "Design, UI/UX",
		},
		{
			Username:   "admin",
			Email:      "admin@example.com",
			Password:   "test1234",
			Role:       "admin",
			University: "Admin University",
			Interests:  "System Administration",
		},
	}

	// Insert users
	for i := range users_data {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users_data[i].Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		// Ensure each user has a UUID
		user := entity.User{
			Username:   users_data[i].Username,
			Email:      users_data[i].Email,
			Password:   string(hashedPassword),
			Role:       users_data[i].Role,
			University: users_data[i].University,
			Interests:  users_data[i].Interests,
		}
		if users_data[i].Organization != nil {
			user.Organization = users_data[i].Organization
		}

		// Create the user
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}
