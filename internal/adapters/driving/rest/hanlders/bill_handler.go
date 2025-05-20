package hanlders

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgsaltarin/SharedBitesBackend/internal/application"
	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BillHandler handles HTTP requests for bill operations
type BillHandler struct {
	billService *application.BillService
}

// NewBillHandler creates a new BillHandler
func NewBillHandler(billService *application.BillService) *BillHandler {
	if billService == nil {
		panic("BillService cannot be nil in NewBillHandler")
	}
	return &BillHandler{billService: billService}
}

// UploadBill godoc
// @Summary Upload a bill image
// @Description Upload a bill image to be stored and later analyzed
// @Tags Bills
// @Accept mpfd
// @Produce json
// @Param image formData file true "Image file of the bill to upload"
// @Success 200 {object} gin.H{"bill_id": string, "message": string} "Successfully uploaded bill"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - e.g., no file or invalid file"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error"
// @Router /bills/upload [post]
func (h *BillHandler) UploadBill(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read image from request: " + err.Error()})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream" // Default content type
	}

	// Create upload request
	uploadReq := domain.UploadBillRequest{
		UserID:      userID,
		File:        file,
		Filename:    header.Filename,
		ContentType: contentType,
	}

	// Upload the bill
	bill, err := h.billService.UploadBill(c, uploadReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload bill: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bill_id": bill.ID.String(),
		"message": "Bill uploaded successfully",
	})
}

// AnalyzeBill godoc
// @Summary Analyze a previously uploaded bill
// @Description Process a bill with Textract to extract information
// @Tags Bills
// @Produce json
// @Param bill_id path string true "Bill ID to analyze"
// @Success 200 {object} gin.H{"message": string} "Successfully analyzed bill"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - e.g., invalid bill ID"
// @Failure 404 {object} gin.H{"error": string} "Not Found - Bill not found"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error"
// @Router /bills/{bill_id}/analyze [post]
func (h *BillHandler) AnalyzeBill(c *gin.Context) {
	billIDStr := c.Param("bill_id")
	billID, err := uuid.Parse(billIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID format"})
		return
	}

	err = h.billService.AnalyzeBill(c, billID)
	if err != nil {
		if err == domain.ErrBillNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze bill: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bill analyzed successfully"})
}

// GetBill godoc
// @Summary Get a bill by ID
// @Description Retrieve bill details including extracted information and line items
// @Tags Bills
// @Produce json
// @Param bill_id path string true "Bill ID to retrieve"
// @Success 200 {object} domain.BillDTO "Bill details"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - e.g., invalid bill ID"
// @Failure 404 {object} gin.H{"error": string} "Not Found - Bill not found"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error"
// @Router /bills/{bill_id} [get]
func (h *BillHandler) GetBill(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	billIDStr := c.Param("bill_id")
	billID, err := uuid.Parse(billIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID format"})
		return
	}

	billWithURL, err := h.billService.GetBill(c, billID, userID)
	if err != nil {
		if err == domain.ErrBillNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bill: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, formatBillResponse(billWithURL))
}

// ListBills godoc
// @Summary List bills for the authenticated user
// @Description List bills with pagination and optional filtering
// @Tags Bills
// @Produce json
// @Param limit query int false "Number of bills to return (default 10)"
// @Param offset query int false "Number of bills to skip (default 0)"
// @Param status query string false "Filter by status (uploaded, pending, processing, analyzed, failed)"
// @Success 200 {object} domain.ListBillsResponseDTO "List of bills"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error"
// @Router /bills [get]
func (h *BillHandler) ListBills(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	status := domain.BillStatus(c.Query("status"))

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	options := domain.ListBillsOptions{
		Limit:  limit,
		Offset: offset,
		Status: status,
	}

	bills, total, err := h.billService.ListBills(c, userID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list bills: " + err.Error()})
		return
	}

	response := domain.ListBillsResponseDTO{
		Bills: make([]domain.BillSummaryDTO, len(bills)),
		Total: total,
	}

	for i, bill := range bills {
		response.Bills[i] = domain.BillSummaryDTO{
			ID:              bill.ID.String(),
			Filename:        bill.Filename,
			Status:          string(bill.Status),
			UploadedAt:      bill.UploadedAt,
			VendorName:      safeString(bill.VendorName),
			TotalAmount:     bill.TotalAmount,
			TransactionDate: bill.TransactionDate,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetBillStatus godoc
// @Summary Get bill status
// @Description Retrieve just the status of a bill
// @Tags Bills
// @Produce json
// @Param bill_id path string true "Bill ID to check"
// @Success 200 {object} gin.H{"status": string} "Bill status"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - e.g., invalid bill ID"
// @Failure 404 {object} gin.H{"error": string} "Not Found - Bill not found"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error"
// @Router /bills/{bill_id}/status [get]
func (h *BillHandler) GetBillStatus(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	billIDStr := c.Param("bill_id")
	billID, err := uuid.Parse(billIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID format"})
		return
	}

	status, err := h.billService.GetBillStatus(c, billID, userID)
	if err != nil {
		if err == domain.ErrBillNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bill status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// DeleteBill godoc
// @Summary Delete a bill
// @Description Delete a bill and all associated data
// @Tags Bills
// @Produce json
// @Param bill_id path string true "Bill ID to delete"
// @Success 200 {object} gin.H{"message": string} "Successfully deleted bill"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - e.g., invalid bill ID"
// @Failure 404 {object} gin.H{"error": string} "Not Found - Bill not found"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error"
// @Router /bills/{bill_id} [delete]
func (h *BillHandler) DeleteBill(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	billIDStr := c.Param("bill_id")
	billID, err := uuid.Parse(billIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID format"})
		return
	}

	err = h.billService.DeleteBill(c, billID, userID)
	if err != nil {
		if err == domain.ErrBillNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bill: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bill deleted successfully"})
}

// Helper functions

func formatBillResponse(billWithURL *domain.BillWithURL) domain.BillDTO {
	bill := billWithURL.Bill
	response := domain.BillDTO{
		ID:              bill.ID.String(),
		Filename:        bill.Filename,
		Status:          string(bill.Status),
		UploadedAt:      bill.UploadedAt.Format(time.RFC3339),
		FileURL:         billWithURL.FileURL,
		VendorName:      safeString(bill.VendorName),
		TotalAmount:     bill.TotalAmount,
		TextTrackOutput: safeString(bill.TextTrackOutput),
	}

	if bill.ProcessedAt != nil {
		processedAt := bill.ProcessedAt.Format(time.RFC3339)
		response.ProcessedAt = &processedAt
	}

	if bill.TransactionDate != nil {
		transactionDate := bill.TransactionDate.Format(time.RFC3339)
		response.TransactionDate = &transactionDate
	}

	response.LineItems = make([]domain.LineItemDTO, len(bill.LineItems))
	for i, item := range bill.LineItems {
		response.LineItems[i] = domain.LineItemDTO{
			ID:          item.ID.String(),
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
		}
	}

	return response
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeFloat64(f *float64) *float64 {
	if f == nil {
		return nil
	}
	return f
}
