package repository

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return userRepository{
		db: db,
	}
}

// FindByEmail implements repository.UserRepository.
func (r userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByUsername implements repository.UserRepository.
func (r userRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindAllRegistration implements repository.UserRepository.
func (r userRepository) FindAllRegistration(id uuid.UUID) ([]*entity.Registration, error) {
	var registrations []*entity.Registration

	result := r.db.Find(&registrations, "user_id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return registrations, nil
}

// Create implements repository.UserRepository.
func (r userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// FindAll implements repository.UserRepository.1
func (r userRepository) FindAll() ([]*entity.User, error) {
	var users []*entity.User
	result := r.db.Preload("OrganizedCompetitions").
		Preload("JoinedCompetitions").
		Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// FindByID implements repository.UserRepository.
func (r userRepository) FindByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	result := r.db.Preload("OrganizedCompetitions").
		Preload("JoinedCompetitions").
		First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Update implements repository.UserRepository.
func (r userRepository) Update(user *entity.User) error {
	result := r.db.Save(user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
