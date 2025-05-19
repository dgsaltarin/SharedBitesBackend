package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// remove json tags
type LineItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
	BillID      uuid.UUID `gorm:"type:uuid;not null;index"` // Foreign key to Bill
	Description string    `gorm:"not null"`
	Quantity    *float64  `gorm:"type:decimal(10,3);default:1;"` // Optional, defaults to 1
	UnitPrice   *float64  `gorm:"type:decimal(10,2);"`           // Price per unit
	TotalPrice  *float64  `gorm:"type:decimal(10,2);"`           // Quantity * UnitPrice (or directly extracted)
	// Consider adding: ProductCode, Category (user-defined or ML-suggested)
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // If line items can be soft-deleted individually
}

func TableName() string {
	return "bill_line_items"
}

func (l *LineItem) BeforeCreate(tx *gorm.DB) (err error) {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}

	now := time.Now().UTC()
	l.CreatedAt = now
	l.UpdatedAt = now
	return
}

func (l *LineItem) BeforeUpdate(tx *gorm.DB) (err error) {
	l.UpdatedAt = time.Now().UTC()
	return
}

func NewLineItem(billID uuid.UUID, description string, quantity, unitPrice, totalPrice *float64) (*LineItem, error) {
	if billID == uuid.Nil {
		return nil, ErrBillIDEmpty
	}
	if description == "" {
		return nil, ErrLineItemDescriptionEmpty
	}

	item := &LineItem{
		BillID:      billID,
		Description: description,
	}

	if quantity != nil {
		item.Quantity = quantity
	} else {
		defaultQuantity := 1.0
		item.Quantity = &defaultQuantity
	}

	if unitPrice != nil {
		item.UnitPrice = unitPrice
	}

	if totalPrice != nil {
		item.TotalPrice = totalPrice
	}

	return item, nil
}
