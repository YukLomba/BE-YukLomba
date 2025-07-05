package service

import (
	"errors"
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	errs "github.com/YukLomba/BE-YukLomba/internal/domain/error"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	GetUser(id uuid.UUID) (*entity.User, error)
	GetAllUsers() ([]*entity.User, error)
	UpdateUser(id uuid.UUID, data *map[string]interface{}) error
	GetAllUserRegistration(id uuid.UUID) ([]*entity.Competition, error)
}

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// GetAllUserRegistration implements UserService.
func (u *UserServiceImpl) GetAllUserRegistration(id uuid.UUID) ([]*entity.Competition, error) {
	_, err := u.userRepo.FindByID(id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrUserNotFound
		default:
			return nil, errs.ErrInternalServer
		}
	}
	registrations, err := u.userRepo.FindAllRegistration(id)
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	return registrations, nil
}

// GetAllUsers implements UserService.
func (u *UserServiceImpl) GetAllUsers() ([]*entity.User, error) {
	users, err := u.userRepo.FindAll()
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	return users, nil
}

// GetUser implements UserService.
func (u *UserServiceImpl) GetUser(id uuid.UUID) (*entity.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrUserNotFound
		default:
			return nil, errs.ErrInternalServer
		}
	}
	return user, nil
}

// UpdateUser implements UserService.
func (u *UserServiceImpl) UpdateUser(id uuid.UUID, data *map[string]interface{}) error {
	_, err := u.userRepo.FindByID(id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrUserNotFound
		default:
			return errs.ErrInternalServer
		}
	}
	if val, ok := (*data)["password"]; ok {
		// log.Println(val)
		(*data)["password"], err = util.HashPassword(val.(string))
		(*data)["password_changed_at"] = time.Now()
		if err != nil {
			return errs.ErrInternalServer
		}
	}
	err = u.userRepo.Update(id, data)
	if err != nil {
		return errs.ErrInternalServer
	}
	return nil
}
