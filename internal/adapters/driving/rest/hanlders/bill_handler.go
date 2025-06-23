package hanlders

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driven/texttrack"
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

// GetBill godoc
// @Summary Retrieve a bill with complete details
// @Description Get comprehensive bill information including metadata, OCR-extracted data, line items, and a pre-signed URL for file access. Only returns bills owned by the authenticated user.
// @Tags Bills
// @Produce json
// @Param bill_id path string true "UUID of the bill to retrieve"
// @Success 200 {object} domain.BillDTO "Complete bill details with file URL and extracted data"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - invalid bill ID format"
// @Failure 401 {object} gin.H{"error": string} "Unauthorized - invalid or missing authentication token"
// @Failure 404 {object} gin.H{"error": string} "Not Found - bill not found or not owned by user"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error - database or file storage error"
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
// @Summary List user's bills with pagination and filtering
// @Description Retrieve a paginated list of bills owned by the authenticated user. Supports filtering by processing status and includes summary information for each bill.
// @Tags Bills
// @Produce json
// @Param limit query int false "Number of bills to return per page (default: 10, max: 100)"
// @Param offset query int false "Number of bills to skip for pagination (default: 0)"
// @Param status query string false "Filter by processing status: uploaded, pending, processing, analyzed, failed"
// @Success 200 {object} domain.ListBillsResponseDTO "Paginated list of bill summaries with total count"
// @Failure 401 {object} gin.H{"error": string} "Unauthorized - invalid or missing authentication token"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error - database query error"
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
// @Summary Get current processing status of a bill
// @Description Retrieve only the current processing status of a specific bill. Useful for polling during OCR processing without fetching complete bill data.
// @Tags Bills
// @Produce json
// @Param bill_id path string true "UUID of the bill to check status for"
// @Success 200 {object} gin.H{"status": string} "Current bill status (uploaded, pending, processing, analyzed, failed)"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - invalid bill ID format"
// @Failure 401 {object} gin.H{"error": string} "Unauthorized - invalid or missing authentication token"
// @Failure 404 {object} gin.H{"error": string} "Not Found - bill not found or not owned by user"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error - database query error"
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
// @Summary Permanently delete a bill and all associated data
// @Description Remove a bill from the system including its file from S3 storage, database record, and all extracted line items. This action cannot be undone. Only the bill owner can delete their bills.
// @Tags Bills
// @Produce json
// @Param bill_id path string true "UUID of the bill to permanently delete"
// @Success 200 {object} gin.H{"message": string} "Bill and all associated data successfully deleted"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - invalid bill ID format"
// @Failure 401 {object} gin.H{"error": string} "Unauthorized - invalid or missing authentication token"
// @Failure 404 {object} gin.H{"error": string} "Not Found - bill not found or not owned by user"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error - database or file storage deletion error"
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

// UploadAndAnalyzeBillWithConfig godoc
// @Summary Upload and analyze a bill image with custom language and analysis configuration
// @Description Upload a bill image file and immediately process it with AWS Textract OCR using custom language settings, confidence thresholds, and currency codes. Ideal for multi-language documents and specific regional requirements.
// @Tags Bills
// @Accept mpfd
// @Produce json
// @Param image formData file true "Image file of the bill to upload and analyze (JPEG, PNG, PDF supported)"
// @Param languages formData string false "Comma-separated list of language codes (e.g., 'es,en' for Spanish and English)"
// @Param min_confidence formData number false "Minimum confidence threshold (0.0 to 1.0, default: 0.7)"
// @Param currency_codes formData string false "Comma-separated list of currency codes (e.g., 'EUR,USD,MXN')"
// @Success 200 {object} domain.BillDTO "Successfully uploaded and analyzed bill with extracted data"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - missing file, invalid file format, or malformed request"
// @Failure 401 {object} gin.H{"error": string} "Unauthorized - invalid or missing authentication token"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error - file storage, Textract processing, or database error"
// @Router /bills/upload-analyze-config [post]
func (h *BillHandler) UploadAndAnalyzeBillWithConfig(c *gin.Context) {
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

	// Parse configuration parameters
	config := texttrack.DefaultConfig() // Start with default config

	// Parse languages
	if languagesStr := c.PostForm("languages"); languagesStr != "" {
		languages := strings.Split(languagesStr, ",")
		for i, lang := range languages {
			languages[i] = strings.TrimSpace(lang)
		}
		config.Languages = languages
	}

	// Parse min_confidence
	if minConfidenceStr := c.PostForm("min_confidence"); minConfidenceStr != "" {
		if minConfidence, err := strconv.ParseFloat(minConfidenceStr, 64); err == nil {
			if minConfidence >= 0.0 && minConfidence <= 1.0 {
				config.MinConfidence = minConfidence
			}
		}
	}

	// Parse currency_codes
	if currencyCodesStr := c.PostForm("currency_codes"); currencyCodesStr != "" {
		currencyCodes := strings.Split(currencyCodesStr, ",")
		for i, code := range currencyCodes {
			currencyCodes[i] = strings.TrimSpace(code)
		}
		config.CurrencyCodes = currencyCodes
	}

	// Create upload request
	uploadReq := domain.UploadBillRequest{
		UserID:      userID,
		File:        file,
		Filename:    header.Filename,
		ContentType: contentType,
	}

	// Upload and analyze the bill with custom configuration
	billWithURL, err := h.billService.UploadAndAnalyzeBillWithConfig(c, uploadReq, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload and analyze bill: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, formatBillResponse(billWithURL))
}

// GetAnalysisConfigs godoc
// @Summary Get available analysis configurations for different languages and regions
// @Description Retrieve information about pre-configured analysis settings for different languages, regions, and use cases. This helps users choose the appropriate configuration for their documents.
// @Tags Bills
// @Produce json
// @Success 200 {object} gin.H{"configs": []gin.H} "Available analysis configurations with descriptions"
// @Failure 401 {object} gin.H{"error": string} "Unauthorized - invalid or missing authentication token"
// @Router /bills/analysis-configs [get]
func (h *BillHandler) GetAnalysisConfigs(c *gin.Context) {
	configs := []gin.H{
		{
			"name":        "default",
			"description": "Default configuration optimized for Spanish and English bills",
			"languages":   []string{"es", "en"},
			"confidence":  0.7,
			"currencies":  []string{"EUR", "USD", "MXN", "ARS", "CLP", "COP", "PEN", "UYU"},
		},
		{
			"name":        "spanish",
			"description": "Optimized for Spanish bills with Spanish-specific keywords and formatting",
			"languages":   []string{"es"},
			"confidence":  0.6,
			"currencies":  []string{"EUR", "MXN", "ARS", "CLP", "COP", "PEN", "UYU", "USD"},
		},
		{
			"name":        "english",
			"description": "Optimized for English bills with English-specific keywords and formatting",
			"languages":   []string{"en"},
			"confidence":  0.7,
			"currencies":  []string{"USD", "EUR", "GBP", "CAD", "AUD", "JPY"},
		},
		{
			"name":        "french",
			"description": "Optimized for French bills with French-specific keywords and formatting",
			"languages":   []string{"fr"},
			"confidence":  0.6,
			"currencies":  []string{"EUR", "USD", "CAD", "CHF", "GBP"},
		},
		{
			"name":        "german",
			"description": "Optimized for German bills with German-specific keywords and formatting",
			"languages":   []string{"de"},
			"confidence":  0.6,
			"currencies":  []string{"EUR", "USD", "CHF", "GBP"},
		},
		{
			"name":        "portuguese",
			"description": "Optimized for Portuguese bills with Portuguese-specific keywords and formatting",
			"languages":   []string{"pt"},
			"confidence":  0.6,
			"currencies":  []string{"EUR", "USD", "BRL", "MZN", "AOA"},
		},
		{
			"name":        "italian",
			"description": "Optimized for Italian bills with Italian-specific keywords and formatting",
			"languages":   []string{"it"},
			"confidence":  0.6,
			"currencies":  []string{"EUR", "USD", "CHF", "GBP"},
		},
		{
			"name":        "latin_american",
			"description": "Optimized for Latin American Spanish bills with regional currency support",
			"languages":   []string{"es"},
			"confidence":  0.5,
			"currencies":  []string{"MXN", "ARS", "CLP", "COP", "PEN", "UYU", "BRL", "USD", "EUR"},
		},
		{
			"name":        "european",
			"description": "Multi-language configuration for European bills (English, Spanish, French, German, Italian, Portuguese)",
			"languages":   []string{"en", "es", "fr", "de", "it", "pt"},
			"confidence":  0.6,
			"currencies":  []string{"EUR", "USD", "GBP", "CHF", "SEK", "NOK", "DKK"},
		},
		{
			"name":        "high_accuracy",
			"description": "High accuracy configuration with strict confidence requirements (85%)",
			"languages":   []string{"es", "en"},
			"confidence":  0.85,
			"currencies":  []string{"EUR", "USD", "MXN", "ARS", "CLP", "COP", "PEN", "UYU"},
		},
		{
			"name":        "low_confidence",
			"description": "Low confidence configuration for poor quality images or unclear text (40%)",
			"languages":   []string{"es", "en"},
			"confidence":  0.4,
			"currencies":  []string{"EUR", "USD", "MXN", "ARS", "CLP", "COP", "PEN", "UYU"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"configs": configs,
		"usage": gin.H{
			"endpoint": "/bills/upload-analyze-config",
			"parameters": gin.H{
				"languages":      "Comma-separated list of language codes (e.g., 'es,en')",
				"min_confidence": "Minimum confidence threshold (0.0 to 1.0)",
				"currency_codes": "Comma-separated list of currency codes (e.g., 'EUR,USD,MXN')",
			},
		},
	})
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
