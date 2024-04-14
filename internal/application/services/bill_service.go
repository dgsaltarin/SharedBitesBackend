package services

type BillService struct{}

// NewBillService creates a new instance of BillService
func NewBillService() BillService {
	return BillService{}
}

// CreateBill creates a new bill
func (b *BillService) GetBillByID() error {
	return nil
}

// SplitBill splits a bill
func (b *BillService) SplitBill() error {
	return nil
}

// GetBillByID gets a bill by its ID
func (b *BillService) DeleteBillByID() error {
	return nil
}
