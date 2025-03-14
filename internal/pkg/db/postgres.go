package db

import (
	"github.com/ladderseeker/gin-crud-starter/internal/pkg/logger"
	"time"

	"github.com/ladderseeker/gin-crud-starter/configs"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// NewPostgresDB establishes a connection to the database
func NewPostgresDB(config *configs.DatabaseConfig) (*gorm.DB, error) {
	// Configure GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		PrepareStmt: true,
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(config.GetDSN()), gormConfig)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Check connection
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	logger.Info("Connected to database",
		zap.String("host", config.Host),
		zap.String("database", config.DBName))

	return db, nil
}

// CloseDatabaseConnection closes the database connection
func CloseDatabaseConnection(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			logger.Error("Error getting SQL DB instance", zap.Error(err))
			return
		}

		if err := sqlDB.Close(); err != nil {
			logger.Error("Error closing database connection", zap.Error(err))
			return
		}

		logger.Info("Database connection closed")
	}
}
