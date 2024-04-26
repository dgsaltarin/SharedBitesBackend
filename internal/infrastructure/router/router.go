package router

import (
	handlders "github.com/dgsaltarin/SharedBitesBackend/internal/infrastructure/handlers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	healthcheckHandler *handlders.HealthCheckHandler
	hanlders           *handlders.Handler
}

func NewRouter(healthcheckHandler *handlders.HealthCheckHandler, hanlder *handlders.Handler) *Router {
	return &Router{
		healthcheckHandler: healthcheckHandler,
		hanlders:           hanlder,
	}
}

func (r *Router) SetupRouter() {
	router := gin.Default()

	router.GET("/healthcheck", r.healthcheckHandler.HealthCheck)
	router.POST("/bills", r.hanlders.SplitBill)

	router.Run(":8080")
}
