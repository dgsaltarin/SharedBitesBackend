package services

import (
	"fmt"

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
	newUser, err := entity.NewUser(user)
	if err != nil {
		fmt.Println("Error creating user: ", err)
		return err
	}
	err = u.repository.UpsertUser(newUser)
	if err != nil {
		fmt.Println("Error signing up: ", err)
		return err
	}

	return nil
}

func (u *userService) Login(username, password string) (string, error) {
	return "", nil
}
