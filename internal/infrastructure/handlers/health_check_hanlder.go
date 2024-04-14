package handlers

import (
	"github.com/dgsaltarin/SharedBitesBackend/internal/application/services"
	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct {
	HealthCheckService services.HealthCheckService
}

func NewHealthCheckHandler() HealthCheckHandler {
	return HealthCheckHandler{
		HealthCheckService: services.HealthCheckService{},
	}
}

func (h *HealthCheckHandler) HealthCheck() (c *gin.Context) {
	c.JSON(200, gin.H{
		"message": h.HealthCheckService.Check(),
	})
	return
}
