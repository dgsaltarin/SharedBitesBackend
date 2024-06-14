package router

import (
	handlders "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/rest/gin/handlers"
	"github.com/gin-gonic/gin"
)

type usersRouter struct {
	group    *gin.RouterGroup
	handlers *handlders.UserHandler
}

// NewRouter creates a new router for the users vertical
func NewUserRoutes(group *gin.RouterGroup, handlers *handlders.UserHandler) *usersRouter {
	usersRouter := &usersRouter{
		group:    group,
		handlers: handlers,
	}

	usersRouter.register()

	return usersRouter
}

// register defines the routes for the users vertical
func (ur *usersRouter) register() {
	ur.group.POST("/signup", ur.handlers.SignUp)
}
