package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgsaltarin/SharedBitesBackend/db"
	"github.com/dgsaltarin/SharedBitesBackend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authorize(c *gin.Context) {
	// Get token from authorization cookie
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}

	// Parse token and validate token that is encrypted using HS256
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used for signing
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	// Validate that token is not expired and is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Validate token is not expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user *models.User
		fmt.Println(claims["sub"].(string))

		// new dynamodb database
		dynamodb := db.GetDynamoDBInstance()

		// new user database
		userdb := db.NewUserDB(dynamodb)
		user, err = userdb.GetUserByUsername(claims["sub"].(string))

		if user.Username == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		fmt.Println("claims error")
		return
	}
}
