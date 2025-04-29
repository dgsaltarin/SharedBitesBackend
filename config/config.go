package config

import (
	"errors"
	"log"
	"os"
)

type AWSConfig struct {
	Region           string `env:"AWS_REGION"`
	AssumeRoleARN    string `env:"AWS_ASSUME_ROLE_ARN,omitempty"`
	TextTrackRoleARN string `env:"AWS_TEXTTRACK_ROLE_ARN,omitempty"` // Example if TextTrack needs specific role
	S3Bucket         string `env:"AWS_S3_BUCKET,omitempty"`
	SQSQueueURL      string `env:"AWS_SQS_QUEUE_URL,omitempty"`
	// ... other AWS settings
}

type FirebaseConfig struct {
	ServiceAccountKeyPath string `env:"FIREBASE_SERVICE_ACCOUNT_KEY_PATH,required"`
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
	// ... other server settings
}

type DatabaseConfig struct {
	DSN string `env:"DB_DSN,required"`
	// ... other db settings
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	AWS      AWSConfig
	Firebase FirebaseConfig // NEW
}

func Load() (*Config, error) {
	// ... existing loading logic ...
	// Ensure FirebaseConfig fields are populated
	// Example using github.com/kelseyhightower/envconfig:
	// var cfg Config
	// err := envconfig.Process("", &cfg)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to load config: %w", err)
	// }
	// return &cfg, nil

	// Placeholder for manual loading if not using a library
	cfg := Config{
		// ... populate other fields from env vars or defaults ...
		Firebase: FirebaseConfig{
			ServiceAccountKeyPath: os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY_PATH"),
		},
	}
	if cfg.Firebase.ServiceAccountKeyPath == "" {
		return nil, errors.New("FIREBASE_SERVICE_ACCOUNT_KEY_PATH must be set")
	}
	// ... validation ...
	return &cfg, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	return cfg
}
