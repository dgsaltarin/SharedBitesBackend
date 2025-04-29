package ports

import (
	"context"
	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, email, name, password string) (*domain.User, error)
	UpdateUserProfile(ctx context.Context, firebaseUID uuid.UUID, name *string) (*domain.User, error)
	GetUserByFirebaseUID(ctx context.Context, firebaseUID uuid.UUID) (*domain.User, error)
}
