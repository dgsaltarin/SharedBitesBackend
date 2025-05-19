package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillStatus string

const (
	BillStatusUploaded   BillStatus = "uploaded"
	BillStatusPending    BillStatus = "pending"
	BillStatusProcessing BillStatus = "processing"
	BillStatusAnalyzed   BillStatus = "analyzed"
	BillStatusFailed     BillStatus = "failed"
)

type Bill struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	//User            User      `gorm:"foreignKey:UserID"` not sure if we need this
	Filename        string `gorm:"size:255;not null"`
	FileStoragePath string `gorm:"size:255;not null"`
	FileType        string `gorm:"size:25;not null"`
	Status          BillStatus
	UploadedAt      time.Time
	UpdatedAt       time.Time
	ProcessedAt     *time.Time

	// fields from textract
	VendorName      *string
	TransactionDate *time.Time
	TotalAmount     *float64
	LineItems       []LineItem `gorm:"foreignKey:BillID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TextTrackOutput *string    `gorm:"type:text"`
}

func (b *Bill) TableName() string {
	return "bills"
}

func (b *Bill) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	if b.UploadedAt.IsZero() {
		b.UploadedAt = time.Now().UTC()
	}

	if b.Status == "" {
		b.Status = BillStatusUploaded
	}

	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = time.Now().UTC()
	}

	return
}

func (b *Bill) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now().UTC()
	return
}

func NewBill(userID uuid.UUID, filename, fileStoragePath, fileType string) (*Bill, error) {
	if userID == uuid.Nil {
		return nil, ErrUserIDEmpty
	}

	if filename == "" || fileStoragePath == "" || fileType == "" {
		return nil, ErrInvalidInput
	}

	return &Bill{
		UserID:          userID,
		Filename:        filename,
		FileStoragePath: fileStoragePath,
		FileType:        fileType,
		Status:          BillStatusUploaded,
	}, nil
}
