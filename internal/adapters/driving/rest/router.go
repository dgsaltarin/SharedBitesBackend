package rest

import (
	"firebase.google.com/go/v4/auth"
	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driving/rest/hanlders"
	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driving/rest/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetupUserRouter configures all user-related routes with appropriate middleware
func SetupUserRouter(router *gin.Engine, userHandler *hanlders.UserHandler, authClient *auth.Client) {
	// Adapt Firebase middleware to work with Gin
	firebaseAuth := func(c *gin.Context) {
		// Create a middleware handler that adapts the Firebase middleware to Gin
		handler := middleware.FirebaseAuthMiddleware(authClient)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Copy the context with the authenticated user to Gin's context
			c.Request = r
			c.Next()
		}))

		// Call the middleware and short-circuit if necessary
		responseWriter := &responseWriterInterceptor{ResponseWriter: c.Writer}
		handler.ServeHTTP(responseWriter, c.Request)

		// If middleware wrote a response, abort the request
		if responseWriter.statusWritten {
			c.Abort()
		}
	}

	// Group for user routes
	userRoutes := router.Group("/api/v1/users")
	{
		// Public routes (no auth required)
		userRoutes.POST("", userHandler.HandleCreateUser) // Register new user

		// Protected routes (require authentication)
		protected := userRoutes.Group("")
		protected.Use(firebaseAuth)
		{
			protected.GET("/:userID", userHandler.HandleGetUser)
			protected.PUT("/:userID", userHandler.HandleUpdateUserProfile)
		}
	}
}

// responseWriterInterceptor helps us determine if the middleware wrote a response
type responseWriterInterceptor struct {
	gin.ResponseWriter
	statusWritten bool
}

func (w *responseWriterInterceptor) WriteHeader(code int) {
	w.statusWritten = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriterInterceptor) Write(b []byte) (int, error) {
	w.statusWritten = true
	return w.ResponseWriter.Write(b)
}
