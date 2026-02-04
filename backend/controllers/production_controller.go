package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"factory-api/database"
	"factory-api/models"
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

// [Tambahkan di backend/controllers/production_controller.go]

// Struct response khusus untuk Dashboard Frontend
type DashboardData struct {
	TotalEmployees  int64               `json:"totalEmployees"`
	TotalOutput     int                 `json:"totalOutput"` // Total lot/WO selesai
	RejectRate      string              `json:"rejectRate"`  // Hardcode dulu atau hitung
	ActiveShift     string              `json:"activeShift"` // Hardcode "Shift 1"
	CuttingOutput   int64               `json:"cuttingOutput"`
	PressingOutput  int64               `json:"pressingOutput"`
	FinishingOutput int64               `json:"finishingOutput"`
	Activities      []ActivityResponse  `json:"activities"`
}

type ActivityResponse struct {
	User   string `json:"user"`
	Action string `json:"action"`
	Time   string `json:"time"`
	Status string `json:"status"`
}

func GetDashboardStats(c *gin.Context) {
	// 1. Hitung Total Employees (User)
	var totalEmp int64
	database.DB.Model(&models.User{}).Count(&totalEmp)

	// 2. Hitung Output per Divisi (Berdasarkan WorkOrder yang 'DONE' atau sesuai status)
	// Catatan: Karena model WorkOrder sederhana, kita hitung kasar dulu
	var cuttingCount int64
	database.DB.Model(&models.WorkOrder{}).Where("status = ?", "CUTTING").Count(&cuttingCount)
	
	var pressingCount int64
	database.DB.Model(&models.WorkOrder{}).Where("status = ?", "PRESSING").Count(&pressingCount)

	// 3. Ambil 5 Aktivitas Terakhir (AuditLog)
	var logs []models.AuditLog
	// Preload User jika ada relasi, tapi karena di AuditLog sudah simpan Username, langsung saja
	database.DB.Order("created_at desc").Limit(5).Find(&logs)

	var activities []ActivityResponse
	for _, log := range logs {
		activities = append(activities, ActivityResponse{
			User:   log.Username,
			Action: log.Action, // Misal: "START_CUTTING"
			Time:   log.CreatedAt.Format("15:04"), // Format Jam:Menit
			Status: "Completed", // Default status untuk log
		})
	}

	// 4. Return JSON
	c.JSON(http.StatusOK, DashboardData{
		TotalEmployees:  totalEmp,
		TotalOutput:     int(cuttingCount + pressingCount), // Contoh penjumlahan
		RejectRate:      "0.5%", // Mock data dulu
		ActiveShift:     "Shift 1 (Pagi)",
		CuttingOutput:   cuttingCount,
		PressingOutput:  pressingCount,
		FinishingOutput: 0,
		Activities:      activities,
	})
}