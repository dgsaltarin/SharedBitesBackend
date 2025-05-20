package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	appconfig "github.com/dgsaltarin/SharedBitesBackend/config"
	"github.com/dgsaltarin/SharedBitesBackend/internal/ports"
)

type s3FileStore struct {
	client        *s3.Client
	uploader      *manager.Uploader
	presignClient *s3.PresignClient
	bucketName    string
	awsRegion     string
}

// NewS3FileStore creates a new S3 file store.
func NewS3FileStore(ctx context.Context, appCfg appconfig.AWSConfig) (ports.FileStore, error) {
	var cfg aws.Config
	var err error

	if appCfg.Region == "" {
		return nil, fmt.Errorf("AWS region must be specified in config for S3FileStore")
	}
	if appCfg.S3Bucket == "" {
		return nil, fmt.Errorf("S3 bucket name must be specified in AWS config")
	}

	cfg, err = awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(appCfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for S3: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client)
	presignClient := s3.NewPresignClient(s3Client)

	return &s3FileStore{
		client:        s3Client,
		uploader:      uploader,
		presignClient: presignClient,
		bucketName:    appCfg.S3Bucket,
		awsRegion:     appCfg.Region,
	}, nil
}

// UploadFile uploads a file to S3 and returns its storage path (key).
func (s *s3FileStore) UploadFile(ctx context.Context, file io.Reader, destinationPath, contentType string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}
	if s.bucketName == "" {
		return "", fmt.Errorf("S3 bucket name not configured")
	}

	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(destinationPath),
		Body:        file,
		ContentType: aws.String(contentType),
	}

	// Using the manager.Uploader for potentially large files and multipart uploads
	_, err := s.uploader.Upload(ctx, uploadInput)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3 (bucket: %s, key: %s): %w", s.bucketName, destinationPath, err)
	}

	return destinationPath, nil // The destinationPath is the storage path (key)
}

// GetFileURL generates a pre-signed URL for accessing an S3 object.
// This URL will be valid for a short duration (e.g., 15 minutes).
func (s *s3FileStore) GetFileURL(ctx context.Context, storagePath string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}
	if s.bucketName == "" {
		return "", fmt.Errorf("S3 bucket name not configured")
	}

	presignDuration := func(po *s3.PresignOptions) {
		po.Expires = 15 * time.Minute // URL expires in 15 minutes
	}

	presignedURL, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(storagePath),
	}, presignDuration)

	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL for S3 object (bucket: %s, key: %s): %w", s.bucketName, storagePath, err)
	}

	return presignedURL.URL, nil
}

// DeleteFile deletes a file from S3.
func (s *s3FileStore) DeleteFile(ctx context.Context, storagePath string) error {
	if s.client == nil {
		return fmt.Errorf("S3 client not initialized")
	}
	if s.bucketName == "" {
		return fmt.Errorf("S3 bucket name not configured")
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(storagePath),
	})

	if err != nil {
		// Check if the error is because the object doesn't exist, which might not be an error for a delete operation.
		var nsk *types.NoSuchKey
		if strings.Contains(err.Error(), "NoSuchKey") || errors.As(err, &nsk) { // More robust check
			// Consider logging this but not returning an error, as the file is already gone.
			// For now, let's return nil to indicate the desired state (file doesn't exist) is achieved.
			return nil
		}
		return fmt.Errorf("failed to delete file from S3 (bucket: %s, key: %s): %w", s.bucketName, storagePath, err)
	}

	return nil
}

func (s *s3FileStore) constructS3URL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.awsRegion, key)
}
