package dto

import (
	"time"

	"github.com/google/uuid"
)

// OrganizationCreateRequest represents the data needed to create an organization
type OrganizationCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Logo        string `json:"logo" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// OrganizationUpdateRequest represents the data needed to update an organization
type OrganizationUpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Logo        string `json:"logo" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// OrganizationResponse represents the organization data returned to the client
type OrganizationResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Logo        string    `json:"logo"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// OrganizationShort represents a simplified organization for response
type OrganizationShort struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// OrganizationListResponse represents a list of organizations
type OrganizationListResponse struct {
	Organizations []OrganizationResponse `json:"organizations"`
	Total         int                    `json:"total"`
}
