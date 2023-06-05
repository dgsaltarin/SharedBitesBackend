package controllers

import (
	"fmt"
	"net/http"

	"github.com/dgsaltarin/SharedBitesBackend/db"
	"github.com/dgsaltarin/SharedBitesBackend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// bin user info from request
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// new dyanmodb database
		dyanmodb, err := db.Connect()
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Error connecting to database",
			})
			fmt.Println(err)
			return
		}

		// add uuid to user
		user.ID = uuid.New().String()

		// new user database
		userdb := db.NewUserDB(dyanmodb)

		// create user in dynamodb
		err = userdb.UpserUser(&user)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Error creating user",
			})
			fmt.Println(err)
			return
		}

		c.JSON(200, gin.H{
			"message": "User created successfully",
		})
	}
}
