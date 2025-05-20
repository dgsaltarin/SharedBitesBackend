package domain

import (
	"io"
	"time"

	"github.com/google/uuid"
)

// BillWithURL combines a bill with its file URL
type BillWithURL struct {
	Bill    *Bill
	FileURL string
}

// UploadBillRequest represents the data needed to upload a bill
type UploadBillRequest struct {
	UserID      uuid.UUID
	File        io.Reader
	Filename    string
	ContentType string
}

// ListBillsOptions represents options for listing bills
type ListBillsOptions struct {
	Limit  int
	Offset int
	Status BillStatus
}

// BillDTO represents a bill data transfer object
type BillDTO struct {
	ID              string        `json:"id"`
	Filename        string        `json:"filename"`
	Status          string        `json:"status"`
	UploadedAt      string        `json:"uploaded_at"`
	ProcessedAt     *string       `json:"processed_at,omitempty"`
	FileURL         string        `json:"file_url"`
	VendorName      string        `json:"vendor_name,omitempty"`
	TransactionDate *string       `json:"transaction_date,omitempty"`
	TotalAmount     *float64      `json:"total_amount,omitempty"`
	LineItems       []LineItemDTO `json:"line_items,omitempty"`
	TextTrackOutput string        `json:"text_track_output,omitempty"`
}

// LineItemDTO represents a line item data transfer object
type LineItemDTO struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Quantity    *float64 `json:"quantity,omitempty"`
	UnitPrice   *float64 `json:"unit_price,omitempty"`
	TotalPrice  *float64 `json:"total_price,omitempty"`
}

// BillSummaryDTO represents a summarized bill for listing
type BillSummaryDTO struct {
	ID              string     `json:"id"`
	Filename        string     `json:"filename"`
	Status          string     `json:"status"`
	UploadedAt      time.Time  `json:"uploaded_at"`
	VendorName      string     `json:"vendor_name,omitempty"`
	TotalAmount     *float64   `json:"total_amount,omitempty"`
	TransactionDate *time.Time `json:"transaction_date,omitempty"`
}

// ListBillsResponseDTO represents a response for listing bills
type ListBillsResponseDTO struct {
	Bills []BillSummaryDTO `json:"bills"`
	Total int64            `json:"total"`
}
