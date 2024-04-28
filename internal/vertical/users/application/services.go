package application

type UserService interface {
	SignUp(username, email, password string) error
	Login(username, password string) (string, error)
}

type HealthCheckService interface {
	Check() string
}
