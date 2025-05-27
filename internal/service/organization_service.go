package service

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/google/uuid"
)

type OrganizationService interface {
	GetOrganization(id uuid.UUID) (*dto.OrganizationResponse, error)
	GetAllOrganizations() (*dto.OrganizationListResponse, error)
	CreateOrganization(org *dto.OrganizationCreateRequest) error
	UpdateOrganization(id uuid.UUID, org *dto.OrganizationUpdateRequest) error
	DeleteOrganization(id uuid.UUID) error
}

type OrganizationServiceImpl struct {
	orgRepo repository.OrganizationRepository
}

func NewOrganizationService(orgRepo repository.OrganizationRepository) OrganizationService {
	return &OrganizationServiceImpl{
		orgRepo: orgRepo,
	}
}

func (s *OrganizationServiceImpl) GetOrganization(id uuid.UUID) (*dto.OrganizationResponse, error) {
	org, err := s.orgRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toOrganizationResponse(org), nil
}

func (s *OrganizationServiceImpl) GetAllOrganizations() (*dto.OrganizationListResponse, error) {
	orgs, err := s.orgRepo.FindAll()
	if err != nil {
		return nil, err
	}

	response := &dto.OrganizationListResponse{
		Total: len(orgs),
	}
	for _, org := range orgs {
		response.Organizations = append(response.Organizations, *s.toOrganizationResponse(org))
	}
	return response, nil
}

func (s *OrganizationServiceImpl) CreateOrganization(org *dto.OrganizationCreateRequest) error {
	entity := &entity.Organization{
		Name:        org.Name,
		Logo:        org.Logo,
		Description: org.Description,
	}
	return s.orgRepo.Create(entity)
}

func (s *OrganizationServiceImpl) UpdateOrganization(id uuid.UUID, org *dto.OrganizationUpdateRequest) error {
	existing, err := s.orgRepo.FindByID(id)
	if err != nil {
		return err
	}

	existing.Name = org.Name
	existing.Logo = org.Logo
	existing.Description = org.Description

	return s.orgRepo.Update(existing)
}

func (s *OrganizationServiceImpl) DeleteOrganization(id uuid.UUID) error {
	return s.orgRepo.Delete(id)
}

func (s *OrganizationServiceImpl) toOrganizationResponse(org *entity.Organization) *dto.OrganizationResponse {
	return &dto.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Logo:        org.Logo,
		Description: org.Description,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}
}
