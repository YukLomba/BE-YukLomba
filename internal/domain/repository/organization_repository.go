package repository

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/google/uuid"
)

type OrganizationRepository interface {
	FindByID(id uuid.UUID) (*entity.Organization, error)
	FindAll() ([]*entity.Organization, error)
	Create(org *entity.Organization) error
	Update(id uuid.UUID, data *map[string]interface{}) error
	Delete(id uuid.UUID) error
}
