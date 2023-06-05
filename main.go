package main

import (
	"github.com/dgsaltarin/SharedBitesBackend/controllers"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Image string
}

func main() {
	router := gin.Default()

	router.GET("/hello", controllers.HelloWorld())
	router.POST("/texttrack", controllers.UploadImage())
	router.POST("/users", controllers.CreateUser())

	router.Run()
}
