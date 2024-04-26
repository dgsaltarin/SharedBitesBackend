package application

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/dgsaltarin/SharedBitesBackend/internal/domain/entity"
)

type UserService interface {
	SignUp(username, email, password string) error
	Login(username, password string) (string, error)
}

type BillService interface {
	SplitBill(bill entity.Bill) (entity.Bill, error)
	GetBillByID(id string) (entity.Bill, error)
	DeleteBillByID(id string) error
}

type HealthCheckService interface {
	Check() string
}

type TextTrackService interface {
	CreateSesion(aws_session *session.Session)
	DelectItmems(session *textract.Textract, path string) error
	ExtractExpensesFromResults(itemsGroup []*textract.LineItemGroup) []entity.Item
}
