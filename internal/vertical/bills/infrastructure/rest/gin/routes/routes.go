package routes

import (
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/bills/infrastructure/rest/gin/handlers"
	"github.com/gin-gonic/gin"
)

type BillsRoutes struct {
	group   *gin.RouterGroup
	hanlder *handlers.BillsHandler
}

// NewBillsRoutes creates a new router for the bills vertical
func NewBillsRoutes(group *gin.RouterGroup, handler *handlers.BillsHandler) *BillsRoutes {
	return &BillsRoutes{
		group:   group,
		hanlder: handler,
	}
}

// register defines the routes for the bills vertical
func (br *BillsRoutes) register() {
	group := br.group.Group("/bills")
	group.POST("/")
}
