package entity

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
	Email    string
	Password string
}

func NewUser(u User) (*User, error) {
	user := &User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
	err := user.GeneratePasswordHash()
	if err != nil {
		fmt.Println("Error generating password hash")
		return nil, err
	}

	user.ID = uuid.New().String()
	return user, nil
}

// GeneratePasswordHash generates a hash for the password
func (u *User) GeneratePasswordHash() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		fmt.Println("Error generating password hash")
		return err
	}

	u.Password = string(bytes)
	return nil
}

// CheckPassword checks if the password is correct
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
