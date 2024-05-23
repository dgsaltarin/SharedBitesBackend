package services

import (
	services "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/application"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/domain/entity"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/domain/repository"
)

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) services.UserService {
	return &userService{
		repository: repository,
	}
}

func (u *userService) SignUp(user entity.User) error {
	err := u.repository.UpsertUser(&user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) Login(username, password string) (string, error) {
	return "", nil
}
