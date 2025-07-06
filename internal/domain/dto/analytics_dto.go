package dto

import "github.com/google/uuid"

type DashboardResponse struct {
	TotalCompetitions  int     `json:"total_competitions"`
	ActiveStudents     int     `json:"active_students"`
	ActiveOrganizers   int     `json:"active_organizers"`
	TotalRegistrations int     `json:"total_registrations"`
	AverageRating      float32 `json:"average_rating"`
}

type CompetitionAnalytics struct {
	CompetitionID      uuid.UUID `json:"competition_id"`
	CompetitionName    string    `json:"competition_name"`
	TotalReviews       int       `json:"total_reviews"`
	AverageRating      float32   `json:"average_rating"`
	TotalRegistrations int       `json:"total_registrations"`
}
