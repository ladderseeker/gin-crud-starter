package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/ladderseeker/gin-crud-starter/handlers"
	"github.com/ladderseeker/gin-crud-starter/models"
	"github.com/ladderseeker/gin-crud-starter/pkg/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// IntegrationTestSuite is our test suite structure
type IntegrationTestSuite struct {
	suite.Suite
	DB     *gorm.DB
	Router *gin.Engine
	// Test data
	testItems []models.Item
}

// SetupSuite runs once before all tests
func (suite *IntegrationTestSuite) SetupSuite() {
	// Initialize logger
	logger.Init()

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test database (SQLite in-memory)
	db, err := gorm.Open(sqlite.Dialector{DriverName: "sqlite", DSN: ":memory:"}, &gorm.Config{})
	if err != nil {
		suite.T().Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Item{})
	if err != nil {
		suite.T().Fatal("Failed to migrate database:", err)
	}

	suite.DB = db

	// Setup the router
	router := gin.Default()
	suite.setupRoutes(router)
	suite.Router = router

	// Add test data
	suite.seedTestData()
}

// TearDownSuite runs after all tests are complete
func (suite *IntegrationTestSuite) TearDownSuite() {
	// Close database connection
	sqlDB, err := suite.DB.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// SetupTest runs before each test
func (suite *IntegrationTestSuite) SetupTest() {
	// Clean up the database before each test
	suite.DB.Exec("DELETE FROM items")
	// Re-seed test data
	suite.seedTestData()
}

// seedTestData adds initial test data to the database
func (suite *IntegrationTestSuite) seedTestData() {
	// Create test items
	suite.testItems = []models.Item{
		{Name: "Test Item 1", Description: "Description 1", Price: 10.99},
		{Name: "Test Item 2", Description: "Description 2", Price: 20.99},
		{Name: "Test Item 3", Description: "Description 3", Price: 30.99},
	}

	// Insert into database
	for i := range suite.testItems {
		suite.DB.Create(&suite.testItems[i])
	}
}

// setupRoutes configures the API routes for testing
func (suite *IntegrationTestSuite) setupRoutes(router *gin.Engine) {
	router.GET("/items", handlers.GetItems(suite.DB))
	router.GET("/items/:id", handlers.GetItemByID(suite.DB))
	router.POST("/items", handlers.CreateItem(suite.DB))
	router.PUT("/items/:id", handlers.UpdateItem(suite.DB))
	router.DELETE("/items/:id", handlers.DeleteItem(suite.DB))
}

// Helper method to perform requests and get responses
func (suite *IntegrationTestSuite) performRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()

	var req *http.Request
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}

	suite.Router.ServeHTTP(w, req)
	return w
}

// Rest of tests remain the same as before
// TestGetItems tests retrieving all items
func (suite *IntegrationTestSuite) TestGetItems() {
	// Perform GET request
	w := suite.performRequest("GET", "/items", nil)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Parse response body
	var items []models.Item
	err := json.Unmarshal(w.Body.Bytes(), &items)
	assert.NoError(suite.T(), err)

	// We should have 3 items from the seed data
	assert.Len(suite.T(), items, 3)
	assert.Equal(suite.T(), "Test Item 1", items[0].Name)
}

// TestGetItemByID tests retrieving a single item
func (suite *IntegrationTestSuite) TestGetItemByID() {
	// Get the first test item's ID
	firstItemID := suite.testItems[0].ID

	// Perform GET request
	w := suite.performRequest("GET", "/items/"+strconv.Itoa(int(firstItemID)), nil)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Parse response body
	var item models.Item
	err := json.Unmarshal(w.Body.Bytes(), &item)
	assert.NoError(suite.T(), err)

	// Verify the correct item was returned
	assert.Equal(suite.T(), "Test Item 1", item.Name)
	assert.Equal(suite.T(), 10.99, item.Price)
}

// TestGetItemByIDNotFound tests retrieving a non-existent item
func (suite *IntegrationTestSuite) TestGetItemByIDNotFound() {
	// Perform GET request for non-existent ID
	w := suite.performRequest("GET", "/items/999", nil)

	// Assert response
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Parse response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	// Verify error message
	assert.Equal(suite.T(), "Item not found", response["error"])
}

// TestCreateItem tests creating a new item
func (suite *IntegrationTestSuite) TestCreateItem() {
	// Create a new item
	newItem := models.Item{
		Name:        "New Item",
		Description: "New Description",
		Price:       40.99,
	}

	// Perform POST request
	w := suite.performRequest("POST", "/items", newItem)

	// Assert response
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// Parse response body
	var createdItem models.Item
	err := json.Unmarshal(w.Body.Bytes(), &createdItem)
	assert.NoError(suite.T(), err)

	// Verify item was created
	assert.NotZero(suite.T(), createdItem.ID)
	assert.Equal(suite.T(), "New Item", createdItem.Name)

	// Verify item exists in database
	var dbItem models.Item
	result := suite.DB.First(&dbItem, createdItem.ID)
	assert.NoError(suite.T(), result.Error)
	assert.Equal(suite.T(), "New Item", dbItem.Name)
}

// TestUpdateItem tests updating an existing item
func (suite *IntegrationTestSuite) TestUpdateItem() {
	// Get the second test item's ID
	itemID := suite.testItems[1].ID

	// Update data
	updatedItem := models.Item{
		Name:        "Updated Item",
		Description: "Updated Description",
		Price:       99.99,
	}

	// Perform PUT request
	w := suite.performRequest("PUT", "/items/"+strconv.Itoa(int(itemID)), updatedItem)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Verify item was updated in database
	var dbItem models.Item
	result := suite.DB.First(&dbItem, itemID)
	assert.NoError(suite.T(), result.Error)
	assert.Equal(suite.T(), "Updated Item", dbItem.Name)
	assert.Equal(suite.T(), 99.99, dbItem.Price)
}

// TestDeleteItem tests deleting an item
func (suite *IntegrationTestSuite) TestDeleteItem() {
	// Get the third test item's ID
	itemID := suite.testItems[2].ID

	// Perform DELETE request
	w := suite.performRequest("DELETE", "/items/"+strconv.Itoa(int(itemID)), nil)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Verify item was deleted from database
	var dbItem models.Item
	result := suite.DB.First(&dbItem, itemID)
	assert.Error(suite.T(), result.Error) // Should get error because item no longer exists
}

// TestCreateItemInvalidJSON tests error handling for invalid JSON
func (suite *IntegrationTestSuite) TestCreateItemInvalidJSON() {
	// Perform POST request with invalid JSON
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer([]byte(`{"name": "Invalid`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// TestMain to run the test suite
func TestIntegrationTestSuite(t *testing.T) {
	// Skip integration tests if environment flag is set
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(IntegrationTestSuite))
}
