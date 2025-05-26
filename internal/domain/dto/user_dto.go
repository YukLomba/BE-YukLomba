package dto

import (
	"github.com/google/uuid"
)

type UserProfileUpdate struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	Username   string    `json:"username" validate:"required"`
	Email      string    `json:"email" validate:"required"`
	Role       string    `json:"role" validate:"required"`
	University string    `json:"university" validate:"required"`
	Interests  string    `json:"interests" validate:"required"`
}
