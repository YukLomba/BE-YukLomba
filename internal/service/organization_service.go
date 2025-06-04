package service

import (
	"errors"
	"slices"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	errs "github.com/YukLomba/BE-YukLomba/internal/domain/error"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrOrganizationExists   = errors.New("organization already exists")
)

type OrganizationService interface {
	GetOrganization(id uuid.UUID) (*entity.Organization, error)
	GetAllOrganizations() ([]*entity.Organization, error)
	CreateOrganization(org *entity.Organization, authInfo *dto.AuthInfo) error
	UpdateOrganization(org *entity.Organization, authInfo *dto.AuthInfo) error
	DeleteOrganization(id uuid.UUID, authInfo *dto.AuthInfo) error
}

type OrganizationServiceImpl struct {
	orgRepo repository.OrganizationRepository
}

func NewOrganizationService(orgRepo repository.OrganizationRepository) OrganizationService {
	return &OrganizationServiceImpl{
		orgRepo: orgRepo,
	}
}

func (s *OrganizationServiceImpl) GetOrganization(id uuid.UUID) (*entity.Organization, error) {
	org, err := s.orgRepo.FindByID(id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrOrganizationNotFound
		default:
			return nil, errs.ErrInternalServer
		}
	}
	return org, nil
}

func (s *OrganizationServiceImpl) GetAllOrganizations() ([]*entity.Organization, error) {
	orgs, err := s.orgRepo.FindAll()
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	return orgs, nil
}

func (s *OrganizationServiceImpl) CreateOrganization(org *entity.Organization, authInfo *dto.AuthInfo) error {
	authorizedRoles := []string{"admin", "organizer"}
	if slices.Contains(authorizedRoles, authInfo.Role) {
		return errs.ErrUnauthorized
	}
	if err := s.orgRepo.Create(org); err != nil {
		return errs.ErrInternalServer
	}
	return nil
}

func (s *OrganizationServiceImpl) UpdateOrganization(org *entity.Organization, authInfo *dto.AuthInfo) error {
	authorizedRoles := []string{"admin", "organizer"}
	if slices.Contains(authorizedRoles, authInfo.Role) {
		return errs.ErrUnauthorized
	}
	if *authInfo.OrganizationID != org.ID || authInfo.Role != "admin" {
		return errs.ErrUnauthorized
	}
	org, err := s.orgRepo.FindByID(org.ID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrOrganizationNotFound
		default:
			return errs.ErrInternalServer
		}
	}

	if err := s.orgRepo.Update(org); err != nil {
		return errs.ErrInternalServer
	}
	return nil
}

func (s *OrganizationServiceImpl) DeleteOrganization(id uuid.UUID, authInfo *dto.AuthInfo) error {
	authorizedRoles := []string{"admin", "organizer"}
	if slices.Contains(authorizedRoles, authInfo.Role) {
		return errs.ErrUnauthorized
	}
	if *authInfo.OrganizationID != id || authInfo.Role != "admin" {
		return errs.ErrUnauthorized
	}

	_, err := s.orgRepo.FindByID(id)

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrOrganizationNotFound
		default:
			return errs.ErrInternalServer
		}

	}

	if err := s.orgRepo.Delete(id); err != nil {
		return errs.ErrInternalServer
	}
	return nil
}
