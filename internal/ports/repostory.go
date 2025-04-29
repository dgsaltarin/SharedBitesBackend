package ports

import (
	"context"
	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error // Creates or updates
	FindByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
