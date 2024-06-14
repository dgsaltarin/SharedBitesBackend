package mappers

import (
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/domain/entity"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/repository/dynamodb/models"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/rest/request"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/rest/response"
)

type Mappers struct{}

func NewMappers() Mappers {
	return Mappers{}
}

func (m *Mappers) MapSignUpRequestToUser(signUpRequest request.SignUpRequest) *entity.User {
	var user entity.User

	user.Username = signUpRequest.Username
	user.Email = signUpRequest.Email
	user.Password = signUpRequest.Password

	return &user
}

func (m *Mappers) MapUserToSignUpResponse(user entity.User) *response.SignUpResponse {
	var response response.SignUpResponse

	response.Username = user.Username
	response.Email = user.Email

	return &response
}

func (m *Mappers) MapUserRepositoryToUser(user *models.User) *entity.User {
	var entity entity.User

	entity.Username = user.Username
	entity.Email = user.Email

	return &entity
}
