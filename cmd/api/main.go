package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/dgsaltarin/SharedBitesBackend/config"
	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driven/firebaseauth"
	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driven/sql"
	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driving/rest"
	"github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driving/rest/hanlders"
	appmiddleware "github.com/dgsaltarin/SharedBitesBackend/internal/adapters/driving/rest/middlewares"
	"github.com/dgsaltarin/SharedBitesBackend/internal/application"
	"github.com/dgsaltarin/SharedBitesBackend/platform/aws"
	"github.com/dgsaltarin/SharedBitesBackend/platform/database"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/api/option"
)

// HandlersContainer remains useful for organizing handlers within main.
type HandlersContainer struct {
	UserHandler     *hanlders.UserHandler
	TextractHandler *hanlders.TextractHandler
}

// --- Main Application Setup ---
func main() {
	tempLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.MustLoad(tempLogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	firebaseApp, err := initFirebase(ctx, cfg.Firebase.ServiceAccountKeyPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Auth client: %v", err)
	}
	firebaseAuthProvider := firebaseauth.NewFirebaseAuthProvider(authClient)

	db := database.MustConnectGORM(cfg.Database)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	textractSvc, err := aws.NewTextractClient(ctx, cfg.AWS)
	if err != nil {
		log.Printf("WARN: Failed to initialize AWS Textract client: %v. Textract features unavailable.", err)
	}

	userRepo := sql.NewGORMUserRepository(db)
	userService := application.NewUserService(userRepo, firebaseAuthProvider)

	userHandler := hanlders.NewUserHandler(*userService)
	var textractHandler *hanlders.TextractHandler
	if textractSvc != nil {
		textractHandler = hanlders.NewTextractHandler(textractSvc)
	} else {
		log.Println("INFO: TextractHandler not initialized.")
	}

	router := setupRouter(userHandler, textractHandler, authClient)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if errSrv := server.ListenAndServe(); errSrv != nil && errSrv != http.ErrServerClosed {
			log.Fatalf("Server error: %v", errSrv)
		}
	}()

	<-signalChan
	log.Println("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped successfully")
}

func initFirebase(ctx context.Context, serviceAccountKeyPath string) (*firebase.App, error) {
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, errFirebase := firebase.NewApp(ctx, nil, opt) // Renamed err to avoid conflict
	if errFirebase != nil {
		return nil, fmt.Errorf("error initializing Firebase app: %w", errFirebase)
	}
	return app, nil
}

func setupRouter(userHandler *hanlders.UserHandler, textractHandler *hanlders.TextractHandler, authClient *auth.Client) *gin.Engine {
	router := gin.Default()

	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group("/api/v1")

	publicApiV1 := apiV1.Group("/")

	protectedApiV1 := apiV1.Group("/")
	// Use the Gin-native Firebase auth middleware directly
	protectedApiV1.Use(appmiddleware.FirebaseAuthMiddleware(authClient))

	rest.SetupAppRoutes(publicApiV1, protectedApiV1, userHandler, textractHandler)

	return router
}
