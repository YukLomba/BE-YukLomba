package service

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/google/uuid"
)

type AnalyticsService interface {
	GetDashboard(authInfo *dto.AuthInfo) (*dto.DashboardResponse, error)
	GetCompetitionAnalytics(competitionID uuid.UUID) (*dto.CompetitionAnalytics, error)
}

type AnalyticsServiceImpl struct {
	userRepo        repository.UserRepository
	competitionRepo repository.CompetitionRepository
	reviewRepo      repository.ReviewRepository
}

func NewAnalyticsService(
	userRepo repository.UserRepository,
	competitionRepo repository.CompetitionRepository,
	reviewRepo repository.ReviewRepository,
) AnalyticsService {
	return &AnalyticsServiceImpl{
		userRepo:        userRepo,
		competitionRepo: competitionRepo,
		reviewRepo:      reviewRepo,
	}
}

func (s *AnalyticsServiceImpl) GetDashboard(authInfo *dto.AuthInfo) (*dto.DashboardResponse, error) {
	var competitions []*entity.Competition
	var err error

	if authInfo.Role == "organizer" {
		competitions, err = s.competitionRepo.FindByOrganizerID(*authInfo.OrganizationID)
	} else {
		competitions, err = s.competitionRepo.FindAll()
	}
	if err != nil {
		return nil, err
	}

	students, err := s.userRepo.CountByRole("student")
	if err != nil {
		return nil, err
	}

	organizers, err := s.userRepo.CountByRole("organizer")
	if err != nil {
		return nil, err
	}

	registrations, err := s.competitionRepo.CountAllRegistrations()
	if err != nil {
		return nil, err
	}

	avgRating, err := s.reviewRepo.GetAverageRatingAll()
	if err != nil {
		return nil, err
	}

	return &dto.DashboardResponse{
		TotalCompetitions:  len(competitions),
		ActiveStudents:     students,
		ActiveOrganizers:   organizers,
		TotalRegistrations: registrations,
		AverageRating:      avgRating,
	}, nil
}

func (s *AnalyticsServiceImpl) GetCompetitionAnalytics(competitionID uuid.UUID) (*dto.CompetitionAnalytics, error) {
	competition, err := s.competitionRepo.FindByID(competitionID)
	if err != nil {
		return nil, err
	}

	reviews, err := s.reviewRepo.GetByCompetition(competitionID)
	if err != nil {
		return nil, err
	}

	registrations, err := s.competitionRepo.CountRegistrations(competitionID)
	if err != nil {
		return nil, err
	}

	var totalRating float32
	for _, review := range reviews {
		totalRating += float32(review.Rating)
	}
	avgRating := float32(0)
	if len(reviews) > 0 {
		avgRating = totalRating / float32(len(reviews))
	}

	return &dto.CompetitionAnalytics{
		CompetitionID:      competitionID,
		CompetitionName:    competition.Title,
		TotalReviews:       len(reviews),
		AverageRating:      avgRating,
		TotalRegistrations: registrations,
	}, nil
}
