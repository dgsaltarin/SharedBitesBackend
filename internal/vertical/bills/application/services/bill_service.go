package services

import (
	"encoding/json"
	"mime/multipart"

	"github.com/dgsaltarin/SharedBitesBackend/internal/infrastructure/helpers"
)

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
func (b *BillService) SplitBill(image *multipart.FileHeader) ([]byte, error) {
	imageData, err := helpers.DecodeImage(image)

	if err != nil {
		return nil, err
	}

	awsSession := helpers.AWSSession()

	helpers.UploadImages3(awsSession, imageData, image.Filename)
	session := helpers.TextTrackSesson(awsSession)
	result := helpers.Detectitems(session, image.Filename)

	output, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// GetBillByID gets a bill by its ID
func (b *BillService) DeleteBillByID() error {
	return nil
}
