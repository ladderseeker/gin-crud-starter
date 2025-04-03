package main

import (
	"github.com/ladderseeker/gin-crud-starter/config"
	"github.com/ladderseeker/gin-crud-starter/internal/database"
	"github.com/ladderseeker/gin-crud-starter/internal/model"
	"github.com/ladderseeker/gin-crud-starter/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	conf, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	logger.Initialize(conf.Logging.Level)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Failed to sync logger", zap.Error(err))
		}
	}(logger.GetLogger())

	// Connect to database
	db, err := database.NewPostgresDB(&conf.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto migrate database schemas
	if err := autoMigrate(db); err != nil {
		logger.Fatal("Failed to migrate database schemas", zap.Error(err))
	}

	// Create and start server
	server := NewServer(conf, db)
	if err := server.Start(); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}

// autoMigrate migrates database schemas
func autoMigrate(db *gorm.DB) error {
	// List of entities to migrate
	entities := []interface{}{
		&model.User{},
		// Add more entities here
	}

	// Migrate entities
	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			return err
		}
	}

	return nil
}
