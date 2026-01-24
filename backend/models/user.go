package models

import (
	"gorm.io/gorm"
	"time"
)

// User represents the user model
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	NIK       string         `gorm:"uniqueIndex;not null" json:"nik" binding:"required"`
	Name      string         `gorm:"not null" json:"name" binding:"required"`
	Role      string         `gorm:"not null;check:role IN ('admin', 'cutting', 'pressing')" json:"role" binding:"required"`
	Password  string         `gorm:"not null" json:"password" binding:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// IsValidRole checks if the given role is valid
func IsValidRole(role string) bool {
	validRoles := []string{"admin", "cutting", "pressing"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// UserLoginRequest represents the login request payload
type UserLoginRequest struct {
	NIK      string `json:"nik" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the user response without password
type UserResponse struct {
	ID        uint      `json:"id"`
	NIK       string    `json:"nik"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}