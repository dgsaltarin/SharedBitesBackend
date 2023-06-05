package controllers

import (
	"github.com/dgsaltarin/SharedBitesBackend/db"
	"github.com/dgsaltarin/SharedBitesBackend/models"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	// bin user info from request
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"message": "Error binding user info",
		})
		return
	}

	// new dyanmodb database
	dyanmodb, err := db.Connect()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error connecting to database",
		})
		return
	}

	// new user database
	userdb := db.NewUserDB(dyanmodb)

	// create user in dynamodb
	err = userdb.UpserUser(&user)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error creating user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User created successfully",
	})
}
