package hanlders

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/dgsaltarin/SharedBitesBackend/internal/application"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TextractHandler handles HTTP requests for Textract operations.
type TextractHandler struct {
	BillService *application.BillService
}

// NewTextractHandler creates a new TextractHandler.
func NewTextractHandler(bs *application.BillService) *TextractHandler {
	if bs == nil {
		// Depending on project logging strategy, log a fatal error or panic.
		// For consistency with other handlers, we could use log.Fatal or return an error if preferred.
		panic("BillService cannot be nil in NewTextractHandler")
	}
	return &TextractHandler{BillService: bs}
}

// AnalyzeBill godoc
// @Summary Analyze a bill image for text
// @Description Upload an image of a bill to detect text using AWS Textract.
// @Tags Textract
// @Accept mpfd
// @Produce json
// @Param image formData file true "Image file of the bill to process"
// @Success 200 {object} gin.H{"detected_lines": []string} "Successfully detected text lines"
// @Failure 400 {object} gin.H{"error": string} "Bad Request - e.g., no file or invalid file"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error"
// @Router /textract/analyze-bill [post]
func (h *TextractHandler) AnalyzeBill(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read image from request: " + err.Error()})
		return
	}
	defer file.Close()

	fmt.Printf("Received file: %s, size: %d, MIME: %s\n", header.Filename, header.Size, header.Header.Get("Content-Type"))

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file into buffer: " + err.Error()})
		return
	}
	//documentBytes := buf.Bytes()

	err = h.BillService.AnalyzeBill(c, uuid.New()) // Using DetectText from the client for now
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Textract failed to process document: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bill analyzed successfully"})
}
