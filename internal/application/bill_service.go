package application

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/dgsaltarin/SharedBitesBackend/internal/ports"
	"github.com/dgsaltarin/SharedBitesBackend/platform/aws"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillService struct {
	textractClient *aws.TextractClient
	fileStore      ports.FileStore
	textProcessor  ports.TextProcessor
	db             *gorm.DB
}

func NewBillService(
	textractClient *aws.TextractClient,
	fileStore ports.FileStore,
	textProcessor ports.TextProcessor,
	db *gorm.DB,
) *BillService {
	return &BillService{
		textractClient: textractClient,
		fileStore:      fileStore,
		textProcessor:  textProcessor,
		db:             db,
	}
}

// UploadBill handles uploading a bill file to storage and saving its metadata
func (s *BillService) UploadBill(ctx context.Context, req domain.UploadBillRequest) (*domain.Bill, error) {
	if req.UserID == uuid.Nil {
		return nil, domain.ErrUserIDEmpty
	}
	if req.File == nil || req.Filename == "" || req.ContentType == "" {
		return nil, domain.ErrInvalidInput
	}

	// Create a storage path for the file
	storagePath := fmt.Sprintf("bills/%s/%s", req.UserID, filepath.Base(req.Filename))

	// Upload the file to storage
	storedPath, err := s.fileStore.UploadFile(ctx, req.File, storagePath, req.ContentType)
	if err != nil {
		return nil, fmt.Errorf("error uploading bill file: %w", err)
	}

	// Create bill record in database
	bill, err := domain.NewBill(req.UserID, req.Filename, storedPath, req.ContentType)
	if err != nil {
		// If there was an error creating the bill record, try to delete the uploaded file
		s.fileStore.DeleteFile(ctx, storedPath)
		return nil, fmt.Errorf("error creating bill record: %w", err)
	}

	// Save bill to database
	if err := s.db.Create(bill).Error; err != nil {
		// If there was an error saving to the database, try to delete the uploaded file
		s.fileStore.DeleteFile(ctx, storedPath)
		return nil, fmt.Errorf("error saving bill to database: %w", err)
	}

	return bill, nil
}

func (s *BillService) AnalyzeBill(ctx context.Context, billID uuid.UUID) error {
	if billID == uuid.Nil {
		return domain.ErrInvalidInput
	}

	// Retrieve the bill from the database
	var bill domain.Bill
	if err := s.db.First(&bill, "id = ?", billID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrBillNotFound
		}
		return fmt.Errorf("error retrieving bill: %w", err)
	}

	// Update bill status to processing
	if err := s.db.Model(&bill).Update("status", domain.BillStatusProcessing).Error; err != nil {
		return fmt.Errorf("error updating bill status to processing: %w", err)
	}

	// Analyze the document with Textract
	result, err := s.textProcessor.AnalyzeDocument(ctx, bill.FileStoragePath)
	if err != nil {
		// Update bill status to failed
		s.db.Model(&bill).Updates(map[string]interface{}{
			"status": domain.BillStatusFailed,
		})
		return fmt.Errorf("error analyzing bill with Textract: %w", err)
	}

	// Start a transaction to update the bill and create line items
	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("error starting transaction: %w", tx.Error)
	}

	// Update bill with extracted information
	billUpdates := map[string]interface{}{
		"vendor_name":       result.VendorName,
		"transaction_date":  result.TransactionDate,
		"total_amount":      result.TotalAmount,
		"text_track_output": result.RawTextOutput,
		"status":            domain.BillStatusAnalyzed,
	}

	if err := tx.Model(&bill).Updates(billUpdates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating bill with extracted data: %w", err)
	}

	// Clear existing line items (if any) and add new ones
	if err := tx.Where("bill_id = ?", billID).Delete(&domain.LineItem{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error clearing existing line items: %w", err)
	}

	// Create line items from extracted data
	for _, item := range result.LineItems {
		lineItem, err := domain.NewLineItem(
			billID,
			item.Description,
			item.Quantity,
			item.UnitPrice,
			item.TotalPrice,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error creating line item: %w", err)
		}

		if err := tx.Create(lineItem).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error saving line item: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// GetBill retrieves a bill and its line items by ID
func (s *BillService) GetBill(ctx context.Context, billID, userID uuid.UUID) (*domain.BillWithURL, error) {
	if billID == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}

	var bill domain.Bill
	err := s.db.Preload("LineItems").First(&bill, "id = ? AND user_id = ?", billID, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBillNotFound
		}
		return nil, fmt.Errorf("error retrieving bill: %w", err)
	}

	// Generate a pre-signed URL for the bill file
	fileURL, err := s.fileStore.GetFileURL(ctx, bill.FileStoragePath)
	if err != nil {
		// Log the error but continue as this is not critical
		fmt.Printf("Warning: Failed to generate pre-signed URL for bill %s: %v\n", billID, err)
		fileURL = "" // Empty URL if we couldn't generate one
	}

	return &domain.BillWithURL{
		Bill:    &bill,
		FileURL: fileURL,
	}, nil
}

// ListBills retrieves all bills for a user with pagination
func (s *BillService) ListBills(ctx context.Context, userID uuid.UUID, options domain.ListBillsOptions) ([]domain.Bill, int64, error) {
	if userID == uuid.Nil {
		return nil, 0, domain.ErrUserIDEmpty
	}

	// Set default values for pagination
	if options.Limit <= 0 {
		options.Limit = 10
	}
	if options.Offset < 0 {
		options.Offset = 0
	}

	// Build query
	query := s.db.Model(&domain.Bill{}).Where("user_id = ?", userID)

	// Apply status filter if provided
	if options.Status != "" {
		query = query.Where("status = ?", options.Status)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting bills: %w", err)
	}

	// Get bills with pagination
	var bills []domain.Bill
	err := query.
		Order("uploaded_at DESC"). // Most recent first
		Limit(options.Limit).
		Offset(options.Offset).
		Find(&bills).Error
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving bills: %w", err)
	}

	return bills, total, nil
}

// GetBillStatus retrieves just the status of a bill
func (s *BillService) GetBillStatus(ctx context.Context, billID, userID uuid.UUID) (domain.BillStatus, error) {
	if billID == uuid.Nil {
		return "", domain.ErrInvalidInput
	}

	var status domain.BillStatus
	err := s.db.Model(&domain.Bill{}).
		Select("status").
		Where("id = ? AND user_id = ?", billID, userID).
		First(&status).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domain.ErrBillNotFound
		}
		return "", fmt.Errorf("error retrieving bill status: %w", err)
	}

	return status, nil
}

// DeleteBill deletes a bill and its associated data
func (s *BillService) DeleteBill(ctx context.Context, billID, userID uuid.UUID) error {
	if billID == uuid.Nil {
		return domain.ErrInvalidInput
	}

	// Get the bill to check ownership and to know the file path
	var bill domain.Bill
	if err := s.db.First(&bill, "id = ? AND user_id = ?", billID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrBillNotFound
		}
		return fmt.Errorf("error retrieving bill: %w", err)
	}

	// Start a transaction to delete the bill and its line items
	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("error starting transaction: %w", tx.Error)
	}

	// Delete line items first (this should use cascading delete if set up in the database)
	if err := tx.Where("bill_id = ?", billID).Delete(&domain.LineItem{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting line items: %w", err)
	}

	// Then delete the bill
	if err := tx.Delete(&bill).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting bill: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	// Delete the file from storage (after successful DB deletion)
	// This is done outside the transaction as it's a separate system
	if err := s.fileStore.DeleteFile(ctx, bill.FileStoragePath); err != nil {
		// Log but don't fail the operation if file deletion fails
		fmt.Printf("Warning: Failed to delete file for bill %s: %v\n", billID, err)
	}

	return nil
}
