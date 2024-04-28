package repository

import (
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/domain/entity"
)

type UserRepository interface {
	GetUser(id string) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	UpserUser(user *entity.User) error
}
