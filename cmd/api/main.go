package main

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/dgsaltarin/SharedBitesBackend/config"
	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driven/firebaseauth"
	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driven/sql"
	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driving/rest"
	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driving/rest/hanlders"
	"github.com/dgsaltarin/SharedBitesBackend/internal/application"
	"github.com/dgsaltarin/SharedBitesBackend/platform/database"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/api/option"
)

func main() {
	// Load configuration
	tempLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.MustLoad(tempLogger)

	// Initialize context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Initialize Firebase
	firebaseApp, err := initFirebase(ctx, cfg.Firebase.ServiceAccountKeyPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Auth client: %v", err)
	}

	// New Firebase Auth Provider
	firebaseAuthProvider := firebaseauth.NewFirebaseAuthProvider(authClient)

	// Initialize database connection with GORM
	db := database.MustConnectGORM(cfg.Database)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Initialize repositories
	userRepo := sql.NewGORMUserRepository(db)

	// Initialize services
	userService := application.NewUserService(userRepo, firebaseAuthProvider)

	// Initialize handlers
	userHandler := hanlders.NewUserHandler(*userService)

	// Create a container for handlers
	handlers := HandlersContainer{
		UserHandler: userHandler,
	}

	// Set up router and routes
	router := setupRouter(handlers, authClient)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for termination signal
	<-signalChan
	log.Println("Shutdown signal received")

	// Create shutdown timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown server gracefully
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped successfully")
}

// initFirebase initializes the Firebase app
func initFirebase(ctx context.Context, serviceAccountKeyPath string) (*firebase.App, error) {
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// setupRouter configures all routes and middleware
func setupRouter(handlers HandlersContainer, authClient *auth.Client) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup user routes with Firebase auth
	rest.SetupUserRouter(router, handlers.UserHandler, authClient)

	// Setup other routes as needed

	return router
}

// HandlersContainer holds all handlers
type HandlersContainer struct {
	UserHandler *hanlders.UserHandler
}
