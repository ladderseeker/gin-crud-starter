package config

import (
	"github.com/ladderseeker/gin-crud-starter/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDB - Clearly sets up and returns GORM database connection
func ConnectDB() (*gorm.DB, error) {
	dsn := "root:password@tcp(localhost:3306)/gin_crud_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Get().Error("Database connection failed: %v", zap.Error(err))
		return nil, err
	}
	return db, nil
}
