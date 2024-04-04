package db

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/dgsaltarin/SharedBitesBackend/services"
)

var dynamodbInstance *Dynamodb
var lock sync.Mutex

// DynamodbInterface interface for the data type of Db field
type DynamodbInterface interface {
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

// Dynamodb structure for dynamodb utils
type Dynamodb struct {
	Db DynamodbInterface
}

// Connect is the method in charge of connect to dynamodb
func Connect() *Dynamodb {
	sess := services.AWSSession()

	svc := dynamodb.New(sess)
	Db := dynamodbiface.DynamoDBAPI(svc)

	return &Dynamodb{
		Db: Db,
	}
}

func GetDynamoDBInstance() *Dynamodb {
	// lock instance for different go routines
	lock.Lock()
	defer lock.Unlock()

	if dynamodbInstance == nil {
		dynamodbInstance = &Dynamodb{}
		dynamodbInstance = Connect()
	} else {
		fmt.Println("Dynamodb already connected")
	}
	return dynamodbInstance
}

// Upsert is the method in charge of create or update objects
func (m Dynamodb) Upsert(table string, model interface{}) error {
	item, err := dynamodbattribute.MarshalMap(model)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(table),
	}

	_, err = m.Db.PutItem(params)
	if err != nil {
		return err
	}

	return nil
}

// GetItem is the method in charge of search one object by its partition key
func (m Dynamodb) GetItem(table string, model interface{}) (*dynamodb.GetItemOutput, error) {
	paramKey, err := dynamodbattribute.MarshalMap(model)
	if err != nil {
		return nil, err
	}

	params := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key:       paramKey,
	}

	return m.Db.GetItem(params)
}

// Query is the method in charge of search objects by its partition key
func (m Dynamodb) Query(params *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return m.Db.Query(params)
}
