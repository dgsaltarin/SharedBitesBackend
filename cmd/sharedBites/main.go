package main

import (
	services "github.com/dgsaltarin/SharedBitesBackend/internal/application/services"
	"github.com/dgsaltarin/SharedBitesBackend/internal/dependencies"
	"github.com/dgsaltarin/SharedBitesBackend/internal/infrastructure/handlers"
	mainRouter "github.com/dgsaltarin/SharedBitesBackend/internal/infrastructure/router"
	"go.uber.org/dig"
)

type RequestBody struct {
	Image string
}

func main() {
	// Dependency Injection
	healtcheckService := services.NewHealthCheckService()
	healthcheckHandler := handlers.NewHealthCheckHandler(&healtcheckService)
	billService := services.NewBillService()
	handlers := handlers.NewHanlder(&billService)

	router := mainRouter.NewRouter(&healthcheckHandler, &handlers)
	container := dependencies.NewWire()
	if err := invokeDependencyInjection(container, router); err != nil {
		panic(err)
	}

	router.SetupRouter()
}

func invokeDependencyInjection(container *dig.Container, api *mainRouter.Router) error {
	return container.Invoke(func(h *dependencies.HandlersContainer) {
		api.HealthCheckHandler = h.HealthCheckHandler
		api.Handler = h.Handler
	})
}
