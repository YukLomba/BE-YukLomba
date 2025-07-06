package repository

import (
	"log/slog"

	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewRepository = repository.ReviewRepository

type ReviewRepositoryImpl struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &ReviewRepositoryImpl{db: db}
}

func (r *ReviewRepositoryImpl) Create(review *entity.Review) error {
	err := r.db.Create(review).Error
	if err != nil {
		slog.Error("Error creating review:",
			"error", err,
		)
		return err
	}
	return nil
}

func (r *ReviewRepositoryImpl) Update(reviewID uuid.UUID, data map[string]interface{}) error {
	return r.db.Model(&entity.Review{}).Where("id = ?", reviewID).Updates(data).Error
}

func (r *ReviewRepositoryImpl) GetByCompetition(competitionID uuid.UUID) ([]*entity.Review, error) {
	var reviews []*entity.Review
	err := r.db.Preload("User").
		Where("competition_id = ?", competitionID).
		Find(&reviews).Error
	if err != nil {
		slog.Error("Error getting reviews by competition:",
			"competitionID", competitionID,
			"error", err,
		)
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewRepositoryImpl) GetByUserAndCompetition(userID, competitionID uuid.UUID) (*entity.Review, error) {
	var review entity.Review
	err := r.db.Where("user_id = ? AND competition_id = ?", userID, competitionID).
		First(&review).Error
	if err != nil {
		slog.Error("Error getting review by user and competition:",
			"userID", userID,
			"competitionID", competitionID,
			"error", err,
		)
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepositoryImpl) GetAverageRating(competitionID uuid.UUID) (float32, error) {
	var avg float32
	err := r.db.Model(&entity.Review{}).
		Where("competition_id = ?", competitionID).
		Select("AVG(rating)").
		Scan(&avg).Error
	if err != nil {
		slog.Error("Error getting average rating:",
			"competitionID", competitionID,
			"error", err,
		)
		return 0, err
	}
	return avg, nil
}

func (r *ReviewRepositoryImpl) GetAverageRatingAll() (float32, error) {
	var avg float32
	err := r.db.Model(&entity.Review{}).
		Select("AVG(rating)").
		Scan(&avg).Error
	if err != nil {
		slog.Error("Error getting average rating:",
			"error", err,
		)
		return 0, err
	}
	return avg, nil
}
