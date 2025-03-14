package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ladderseeker/gin-crud-starter/handlers"
)

// SetupRouter initializes all your API routes
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	itemGroup := router.Group("/items")
	{
		itemGroup.GET("/", handlers.GetItems(db))
		itemGroup.GET("/:id", handlers.GetItemByID(db))
		itemGroup.POST("/", handlers.CreateItem(db))
		itemGroup.PUT("/:id", handlers.UpdateItem(db))
		itemGroup.DELETE("/:id", handlers.DeleteItem(db))
	}

	return router
}
