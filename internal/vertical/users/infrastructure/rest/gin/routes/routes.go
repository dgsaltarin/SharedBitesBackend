package router

import (
	handlders "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/rest/gin/handlers"
	"github.com/gin-gonic/gin"
)

type usersRouter struct {
	group              *gin.RouterGroup
	healthcheckHandler *handlders.HealthCheckHandler
}

// NewRouter creates a new router for the users vertical
func NewUserRoutes(group *gin.RouterGroup, healthcheckHandler *handlders.HealthCheckHandler) *usersRouter {
	usersRouter := &usersRouter{
		group:              group,
		healthcheckHandler: healthcheckHandler,
	}

	usersRouter.register()

	return usersRouter
}

// register defines the routes for the users vertical
func (ur *usersRouter) register() {
	ur.group.GET("/healthcheck", ur.healthcheckHandler.HealthCheck)
}
