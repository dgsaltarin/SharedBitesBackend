package application

import "github.com/dgsaltarin/SharedBitesBackend/internal/domain/entity"

type UserService interface {
	SignUp(username, email, password string) error
	Login(username, password string) (string, error)
}

type BillService interface {
	SplitBill(bill entity.Bill) (entity.Bill, error)
	GetBillByID(id string) (entity.Bill, error)
	DeleteBillByID(id string) error
}

type S3Service interface {
	UploadImage(image []byte) (string, error)
}
