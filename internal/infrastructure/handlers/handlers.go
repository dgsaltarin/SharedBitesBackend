package handlers

import (
	services "github.com/dgsaltarin/SharedBitesBackend/internal/application/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	billService services.BillService
}

func NewHanlder(billservice *services.BillService) Handler {
	return Handler{
		billService: *billservice,
	}
}

func (h *Handler) SplitBill(c *gin.Context) {
	image, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	result, err := h.billService.SplitBill(image)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error splitting bill"})
		return
	}
	c.JSON(200, gin.H{"message": result})
}
