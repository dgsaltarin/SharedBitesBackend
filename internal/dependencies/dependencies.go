package dependencies

import (
	awssession "github.com/dgsaltarin/SharedBitesBackend/internal/common/aws/session"
	billsHandler "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/bills/infrastructure/rest/gin/handlers"
	userService "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/application/services"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/rest/gin/handlers"
	userHandler "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/rest/gin/handlers"

	"go.uber.org/dig"
)

type HandlersContainer struct {
	HealthCheckHandler *userHandler.HealthCheckHandler
	BillsHandler       *billsHandler.BillsHandler
}

// NewWire creates a new container with all the dependencies
func NewWire() *dig.Container {
	container := dig.New()

	// aws dependencies
	container.Provide(awssession.NewAWSSession)

	// user dependencies
	container.Provide(userService.NewUserService)
	container.Provide(userService.NewHealthCheckService)
	container.Provide(userHandler.NewUserHandler)
	container.Provide(userHandler.NewHealthCheckHandler)
	container.Provide(billsHandler.NewBillsHandler)

	// handlers dependencies
	container.Provide(
		func(
			healthcheckHandler *handlers.HealthCheckHandler,
		) *HandlersContainer {
			return &HandlersContainer{
				HealthCheckHandler: healthcheckHandler,
			}
		})

	return container
}
