package repository

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindByID(id uuid.UUID) (*entity.User, error)
	FindAll() ([]*entity.User, error)
	Create(user *entity.User) error
	Update(userID uuid.UUID, data *map[string]interface{}) error
	FindAllRegistration(id uuid.UUID) ([]*entity.Competition, error)
	FindByEmail(email string) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
}
