package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

func RekognitionSession() *rekognition.Rekognition {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		fmt.Println(err)
	}

	svc := rekognition.New(sess)

	return svc
}

func DetectLabels(svc *rekognition.Rekognition, decodedImage []byte) *rekognition.DetectLabelsOutput {
	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			Bytes: decodedImage,
		},
	}

	result, err := svc.DetectLabels(input)
	if err != nil {
		fmt.Println(err)
	}

	return result
}
