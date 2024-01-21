package main

import (
	"github.com/dgsaltarin/SharedBitesBackend/controllers"
	"github.com/dgsaltarin/SharedBitesBackend/middlewares"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Image string
}

func main() {
	router := gin.Default()

	router.GET("/hello", middlewares.Authorize, controllers.HelloWorld())
	router.POST("/texttrack", controllers.UploadImage())
	router.POST("/users", controllers.CreateUser())
	router.GET("/users", controllers.GetUserByUsername())
	router.POST("/login", controllers.Login())
	router.POST("/signup", controllers.SignUp())
	router.GET("/healthcheck", controllers.HealthCheck())

	router.Run()
}
