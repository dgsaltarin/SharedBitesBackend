package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dgsaltarin/SharedBitesBackend/services"

	"github.com/gin-gonic/gin"
)

func UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		image, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		src, err := image.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		defer src.Close()

		data := make([]byte, image.Size)
		_, err = src.Read(data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

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
	}
}
