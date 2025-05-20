package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	appconfig "github.com/dgsaltarin/SharedBitesBackend/config"
)

// LoadAWSConfig loads AWS configuration using the provided app config.
// It returns the AWS Config object that can be used to initialize AWS service clients.
func LoadAWSConfig(ctx context.Context, appCfg appconfig.AWSConfig) (aws.Config, error) {
	var cfg aws.Config
	var err error

	if appCfg.Region != "" {
		cfg, err = awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(appCfg.Region))
	} else {
		// If no region is specified in our app config, load default config.
		// This will attempt to find region from environment variables, shared config, etc.
		cfg, err = awsconfig.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// You can add custom credential providers here if needed
	// For example, if AccessKeyID and SecretAccessKey are set in appCfg and you want to use them explicitly.

	// By default, LoadDefaultConfig will try to find credentials in the standard chain:
	// 1. Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN)
	// 2. Shared credentials file (~/.aws/credentials)
	// 3. Shared configuration file (~/.aws/config)
	// 4. EC2 Instance Metadata Service (if running on an EC2 instance)

	return cfg, nil
}
