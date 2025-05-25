package middleware

import (
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/dgsaltarin/SharedBitesBackend/internal/application"
	"github.com/gin-gonic/gin"
)

type contextKey string

// UserContextKey is the key used to store authenticated user info in the context.
const UserContextKey string = "authenticatedUser"

// AuthenticatedUser holds information about the verified user.
type AuthenticatedUser struct {
	UID string // Firebase User ID
	// Email string
	// Name string
}

// FirebaseAuthMiddleware creates a Gin middleware handler that verifies Firebase ID tokens.
func FirebaseAuthMiddleware(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Auth Middleware: Missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Println("Auth Middleware: Invalid Authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		idToken := parts[1]

		token, err := authClient.VerifyIDToken(c.Request.Context(), idToken) // Use c.Request.Context()
		if err != nil {
			log.Printf("Auth Middleware: Error verifying Firebase ID token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired authentication token"})
			return
		}

		authUser := AuthenticatedUser{
			UID: token.UID,
		}

		c.Set(UserContextKey, authUser)

		c.Next()
	}
}

func UserLookupMiddleware(userService *application.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authUser, exists := GetUserFromGinContext(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		user, err := userService.GetUserByFirebaseUID(c.Request.Context(), authUser.UID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set("userID", user.ID)

		c.Next()
	}
}

// GetUserFromGinContext retrieves the authenticated user from the Gin request context.
func GetUserFromGinContext(c *gin.Context) (AuthenticatedUser, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return AuthenticatedUser{}, false
	}
	authUser, ok := user.(AuthenticatedUser)
	return authUser, ok
}
