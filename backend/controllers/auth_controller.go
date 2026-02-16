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
	// var input models.RegisterInput

	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// user := models.User{
	// 	Username: input.Username,
	// 	Password: string(hashedPassword),
	// 	Role:     input.Role,
	// }

	// if err := database.DB.Create(&user).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah user"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"message": "Operator berhasil dibuat", "data": user})
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
	// Kita panggil fungsi helper di atas
	userRole := determineRole(employee.IDEmploy)

	// 4. Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  employee.NIK,      // ID pakai NIK
		"username": employee.NIK,      // Username pakai NIK
		"nama":     employee.Nama,     // Nama Asli
		"role":     userRole,          // Role hasil hardcode
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString(jwtKey)

	fmt.Printf("[LOGIN] User: %s | ID: %d | Role Assigned: %s\n", employee.Nama, employee.IDEmploy, userRole)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"role":  userRole,
		"nama":  employee.Nama,
		"nik":   employee.NIK,
	})
}

func GetUserProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("userRole")
	username, _ := c.Get("username")

	var emp models.TmEmploy
	var namaDisplay string
	
	// Cek context dulu
	if val, ok := c.Get("nama"); ok {
		namaDisplay = val.(string)
	} else {
		// Fallback query
		if err := database.MySQL.Where("nik = ?", userID).First(&emp).Error; err == nil {
			namaDisplay = emp.Nama
		} else {
			namaDisplay = username.(string)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       userID,
			"username": username,
			"role":     role,
			"name":     namaDisplay,
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