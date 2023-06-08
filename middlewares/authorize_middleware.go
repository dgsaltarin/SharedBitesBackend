package middlewares

import (
	"fmt"
	"net/http"
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

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return token, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user *models.User
		fmt.Println(claims["sub"].(string))

		// new dynamodb database
		dynamodb, err := db.Connect()
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Error connecting to database",
			})
			fmt.Println(err)
			return
		}

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
		return
	}
}
