package services

import (
	services "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/application"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/domain/entity"
)

type userService struct{}

func NewUserService() services.UserService {
	return &userService{}
}

func (u *userService) SignUp(user entity.User) error {
	return nil
}

func (u *userService) Login(username, password string) (string, error) {
	return "", nil
}
