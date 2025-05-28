package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Competition represents a competition entity in the system
type Organization struct {
	ID                    uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;"`
	Name                  string        `json:"title" gorm:"not null"`
	Logo                  string        `json:"logo" gorm:"type:text"`
	Description           string        `json:"description" gorm:"type:text"`
	CreatedAt             time.Time     `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt             time.Time     `json:"updatedAt" gorm:"autoUpdateTime"`
	OrganizedCompetitions []Competition `json:"competitions" gorm:"foreignKey:OrganizerID;constraint:OnDelete:SET NULL"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (O *Organization) BeforeCreate(tx *gorm.DB) error {
	if O.ID == uuid.Nil {
		O.ID = uuid.New()
	}
	return nil
}
