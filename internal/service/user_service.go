package service

import (
	"context"
	"github.com/ladderseeker/gin-crud-starter/internal/model"
	"github.com/ladderseeker/gin-crud-starter/pkg/errors"
	"github.com/ladderseeker/gin-crud-starter/pkg/logger"
	"time"

	"github.com/ladderseeker/gin-crud-starter/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user service
type UserService interface {
	GetAllUsers(ctx context.Context) ([]model.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*model.UserResponse, error)
	CreateUser(ctx context.Context, input model.UserCreate) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, id uint, input model.UserUpdate) (*model.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
}

// userServiceImpl implements the UserService interface
type userServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

// GetAllUsers retrieves all users
func (s *userServiceImpl) GetAllUsers(ctx context.Context) ([]model.UserResponse, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		logger.Error("Failed to get all users", zap.Error(err))
		return nil, err
	}

	// Convert users to response format
	var response []model.UserResponse
	for _, user := range users {
		response = append(response, user.ToResponse())
	}

	return response, nil
}

// GetUserByID retrieves a user by ID
func (s *userServiceImpl) GetUserByID(ctx context.Context, id uint) (*model.UserResponse, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get user by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// CreateUser creates a new user
func (s *userServiceImpl) CreateUser(ctx context.Context, input model.UserCreate) (*model.UserResponse, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return nil, errors.NewInternalError("Failed to process password", err)
	}

	// Create user entity
	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
		Active:   true,
	}

	// Default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	// Create user
	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error("Failed to create user", zap.String("email", input.Email), zap.Error(err))
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateUser updates a user
func (s *userServiceImpl) UpdateUser(ctx context.Context, id uint, input model.UserUpdate) (*model.UserResponse, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Retrieve user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("Failed to retrieve user for update", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	// Update user fields if provided
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("Failed to hash password during update", zap.Error(err))
			return nil, errors.NewInternalError("Failed to process password", err)
		}
		user.Password = string(hashedPassword)
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.Active != nil {
		user.Active = *input.Active
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Error("Failed to update user", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// DeleteUser deletes a user
func (s *userServiceImpl) DeleteUser(ctx context.Context, id uint) error {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Delete user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete user", zap.Uint("id", id), zap.Error(err))
		return err
	}

	return nil
}
