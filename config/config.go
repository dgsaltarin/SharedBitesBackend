package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AWSConfig struct {
	Region          string `envconfig:"AWS_REGION" required:"true"`
	AccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY"`
	S3Bucket        string `envconfig:"AWS_S3_BUCKET"`
}

type FirebaseConfig struct {
	ServiceAccountKeyPath string `envconfig:"FIREBASE_SERVICE_ACCOUNT_KEY_PATH" required:"true"`
}

type ServerConfig struct {
	Port     string `envconfig:"SERVER_PORT" default:"8080"`
	GinMode  string `envconfig:"GIN_MODE" default:"release"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

type DatabaseConfig struct {
	DSN string `envconfig:"DB_DSN" required:"true"`
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	AWS      AWSConfig
	Firebase FirebaseConfig
}

func Load(logger *slog.Logger) (*Config, error) {
	err := godotenv.Load() // Loads .env from current directory by default
	if err != nil {
		if !os.IsNotExist(err) { // Log only if it's an actual error loading an existing .env
			logger.Warn("Error loading .env file", "error", err)
		} else {
			logger.Info(".env file not found, proceeding with environment variables")
		}
	} else {
		logger.Info(".env file loaded successfully")
	}

	var cfg Config

	err = envconfig.Process("", &cfg)
	if err != nil {
		logger.Error("Failed to process configuration from environment variables", "error", err)
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}

	// envconfig handles 'required' and 'default' tags.
	// Additional custom validation can be done here if needed.
	// For example, check if ServiceAccountKeyPath file exists:
	// if _, err := os.Stat(cfg.Firebase.ServiceAccountKeyPath); os.IsNotExist(err) {
	// 	logger.Error("Firebase service account key file not found", "path", cfg.Firebase.ServiceAccountKeyPath)
	// 	return nil, fmt.Errorf("firebase service account key file not found at %s: %w", cfg.Firebase.ServiceAccountKeyPath, err)
	// }

	logger.Info("Configuration loaded successfully from environment variables")
	return &cfg, nil
}

func MustLoad(logger *slog.Logger) *Config {
	cfg, err := Load(logger)
	if err != nil {
		log.Fatalf("CRITICAL: failed to load configuration: %v", err)
	}
	return cfg
}
