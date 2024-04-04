package entity

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
	Email    string
	Password string
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

type Item struct {
	Name  string
	Price float64
}

type Bill struct {
	ID          string
	Items       []Item
	Total       float64
	People      int
	SplitEqualy bool
	UserID      string
	Date        time.Time
}
