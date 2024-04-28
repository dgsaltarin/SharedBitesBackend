package dependencies

import (
	awssession "github.com/dgsaltarin/SharedBitesBackend/internal/common/aws/session"
	billsHandler "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/bills/infrastructure/rest/gin/handlers"
	userHandler "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/rest/gin/handlers"
	"go.uber.org/dig"
)

type HandlersContainer struct {
	userHandler  *userHandler.UserHandler
	billsHandler *billsHandler.BillsHandler
}

func NewWire() *dig.Container {
	container := dig.New()

	//aws dependencies
	container.Provide(awssession.NewAWSSession)

	//handlers dependencies
	container.Provide(userHandler.NewUserHandler)
	container.Provide(billsHandler.NewBillsHandler)

	return container
}
