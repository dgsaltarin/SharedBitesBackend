package ports

import (
	"context"

	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error // Creates or updates
	FindByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type BillRepository interface {
	CreateBill(ctx context.Context, bill *domain.Bill) error
	GetBillByID(ctx context.Context, billID uuid.UUID) (*domain.Bill, error)
	GetBillsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Bill, error)
	UpdateBill(ctx context.Context, bill *domain.Bill) error
	SaveBillWithLineItems(ctx context.Context, bill *domain.Bill, lineItems []*domain.LineItem) error
}

// GroupRepository defines the interface for group data access operations
type GroupRepository interface {
	Create(ctx context.Context, group *domain.Group) error
	GetByID(ctx context.Context, groupID uuid.UUID) (*domain.Group, error)
	GetByIDAndOwner(ctx context.Context, groupID, ownerID uuid.UUID) (*domain.Group, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID, options domain.ListGroupsOptions) ([]domain.Group, int64, error)
	Update(ctx context.Context, group *domain.Group) error
	Delete(ctx context.Context, groupID uuid.UUID) error
}
