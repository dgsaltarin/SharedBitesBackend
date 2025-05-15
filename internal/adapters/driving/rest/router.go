package rest

import (
	"log"

	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driving/rest/hanlders"
	"github.com/gin-gonic/gin"
)

// SetupAppRoutes configures application routes on the provided public and protected router groups.
func SetupAppRoutes(
	publicRoutes *gin.RouterGroup,
	protectedRoutes *gin.RouterGroup,
	userHandler *hanlders.UserHandler,
	textractHandler *hanlders.TextractHandler,
) {
	// --- User Routes --- //
	if userHandler != nil {
		// Public user routes
		userPublic := publicRoutes.Group("/users")
		{
			userPublic.POST("", userHandler.HandleCreateUser) // e.g., User registration
		}

		// Protected user routes
		userProtected := protectedRoutes.Group("/users")
		{
			userProtected.GET("/:userID", userHandler.HandleGetUser)
			userProtected.PUT("/:userID", userHandler.HandleUpdateUserProfile)
			// Add other user routes that need protection here
		}
	} else {
		log.Println("WARN: UserHandler is nil, user routes not configured in SetupAppRoutes.")
	}

	// --- Textract Routes --- //
	if textractHandler != nil {
		// Example: Making Textract routes protected by default
		// If you want some Textract routes to be public, create a group on 'publicRoutes'
		textractProtected := protectedRoutes.Group("/textract")
		{
			textractProtected.POST("/analyze-bill", textractHandler.AnalyzeBill)
			// Add other textract routes that need protection here
		}
	} else {
		log.Println("WARN: TextractHandler is nil, Textract routes not configured in SetupAppRoutes.")
	}
}
