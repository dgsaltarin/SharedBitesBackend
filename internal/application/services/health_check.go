package services

type HealthCheckService struct{}

func (h *HealthCheckService) Check() string {
	return "OK"
}
