package session

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AWSSession struct{}

func NewAWSSession() *AWSSession {
	return &AWSSession{}
}

func (as *AWSSession) CreateSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		fmt.Println(err)
	}

	return sess
}
