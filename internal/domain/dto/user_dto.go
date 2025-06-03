package dto

import "time"

type UserProfileUpdate struct {
	Username   *string `json:"username" validate:"required"`
	Password   *string `json:"password" validate:"required"`
	Email      *string `json:"email" validate:"required"`
	Role       *string `json:"role" validate:"required"`
	University *string `json:"university" validate:"required"`
	Interests  *string `json:"interests" validate:"required"`
}

type UserResponse struct {
	ID                 string              `json:"id" `
	Username           string              `json:"username" `
	Email              string              `json:"email" `
	University         string              `json:"university"`
	Interests          string              `json:"interests"`
	CreatedAt          time.Time           `json:"created_at"`
	Organization       *OrganizationShort  `json:"organization,omitempty" `
	JoinedCompetitions []*CompetitionShort `json:"joined_competitions" `
}
