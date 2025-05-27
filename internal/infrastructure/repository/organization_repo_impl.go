package repository

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) repository.OrganizationRepository {
	return &organizationRepository{
		db: db,
	}
}

func (r *organizationRepository) FindByID(id uuid.UUID) (*entity.Organization, error) {
	var organization entity.Organization
	result := r.db.First(&organization, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &organization, nil
}

func (r *organizationRepository) FindAll() ([]*entity.Organization, error) {
	var organizations []*entity.Organization
	result := r.db.Find(&organizations)
	if result.Error != nil {
		return nil, result.Error
	}
	return organizations, nil
}

func (r *organizationRepository) Create(org *entity.Organization) error {
	return r.db.Create(org).Error
}

func (r *organizationRepository) Update(org *entity.Organization) error {
	result := r.db.Save(org)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *organizationRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&entity.Organization{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
