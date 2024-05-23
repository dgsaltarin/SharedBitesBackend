package response

type SignUpResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
