package service

import (
	"errors"
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	errs "github.com/YukLomba/BE-YukLomba/internal/domain/error"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/YukLomba/BE-YukLomba/internal/infrastructure/util"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserService interface {
	GetUser(id uuid.UUID) (*entity.User, error)
	GetAllUsers() ([]*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(id uuid.UUID, data *map[string]interface{}) error
	GetAllUserRegistration(id uuid.UUID) ([]*entity.Registration, error)
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
func (u *UserServiceImpl) GetAllUserRegistration(id uuid.UUID) ([]*entity.Registration, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	registrations, err := u.userRepo.FindAllRegistration(id)
	if err != nil {
		return nil, errs.ErrInternalServer
	}
	return registrations, nil
}

// CreateUser implements UserService.
func (u *UserServiceImpl) CreateUser(user *entity.User) error {
	user, err := u.userRepo.FindByEmail(user.Email)
	if err != nil {
		return errs.ErrInternalServer
	}
	if user != nil {
		return ErrUserAlreadyExists
	}
	return u.userRepo.Create(user)
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
		return nil, errs.ErrInternalServer
	}
	if user == nil {
		return nil, errs.ErrNotFound
	}
	return user, nil
}

// UpdateUser implements UserService.
func (u *UserServiceImpl) UpdateUser(id uuid.UUID, data *map[string]interface{}) error {
	existing, err := u.userRepo.FindByID(id)
	if err != nil {
		return errs.ErrInternalServer
	}
	if existing == nil {
		return errs.ErrNotFound
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
