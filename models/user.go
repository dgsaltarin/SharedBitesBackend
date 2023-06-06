package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// InvalidPassword returns true if the password is invalid
func (u *User) InvalidPassword(password string) bool {
	return u.Password != password
}
