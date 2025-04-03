package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ladderseeker/gin-crud-starter/internal/controller/v1"
	"github.com/ladderseeker/gin-crud-starter/internal/middleware"
	"github.com/ladderseeker/gin-crud-starter/internal/repository"
	"github.com/ladderseeker/gin-crud-starter/internal/service"
	"gorm.io/gorm"
)

// SetupRoutes configures all the router for the application
func SetupRoutes(router *gin.Engine, db *gorm.DB) {

	// Initialize user related instance
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := v1.NewUserController(userService)

	// Setup middleware
	middleware.SetupMiddleware(router)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API router
	api := router.Group("/api/v1")
	{
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
