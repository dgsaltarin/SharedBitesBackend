package handlers

import (
	services "github.com/dgsaltarin/SharedBitesBackend/internal/application/services"
)

type BillHandler struct {
	billService services.BillService
}
