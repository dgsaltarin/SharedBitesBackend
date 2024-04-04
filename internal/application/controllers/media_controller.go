package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dgsaltarin/SharedBitesBackend/helpers"
	"github.com/dgsaltarin/SharedBitesBackend/services"

	"github.com/gin-gonic/gin"
)

// default function
func HelloWorld() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	}
}

// UploadImage function upload image to s3 and start analyze to detect items
func UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		image, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		data, err := helpers.DecodeImage(image)

		aws_session := services.AWSSession()

		services.UploadImages3(aws_session, data, image.Filename)
		session := services.TextTrackSesson(aws_session)
		result := services.Detectitems(session, image.Filename)

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
	}
}
