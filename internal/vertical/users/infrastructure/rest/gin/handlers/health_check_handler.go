package handlers

import (
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/application/services"
	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct {
	HealthCheckService services.HealthCheckService
}

func NewHealthCheckHandler(healthCheckService services.HealthCheckService) *HealthCheckHandler {
	return &HealthCheckHandler{
		HealthCheckService: healthCheckService,
	}
}

func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": h.HealthCheckService.Check(),
	})
}
