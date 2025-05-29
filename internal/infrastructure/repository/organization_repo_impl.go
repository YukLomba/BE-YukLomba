package repository

import (
	"log/slog"

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
		slog.Error("Error finding organization by ID:",
			"id", id,
			"error", result.Error,
		)
		return nil, result.Error
	}
	return &organization, nil
}

func (r *organizationRepository) FindAll() ([]*entity.Organization, error) {
	var organizations []*entity.Organization
	result := r.db.Find(&organizations)
	if result.Error != nil {
		slog.Error("Error finding all organizations:",
			"error", result.Error,
		)
		return nil, result.Error
	}
	return organizations, nil
}

func (r *organizationRepository) Create(org *entity.Organization) error {
	result := r.db.Create(org)
	if result.Error != nil {
		slog.Error("Error creating organization:",
			"organization", org,
			"error", result.Error,
		)
		return result.Error
	}
	return nil
}

func (r *organizationRepository) Update(org *entity.Organization) error {
	result := r.db.Save(org)
	if result.Error != nil {
		slog.Error("Error updating organization:",
			"organization", org,
			"error", result.Error,
		)
		return result.Error
	}
	if result.RowsAffected == 0 {
		slog.Warn("No organization updated")
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *organizationRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&entity.Organization{}, "id = ?", id)
	if result.Error != nil {
		slog.Error("Error deleting organization:",
			"id", id,
			"error", result.Error,
		)
		return result.Error
	}
	if result.RowsAffected == 0 {
		slog.Warn("No organization deleted")
		return gorm.ErrRecordNotFound
	}
	return nil
}
