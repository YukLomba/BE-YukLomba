package dto

import (
	"time"

	"github.com/google/uuid"
)

// CompetitionCreateRequest represents the data needed to create a competition
type CompetitionCreateRequest struct {
	Title       string    `json:"title" validate:"required"`
	Type        string    `json:"type" validate:"required"`
	Description string    `json:"description" validate:"required"`
	OrganizerID uuid.UUID `json:"organizer_id" validate:"required"`
	Deadline    time.Time `json:"deadline" validate:"required,future"`
	Category    string    `json:"category" validate:"required"`
	Rules       string    `json:"rules"`
	EventLink   string    `json:"eventLink" validate:"omitempty,url"`
}

// CompetitionUpdateRequest represents the data needed to update a competition
type CompetitionUpdateRequest struct {
	Title       string    `json:"title" validate:"required"`
	Type        string    `json:"type" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Deadline    time.Time `json:"deadline" validate:"required,future"`
	Category    string    `json:"category" validate:"required"`
	Rules       string    `json:"rules"`
	EventLink   string    `json:"eventLink" validate:"omitempty,url"`
	Results     string    `json:"results"`
}

// CompetitionResponse represents the competition data returned to the client
type CompetitionResponse struct {
	ID          uuid.UUID         `json:"id"`
	Title       string            `json:"title"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Organizer   OrganizationShort `json:"organizer"`
	Deadline    time.Time         `json:"deadline"`
	Category    string            `json:"category"`
	Rules       string            `json:"rules"`
	EventLink   string            `json:"eventLink"`
	Results     string            `json:"results"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// CompetitionListResponse represents a list of competitions
type CompetitionListResponse struct {
	Competitions []CompetitionResponse `json:"competitions"`
	Total        int                   `json:"total"`
}
