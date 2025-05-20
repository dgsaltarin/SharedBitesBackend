package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/textract"
	appconfig "github.com/dgsaltarin/SharedBitesBackend/config"
)

// TextractClient is a wrapper around the AWS Textract service client.
type TextractClient struct {
	client *textract.Client
}

// NewTextractClient creates a new TextractClient.
// It loads AWS configuration using the region from appConfig.
// If appConfig.Region is empty, it attempts to load the default AWS configuration.
func NewTextractClient(ctx context.Context, appCfg appconfig.AWSConfig) (*TextractClient, error) {
	cfg, err := LoadAWSConfig(ctx, appCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for Textract: %w", err)
	}

	svc := textract.NewFromConfig(cfg)
	return &TextractClient{client: svc}, nil
}
