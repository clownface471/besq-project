package controllers

import (
	"besq-backend/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthController handles authentication-related requests
type AuthController struct {
	db *gorm.DB
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

// JWT secret key - in production, this should be in environment variables
var jwtSecret = []byte("your-secret-key-change-this-in-production")

// Claims represents the JWT claims structure
type Claims struct {
	Sub  string `json:"sub"`  // User ID
	Role string `json:"role"` // User Role
	jwt.RegisteredClaims
}

// Login handles user login and JWT token generation
func (ac *AuthController) Login(c *gin.Context) {
	var loginReq models.UserLoginRequest

	// Bind JSON input
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Normalize NIK to lowercase for case-insensitive login
	loginReq.NIK = strings.ToLower(loginReq.NIK)

	// Find user by NIK using GORM
	var user models.User
	if err := ac.db.Where("nik = ?", loginReq.NIK).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid NIK or password",
		})
		return
	}

	// Compare hashed password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid NIK or password",
		})
		return
	}

	// Generate JWT token
	token, err := ac.generateJWTToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Return success response with token
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// generateJWTToken creates a JWT token with sub (UserID) and role (UserRole), expiring in 24 hours
func (ac *AuthController) generateJWTToken(userID uint, role string) (string, error) {
	// Create claims with sub (UserID) and role (UserRole)
	claims := &Claims{
		Sub:  strconv.FormatUint(uint64(userID), 10), // Convert userID to string
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours expiration
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "besq-backend",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWTToken validates and parses a JWT token (utility function for future middleware)
func ValidateJWTToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
