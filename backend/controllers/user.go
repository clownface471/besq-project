package controllers

import (
	"besq-backend/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserController handles user-related CRUD operations
type UserController struct {
	db *gorm.DB
}

// NewUserController creates a new instance of UserController
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

// CreateUser creates a new user (Admin only)
// Hashes password before saving to database
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("Validation Error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize NIK to lowercase for case-insensitive user handling
	user.NIK = strings.ToLower(user.NIK)

	// Validate role
	if !models.IsValidRole(user.Role) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role. Valid roles are: admin, cutting, pressing",
		})
		return
	}

	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	user.Password = string(hashedPassword)

	// Create user in database
	if err := uc.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "NIK already exists or database error",
		})
		return
	}

	// Return response without password
	userResponse := models.UserResponse{
		ID:        user.ID,
		NIK:       user.NIK,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    userResponse,
	})
}

// GetUsers returns a list of all users
func (uc *UserController) GetUsers(c *gin.Context) {
	var users []models.User
	if err := uc.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}

	// Convert to response format (without passwords)
	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			ID:        user.ID,
			NIK:       user.NIK,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": userResponses,
		"total": len(userResponses),
	})
}

// GetUser returns a single user by ID
func (uc *UserController) GetUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := uc.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// Return user without password
	userResponse := models.UserResponse{
		ID:        user.ID,
		NIK:       user.NIK,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userResponse,
	})
}

// UpdateUser edits user name, role, and/or password (Admin only)
func (uc *UserController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var updateData struct {
		Name     string `json:"name"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validate role if provided
	if updateData.Role != "" && !models.IsValidRole(updateData.Role) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role. Valid roles are: admin, cutting, pressing",
		})
		return
	}

	// Find user to update
	var user models.User
	if err := uc.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// Normalize NIK to lowercase for case-insensitive user handling
	user.NIK = strings.ToLower(user.NIK)

	// Update user fields if provided
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Role != "" {
		user.Role = updateData.Role
	}
	if updateData.Password != "" {
		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password",
			})
			return
		}
		user.Password = string(hashedPassword)
	}

	// Save updated user
	if err := uc.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	// Return updated user without password
	userResponse := models.UserResponse{
		ID:        user.ID,
		NIK:       user.NIK,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    userResponse,
	})
}

// DeleteUser deletes a user by ID (Admin only)
func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	// Find user to delete
	var user models.User
	if err := uc.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// Prevent deletion of the last admin user
	if user.Role == "admin" {
		var adminCount int64
		uc.db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot delete the last admin user",
			})
			return
		}
	}

	// Delete user
	if err := uc.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// GetProfile returns the current user's profile (requires authentication)
func (uc *UserController) GetProfile(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	var user models.User
	if err := uc.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// Return user profile without password
	userResponse := models.UserResponse{
		ID:        user.ID,
		NIK:       user.NIK,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userResponse,
	})
}
