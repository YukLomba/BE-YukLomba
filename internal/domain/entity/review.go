package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID `gorm:"type:uuid;not null"`
	CompetitionID uuid.UUID `gorm:"type:uuid;not null"`
	Rating        int       `gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment       string    `gorm:"type:text"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`

	User        User        `gorm:"foreignKey:UserID"`
	Competition Competition `gorm:"foreignKey:CompetitionID"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
