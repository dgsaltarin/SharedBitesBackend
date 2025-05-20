package ports

import (
	"context"
	"io"

	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, email, name, password string) (*domain.User, error)
	UpdateUserProfile(ctx context.Context, firebaseUID uuid.UUID, name *string) (*domain.User, error)
	GetUserByFirebaseUID(ctx context.Context, firebaseUID uuid.UUID) (*domain.User, error)
}

type UploadBillParamsDTO struct {
	UserID      uuid.UUID // Internal UUID of the user uploading the bill
	FileName    string    // Original name of the file
	ContentType string    // MIME type of the file
	FileContent io.Reader // The actual file content stream
	FileSize    int64     // Optional: size of the file, for validation or progress
}

type BillService interface {
	UploadBill(ctx context.Context, req domain.UploadBillRequest) (*domain.Bill, error)
	AnalyzeBill(ctx context.Context, billID uuid.UUID) error
	GetBill(ctx context.Context, billID, userID uuid.UUID) (*domain.BillWithURL, error)
	ListBills(ctx context.Context, userID uuid.UUID, options domain.ListBillsOptions) ([]domain.Bill, int64, error)
	GetBillStatus(ctx context.Context, billID uuid.UUID, userID uuid.UUID) (domain.BillStatus, error)
	DeleteBill(ctx context.Context, billID uuid.UUID, userID uuid.UUID) error
}
