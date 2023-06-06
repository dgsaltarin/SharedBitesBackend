package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/dgsaltarin/SharedBitesBackend/models"
)

const userstable = "sharedbites-users"

// DatabaseInterface interface for the data type of db field
type DatabaseInterface interface {
	Upsert(table string, model interface{}) error
	GetItem(table string, model interface{}) (*dynamodb.GetItemOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

// UserDB storage structure
type UsersDb struct {
	db DatabaseInterface
}

// GetUser is the method in charge of obtain the product by the sku
func (u *UsersDb) GetUser(id string) (*models.User, error) {
	var user *models.User
	response, err := u.db.GetItem(userstable, map[string]string{"id": id})
	if err != nil {
		return user, err
	}

	err = dynamodbattribute.UnmarshalMap(response.Item, &user)
	if err != nil {
		return user, err
	}

	return user, err
}

// GetUserByUsername is the method in charge of obtain the user by the username
func (u *UsersDb) GetUserByUsername(username string) (*models.User, error) {
	var user *models.User
	// Define the input parameters for the Query operation
	params := &dynamodb.QueryInput{
		TableName: aws.String(userstable),       // Replace with your DynamoDB table name
		IndexName: aws.String("username-index"), // Replace with your secondary index name
		KeyConditions: map[string]*dynamodb.Condition{
			"username": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(username),
					},
				},
			},
		},
	}

	// Execute the Query operation
	result, err := u.db.Query(params)
	if err != nil {
		return nil, err
	}

	// Parse the query results
	var users []models.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	user = &users[0]

	return user, err
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
