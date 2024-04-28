package dependencies

import (
	awssession "github.com/dgsaltarin/SharedBitesBackend/internal/common/aws/session"
	HealthCheckHandler "github.com/dgsaltarin/SharedBitesBackend/internal/infrastructure/handlers"
	"go.uber.org/dig"
)

type HandlersContainer struct {
	HealthCheckHandler *HealthCheckHandler.HealthCheckHandler
}

func NewWire() *dig.Container {
	container := dig.New()

	//aws dependencies
	container.Provide(awssession.NewAWSSession)

	container.Provide(HealthCheckHandler.NewHealthCheckHandler)

	return container
}
