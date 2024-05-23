package main

import (
	"github.com/dgsaltarin/SharedBitesBackend/internal/dependencies"
	billsRoutes "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/bills/infrastructure/rest/gin/routes"
	userRoutes "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/rest/gin/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/dig"
)

type RequestBody struct {
	Image string
}

func main() {
	ginInstance := setupGin()

	routerGrup := ginInstance.Group("/api/v1")
	container := dependencies.NewWire()
	if err := invokeDependencyInjection(container, routerGrup); err != nil {
		panic(err)
	}

	if err := invokeDependencyInjectionBills(container, routerGrup); err != nil {
		panic(err)
	}

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

// invokeDependencyInjection invokes the dependency injection for the users vertical
func invokeDependencyInjection(container *dig.Container, api *gin.RouterGroup) error {
	return container.Invoke(func(h *dependencies.HandlersContainer) {
		userRoutes.NewUserRoutes(api.Group("/users"), h.HealthCheckHandler)
	})
}

func invokeDependencyInjectionBills(container *dig.Container, api *gin.RouterGroup) error {
	return container.Invoke(func(h *dependencies.HandlersContainer) {
		billsRoutes.NewBillsRoutes(api.Group("/bills"), h.BillsHandler)
	})
}
