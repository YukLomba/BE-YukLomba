package repository

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/google/uuid"
)

type CompetitionRepository interface {
	FindByID(id uuid.UUID) (*entity.Competition, error)
	FindAll() ([]*entity.Competition, error)
	Create(competition *entity.Competition) error
	CreateMany(competitions *[]entity.Competition) error
	Update(id uuid.UUID, data *map[string]interface{}) error
	Delete(id uuid.UUID) error
	CreateUserRegistration(registration *entity.Registration) error
	FindByOrganizerID(organizerID uuid.UUID) ([]*entity.Competition, error)
	FindWithFilter(filter *dto.CompetitionFilter) ([]*entity.Competition, error)
	FindUserRegistration(competitionID uuid.UUID, userID uuid.UUID) (*entity.Registration, error)
}
