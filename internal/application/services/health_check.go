package services

type HealthCheckService struct{}

func NewHealthCheckService() HealthCheckService {
	return HealthCheckService{}
}

func (h *HealthCheckService) Check() string {
	return "OK"
}
