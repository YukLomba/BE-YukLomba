package service

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/google/uuid"
)

type CompetitionService interface {
	GetCompetition(id uuid.UUID) (*entity.Competition, error)
	GetAllCompetitions() ([]*entity.Competition, error)
	CreateCompetition(competition *entity.Competition) error
	UpdateCompetition(competition *entity.Competition) error
	DeleteCompetition(id uuid.UUID) error
	GetCompetitionsByOrganizer(organizerID uuid.UUID) ([]*entity.Competition, error)
}

type CompetitionServiceImpl struct {
	competitionRepo repository.CompetitionRepository
}

func NewCompetitionService(competitionRepo repository.CompetitionRepository) CompetitionService {
	return &CompetitionServiceImpl{
		competitionRepo: competitionRepo,
	}
}

// GetCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) GetCompetition(id uuid.UUID) (*entity.Competition, error) {
	return s.competitionRepo.FindByID(id)
}

// GetAllCompetitions implements CompetitionService.
func (s *CompetitionServiceImpl) GetAllCompetitions() ([]*entity.Competition, error) {
	return s.competitionRepo.FindAll()
}

// CreateCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) CreateCompetition(competition *entity.Competition) error {
	return s.competitionRepo.Create(competition)
}

// UpdateCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) UpdateCompetition(competition *entity.Competition) error {
	return s.competitionRepo.Update(competition)
}

// DeleteCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) DeleteCompetition(id uuid.UUID) error {
	return s.competitionRepo.Delete(id)
}

// GetCompetitionsByOrganizer implements CompetitionService.
func (s *CompetitionServiceImpl) GetCompetitionsByOrganizer(organizerID uuid.UUID) ([]*entity.Competition, error) {
	return s.competitionRepo.FindByOrganizerID(organizerID)
}
