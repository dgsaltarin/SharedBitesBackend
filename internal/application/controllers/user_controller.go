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
		dyanmodb := db.GetDynamoDBInstance()

		// add uuid to user
		user.ID = uuid.New().String()

		// new user database
		userdb := db.NewUserDB(dyanmodb)

		// create user in dynamodb
		err := userdb.UpserUser(&user)
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

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the query parameter value from the request
		userID := c.Query("id")

		// new dyanmodb database
		dyanmodb := db.GetDynamoDBInstance()

		// new user database
		userdb := db.NewUserDB(dyanmodb)

		// get users from dynamodb
		users, err := userdb.GetUser(userID)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Error getting users",
			})
			fmt.Println(err)
			return
		}

		c.JSON(200, gin.H{
			"users": users,
		})
	}
}

// GetUserByUsername is the method in charge of obtain the user by the username
func GetUserByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		// bin user info from request
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// new dyanmodb database
		dyanmodb := db.GetDynamoDBInstance()

		// new user database
		userdb := db.NewUserDB(dyanmodb)

		// get users from dynamodb
		users, err := userdb.GetUserByUsername(user.Username)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Error getting users",
			})
			fmt.Println(err)
			return
		}

		c.JSON(200, gin.H{
			"users": users,
		})
	}
}
