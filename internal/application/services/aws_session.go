package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AWSSession creates a new session for aws
func AWSSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_S3_REGION),
	})

	if err != nil {
		fmt.Println(err)
	}

	return sess
}
