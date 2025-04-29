package api

import (
	"github.com/dgsaltarin/SharedBitesBackend/internal/dependencies"
	billsRoutes "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/bills/infrastructure/rest/gin/routes"
	userRoutes "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/rest/gin/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RequestBody struct {
	Image string
}

func main() {
	ginInstance := setupGin()

	routerGrup := ginInstance.Group("/api/v1")

	// Create all handlers directly instead of using dig
	handlers := dependencies.NewHandlers()

	// Setup routes directly
	setupUserRoutes(routerGrup, handlers)
	setupBillsRoutes(routerGrup, handlers)

	ginInstance.Run(":8080")
}

// setupGin creates a new gin instance
func setupGin() *gin.Engine {
	ginInstance := gin.Default()

	ginInstance.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginInstance.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	return ginInstance
}

// setupUserRoutes sets up the routes for the users vertical
func setupUserRoutes(api *gin.RouterGroup, handlers *dependencies.HandlersContainer) {
	userRoutes.NewUserRoutes(api.Group("/users"), handlers.UserHandler)
}

// setupBillsRoutes sets up the routes for the bills vertical
func setupBillsRoutes(api *gin.RouterGroup, handlers *dependencies.HandlersContainer) {
	billsRoutes.NewBillsRoutes(api.Group("/bills"), handlers.BillsHandler)
}
