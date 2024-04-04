package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	services "github.com/dgsaltarin/SharedBitesBackend/internal/application"
	"github.com/dgsaltarin/SharedBitesBackend/internal/infrastructure/rest/request"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return UserHandler{
		userService: userService,
	}
}

func (uh *UserHandler) SignUp() (c *gin.Context) {
	// bin user info from request
	var user request.SignUpRequest
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

func (uh *UserHandler) Login() (c *gin.Context) {
	// get user from request
	var userRequest *models.User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// new dynamodb database
	dyanmodb := db.GetDynamoDBInstance()

	// new user database
	userdb := db.NewUserDB(dyanmodb)

	user, err := userdb.GetUserByUsername(userRequest.Username)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Invalid username or password",
		})
		return
	}

	// check if password is correct
	if user.CheckPassword(userRequest.Password) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Username,
			"exp": time.Now().Add(time.Minute * 10).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})

	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid email or password",
		})
	}

}
