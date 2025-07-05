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
	ErrMaximumOrganizations = errors.New("maximum organizations reached")
)

type OrganizationService interface {
	GetOrganization(id uuid.UUID) (*entity.Organization, error)
	GetAllOrganizations() ([]*entity.Organization, error)
	CreateOrganization(org *entity.Organization, authInfo *dto.AuthInfo) error
	UpdateOrganization(authInfo *dto.AuthInfo, id uuid.UUID, data *map[string]interface{}) error
	DeleteOrganization(id uuid.UUID, authInfo *dto.AuthInfo) error
}

type OrganizationServiceImpl struct {
	orgRepo  repository.OrganizationRepository
	userRepo repository.UserRepository
}

func NewOrganizationService(orgRepo repository.OrganizationRepository, userRepo repository.UserRepository) OrganizationService {
	return &OrganizationServiceImpl{
		orgRepo:  orgRepo,
		userRepo: userRepo,
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
	if !slices.Contains(authorizedRoles, authInfo.Role) {
		return errs.ErrUnauthorized
	}

	if (*authInfo).OrganizationID != nil && (*authInfo).Role == "organizer" {
		return ErrMaximumOrganizations
	}
	if err := s.orgRepo.Create(org); err != nil {
		return errs.ErrInternalServer
	}
	if (*authInfo).Role == "organizer" && (*authInfo).OrganizationID == nil {
		data := map[string]interface{}{
			"organization_id": org.ID,
		}
		if err := s.userRepo.Update((*authInfo).ID, &data); err != nil {
			return errs.ErrInternalServer
		}
	}
	return nil
}

func (s *OrganizationServiceImpl) UpdateOrganization(authInfo *dto.AuthInfo, id uuid.UUID, data *map[string]interface{}) error {
	org, err := s.orgRepo.FindByID(id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrOrganizationNotFound
		default:
			return errs.ErrInternalServer
		}
	}
	if *(*authInfo).OrganizationID != org.ID && (*authInfo).Role == "organizer" {
		return errs.ErrUnauthorized
	}

	if err := s.orgRepo.Update(id, data); err != nil {
		return errs.ErrInternalServer
	}
	return nil
}

func (s *OrganizationServiceImpl) DeleteOrganization(id uuid.UUID, authInfo *dto.AuthInfo) error {
	if *(*authInfo).OrganizationID != id && authInfo.Role == "organizer" {
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

	err = s.orgRepo.Delete(id)
	if err != nil {
		return errs.ErrInternalServer
	}
	// delete user organization_id if organizer
	if (*authInfo).Role == "organizer" && (*authInfo).OrganizationID == nil {
		data := map[string]interface{}{
			"organization_id": nil,
		}
		if err := s.userRepo.Update((*authInfo).ID, &data); err != nil {
			return errs.ErrInternalServer
		}
	}
	return nil
}
