package firebaseauth

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/dgsaltarin/SharedBitesBackend/internal/ports"
	"log"

	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/errorutils"
)

type firebaseAuthProvider struct {
	client *auth.Client
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
		EmailVerified(false). // let's not verify email for now
		Disabled(false)

	firebaseUser, err := p.client.CreateUser(ctx, params)
	if err != nil {
		if auth.IsEmailAlreadyExists(err) {
			log.Printf("Firebase createUser failed: email %s already exists", email)
			return "", fmt.Errorf("firebase: %w", domain.ErrUserAlreadyExists)
		}
		log.Printf("Firebase createUser failed for email %s: %v", email, err)
		return "", fmt.Errorf("%w: %v", domain.ErrFirebaseUserCreationFailed, err)
	}

	log.Printf("Successfully created user in Firebase: UID=%s, Email=%s", firebaseUser.UID, email)
	return firebaseUser.UID, nil
}

func (p *firebaseAuthProvider) DeleteUser(ctx context.Context, uid string) error {
	err := p.client.DeleteUser(ctx, uid)
	if err != nil {
		// Check if the user simply doesn't exist (might happen in complex rollback scenarios)
		if errorutils.IsNotFound(err) {
			log.Printf("Firebase deleteAuthUser failed: uid %s does not exist", uid)
			return nil
		}

		log.Printf("Firebase deleteAuthUser failed for uid %s: %v", uid, err)
		return errors.Join(errors.New("failed to delete firebase user"), err) // Wrap the error
	}

	log.Printf("Successfully deleted user in Firebase: UID=%s", uid)
	return nil
}
