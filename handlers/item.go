package handlers

import (
	"github.com/ladderseeker/gin-crud-starter/models"
	"github.com/ladderseeker/gin-crud-starter/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetItems retrieves all items clearly from DB
func GetItems(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []models.Item
		if err := db.Find(&items).Error; err != nil {
			logger.Get().Error("Failed to retrieve items", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

// GetItemByID retrieves an item clearly by ID from DB
func GetItemByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			logger.Get().Error("Invalid ID provided", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var item models.Item
		if err := db.First(&item, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				logger.Get().Warn("Item not found", zap.Int("id", id))
				c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			} else {
				logger.Get().Error("Database error", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

// CreateItem inserts a new item clearly into DB
func CreateItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newItem models.Item
		if err := c.ShouldBindJSON(&newItem); err != nil {
			logger.Get().Error("JSON binding error", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&newItem).Error; err != nil {
			logger.Get().Error("Failed to create item", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusCreated, newItem)
	}
}

// UpdateItem modifies an existing item clearly in DB
func UpdateItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			logger.Get().Error("Invalid ID", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var item models.Item
		if err := c.ShouldBindJSON(&item); err != nil {
			logger.Get().Error("JSON binding error", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item.ID = uint(id)
		if err := db.Model(&item).Updates(item).Error; err != nil {
			logger.Get().Error("Failed to update item", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

// DeleteItem removes an item clearly from DB
func DeleteItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			logger.Get().Error("Invalid ID", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		if err := db.Delete(&models.Item{}, id).Error; err != nil {
			logger.Get().Error("Failed to delete item", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deleted": id})
	}
}
