package firebaseauth

import (
	"context"
	"fmt"
	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/dgsaltarin/SharedBitesBackend/internal/ports"
	"log"

	"firebase.google.com/go/v4/auth"
)

type firebaseAuthProvider struct {
	client *auth.Client // Injected Firebase Auth client
}

// NewFirebaseAuthProvider creates a new adapter for Firebase Auth interactions.
func NewFirebaseAuthProvider(client *auth.Client) ports.FirebaseAuthProvider {
	if client == nil {
		log.Fatal("Firebase Auth client cannot be nil") // Or return an error
	}
	return &firebaseAuthProvider{client: client}
}

// CreateUser implements ports.FirebaseAuthProvider.
func (p *firebaseAuthProvider) CreateUser(ctx context.Context, email, password, displayName string) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(displayName).
		EmailVerified(false). // You might want email verification flow later
		Disabled(false)

	firebaseUser, err := p.client.CreateUser(ctx, params)
	if err != nil {
		// Map Firebase specific errors to domain errors
		if auth.IsEmailAlreadyExists(err) {
			log.Printf("Firebase createUser failed: email %s already exists", email)
			return "", fmt.Errorf("firebase: %w", domain.ErrUserAlreadyExists)
		}
		// Add other specific Firebase error mappings if needed
		log.Printf("Firebase createUser failed for email %s: %v", email, err)
		return "", fmt.Errorf("%w: %v", domain.ErrFirebaseUserCreationFailed, err)
	}

	log.Printf("Successfully created user in Firebase: UID=%s, Email=%s", firebaseUser.UID, email)
	return firebaseUser.UID, nil
}
