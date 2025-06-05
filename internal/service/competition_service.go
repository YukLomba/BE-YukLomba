package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	errs "github.com/YukLomba/BE-YukLomba/internal/domain/error"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrCompetitionNotFound          = errors.New("competition not found")
	ErrCompetitionAlreadyExists     = errors.New("competition already exists")
	ErrCompetitionNotBelongsToOrg   = errors.New("competition does not belong to organization")
	ErrCompetitionAlreadyRegistered = errors.New("user already registered to competition")
	ErrCompetitionDeadlinePassed    = errors.New("competition deadline has passed")
)

type CompetitionService interface {
	GetCompetition(id uuid.UUID) (*entity.Competition, error)
	GetAllCompetitions(filter *dto.CompetitionFilter) ([]*entity.Competition, error)
	CreateCompetition(authInfo *dto.AuthInfo, competition *entity.Competition) error
	CreateManyCompetitition(authInfo *dto.AuthInfo, competitions []*entity.Competition) (*[]string, error)
	UpdateCompetition(authInfo *dto.AuthInfo, id uuid.UUID, data *map[string]interface{}) error
	DeleteCompetition(authInfo *dto.AuthInfo, id uuid.UUID) error
	RegisterUserToCompetition(authInfo *dto.AuthInfo, competitionID uuid.UUID) error
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
	competition, err := s.competitionRepo.FindByID(id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrCompetitionNotFound
		default:
			return nil, errs.ErrInternalServer
		}
	}
	return competition, nil
}

// GetAllCompetitions implements CompetitionService.
func (s *CompetitionServiceImpl) GetAllCompetitions(filter *dto.CompetitionFilter) ([]*entity.Competition, error) {
	competitions, err := s.competitionRepo.FindWithFilter(filter)
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	return competitions, nil
}

// CreateCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) CreateCompetition(authInfo *dto.AuthInfo, competition *entity.Competition) error {
	if competition.Deadline.Before(time.Now()) {
		return ErrCompetitionDeadlinePassed
	}
	if authInfo.OrganizationID == nil {
		return errs.ErrUnauthorized
	}
	competition.OrganizerID = *authInfo.OrganizationID
	return s.competitionRepo.Create(competition)
}

func (s *CompetitionServiceImpl) CreateManyCompetitition(authInfo *dto.AuthInfo, competitions []*entity.Competition) (*[]string, error) {
	Competitions := new([]entity.Competition)
	var notValidMessage []string
	if authInfo.Role != "admin" {
		return nil, errs.ErrUnauthorized
	}

	for _, comp := range competitions {
		if comp.Deadline.Before(time.Now()) {
			notValidMessage = append(notValidMessage, fmt.Sprintf("Deadline must be after %s for competition with title %s", time.Now().Format("2006-01-02"), comp.Title))
			continue
		}
		*Competitions = append(*Competitions, *comp)
	}
	err := s.competitionRepo.CreateMany(Competitions)
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	if len(notValidMessage) > 0 {
		return &notValidMessage, nil
	}
	return nil, nil
}

// UpdateCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) UpdateCompetition(authInfo *dto.AuthInfo, id uuid.UUID, data *map[string]interface{}) error {
	competition, err := s.competitionRepo.FindByID(id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrCompetitionNotFound
		default:
			return errs.ErrInternalServer
		}
	}

	if *(*authInfo).OrganizationID != competition.OrganizerID {
		return ErrCompetitionNotBelongsToOrg
	}

	err = s.competitionRepo.Update(id, data)
	if err != nil {
		return errs.ErrInternalServer
	}
	return nil
}

// DeleteCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) DeleteCompetition(authInfo *dto.AuthInfo, id uuid.UUID) error {
	competition, err := s.competitionRepo.FindByID(id)
	if err != nil {
		return ErrCompetitionNotFound
	}
	if *authInfo.OrganizationID != competition.OrganizerID {
		return ErrCompetitionNotBelongsToOrg
	}
	err = s.competitionRepo.Delete(id)
	if err != nil {
		return errs.ErrInternalServer
	}
	return nil
}

// GetCompetitionsByOrganizer implements CompetitionService.
func (s *CompetitionServiceImpl) GetCompetitionsByOrganizer(organizerID uuid.UUID) ([]*entity.Competition, error) {
	competitions, err := s.competitionRepo.FindByOrganizerID(organizerID)
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	return competitions, nil
}

// RegisterUserToCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) RegisterUserToCompetition(authInfo *dto.AuthInfo, competitionID uuid.UUID) error {
	competition, err := s.competitionRepo.FindByID(competitionID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrCompetitionNotFound
		default:
			return errs.ErrInternalServer
		}
	}

	if competition.Deadline.Before(time.Now()) {
		return ErrCompetitionDeadlinePassed
	}

	_, err = s.competitionRepo.FindUserRegistration(competitionID, authInfo.ID)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrInternalServer
		}
	} else {
		// User registration found, already registered
		return ErrCompetitionAlreadyRegistered
	}

	registration := &entity.Registration{
		UserID:        authInfo.ID,
		CompetitionID: competitionID,
	}

	err = s.competitionRepo.CreateUserRegistration(registration)

	if err != nil {
		return errs.ErrInternalServer
	}
	return nil
}
