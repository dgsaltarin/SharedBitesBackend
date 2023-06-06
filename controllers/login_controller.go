package controllers

import (
	"fmt"
	"net/http"

	"github.com/dgsaltarin/SharedBitesBackend/db"
	"github.com/dgsaltarin/SharedBitesBackend/models"
	"github.com/gin-gonic/gin"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user from request
		var userRequest *models.User
		if err := c.ShouldBindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			fmt.Println(err)
			return
		}

		// new dynamodb database
		dyanmodb, err := db.Connect()
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Error connecting to database",
			})
			fmt.Println(err)
			return
		}

		// new user database
		userdb := db.NewUserDB(dyanmodb)

		user, err := userdb.GetUser(userRequest.ID)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Error getting user",
			})
			fmt.Println(err)
			return
		}

		// check if password is correct
		if user.Username != userRequest.Username || user.Password != userRequest.Password {
			c.JSON(500, gin.H{
				"message": "Incorrect username or password",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "User logged in successfully",
		})
	}
}
