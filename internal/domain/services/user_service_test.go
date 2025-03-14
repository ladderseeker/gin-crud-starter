package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ladderseeker/gin-crud-starter/internal/domain/entities"
	apperrors "github.com/ladderseeker/gin-crud-starter/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of repositories.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]entities.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetAllUsers(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Create sample users
	users := []entities.User{
		{ID: 1, Name: "User 1", Email: "user1@example.com", Role: "user"},
		{ID: 2, Name: "User 2", Email: "user2@example.com", Role: "admin"},
	}

	// Set expectations
	mockRepo.On("FindAll", mock.Anything).Return(users, nil)

	// Create service with mock repository
	service := NewUserService(mockRepo)

	// Call the service method
	result, err := service.GetAllUsers(context.Background())

	// Assert results
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "User 1", result[0].Name)
	assert.Equal(t, "user2@example.com", result[1].Email)

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	// Create sample user
	user := &entities.User{
		ID:     1,
		Name:   "John Doe",
		Email:  "john@example.com",
		Role:   "user",
		Active: true,
	}

	// Test cases
	testCases := []struct {
		name          string
		id            uint
		mockReturn    *entities.User
		mockError     error
		expectedError bool
	}{
		{
			name:          "Success",
			id:            1,
			mockReturn:    user,
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "UserNotFound",
			id:            2,
			mockReturn:    nil,
			mockError:     apperrors.NewResourceNotFoundError("User not found", nil, nil),
			expectedError: true,
		},
		{
			name:          "DatabaseError",
			id:            3,
			mockReturn:    nil,
			mockError:     errors.New("database error"),
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a fresh mock for each test case
			mockRepo := new(MockUserRepository)

			// Set expectations
			mockRepo.On("FindByID", mock.Anything, tc.id).Return(tc.mockReturn, tc.mockError)

			// Create service with mock repository
			service := NewUserService(mockRepo)

			// Call the service method
			result, err := service.GetUserByID(context.Background(), tc.id)

			// Assert results
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, user.Name, result.Name)
				assert.Equal(t, user.Email, result.Email)
			}

			// Verify expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateUser(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Create sample user input
	userInput := entities.UserCreate{
		Name:     "New User",
		Email:    "newuser@example.com",
		Password: "password123",
		Role:     "user",
	}

	// Capture the created user for validation
	var capturedUser *entities.User
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *entities.User) bool {
		capturedUser = u
		return u.Name == userInput.Name && u.Email == userInput.Email
	})).Return(nil)

	// Create service with mock repository
	service := NewUserService(mockRepo)

	// Call the service method
	result, err := service.CreateUser(context.Background(), userInput)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userInput.Name, result.Name)
	assert.Equal(t, userInput.Email, result.Email)

	// Assert password was hashed (not stored as plaintext)
	assert.NotEqual(t, userInput.Password, capturedUser.Password)

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Set expectations
	mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)
	mockRepo.On("Delete", mock.Anything, uint(2)).Return(apperrors.NewResourceNotFoundError("User not found", nil, nil))

	// Create service with mock repository
	service := NewUserService(mockRepo)

	// Test successful deletion
	err := service.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)

	// Test deletion of non-existent user
	err = service.DeleteUser(context.Background(), 2)
	assert.Error(t, err)
	assert.True(t, apperrors.IsNotFound(err))

	// Verify expectations
	mockRepo.AssertExpectations(t)
}
