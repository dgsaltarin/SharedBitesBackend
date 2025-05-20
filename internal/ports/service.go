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
	UploadAndInitiateProcessing(ctx context.Context, params UploadBillParamsDTO) (*domain.Bill, error)
	GetBillDetails(ctx context.Context, billID uuid.UUID, userID uuid.UUID) (*domain.Bill, error)
	ListUserBills(ctx context.Context, userID uuid.UUID, offset, limit int) ([]domain.Bill, int64, error) // Returns bills and total count
	DeleteBill(ctx context.Context, billID uuid.UUID, userID uuid.UUID) error
	GetBillStatus(ctx context.Context, billID uuid.UUID, userID uuid.UUID) (domain.BillStatus, error)
}
