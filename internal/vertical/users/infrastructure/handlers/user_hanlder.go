package handlers

import (
	services "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/application"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return UserHandler{
		userService: userService,
	}
}

func (uh *UserHandler) SignUp(c *gin.Context) {
	c.JSON(200, gin.H{"message": "SignUp"})
}

func (uh *UserHandler) Login(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Login"})
}
