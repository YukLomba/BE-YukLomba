package dto

import (
	"time"

	"github.com/google/uuid"
)

// CompetitionCreateRequest represents the data needed to create a competition
type CompetitionCreateRequest struct {
	Title       string     `json:"title" validate:"required"`
	Type        string     `json:"type" validate:"required"`
	Description string     `json:"description" validate:"required"`
	Image       *[]string  `json:"image" validate:"dive,url"`
	OrganizerID *uuid.UUID `json:"organizer_id" validate:"required"`
	Deadline    time.Time  `json:"deadline" validate:"required,future"`
	Category    string     `json:"category" validate:"required"`
	EventLink   string     `json:"eventLink" validate:"required,url"`
}

type MultiCompetitionCreateRequest struct {
	Competitions []*CompetitionCreateRequest `json:"competitions" validate:"required,dive"`
}

// CompetitionUpdateRequest represents the data needed to update a competition
type CompetitionUpdateRequest struct {
	Title       *string    `json:"title" validate:"required"`
	Type        *string    `json:"type" validate:"required"`
	Description *string    `json:"description" validate:"required"`
	Image       *[]string  `json:"image" validate:"dive,url"`
	Deadline    *time.Time `json:"deadline" validate:"required,future"`
	Category    *string    `json:"category" validate:"required"`
	EventLink   *string    `json:"eventLink" validate:"url"`
	Results     *string    `json:"results"`
}

// CompetitionResponse represents the competition data returned to the client
type CompetitionResponse struct {
	ID          uuid.UUID         `json:"id"`
	Title       string            `json:"title"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Image       *[]string         `json:"image"`
	Organizer   OrganizationShort `json:"organizer"`
	Deadline    time.Time         `json:"deadline"`
	Category    string            `json:"category"`
	EventLink   string            `json:"eventLink"`
	Results     string            `json:"results"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// // CompetitionListResponse represents a list of competitions
// type CompetitionListResponse struct {
// 	Competitions []CompetitionResponse `json:"competitions"`
// }

type CompetitionFilter struct {
	Title    *string    `query:"title" validate:"min=3"`
	Type     *string    `query:"type"`
	Category *string    `query:"category"`
	Before   *time.Time `query:"before"` // untuk deadline sebelum tanggal tertentu
	After    *time.Time `query:"after"`  // untuk deadline setelah tanggal tertentu
}

type CompetitionShort struct {
	ID          uuid.UUID         `json:"id"`
	Title       string            `json:"title"`
	Type        string            `json:"type"`
	Organizer   OrganizationShort `json:"organization"`
	Description string            `json:"description"`
	Deadline    time.Time         `json:"deadline"`
	Category    string            `json:"category"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}
