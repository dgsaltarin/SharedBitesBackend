package domain

import (
	"github.com/google/uuid" // Use UUIDs for internal user IDs
	"time"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"` // Internal DB ID
	Name        string    `gorm:"size:255;not null"`
	Email       string    `gorm:"size:255;not null;uniqueIndex"`
	FirebaseUID string    `gorm:"size:255;not null;uniqueIndex"`
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time
}

func NewUser(name string, email string, firebaseUID string) (*User, error) {
	if name == "" {
		return nil, ErrUserNameEmpty
	}
	if email == "" {
		return nil, ErrUserEmailEmpty
	}

	if firebaseUID == "" {
		return nil, ErrFirebaseIDEmpty
	}

	now := time.Now().UTC()
	return &User{
		Name:        name,
		Email:       email,
		FirebaseUID: firebaseUID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
