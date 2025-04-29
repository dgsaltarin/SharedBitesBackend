package firebase

import (
	"context"
	"fmt"
	"github.com/dgsaltarin/SharedBitesBackend/config"
	"log" // Use your preferred logger

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// NewFirebaseAuthClient initializes the Firebase App and returns an Auth client.
func NewFirebaseAuthClient(ctx context.Context, cfg config.FirebaseConfig) (*auth.Client, error) {
	if cfg.ServiceAccountKeyPath == "" {
		return nil, fmt.Errorf("firebase service account key path is required")
	}

	opt := option.WithCredentialsFile(cfg.ServiceAccountKeyPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Firebase Auth client: %w", err)
	}

	log.Println("Firebase Auth client initialized successfully.") // Use logger
	return client, nil
}

// MustNewFirebaseAuthClient is like NewFirebaseAuthClient but panics on error.
func MustNewFirebaseAuthClient(ctx context.Context, cfg config.FirebaseConfig) *auth.Client {
	client, err := NewFirebaseAuthClient(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize Firebase Auth client: %v", err)
	}
	return client
}
