package services

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	AWS_S3_REGION = "us-east-1"
	AWS_S3_BUCKET = "test-dgsaltarin"
)

// UploadImageS3 upload an image into a s3 bucket
func UploadImages3(session *session.Session, data []byte, filename string) error {
	svc := s3.New(session)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return err
	}

	return nil
}
