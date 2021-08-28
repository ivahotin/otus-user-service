package service

import "example.com/arch/user-service/internal/users/models"

type UserRepository interface {
	CreateUser(*models.User) (models.UserId, error)
	GetUser(models.UserId) (models.User, error)
	UpdateUser(*models.User) error
	DeleteUser(models.UserId, string) error
}
