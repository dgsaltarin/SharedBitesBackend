package router

import (
	handlders "github.com/dgsaltarin/SharedBitesBackend/internal/infrastructure/handlers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	healthcheckHandler *handlders.HealthCheckHandler
}

func NewRouter(healthcheckHandler *handlders.HealthCheckHandler) *Router {
	return &Router{
		healthcheckHandler: healthcheckHandler,
	}
}

func (r *Router) SetupRouter() {
	router := gin.Default()

	router.GET("/healthcheck", r.healthcheckHandler.HealthCheck)

	router.Run(":8080")
}
