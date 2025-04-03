package main

import (
	"context"
	"fmt"
	"github.com/ladderseeker/gin-crud-starter/internal/database"
	"github.com/ladderseeker/gin-crud-starter/internal/model"
	"github.com/ladderseeker/gin-crud-starter/pkg/logger"
	"os"
	"time"

	"github.com/ladderseeker/gin-crud-starter/config"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// TestUser represents a user for seeding the database
type TestUser struct {
	Name     string
	Email    string
	Password string
	Role     string
	Active   bool
}

func main() {
	// Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Initialize(config.Logging.Level)
	defer logger.GetLogger().Sync()

	// Connect to database
	database, err := database.NewPostgresDB(&config.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto migrate database schemas
	if err := autoMigrate(database); err != nil {
		logger.Fatal("Failed to migrate database schemas", zap.Error(err))
	}

	// Seed test data
	if err := seedTestData(database); err != nil {
		logger.Fatal("Failed to seed test data", zap.Error(err))
	}

	logger.Info("Test data seeded successfully")
}

// autoMigrate migrates database schemas
func autoMigrate(database *gorm.DB) error {
	// List of entities to migrate
	entities := []interface{}{
		&model.User{},
		// Add more entities here as needed
	}

	// Migrate entities
	for _, entity := range entities {
		if err := database.AutoMigrate(entity); err != nil {
			return err
		}
	}

	return nil
}

// seedTestData seeds the database with test data
func seedTestData(database *gorm.DB) error {
	// Define test users
	testUsers := []TestUser{
		{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: "password123",
			Role:     "admin",
			Active:   true,
		},
		{
			Name:     "Regular User",
			Email:    "user@example.com",
			Password: "password123",
			Role:     "user",
			Active:   true,
		},
		{
			Name:     "Inactive User",
			Email:    "inactive@example.com",
			Password: "password123",
			Role:     "user",
			Active:   false,
		},
		{
			Name:     "John Smith",
			Email:    "john.smith@example.com",
			Password: "password123",
			Role:     "user",
			Active:   true,
		},
		{
			Name:     "Jane Doe",
			Email:    "jane.doe@example.com",
			Password: "password123",
			Role:     "user",
			Active:   true,
		},
	}

	// Clear existing data
	if err := database.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
		return err
	}

	// Create users
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, u := range testUsers {
		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Create user
		user := &model.User{
			Name:     u.Name,
			Email:    u.Email,
			Password: string(hashedPassword),
			Role:     u.Role,
			Active:   u.Active,
		}

		if err := database.WithContext(ctx).Create(user).Error; err != nil {
			return err
		}

		logger.Info("Created test user", zap.String("name", u.Name), zap.String("email", u.Email))
	}

	return nil
}
