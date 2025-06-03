package seeder

import (
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedCompetitions seeds the competitions table with sample data
func SeedCompetitions(db *gorm.DB) error {
	// Check if competitions already exist
	var count int64
	db.Model(&entity.Competition{}).Count(&count)
	if count > 0 {
		return nil // Skip seeding if competitions already exist
	}

	// get organizers
	var orgsList []entity.Organization
	if err := db.Find(&orgsList).Error; err != nil {
		return err
	}

	organization := make(map[string]entity.Organization)
	for _, org := range orgsList {
		organization[org.Name] = org
	}

	// Define sample competitions
	competitions := []entity.Competition{
		{
			Title:       "Web Development Hackathon",
			Type:        "Hackathon",
			Description: "A 48-hour hackathon focused on creating innovative web applications.",
			Image:       &[]string{"https://lombasma.com/wp-content/uploads/2022/12/lomba-sma-768x799.png", "https://th.bing.com/th/id/OIP.Gg4-8VRGkmSTQGBqgaryGgHaKe?cb=iwc2&pid=ImgDet&w=474&h=670&rs=1"},
			Deadline:    time.Now().AddDate(0, 1, 0), // 1 month from now
			Organizer:   organization["Google LLC"],
			Category:    "Technology",
			EventLink:   "https://example.com/hackathon",
		},
		{
			Title:       "Mobile App Design Competition",
			Type:        "Design",
			Description: "Design a mobile app that solves a real-world problem.",
			Deadline:    time.Now().AddDate(0, 2, 0), // 2 months from now
			Category:    "Design",
			Organizer:   organization["University of Indonesia"],
			EventLink:   "https://example.com/design-competition",
		},
		{
			Title:       "Business Plan Competition",
			Type:        "Business",
			Description: "Present a business plan for a startup idea with potential for growth and impact.",
			Deadline:    time.Now().AddDate(0, 3, 0), // 3 months from now
			Category:    "Business",
			Organizer:   organization["University of Indonesia"],
			EventLink:   "https://example.com/business-plan",
		},
		{
			Title:       "AI Research Paper Competition",
			Type:        "Research",
			Description: "Submit original research papers on artificial intelligence and machine learning.",
			Deadline:    time.Now().AddDate(0, 4, 0), // 4 months from now
			Category:    "Technology",
			Organizer:   organization["University of Indonesia"],
			EventLink:   "https://example.com/ai-research",
		},
		{
			Title:       "Sustainable Innovation Challenge",
			Type:        "Innovation",
			Description: "Develop innovative solutions for environmental sustainability challenges.",
			Deadline:    time.Now().AddDate(0, 5, 0), // 5 months from now
			Category:    "Environment",
			Organizer:   organization["Bandung Institute of Technology"],
			EventLink:   "https://example.com/sustainability",
		},
	}

	// Insert competitions
	for i := range competitions {
		// Ensure each competition has a UUID
		competitions[i].ID = uuid.New()

		// Create the competition
		if err := db.Create(&competitions[i]).Error; err != nil {
			return err
		}
	}

	return nil
}
