package entity

import (
	"time"

	"github.com/google/uuid"
)

type Registration struct {
	UserID        uuid.UUID `json:"userId" gorm:"type:uuid;not null;primaryKey"`
	CompetitionID uuid.UUID `json:"competitionId" gorm:"type:uuid;not null;primaryKey"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
