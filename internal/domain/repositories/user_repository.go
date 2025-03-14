package repositories

import (
	"context"
	"github.com/ladderseeker/gin-crud-starter/internal/pkg/errors"

	"github.com/ladderseeker/gin-crud-starter/internal/domain/entities"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user repository
type UserRepository interface {
	FindAll(ctx context.Context) ([]entities.User, error)
	FindByID(ctx context.Context, id uint) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uint) error
}

// userRepository implements the UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// FindAll retrieves all users
func (r *userRepository) FindAll(ctx context.Context) ([]entities.User, error) {
	var users []entities.User
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, errors.NewDatabaseError("Failed to retrieve users", result.Error)
	}
	return users, nil
}

// FindByID retrieves a user by ID
func (r *userRepository) FindByID(ctx context.Context, id uint) (*entities.User, error) {
	var user entities.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewResourceNotFoundError("User not found", map[string]interface{}{"id": id}, result.Error)
		}
		return nil, errors.NewDatabaseError("Failed to retrieve user", result.Error)
	}
	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewResourceNotFoundError("User not found", map[string]interface{}{"email": email}, result.Error)
		}
		return nil, errors.NewDatabaseError("Failed to retrieve user by email", result.Error)
	}
	return &user, nil
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	// Check if user with the same email already exists
	existingUser, err := r.FindByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return errors.NewDuplicateResourceError("User with this email already exists", map[string]interface{}{"email": user.Email}, nil)
	}

	// Create user
	result := r.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return errors.NewDatabaseError("Failed to create user", result.Error)
	}
	return nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	result := r.db.WithContext(ctx).Save(&user)
	if result.Error != nil {
		return errors.NewDatabaseError("Failed to update user", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewResourceNotFoundError("User not found", map[string]interface{}{"id": user.ID}, nil)
	}
	return nil
}

// Delete deletes a user
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entities.User{}, id)
	if result.Error != nil {
		return errors.NewDatabaseError("Failed to delete user", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewResourceNotFoundError("User not found", map[string]interface{}{"id": id}, nil)
	}
	return nil
}
