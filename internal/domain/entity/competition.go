package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Competition represents a competition entity in the system
type Competition struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;" validate:"required"`
	Title       string       `json:"title" gorm:"not null" validate:"required"`
	Type        string       `json:"type" gorm:"not null" validate:"required"`
	Description string       `json:"description" gorm:"type:text" validate:"required"`
	OrganizerID uuid.UUID    `json:"organizer_id" gorm:"type:uuid;not null"`
	Organizer   Organization `json:"organization" gorm:"foreignKey:OrganizerID;constraint:OnDelete:SET NULL"`
	Deadline    time.Time    `json:"deadline" gorm:"not null" validate:"required,future"`
	Category    string       `json:"category" gorm:"not null" validate:"required"`
	Rules       string       `json:"rules" gorm:"type:text"`
	EventLink   string       `json:"eventLink" gorm:"column:event_link" validate:"omitempty,url"`
	Results     string       `json:"results" gorm:"type:text"`
	CreatedAt   time.Time    `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updatedAt" gorm:"autoUpdateTime"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Competition) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
