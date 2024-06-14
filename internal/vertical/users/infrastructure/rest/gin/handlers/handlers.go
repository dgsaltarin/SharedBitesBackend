package handlers

import (
	"fmt"
	"net/http"

	services "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/application"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/mappers"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/rest/request"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
	mapper      mappers.Mappers
}

func NewUserHandler(userService services.UserService, mapper mappers.Mappers) *UserHandler {
	return &UserHandler{
		userService: userService,
		mapper:      mapper,
	}
}

func (uh *UserHandler) SignUp(c *gin.Context) {
	var userReuest request.SignUpRequest

	if err := c.ShouldBindJSON(&userReuest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("Invalid request: %v", err)})
		return
	}

	user := uh.mapper.MapSignUpRequestToUser(userReuest)

	if err := uh.userService.SignUp(*user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Error signing up: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (uh *UserHandler) Login(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Login"})
}
