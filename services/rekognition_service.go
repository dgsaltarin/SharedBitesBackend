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

	fmt.Println("Rekognition session created")

	return svc
}

func DetectLabels(svc *rekognition.Rekognition, decodedImage []byte) *rekognition.DetectTextOutput {
	input := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			Bytes: decodedImage,
		},
	}

	result, err := svc.DetectText(input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Labels detected")

	return result
}
