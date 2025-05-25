package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"unicode/utf8"

	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"
	"github.com/dgsaltarin/SharedBitesBackend/internal/ports"
)

const minPasswordLength = 8 // Example minimum password length

// UserService implements the UserService interface.
type UserService struct {
	userRepo     ports.UserRepository
	firebaseAuth ports.FirebaseAuthProvider
}

// NewUserService creates a new user service instance.
func NewUserService(ur ports.UserRepository, fa ports.FirebaseAuthProvider) *UserService {
	return &UserService{
		userRepo:     ur,
		firebaseAuth: fa,
	}
}

// CreateUser orchestrates creating a user in Firebase and the local database.
func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	if name == "" {
		return nil, domain.ErrUserNameEmpty
	}

	// TODO add email validation
	if email == "" {
		return nil, domain.ErrUserEmailEmpty
	}

	if utf8.RuneCountInString(password) < minPasswordLength {
		return nil, domain.ErrUserPasswordTooShort
	}

	// 2. Check if user exists in local DB (cheaper check first)
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, fmt.Errorf("failed checking existing user: %w", err)
	}

	// --- Start "transaction-like" block (with compensation) ---
	var createdFirebaseUID string = ""

	defer func() {
		// Compensation: If Firebase user was created but DB save failed, delete from Firebase.
		if err != nil && createdFirebaseUID != "" { // Check if an error occurred AND Firebase user was created
			log.Printf("COMPENSATION: Attempting to delete Firebase user %s due to failed DB creation for email %s", createdFirebaseUID, email)
			delErr := s.firebaseAuth.DeleteUser(context.Background(), createdFirebaseUID) // Use background context for cleanup
			if delErr != nil {
				// Log compensation failure, but don't overwrite original error 'err'
				log.Printf("ERROR: Failed to execute compensation (delete Firebase user %s): %v", createdFirebaseUID, delErr)
			} else {
				log.Printf("COMPENSATION: Successfully deleted Firebase user %s", createdFirebaseUID)
			}
		}
	}()

	createdFirebaseUID, err = s.firebaseAuth.CreateUser(ctx, email, password, name)
	if err != nil {
		log.Printf("Failed to create user in Firebase for email %s: %v", email, err)
		return nil, domain.ErrFirebaseUserAlreadyExists
	}

	domainUser, err := domain.NewUser(name, email, createdFirebaseUID)
	if err != nil {
		log.Printf("Error creating domain user object for email %s: %v", email, err)
		return nil, err
	}

	err = s.userRepo.Save(ctx, domainUser)
	if err != nil {
		log.Printf("Failed to save user to database for email %s (Firebase UID: %s): %v", email, createdFirebaseUID, err)
		return nil, fmt.Errorf("failed to save user record: %w", domain.ErrDatabaseUserCreationFailed)
	}

	log.Printf("Successfully created user: DB_ID=%s, FirebaseUID=%s, Email=%s", domainUser.ID, domainUser.FirebaseUID, domainUser.Email)

	return domainUser, nil
}

// UpdateUserProfile updates the profile data for an existing user identified by Firebase UID.
// Only allows updating fields provided (non-nil pointers).
func (s *UserService) UpdateUserProfile(ctx context.Context, firebaseUID string, name *string /* add other updatable fields as pointers */) (*domain.User, error) {
	user, err := s.userRepo.FindByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.Printf("User profile not found for update: %s", firebaseUID)
		} else {
			log.Printf("Failed to get user by Firebase UID for update: %v", err)
		}
		return nil, err
	}

	updated := false
	if name != nil && *name != user.Name {
		user.Name = *name
		updated = true
		log.Printf("Updating user name for user: %s", user.ID)
	}

	if updated {
		user.UpdatedAt = time.Now().UTC()
		err = s.userRepo.Save(ctx, user)
		if err != nil {
			log.Printf("Failed to save updated user profile for userID %s: %v", user.ID, err)
			return nil, err
		}
		log.Printf("User profile updated successfully for userID: %s", user.ID)
	} else {
		log.Printf("No changes detected for user profile update, userID: %s", user.ID)
	}

	return user, nil
}

// GetUserByFirebaseUID retrieves a user profile by their Firebase UID.
func (s *UserService) GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	log.Printf("Getting user by Firebase UID: %s", firebaseUID)
	user, err := s.userRepo.FindByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.Printf("User profile not found by Firebase UID: %s", firebaseUID)
		} else {
			log.Printf("Failed to get user by Firebase UID %s: %v", firebaseUID, err)
		}
		return nil, err
	}
	log.Printf("User profile retrieved successfully, userID: %s", user.ID)
	return user, nil
}
