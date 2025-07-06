package repository

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/google/uuid"
)

type ReviewRepository interface {
	Create(review *entity.Review) error
	Update(reviewID uuid.UUID, data map[string]interface{}) error
	GetByCompetition(competitionID uuid.UUID) ([]*entity.Review, error)
	GetByUserAndCompetition(userID, competitionID uuid.UUID) (*entity.Review, error)
	GetAverageRating(competitionID uuid.UUID) (float32, error)
	GetAverageRatingAll() (float32, error)
}
