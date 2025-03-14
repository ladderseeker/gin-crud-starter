package main

import (
	"github.com/ladderseeker/gin-crud-starter/configs"
	"github.com/ladderseeker/gin-crud-starter/internal/api"
	"github.com/ladderseeker/gin-crud-starter/internal/domain/entities"
	"github.com/ladderseeker/gin-crud-starter/internal/pkg/db"
	"github.com/ladderseeker/gin-crud-starter/internal/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	logger.Initialize(config.Logging.Level)
	defer logger.GetLogger().Sync()

	// Connect to database
	database, err := db.NewPostgresDB(&config.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto migrate database schemas
	if err := autoMigrate(database); err != nil {
		logger.Fatal("Failed to migrate database schemas", zap.Error(err))
	}

	// Create and start server
	server := api.NewServer(config, database)
	if err := server.Start(); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}

// autoMigrate migrates database schemas
func autoMigrate(database *gorm.DB) error {
	// List of entities to migrate
	entities := []interface{}{
		&entities.User{},
		// Add more entities here
	}

	// Migrate entities
	for _, entity := range entities {
		if err := database.AutoMigrate(entity); err != nil {
			return err
		}
	}

	return nil
}
