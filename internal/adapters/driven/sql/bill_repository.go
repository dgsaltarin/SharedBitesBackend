package sql

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/dgsaltarin/SharedBitesBackend/internal/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormBillRepository struct {
	db *gorm.DB
}

// NewGORMBillRepository creates a new GORM bill repository.
func NewGORMBillRepository(db *gorm.DB) ports.BillRepository {
	if db == nil {
		log.Fatal("GORM DB cannot be nil for BillRepository")
	}
	return &gormBillRepository{db: db}
}

// CreateBill creates a new bill record.
func (r *gormBillRepository) CreateBill(ctx context.Context, bill *domain.Bill) error {
	if err := r.db.WithContext(ctx).Create(bill).Error; err != nil {
		log.Printf("Error creating bill (UserID: %s, Filename: %s): %v", bill.UserID, bill.Filename, err)
		return fmt.Errorf("database error creating bill: %w", err)
	}
	return nil
}

// GetBillByID retrieves a bill by its UUID, preloading LineItems.
func (r *gormBillRepository) GetBillByID(ctx context.Context, billID uuid.UUID) (*domain.Bill, error) {
	var bill domain.Bill
	err := r.db.WithContext(ctx).Preload("LineItems").Where("id = ?", billID).First(&bill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBillNotFound
		}
		log.Printf("Error finding bill by ID %s: %v", billID, err)
		return nil, fmt.Errorf("database error finding bill by ID: %w", err)
	}
	return &bill, nil
}

// GetBillsByUserID retrieves all bills for a given user ID, preloading LineItems.
func (r *gormBillRepository) GetBillsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Bill, error) {
	var bills []*domain.Bill
	err := r.db.WithContext(ctx).Preload("LineItems").Where("user_id = ?", userID).Find(&bills).Error
	if err != nil {
		log.Printf("Error finding bills by UserID %s: %v", userID, err)
		return nil, fmt.Errorf("database error finding bills by UserID: %w", err)
	}
	return bills, nil
}

// UpdateBill updates an existing bill record.
func (r *gormBillRepository) UpdateBill(ctx context.Context, bill *domain.Bill) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Save the bill itself
		if err := tx.Save(bill).Error; err != nil {
			log.Printf("Error updating bill (ID: %s): %v", bill.ID, err)
			return fmt.Errorf("database error updating bill: %w", err)
		}

		if err := tx.Where("bill_id = ?", bill.ID).Delete(&domain.LineItem{}).Error; err != nil {
			log.Printf("Error deleting existing line items for bill ID %s: %v", bill.ID, err)
			return fmt.Errorf("database error clearing line items for update: %w", err)
		}

		if len(bill.LineItems) > 0 {
			for i := range bill.LineItems {
				bill.LineItems[i].BillID = bill.ID // Ensure BillID is set
			}
			if err := tx.Create(&bill.LineItems).Error; err != nil {
				log.Printf("Error creating new line items for bill ID %s: %v", bill.ID, err)
				return fmt.Errorf("database error creating line items for update: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// SaveBillWithLineItems creates a new bill and its associated line items in a transaction.
func (r *gormBillRepository) SaveBillWithLineItems(ctx context.Context, bill *domain.Bill, lineItems []*domain.LineItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the bill
		if err := tx.Create(bill).Error; err != nil {
			log.Printf("Error creating bill in transaction (UserID: %s): %v", bill.UserID, err)
			return fmt.Errorf("transaction error creating bill: %w", err)
		}

		if len(lineItems) > 0 {
			for _, item := range lineItems {
				item.BillID = bill.ID // Set the foreign key
			}
			if err := tx.Create(lineItems).Error; err != nil {
				log.Printf("Error creating line items in transaction for bill ID %s: %v", bill.ID, err)
				return fmt.Errorf("transaction error creating line items: %w", err)
			}
		}
		// bill.LineItems = lineItems // Assign created line items back to the bill object for consistency
		// After successful creation, if you need to update bill.LineItems to reflect what's in the DB (including IDs and timestamps)
		// you might need to re-fetch or carefully construct them. For now, we assume the passed lineItems are what's persisted.
		// However, the original bill.LineItems is []domain.LineItem, and lineItems is []*domain.LineItem.
		// Let's assign them back carefully if needed, or ensure the caller updates its reference.
		// For now, we will assign the created items back to the bill object by dereferencing.
		domainLineItems := make([]domain.LineItem, len(lineItems))
		for i, item := range lineItems {
			if item != nil {
				domainLineItems[i] = *item
			}
		}
		bill.LineItems = domainLineItems
		return nil
	})
}
