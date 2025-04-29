package rest

import (
	"firebase.google.com/go/v4/auth"
	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driving/rest/hanlders"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupGinRouter configures the Gin HTTP router.
func SetupGinRouter(
	userHandler *hanlders.UserHandler, // NEW
	authMiddleware gin.HandlerFunc, // Firebase Auth Middleware adapted for Gin
) *gin.Engine {

	// gin.SetMode(gin.ReleaseMode) // Set for production
	r := gin.Default() // Includes Logger and Recovery middleware

	// --- Base Middleware ---
	// Example CORS config (adjust origins etc.)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Be more specific in production!
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// MaxAge: 12 * time.Hour,
	}))
	// Add other global middleware if needed

	// --- API Routes v1 ---
	apiV1 := r.Group("/api/v1")
	{
		// --- Public Routes ---
		apiV1.POST("/users", userHandler.HandleCreateUser) // User creation is public initially
		// Example: Getting public group details might be public
		apiV1.GET("/groups/:groupID", groupHandler.HandleGetGroup) // Adapt handler for Gin context

		// --- Protected Routes (Require Authentication) ---
		protected := apiV1.Group("/")
		protected.Use(authMiddleware) // Apply Firebase Auth middleware
		{
			// Group related endpoints
			protected.POST("/groups", groupHandler.HandleCreateGroup) // Adapt handler for Gin
			protected.PUT("/groups/:groupID", groupHandler.HandleUpdateGroup)
			protected.POST("/groups/:groupID/members", groupHandler.HandleAddMemberByName)
			protected.DELETE("/groups/:groupID/members/:memberID", groupHandler.HandleRemoveMember)

			// Expense related endpoints
			protected.POST("/groups/:groupID/expenses", expenseHandler.HandleCreateExpense)
			protected.GET("/groups/:groupID/expenses", expenseHandler.HandleListExpenses)
			protected.GET("/expenses/:expenseID", expenseHandler.HandleGetExpense)
			protected.PUT("/expenses/:expenseID", expenseHandler.HandleUpdateExpense)
			protected.DELETE("/expenses/:expenseID", expenseHandler.HandleDeleteExpense)

			// Current user endpoint
			// protected.GET("/users/me", userHandler.HandleGetMe) // Needs implementation
		}

		// --- Health Check ---
		r.GET("/health", func(c *gin.Context) {
			// Add DB ping check if needed
			c.Status(http.StatusOK)
		})
	}

	return r
}

// Adapt your Firebase middleware for Gin:
func AdaptFirebaseAuthMiddleware(authClient *auth.Client) gin.HandlerFunc {
	// Wrap your existing http.Handler middleware logic inside a gin.HandlerFunc
	// See previous Firebase middleware example for the core token verification logic.
	// Use c.GetHeader("Authorization"), c.AbortWithStatusJSON(), c.Set(), c.Next()
	// Example structure:
	return func(c *gin.Context) {
		// ... (Get token from c.GetHeader("Authorization")) ...
		// ... (Verify token using authClient) ...
		// if err != nil {
		//     c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		//     return
		// }
		// authUser := middleware.AuthenticatedUser{ UID: token.UID } // From previous example
		// c.Set(string(middleware.UserContextKey), authUser) // Use the same context key
		// c.Next()
	}
}

// In your handlers, get user like this:
// func GetUserFromGinContext(c *gin.Context) (middleware.AuthenticatedUser, bool) {
//     userVal, exists := c.Get(string(middleware.UserContextKey))
//     if !exists { return middleware.AuthenticatedUser{}, false }
//     user, ok := userVal.(middleware.AuthenticatedUser)
//     return user, ok
// }
