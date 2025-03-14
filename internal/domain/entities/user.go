package entities

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user entity
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" binding:"required" gorm:"size:100;not null"`
	Email     string         `json:"email" binding:"required,email" gorm:"size:100;uniqueIndex;not null"`
	Password  string         `json:"-" binding:"required,min=6" gorm:"size:100;not null"`
	Role      string         `json:"role" gorm:"size:20;default:'user'"`
	Active    bool           `json:"active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}

// UserCreate represents the structure for creating a user
type UserCreate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"omitempty,oneof=admin user"`
}

// UserUpdate represents the structure for updating a user
type UserUpdate struct {
	Name     *string `json:"name" binding:"omitempty"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6"`
	Role     *string `json:"role" binding:"omitempty,oneof=admin user"`
	Active   *bool   `json:"active" binding:"omitempty"`
}

// UserResponse represents the response structure for a user
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
