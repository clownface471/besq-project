package controllers

import (
	"factory-api/database"
	"factory-api/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	
	// Ambil parameter dari URL, contoh: /admin/audit-logs?page=1&limit=10
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	
	database.DB.Order("created_at desc").Limit(limit).Offset(offset).Find(&logs)
	
	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"data":  logs,
	})
}

// GET: Ambil semua user (Operator)
func GetAllUsers(c *gin.Context) {
	var users []models.User
	// Ambil semua user kecuali password hashnya agar lebih aman (opsional, tapi good practice)
	if err := database.DB.Select("id, username, role, created_at, updated_at").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// PUT: Update User (Ganti Role atau Reset Password)
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Username string `json:"username"`
		Role     string `json:"role"`
		Password string `json:"password"` // Opsional, hanya jika ingin reset
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Update field
	if input.Username != "" { user.Username = input.Username }
	if input.Role != "" { user.Role = input.Role }
	
	// Jika password diisi, hash ulang
	if input.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
	}

	database.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User berhasil diupdate", "data": user})
}

// DELETE: Hapus User
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}