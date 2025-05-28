package service

import (
	"fmt"
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/google/uuid"
)

type CompetitionService interface {
	GetCompetition(id uuid.UUID) (*dto.CompetitionResponse, error)
	GetAllCompetitions(filter *dto.CompetitionFilter) (*dto.CompetitionListResponse, error)
	CreateCompetition(competition *dto.CompetitionCreateRequest) error
	CreateManyCompetitition(competitions *dto.MultiCompetitionCreateRequest) (*[]string, error)
	UpdateCompetition(id uuid.UUID, competition *dto.CompetitionUpdateRequest) error
	DeleteCompetition(id uuid.UUID) error
	RegisterUserToCompetition(userID uuid.UUID, competitionID uuid.UUID) error
	GetCompetitionsByOrganizer(organizerID uuid.UUID) (*dto.CompetitionListResponse, error)
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
func (s *CompetitionServiceImpl) GetCompetition(id uuid.UUID) (*dto.CompetitionResponse, error) {
	competition, err := s.competitionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toCompetitionResponse(competition), nil
}

// GetAllCompetitions implements CompetitionService.
func (s *CompetitionServiceImpl) GetAllCompetitions(filter *dto.CompetitionFilter) (*dto.CompetitionListResponse, error) {
	var competitions []*entity.Competition
	var err error
	if filter != nil {
		competitions, err = s.competitionRepo.FindWithFilter(filter)
		if err != nil {
			return nil, err
		}
	} else {
		competitions, err = s.competitionRepo.FindAll()
	}
	if err != nil {
		return nil, err
	}

	response := &dto.CompetitionListResponse{
		Total: len(competitions),
	}
	for _, comp := range competitions {
		response.Competitions = append(response.Competitions, *s.toCompetitionResponse(comp))
	}
	return response, nil
}

// CreateCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) CreateCompetition(competition *dto.CompetitionCreateRequest) error {
	entity := &entity.Competition{
		Title:       competition.Title,
		Type:        competition.Type,
		Description: competition.Description,
		OrganizerID: *competition.OrganizerID,
		Deadline:    competition.Deadline,
		Category:    competition.Category,
		Rules:       competition.Rules,
		EventLink:   competition.EventLink,
	}
	return s.competitionRepo.Create(entity)
}

func (s *CompetitionServiceImpl) CreateManyCompetitition(competitions *dto.MultiCompetitionCreateRequest) (*[]string, error) {
	var Competitions []entity.Competition
	var notValidMessage []string

	for _, comp := range competitions.Competitions {

		if comp.OrganizerID == nil {
			notValidMessage = append(notValidMessage, fmt.Sprintf("OrganizerID is required for competition with title %s", comp.Title))
			continue
		}
		if comp.Deadline.Before(time.Now()) {
			notValidMessage = append(notValidMessage, fmt.Sprintf("Deadline must be after %s for competition with title %s", time.Now().Format("2006-01-02"), comp.Title))
			continue
		}
		entity := entity.Competition{
			Title:       comp.Title,
			Type:        comp.Type,
			Description: comp.Description,
			OrganizerID: *comp.OrganizerID,
			Deadline:    comp.Deadline,
			Category:    comp.Category,
			Rules:       comp.Rules,
			EventLink:   comp.EventLink,
		}
		Competitions = append(Competitions, entity)
	}
	err := s.competitionRepo.CreateMany(&Competitions)
	if err != nil {
		return nil, err
	}
	if len(notValidMessage) > 0 {
		return &notValidMessage, nil
	}
	return nil, nil
}

// UpdateCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) UpdateCompetition(id uuid.UUID, competition *dto.CompetitionUpdateRequest) error {
	existing, err := s.competitionRepo.FindByID(id)
	if err != nil {
		return err
	}

	existing.Title = competition.Title
	existing.Type = competition.Type
	existing.Description = competition.Description
	existing.Deadline = competition.Deadline
	existing.Category = competition.Category
	existing.Rules = competition.Rules
	existing.EventLink = competition.EventLink
	existing.Results = competition.Results

	return s.competitionRepo.Update(existing)
}

// DeleteCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) DeleteCompetition(id uuid.UUID) error {
	return s.competitionRepo.Delete(id)
}

// GetCompetitionsByOrganizer implements CompetitionService.
func (s *CompetitionServiceImpl) GetCompetitionsByOrganizer(organizerID uuid.UUID) (*dto.CompetitionListResponse, error) {
	competitions, err := s.competitionRepo.FindByOrganizerID(organizerID)
	if err != nil {
		return nil, err
	}

	response := &dto.CompetitionListResponse{
		Total: len(competitions),
	}
	for _, comp := range competitions {
		response.Competitions = append(response.Competitions, *s.toCompetitionResponse(comp))
	}
	return response, nil
}

func (s *CompetitionServiceImpl) toCompetitionResponse(competition *entity.Competition) *dto.CompetitionResponse {
	return &dto.CompetitionResponse{
		ID:          competition.ID,
		Title:       competition.Title,
		Type:        competition.Type,
		Description: competition.Description,
		Organizer: dto.OrganizationShort{
			ID:   competition.OrganizerID,
			Name: competition.Organizer.Name,
		},
		Deadline:  competition.Deadline,
		Category:  competition.Category,
		Rules:     competition.Rules,
		EventLink: competition.EventLink,
		Results:   competition.Results,
		CreatedAt: competition.CreatedAt,
		UpdatedAt: competition.UpdatedAt,
	}
}

// RegisterUserToCompetition implements CompetitionService.
func (s *CompetitionServiceImpl) RegisterUserToCompetition(userID uuid.UUID, competitionID uuid.UUID) error {
	registration := &entity.Registration{
		UserID:        userID,
		CompetitionID: competitionID,
	}
	return s.competitionRepo.CreateUserRegistration(registration)
}
