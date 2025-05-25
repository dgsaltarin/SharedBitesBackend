package database

import (
	"fmt"
	"log"
	"time"

	"github.com/dgsaltarin/SharedBitesBackend/config"
	"github.com/dgsaltarin/SharedBitesBackend/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MustConnectGORM establishes a GORM connection to PostgreSQL and runs migrations. Panics on error.
func MustConnectGORM(cfg config.DatabaseConfig) *gorm.DB {
	db, err := ConnectGORM(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database with GORM: %v", err)
	}
	return db
}

// ConnectGORM establishes a GORM connection to PostgreSQL and runs migrations.
func ConnectGORM(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// Consider adding more robust DSN parsing if needed
	dsn := cfg.DSN
	if dsn == "" {
		return nil, fmt.Errorf("database DSN is required")
	}

	// Configure GORM logger (adjust level as needed)
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second * 2, // Slow SQL threshold
			LogLevel:                  logger.Warn,     // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: false,           // Do not ignore ErrRecordNotFound errors
			Colorful:                  true,            // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		// PrepareStmt: true, // Cache prepared statements for performance
	})

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Optional: Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully with GORM.")

	// --- Auto Migration ---
	// IMPORTANT: In production, use a dedicated migration tool (e.g., golang-migrate, Atlas).
	// AutoMigrate is suitable for development/simple cases.
	log.Println("Running GORM AutoMigrate...")
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Bill{},
		&domain.LineItem{},
		//&domain.Group{},   // Add other domain models you need tables for
		//&domain.Expense{}, // Add other domain models you need tables for
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run GORM auto-migration: %w", err)
	}
	log.Println("GORM AutoMigrate completed.")

	return db, nil
}
