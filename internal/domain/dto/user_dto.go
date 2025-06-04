package dto

import "time"

type UserProfileUpdate struct {
	Username   *string `json:"username"`
	Password   *string `json:"password"`
	Email      *string `json:"email"`
	Role       *string `json:"role"`
	University *string `json:"university"`
	Interests  *string `json:"interests"`
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
