package controllers

import (
	"factory-api/database"
	"factory-api/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateWorkOrder - Hanya Admin yang bisa membuat WO baru
func CreateWorkOrder(c *gin.Context) {
	var input struct {
		OrderNumber string `json:"order_number" binding:"required"`
		PartName    string `json:"part_name" binding:"required"`
		Quantity    int    `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// Ambil data admin yang sedang login dari token
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	wo := models.WorkOrder{
		OrderNumber: input.OrderNumber,
		PartName:    input.PartName,
		Quantity:    input.Quantity,
		Status:      "PENDING", // Status awal selalu PENDING
		OperatorID:  uint(userID.(float64)),
	}

	if err := database.DB.Create(&wo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat Work Order"})
		return
	}

	// Catat ke Audit Log
	database.RecordActivity(uint(userID.(float64)), username.(string), "CREATE_WORK_ORDER", "/admin/work-order")

	c.JSON(http.StatusCreated, gin.H{"message": "Work Order berhasil dibuat", "data": wo})
}

// GetAllWorkOrders - Semua role bisa melihat daftar pekerjaan
func GetAllWorkOrders(c *gin.Context) {
	var workOrders []models.WorkOrder
	
	// Kita ambil semua data (bisa ditambahkan pagination di sini nanti)
	database.DB.Order("created_at desc").Find(&workOrders)

	c.JSON(http.StatusOK, gin.H{
		"total": len(workOrders),
		"data":  workOrders,
	})
}

// UpdateWOStatus - Operator mengupdate status (Misal: PENDING -> ON_PROGRESS)
func UpdateWOStatus(c *gin.Context) {
	id := c.Param("id")
	var wo models.WorkOrder

	// Cek apakah WO ada
	if err := database.DB.First(&wo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work Order tidak ditemukan"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus diisi"})
		return
	}

	// Ambil data operator yang melakukan update
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	// Update Status
	oldStatus := wo.Status
	wo.Status = input.Status
	wo.OperatorID = uint(userID.(float64)) // Update siapa yang terakhir mengerjakan

	database.DB.Save(&wo)

	// Catat ke Audit Log (Penting untuk pelacakan)
	actionDesc := fmt.Sprintf("UPDATE_STATUS: %s -> %s", oldStatus, input.Status)
	database.RecordActivity(uint(userID.(float64)), username.(string), actionDesc, "/production/work-order/status")

	c.JSON(http.StatusOK, gin.H{
		"message": "Status Work Order diperbarui", 
		"data": wo,
	})
}