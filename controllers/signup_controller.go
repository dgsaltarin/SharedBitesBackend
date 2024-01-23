package controllers

import (
	"fmt"
	"net/http"

	"github.com/dgsaltarin/SharedBitesBackend/db"
	"github.com/dgsaltarin/SharedBitesBackend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SignUp function create a new user in dynamodb
func SignUp() gin.HandlerFunc {
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

		// check if user already exists
		dbuser, err := userdb.GetUserByUsername(user.Username)
		if dbuser != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "User already exists",
			})
			fmt.Println(err)
			return
		}

		// add uuid to user
		user.ID = uuid.New().String()

		// hashnig password
		err = user.GeneratePasswordHash()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error hashing password",
			})
			fmt.Println(err)
			return
		}

		// create user in dynamodb
		err = userdb.UpserUser(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error creating user",
			})
			fmt.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}
