package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ladderseeker/gin-crud-starter/models"
	"github.com/ladderseeker/gin-crud-starter/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	// Initialize logger and set Gin to test mode once for all tests
	logger.Init()
	gin.SetMode(gin.TestMode)
}

// setupTestDB creates and returns a new isolated in-memory database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	// Create unique database name to ensure isolation
	dbName := fmt.Sprintf("file:memdb%d?mode=memory", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Migrate schema
	err = db.AutoMigrate(&models.Item{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// seedTestItems adds test data to the provided database
func seedTestItems(db *gorm.DB) []models.Item {
	// Create test items
	items := []models.Item{
		{Name: "Test Item 1", Description: "Description 1", Price: 10.99},
		{Name: "Test Item 2", Description: "Description 2", Price: 20.99},
	}

	for i := range items {
		db.Create(&items[i])
	}

	return items
}

// createTestContext creates a test context with the given HTTP method, path, and optional body
func createTestContext(w *httptest.ResponseRecorder, method, path string, body interface{}) (*gin.Context, *gin.Engine) {
	r := gin.New()
	c, _ := gin.CreateTestContext(w)

	var req *http.Request
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}

	c.Request = req
	return c, r
}

// cleanupDB handles cleanup of the database connection
func cleanupDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		t.Logf("Warning: couldn't get SQL DB for cleanup: %v", err)
		return
	}
	sqlDB.Close()
}

// TestGetItems tests the GetItems handler
func TestGetItems(t *testing.T) {
	// Setup: Create isolated test database
	db := setupTestDB(t)
	defer cleanupDB(t, db)

	// Seed test data
	seedTestItems(db)

	// Create test request and context
	w := httptest.NewRecorder()
	c, _ := createTestContext(w, "GET", "/items", nil)

	// Call the handler
	handler := GetItems(db)
	handler(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Item
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Test Item 1", response[0].Name)
}

// TestGetItemByID tests the GetItemByID handler
func TestGetItemByID(t *testing.T) {
	// Setup: Create isolated test database
	db := setupTestDB(t)
	defer cleanupDB(t, db)

	// Seed test data
	items := seedTestItems(db)

	// Create test request and context
	w := httptest.NewRecorder()
	c, _ := createTestContext(w, "GET", "/items/1", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: fmt.Sprintf("%d", items[0].ID)}}

	// Call the handler
	handler := GetItemByID(db)
	handler(c)

	// Assertions for success case
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Item
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, items[0].ID, response.ID)
	assert.Equal(t, "Test Item 1", response.Name)

	// Test not found case
	w = httptest.NewRecorder()
	c, _ = createTestContext(w, "GET", "/items/999", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "999"}}

	handler(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestCreateItem tests the CreateItem handler
func TestCreateItem(t *testing.T) {
	// Setup: Create isolated test database
	db := setupTestDB(t)
	defer cleanupDB(t, db)

	// Create a new item to test with
	newItem := models.Item{
		Name:        "New Test Item",
		Description: "New Test Description",
		Price:       15.99,
	}

	// Create test request and context
	w := httptest.NewRecorder()
	c, _ := createTestContext(w, "POST", "/items", newItem)

	// Call the handler
	handler := CreateItem(db)
	handler(c)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Item
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotZero(t, response.ID)
	assert.Equal(t, "New Test Item", response.Name)

	// Verify item is in database
	var dbItem models.Item
	result := db.First(&dbItem, response.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "New Test Item", dbItem.Name)
}

// TestUpdateItem tests the UpdateItem handler
func TestUpdateItem(t *testing.T) {
	// Setup: Create isolated test database
	db := setupTestDB(t)
	defer cleanupDB(t, db)

	// Seed test data
	items := seedTestItems(db)
	itemID := items[0].ID

	// Create updated item data
	updatedItem := models.Item{
		Name:        "Updated Item",
		Description: "Updated Description",
		Price:       25.99,
	}

	// Create test request and context
	w := httptest.NewRecorder()
	c, _ := createTestContext(w, "PUT", "/items/1", updatedItem)
	c.Params = gin.Params{gin.Param{Key: "id", Value: fmt.Sprintf("%d", itemID)}}

	// Call the handler
	handler := UpdateItem(db)
	handler(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify item was updated in database
	var dbItem models.Item
	result := db.First(&dbItem, itemID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Updated Item", dbItem.Name)
	assert.Equal(t, 25.99, dbItem.Price)
}

// TestDeleteItem tests the DeleteItem handler
func TestDeleteItem(t *testing.T) {
	// Setup: Create isolated test database
	db := setupTestDB(t)
	defer cleanupDB(t, db)

	// Seed test data
	items := seedTestItems(db)
	itemID := items[0].ID

	// Create test request and context
	w := httptest.NewRecorder()
	c, _ := createTestContext(w, "DELETE", "/items/1", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: fmt.Sprintf("%d", itemID)}}

	// Call the handler
	handler := DeleteItem(db)
	handler(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify item was deleted from database
	var dbItem models.Item
	result := db.First(&dbItem, itemID)
	assert.Error(t, result.Error) // Should get error because item no longer exists
}

// TestCreateItemInvalidJSON tests error handling for invalid JSON input
func TestCreateItemInvalidJSON(t *testing.T) {
	// Setup: Create isolated test database
	db := setupTestDB(t)
	defer cleanupDB(t, db)

	// Create test request with invalid JSON
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/items", bytes.NewBufferString(`{"name": "Invalid JSON`))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Call the handler
	handler := CreateItem(db)
	handler(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
