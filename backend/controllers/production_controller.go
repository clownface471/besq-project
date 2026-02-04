package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"factory-api/database"
)

// GET: Melihat status mesin Cutting
func GetCuttingStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"section": "Cutting",
		"status":  "Operating",
		"message": "Data ini bisa dilihat oleh semua operator dan admin",
	})
}

// GET: Melihat status mesin Pressing
func GetPressingStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"section": "Pressing",
		"status":  "Idle",
		"message": "Data ini bisa dilihat oleh semua operator dan admin",
	})
}

// POST: Menjalankan mesin Cutting (Hanya Role Cutting & Admin)
func StartCutting(c *gin.Context) {
	// Ambil data user dari JWT claims (kita perlu update middleware sedikit untuk menyimpan userID)
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	// Simpan ke Audit Log
	database.RecordActivity(uint(userID.(float64)), username.(string), "START_CUTTING_MACHINE", c.Request.URL.Path)

	c.JSON(http.StatusOK, gin.H{
		"message": "Mesin Cutting berhasil dijalankan!",
		"logged":  true,
	})
}

// POST: Menjalankan mesin Pressing (Hanya Role Pressing & Admin)
func StartPressing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Mesin Pressing berhasil dijalankan!",
	})
}
