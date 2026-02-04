package controllers

import (
	"factory-api/database"
	"fmt"
	"factory-api/models"
	"net/http"
	"time"              // TAMBAHKAN INI
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // TAMBAHKAN INI
	"golang.org/x/crypto/bcrypt"
)

func CreateOperator(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password sebelum simpan
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Operator berhasil dibuat", "data": user})
}

var jwtKey = []byte("secret_key_besq_2026") // Ganti dengan key yang lebih kuat di production

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
        fmt.Println("Error Binding:", err.Error()) // Lihat log di VS Code terminal
        c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
        return
    }
	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Cek Password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString(jwtKey)

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "role": user.Role})
}

func GetUserProfile(c *gin.Context) {
	// Mengambil data yang diset oleh middleware (userID, role, username)
	// Kita tidak perlu query DB lagi karena data penting sudah ada di token
	userID, _ := c.Get("userID")
	role, _ := c.Get("userRole")
	username, _ := c.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       userID,
			"username": username,
			"role":     role, // Frontend butuh ini untuk redirect
			"name":     username, // Gunakan username sebagai nama tampilan sementara
		},
	})
}