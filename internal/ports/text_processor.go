package ports

import (
	"context"
	"time"
)

type ParsedTextractData struct {
	VendorName      *string
	TransactionDate *time.Time
	TotalAmount     *float64
	LineItems       []ParsedLineItem
	RawTextOutput   string
}

type ParsedLineItem struct {
	Description string
	Quantity    *float64
	UnitPrice   *float64
	TotalPrice  *float64
}

type TextProcessor interface {
	AnalyzeDocument(ctx context.Context, storagePath string) (*ParsedTextractData, error)
}
