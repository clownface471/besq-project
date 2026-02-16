package controllers

import (
	"factory-api/database"
	"factory-api/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	// "golang.org/x/crypto/bcrypt"
)

// ... (CreateOperator tetap sama, tidak perlu diubah) ...
func CreateOperator(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Fitur tambah user dimatikan saat menggunakan DB Perusahaan"})
}

var jwtKey = []byte("secret_key_besq_2026")

// --- FUNGSI LOGIN (YANG DIPERBAIKI) ---
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"` // NIK
		Password string `json:"password" binding:"required"` // Passw
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	var employee models.TmEmploy

	// 1. Cek ke Database MySQL berdasarkan NIK
	if err := database.MySQL.Where("nik = ?", input.Username).First(&employee).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "NIK tidak ditemukan"})
		return
	}

	// 2. Cek Password (Plain Text sesuai instruksi sebelumnya)
	if employee.Passw != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	// 3. Tentukan Role (HARDCODE CHECK)
	userRole := determineRole(employee.IDEmploy)

	// 4. Generate JWT - TAMBAHKAN EMAIL
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  employee.NIK,      // ID pakai NIK
		"username": employee.NIK,      // Username pakai NIK
		"nama":     employee.Nama,     // Nama Asli dari kolom 'nama'
		"email":    employee.EmailAddr, // Email dari kolom 'emailAddr'
		"role":     userRole,          // Role hasil hardcode
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString(jwtKey)

	fmt.Printf("[LOGIN] User: %s | ID: %d | Role Assigned: %s | Email: %s\n", 
		employee.Nama, employee.IDEmploy, userRole, employee.EmailAddr)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"role":  userRole,
		"nama":  employee.Nama,
		"nik":   employee.NIK,
		"email": employee.EmailAddr, // Tambahkan email dalam response
	})
}

func GetUserProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("userRole")
	username, _ := c.Get("username")

	var emp models.TmEmploy
	var namaDisplay string
	var emailDisplay string
	
	// Cek context dulu
	if val, ok := c.Get("nama"); ok {
		namaDisplay = val.(string)
	} else {
		namaDisplay = username.(string)
	}

	if val, ok := c.Get("email"); ok {
		emailDisplay = val.(string)
	} else {
		// Fallback query
		if err := database.MySQL.Where("nik = ?", userID).First(&emp).Error; err == nil {
			namaDisplay = emp.Nama
			emailDisplay = emp.EmailAddr
		} else {
			emailDisplay = "-"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       userID,
			"username": username,
			"role":     role,
			"name":     namaDisplay,
			"email":    emailDisplay, // Tambahkan email
		},
	})
}

func determineRole(idEmploy int) string {
	// Jika nanti DB sudah ada kolom role, fungsi ini bisa diganti
	// menjadi sekedar return employee.Role
	
	switch idEmploy {
	case 698:
		return "ADMIN"
	case 699:
		return "MANAGER"
	case 700:
		return "LEADER"
	case 701:
		return "OPERATOR_CUTTING"
	default:
		// Sisanya dianggap Pressing
		return "OPERATOR_PRESSING"
	}
}