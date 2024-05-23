package application

import (
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/domain/entity"
)

type UserService interface {
	SignUp(user entity.User) error
	Login(username, password string) (string, error)
}

type HealthCheckService interface {
	Check() string
}
