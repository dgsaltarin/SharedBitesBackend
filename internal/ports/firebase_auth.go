package ports

import (
	"context"
)

// FirebaseAuthProvider defines operations for interacting with the external Firebase Auth service.
type FirebaseAuthProvider interface {
	CreateUser(ctx context.Context, email, password, displayName string) (firebaseUID string, err error)
	DeleteUser(ctx context.Context, uid string) error
}
