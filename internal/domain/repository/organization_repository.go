package repository

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/google/uuid"
)

type OrganizationRepository interface {
	FindByID(id uuid.UUID) (*entity.User, error)
	FindAll() ([]*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uuid.UUID) error
}
