package services

import (
	services "github.com/dgsaltarin/SharedBitesBackend/internal/application"
)

type userService struct {
}

func NewUserService() services.UserService {
	return &userService{}
}

func (u *userService) SignUp(username, email, password string) error {
	return nil
}

func (u *userService) Login(username, password string) (string, error) {
	return "", nil
}
