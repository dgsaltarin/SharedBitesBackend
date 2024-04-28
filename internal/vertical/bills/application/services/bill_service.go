package services

import (
	"encoding/json"
	"mime/multipart"

	awssession "github.com/dgsaltarin/SharedBitesBackend/internal/common/aws/session"
	decoder "github.com/dgsaltarin/SharedBitesBackend/internal/common/decoder"
	s3 "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/bills/application/providers/s3"
	textTrack "github.com/dgsaltarin/SharedBitesBackend/internal/vertical/bills/application/providers/text_track"
)

type BillService struct {
	decoder    *decoder.Decoder
	awsSession *awssession.AWSSession
	s3         *s3.S3
	textTrack  *textTrack.TextTrackService
}

// NewBillService creates a new instance of BillService
func NewBillService(decoder *decoder.Decoder, awssession *awssession.AWSSession, s3 *s3.S3, textTrack *textTrack.TextTrackService) BillService {
	return BillService{
		decoder:    decoder,
		awsSession: awssession,
		s3:         s3,
		textTrack:  textTrack,
	}
}

// CreateBill creates a new bill
func (b *BillService) GetBillByID() error {
	return nil
}

// SplitBill splits a bill
func (b *BillService) SplitBill(image *multipart.FileHeader) ([]byte, error) {
	imageData, err := b.decoder.DecodeImage(image)

	if err != nil {
		return nil, err
	}

	awsSession := b.awsSession.CreateSession()

	b.s3.UploadImages3(awsSession, imageData, image.Filename)
	session := b.textTrack.TextTrackSesson(awsSession)
	result := b.textTrack.Detectitems(session, image.Filename)

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
