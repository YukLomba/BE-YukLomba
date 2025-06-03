package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;"`
	Username           string         `json:"username" gorm:"unique;not null"`
	Email              string         `json:"email" gorm:"unique;not null"`
	Password           string         `json:"-" gorm:"not null"`
	Role               string         `json:"role"`
	University         string         `json:"university"`
	Interests          string         `json:"interests"`
	CreatedAt          time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	OrganizationID     *uuid.UUID     `json:"organization_id"`
	Organization       *Organization  `json:"organization" gorm:"foreignKey:OrganizationID; constraint:OnDelete:SET NULL"`
	JoinedCompetitions []*Competition `json:"joined_competitions" gorm:"many2many:registrations"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
