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
