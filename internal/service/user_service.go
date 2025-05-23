package service

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/YukLomba/BE-YukLomba/internal/domain/repository"
	"github.com/google/uuid"
)

type UserService interface {
	GetUser(id uuid.UUID) (*entity.User, error)
	GetAllUsers() ([]*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
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
	return u.userRepo.FindAllRegistration(id)
}

// CreateUser implements UserService.
func (u *UserServiceImpl) CreateUser(user *entity.User) error {
	return u.userRepo.Create(user)
}

// GetAllUsers implements UserService.
func (u *UserServiceImpl) GetAllUsers() ([]*entity.User, error) {
	return u.userRepo.FindAll()
}

// GetUser implements UserService.
func (u *UserServiceImpl) GetUser(id uuid.UUID) (*entity.User, error) {
	return u.userRepo.FindByID(id)
}

// UpdateUser implements UserService.
func (u *UserServiceImpl) UpdateUser(user *entity.User) error {
	return u.userRepo.Update(user)
}
