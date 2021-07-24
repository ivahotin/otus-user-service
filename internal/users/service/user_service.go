package service

import (
	"example.com/arch/user-service/internal/users/models"
	"github.com/go-playground/validator"
)

type UserServiceImpl struct {
	repository    UserRepository
	userValidator *validator.Validate
}

func NewUserService(repository UserRepository) *UserServiceImpl {
	userValidator := validator.New()
	return &UserServiceImpl{repository: repository, userValidator: userValidator}
}

func (service *UserServiceImpl) CreateUser(user *models.User) (models.UserId, error) {
	err := service.userValidator.Struct(user)
	if err != nil {
		return -1, err
	}
	return service.repository.CreateUser(user)
}

func (service *UserServiceImpl) GetUser(userId models.UserId) (models.User, error) {
	return service.repository.GetUser(userId)
}

func (service *UserServiceImpl) UpdateUser(user *models.User) error {
	err := service.userValidator.Struct(user)
	if err != nil {
		return err
	}
	return service.repository.UpdateUser(user)
}

func (service *UserServiceImpl) DeleteUser(userId models.UserId) error {
	return service.repository.DeleteUser(userId)
}
