package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/dgsaltarin/SharedBitesBackend/internal/vertical/users/infrastructure/repository/dynamodb/models"
)

var USER_TABLE string = os.Getenv("USER_TABLE")

type dynamodbUserRepository struct {
	dynamodb *dynamodb.DynamoDB
}

func NewDynamoDBUserRepository() *dynamodbUserRepository {
	return &dynamodbUserRepository{
		dynamodb: dynamodb.New(session.New()),
	}
}

// GetUser is the method in charge of obtain the product by the sku
func (u *dynamodbUserRepository) GetUser(id string) (*models.User, error) {
	// Define the input parameters for the GetItem operation
	input := &dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE), // Replace with your table name
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	// Execute the GetItem operation
	result, err := u.dynamodb.GetItem(input)
	if err != nil {
		return nil, err
	}

	// Parse the query result
	user := models.User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername is the method in charge of obtain the user by the username
func (u *dynamodbUserRepository) GetUserByUsername(username string) (*models.User, error) {
	// Define the input parameters for the GetItem operation
	input := &dynamodb.ScanInput{
		TableName: aws.String(USER_TABLE), // Replace with your table name
	}

	// Execute the GetItem operation
	result, err := u.dynamodb.Scan(input)
	if err != nil {
		return nil, err
	}

	// Parse the query result
	user := models.User{}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpserUser is the method in charge of insert or update one product
func (u *dynamodbUserRepository) UpsertUser(user *models.User) error {
	// Define the input parameters for the PutItem operation
	input := &dynamodb.PutItemInput{
		TableName: aws.String(USER_TABLE), // Replace with your table name
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(user.ID),
			},
			"username": {
				S: aws.String(user.Username),
			},
			"email": {
				S: aws.String(user.Email),
			},
			"password": {
				S: aws.String(user.Password),
			},
			"created_at": {
				S: aws.String(user.CreatedAt),
			},
			"updated_at": {
				S: aws.String(user.UpdatedAt),
			},
		},
	}

	// Execute the PutItem operation
	_, err := u.dynamodb.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
