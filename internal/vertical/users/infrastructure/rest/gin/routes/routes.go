package router

import (
	handlders "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/rest/gin/handlers"
	"github.com/gin-gonic/gin"
)

type UsersRouter struct {
	group              *gin.RouterGroup
	healthcheckHandler *handlders.HealthCheckHandler
}

// NewRouter creates a new router for the users vertical
func NewUserRoutes(group *gin.RouterGroup, healthcheckHandler *handlders.HealthCheckHandler) *UsersRouter {
	return &UsersRouter{
		group:              group,
		healthcheckHandler: healthcheckHandler,
	}
}

// register defines the routes for the users vertical
func (ur *UsersRouter) register() {
	group := ur.group.Group("/users")
	group.GET("/healthcheck", ur.healthcheckHandler.HealthCheck)
}