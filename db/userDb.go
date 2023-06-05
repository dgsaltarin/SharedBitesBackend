package db

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/dgsaltarin/SharedBitesBackend/models"
)

const userstable = "sharedbites-users"

// DatabaseInterface interface for the data type of db field
type DatabaseInterface interface {
	Upsert(table string, model interface{}) error
	GetItem(table string, model interface{}) (*dynamodb.GetItemOutput, error)
}

// UserDB storage structure
type UsersDb struct {
	db DatabaseInterface
}

// GetUser is the method in charge of obtain the product by the sku
func (u *UsersDb) GetUser(id string) (*models.User, error) {
	var product *models.User
	response, err := u.db.GetItem(userstable, map[string]string{"id": id})
	if err != nil {
		return product, err
	}

	err = dynamodbattribute.UnmarshalMap(response.Item, &product)
	if err != nil {
		return product, err
	}

	return product, err
}

// UpserUser is the method in charge of insert or update one product
func (u *UsersDb) UpserUser(user *models.User) error {
	return u.db.Upsert(userstable, user)
}

// NewSiigoDB build repository structure and return it with the settings
func NewUserDB(
	db DatabaseInterface,
) *UsersDb {
	return &UsersDb{
		db: db,
	}
}
