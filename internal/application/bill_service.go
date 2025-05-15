package application

import (
	"context"

	"github.com/dgsaltarin/SharedBitesBackend/platform/aws"
	"github.com/google/uuid"
)

type BillService struct {
	textractClient *aws.TextractClient
}

func NewBillService(textractClient *aws.TextractClient) *BillService {
	return &BillService{textractClient: textractClient}
}

func (s *BillService) AnalyzeBill(ctx context.Context, billID uuid.UUID) error {
	return nil
}
