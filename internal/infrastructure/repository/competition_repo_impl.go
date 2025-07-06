package repository

import (
	"encoding/json"
	"log/slog"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
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
	result := r.db.Preload("Organizer").First(&competition, "id = ?", id)
	if result.Error != nil {
		slog.Error("Error finding competition by ID:",
			"id", id,
			"error", result.Error,
		)
		return nil, result.Error
	}
	return &competition, nil
}

// FindAll implements repository.CompetitionRepository.
func (r *competitionRepository) FindAll() ([]*entity.Competition, error) {
	var competitions []*entity.Competition
	result := r.db.Find(&competitions)
	if result.Error != nil {
		slog.Error("Error finding all competitions:",
			"error", result.Error,
		)
		return nil, result.Error
	}
	return competitions, nil
}

// Create implements repository.CompetitionRepository.
func (r *competitionRepository) Create(competition *entity.Competition) error {
	result := r.db.Create(competition)
	if result.Error != nil {
		slog.Error("Error creating competition:",
			"competition", competition,
			"error", result.Error,
		)
		return result.Error
	}
	return nil
}
func (r *competitionRepository) CreateMany(competition *[]entity.Competition) error {
	len := len(*competition)
	result := r.db.CreateInBatches(competition, len)
	if result.Error != nil {
		slog.Error("Error creating competitions:",
			"competitions", competition,
			"error", result.Error,
		)
		return result.Error
	}
	return nil
}

// Update implements repository.CompetitionRepository.
func (r *competitionRepository) Update(id uuid.UUID, data *map[string]interface{}) error {
	result := r.db.Model(&entity.Competition{}).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		slog.Error("Error updating competition:",
			"competition", data,
			"error", result.Error,
		)
		return result.Error
	}
	if result.RowsAffected == 0 {
		slog.Warn("No competition updated")
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete implements repository.CompetitionRepository.
func (r *competitionRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&entity.Competition{}, "id = ?", id)
	if result.Error != nil {
		slog.Error("Error deleting competition:",
			"id", id,
			"error", result.Error,
		)
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
	result := r.db.Preload("Organizer").Find(&competitions, "organizer_id = ?", organizerID)
	if result.Error != nil {
		slog.Error("Error finding competitions by organizer ID:",
			"organizerID", organizerID,
			"error", result.Error,
		)
		return nil, result.Error
	}
	return competitions, nil
}
func (r *competitionRepository) FindWithFilter(filter *dto.CompetitionFilter) ([]*entity.Competition, error) {
	var competitions []*entity.Competition
	query := r.db.Model(&entity.Competition{})
	if filter != nil {
		if filter.Title != nil && *filter.Title != "" {
			query = query.Where("title LIKE ?", "%"+*filter.Title+"%")
		}
		if filter.Category != nil && *filter.Category != "" {
			query = query.Where("category = ?", *filter.Category)
		}
		if filter.Type != nil {
			query = query.Where("type = ?", *filter.Type)
		}
		if filter.Before != nil {
			query = query.Where("deadline <= ?", *filter.Before)
		}
		if filter.After != nil {
			query = query.Where("deadline >= ?", *filter.After)
		}
	}
	result := query.Preload("Organizer").Find(&competitions)
	if result.Error != nil {
		filterJson, err := json.Marshal(filter)
		if err != nil {
			slog.Error("Error marshalling filter:",
				"error", err,
			)
			return nil, err
		}
		slog.Error("Error finding competitions by filter:",
			"filter", filterJson,
			"error", result.Error,
		)
		return nil, result.Error
	}
	return competitions, nil
}

// RegisterUserToCompetition implements repository.CompetitionRepository.
func (r *competitionRepository) CreateUserRegistration(registration *entity.Registration) error {
	var competition entity.Competition
	user := &entity.User{
		ID: registration.UserID,
	}
	// double insert
	if err := r.db.First(&competition, registration.CompetitionID).Error; err != nil {
		return err
	}
	if err := r.db.Create(registration).Error; err != nil {
		return err
	}
	err := r.db.Model(&user).Association("JoinedCompetitions").Append(&competition)
	if err != nil {
		slog.Error("Error registering user to competition:",
			"registration", registration,
			"error", err,
		)
		return err
	}
	return nil
}
func (r *competitionRepository) DeleteUserRegistration(registration *entity.Registration) error {
	return r.db.Delete(registration).Error
}
func (r *competitionRepository) FindUserRegistration(competitionID uuid.UUID, userID uuid.UUID) (*entity.Registration, error) {
	var registration entity.Registration
	result := r.db.Where("competition_id = ? AND user_id = ?", competitionID, userID).First(&registration)
	if result.Error != nil {
		slog.Error("Error finding registration:",
			"competitionID", competitionID,
			"userID", userID,
			"error", result.Error,
		)
		return nil, result.Error
	}
	return &registration, nil
}

func (r *competitionRepository) CountAllRegistrations() (int, error) {
	var count int64
	err := r.db.Model(&entity.Registration{}).Count(&count).Error
	if err != nil {
		slog.Error("Error counting all registrations:",
			"error", err,
		)
		return 0, err
	}
	return int(count), nil
}

func (r *competitionRepository) CountRegistrations(competitionID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Model(&entity.Registration{}).
		Where("competition_id = ?", competitionID).
		Count(&count).Error
	if err != nil {
		slog.Error("Error counting registrations for competition:",
			"competitionID", competitionID,
			"error", err,
		)
		return 0, err
	}
	return int(count), nil
}
