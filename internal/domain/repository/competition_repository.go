package repository

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/google/uuid"
)

type CompetitionRepository interface {
	FindByID(id uuid.UUID) (*entity.Competition, error)
	FindAll() ([]*entity.Competition, error)
	Create(competition *entity.Competition) error
	Update(competition *entity.Competition) error
	Delete(id uuid.UUID) error
	FindByOrganizerID(organizerID uuid.UUID) ([]*entity.Competition, error)
}
