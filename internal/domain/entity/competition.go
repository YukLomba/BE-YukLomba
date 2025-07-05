package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Competition represents a competition entity in the system
type Competition struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;"`
	Title       string       `json:"title" gorm:"not null"`
	Type        string       `json:"type" gorm:"not null"`
	Description string       `json:"description" gorm:"type:text"`
	Image       *[]string    `json:"image" gorm:"serializer:json" validate:"required,dive,url"`
	OrganizerID uuid.UUID    `json:"organizer_id,omitempty" gorm:"type:uuid;not null"`
	Organizer   Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizerID;constraint:OnDelete:SET NULL"`
	Deadline    time.Time    `json:"deadline" gorm:"not null" validate:"required,future"`
	Category    string       `json:"category" gorm:"not null"`
	EventLink   string       `json:"eventLink" gorm:"column:event_link" validate:"omitempty,url"`
	Results     string       `json:"results" gorm:"type:text"`
	Registrant  []*User      `gorm:"many2many:registrations"`
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
