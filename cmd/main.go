package main

import (
	"github.com/ladderseeker/gin-crud-starter/config"
	"github.com/ladderseeker/gin-crud-starter/models"
	"github.com/ladderseeker/gin-crud-starter/pkg/logger"
	"github.com/ladderseeker/gin-crud-starter/routers"
	"go.uber.org/zap"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		logger.Get().Fatal("Could not connect to DB: %v", zap.Error(err))
	}

	// Auto-migrate Item schema (creates tables automatically clearly)
	if err := db.AutoMigrate(&models.Item{}); err != nil {
		logger.Get().Fatal("Database migration failed: %v", zap.Error(err))
	}

	router := routers.SetupRouter(db)
	router.Run(":8080")
}
