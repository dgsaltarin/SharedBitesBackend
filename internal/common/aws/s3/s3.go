package s3

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct{}

func NewS3() *S3 {
	return &S3{}
}

// UploadImageS3 upload an image into a s3 bucket
func (s *S3) UploadImages3(session *session.Session, data []byte, filename string) error {
	svc := s3.New(session)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return err
	}

	return nil
}
