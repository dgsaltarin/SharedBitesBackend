package main

import (
	services "github.com/dgsaltarin/SharedBitesBackend/internal/application/services"
	"github.com/dgsaltarin/SharedBitesBackend/internal/infrastructure/handlers"
	"github.com/dgsaltarin/SharedBitesBackend/internal/infrastructure/router"
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

	router := router.NewRouter(&healthcheckHandler, &handlers)

	router.SetupRouter()
}
