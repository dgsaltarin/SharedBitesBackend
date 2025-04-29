package hanlders

import (
	"errors"
	"github.com/dgsaltarin/SharedBitesBackend/internal/application"
	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	userService application.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(us application.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

// --- DTOs (Data Transfer Objects) ---

// CreateUserRequest defines the expected JSON body for creating a user.
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"` // Use Gin's binding validator
	Password string `json:"password" binding:"required"`
}

// CreateUserResponse defines the JSON response after successfully creating a user.
type CreateUserResponse struct {
	ID        uuid.UUID `json:"id"` // Internal DB ID
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// --- Handler Methods ---

// HandleCreateUser handles the POST /users request.
func (h *UserHandler) HandleCreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Call the application service to create the user
	createdUser, err := h.userService.CreateUser(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		// Map domain errors to HTTP status codes
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrUserNameEmpty),
			errors.Is(err, domain.ErrUserEmailEmpty),
			errors.Is(err, domain.ErrUserPasswordTooShort):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrFirebaseUserCreationFailed),
			errors.Is(err, domain.ErrDatabaseUserCreationFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user due to an internal issue."})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		}
		return
	}

	response := CreateUserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) HandleUpdateUserProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Call the application service to update the user profile
	updatedUser, err := h.userService.UpdateUserProfile(c.Request.Context(), userID, &req.Name)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrUserNameEmpty),
			errors.Is(err, domain.ErrUserEmailEmpty):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		}
		return
	}

	response := CreateUserResponse{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) HandleGetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUserByFirebaseUID(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		}
		return
	}

	response := CreateUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}
