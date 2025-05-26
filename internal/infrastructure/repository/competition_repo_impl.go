package repository

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type competitionRepository struct {
	db *gorm.DB
}

func NewCompetitionRepository(db *gorm.DB) repository.CompetitionRepository {
	return &competitionRepository{
		db: db,
	}
}

// FindByID implements repository.CompetitionRepository.
func (r *competitionRepository) FindByID(id uuid.UUID) (*entity.Competition, error) {
	var competition entity.Competition
	result := r.db.First(&competition, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &competition, nil
}

// FindAll implements repository.CompetitionRepository.
func (r *competitionRepository) FindAll() ([]*entity.Competition, error) {
	var competitions []*entity.Competition
	result := r.db.Find(&competitions)
	if result.Error != nil {
		return nil, result.Error
	}
	return competitions, nil
}

// Create implements repository.CompetitionRepository.
func (r *competitionRepository) Create(competition *entity.Competition) error {
	return r.db.Create(competition).Error
}

// Update implements repository.CompetitionRepository.
func (r *competitionRepository) Update(competition *entity.Competition) error {
	result := r.db.Save(competition)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete implements repository.CompetitionRepository.
func (r *competitionRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&entity.Competition{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindByOrganizerID implements repository.CompetitionRepository.
func (r *competitionRepository) FindByOrganizerID(organizerID uuid.UUID) ([]*entity.Competition, error) {
	var competitions []*entity.Competition
	result := r.db.Find(&competitions, "organizer_id = ?", organizerID)
	if result.Error != nil {
		return nil, result.Error
	}
	return competitions, nil
}
