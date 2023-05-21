package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgsaltarin/SharedBitesBackend/services"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Image string
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/rekognition", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		defer src.Close()

		data := make([]byte, file.Size)
		_, err = src.Read(data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		fmt.Println(file.Filename)

		session := services.RekognitionSession()

		result := services.DetectLabels(session, data)

		output, err := json.Marshal(result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": string(output),
		})

	})

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
