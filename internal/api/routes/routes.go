package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ladderseeker/gin-crud-starter/internal/api/controllers"
	"github.com/ladderseeker/gin-crud-starter/internal/api/middleware"
	"github.com/ladderseeker/gin-crud-starter/internal/domain/repositories"
	"github.com/ladderseeker/gin-crud-starter/internal/domain/services"
	"gorm.io/gorm"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Setup middleware
	middleware.SetupMiddleware(router)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Initialize repositories
		userRepo := repositories.NewUserRepository(db)

		// Initialize services
		userService := services.NewUserService(userRepo)

		// Initialize controllers
		userController := controllers.NewUserController(userService)

		// Register controller routes
		userController.Register(api)
	}

	// Handle 404 Not Found
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    "NOT_FOUND",
			"message": "The requested resource was not found",
		})
	})
}
