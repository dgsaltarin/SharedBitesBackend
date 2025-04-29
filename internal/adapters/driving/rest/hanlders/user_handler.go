package hanlders

import (
	"errors"
	"github.com/dgsaltarin/SharedBitesBackend/internal/application"
	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"log" // Use your logger
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
// IMPORTANT: Never return sensitive info like passwords or Firebase UIDs unless absolutely necessary.
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

	// Bind JSON request to the struct and validate using 'binding' tags
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("User creation bind/validation error: %v", err)
		// Provide more specific error messages based on validation failures if desired
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Call the application service to create the user
	createdUser, err := h.userService.CreateUser(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		log.Printf("Error calling CreateUser service for email %s: %v", req.Email, err)

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
			// Log the internal error but return a generic server error message
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user due to an internal issue."})
		default:
			// Catch-all for other unexpected errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		}
		return
	}

	// Map the domain user to the response DTO
	response := CreateUserResponse{
		ID:        createdUser.ID, // ID is populated after successful DB save
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}
