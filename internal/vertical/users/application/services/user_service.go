package services

import (
	"fmt"

	services "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/application"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/domain/entity"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/domain/repository"
	"github.com/google/uuid"
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
	user.ID = uuid.New().String()
	err := u.repository.UpsertUser(&user)
	if err != nil {
		fmt.Println("Error signing up: ", err)
		return err
	}

	return nil
}

func (u *userService) Login(username, password string) (string, error) {
	return "", nil
}
