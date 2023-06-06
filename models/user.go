package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
