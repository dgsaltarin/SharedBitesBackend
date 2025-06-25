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
	billHandler *hanlders.BillHandler,
	groupHandler *hanlders.GroupHandler,
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

	// --- Bill Routes --- //
	if billHandler != nil {
		billProtected := protectedRoutes.Group("/bills")
		{
			billProtected.POST("/upload-analyze-config", billHandler.UploadAndAnalyzeBillWithConfig)
			billProtected.GET("/analysis-configs", billHandler.GetAnalysisConfigs)
			billProtected.GET("", billHandler.ListBills)
			billProtected.GET("/:bill_id", billHandler.GetBill)
			billProtected.GET("/:bill_id/status", billHandler.GetBillStatus)
			billProtected.DELETE("/:bill_id", billHandler.DeleteBill)
		}
	} else {
		log.Println("WARN: BillHandler is nil, Bill routes not configured in SetupAppRoutes.")
	}

	// --- Group Routes --- //
	if groupHandler != nil {
		groupProtected := protectedRoutes.Group("/groups")
		{
			groupProtected.POST("", groupHandler.CreateGroup)
			groupProtected.GET("", groupHandler.ListGroups)
			groupProtected.GET("/:group_id", groupHandler.GetGroup)
			groupProtected.PUT("/:group_id", groupHandler.UpdateGroup)
			groupProtected.DELETE("/:group_id", groupHandler.DeleteGroup)
		}
	} else {
		log.Println("WARN: GroupHandler is nil, Group routes not configured in SetupAppRoutes.")
	}
}
