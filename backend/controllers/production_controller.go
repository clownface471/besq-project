package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"factory-api/database"
	"factory-api/models"
	"time"
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

func RecordPressingCycle(c *gin.Context) {
	var input models.PerCycle

	// Validasi input JSON dari Frontend
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Opsi Tambahan: Jika frontend tidak mengirim nama_operator, 
	// kita bisa ambil otomatis dari Token login (Context)
	if input.NamaOperator == "" {
		username, _ := c.Get("username")
		input.NamaOperator = username.(string)
	}

	// Simpan ke Database
	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data cycle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cycle recorded successfully",
		"data":    input,
	})
}
type PressingDashboardData struct {
	Stats      PressingStats      `json:"stats"`
	Scans      []PressingScan     `json:"scans"`
	LwpRecords []models.PerCycle  `json:"lwpRecords"` // Kita gunakan model cycle sederhana dulu
}

type PressingStats struct {
	Completed      int `json:"completed"`
	Target         int `json:"target"`
	Efficiency     int `json:"efficiency"`
	TodayCompleted int `json:"todayCompleted"`
	TodayTarget    int `json:"todayTarget"`
}

type PressingScan struct {
	ID        uint   `json:"id"`
	Lot       string `json:"lot"`
	TimeStart string `json:"timeStart"`
	TimeEnd   string `json:"timeEnd"`
	Status    string `json:"status"`
	Pic       string `json:"pic"`
}

// GET: /api/pressing/today
func GetPressingDashboard(c *gin.Context) {
	// 1. Data Dummy / Mockup (Nanti bisa diganti query DB beneran)
	stats := PressingStats{
		Completed:      125,
		Target:         200,
		Efficiency:     85,
		TodayCompleted: 45, // Data untuk Donut Chart
		TodayTarget:    60,
	}

	// 2. Data Scans Dummy
	scans := []PressingScan{
		{ID: 1, Lot: "LOT-2026-001", TimeStart: "08:00", TimeEnd: "08:45", Status: "Selesai", Pic: "Budi"},
		{ID: 2, Lot: "LOT-2026-002", TimeStart: "09:00", TimeEnd: "09:45", Status: "Selesai", Pic: "Budi"},
		{ID: 3, Lot: "LOT-2026-003", TimeStart: "10:00", TimeEnd: "-", Status: "Proses", Pic: "Budi"},
	}

	// 3. Data LWP (Ambil dari tabel per_cycles yang baru kita buat)
	var cycles []models.PerCycle
	database.DB.Order("created_at desc").Limit(10).Find(&cycles)

	c.JSON(http.StatusOK, gin.H{
		"stats":      stats,
		"scans":      scans,
		"lwpRecords": cycles,
	})
}
// [backend/controllers/production_controller.go]

// Struct untuk input Scan Mesin
type ScanMachineInput struct {
	MachineCode string `json:"machineCode" binding:"required"`
	Timestamp   string `json:"timestamp"`
}

// POST: /api/scan-machine
func ScanMachine(c *gin.Context) {
	var input ScanMachineInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Di sini kamu bisa cek ke DB apakah kode mesin valid
	// Untuk sekarang kita anggap semua kode yang diawali "PRESS" atau "MESIN" valid
	// Simpan log atau update status mesin di DB jika perlu

	c.JSON(http.StatusOK, gin.H{
		"message": "Mesin berhasil divalidasi",
		"machine": input.MachineCode,
	})
}
// [backend/controllers/production_controller.go]

// Struct untuk input LWP (sesuaikan dengan frontend)
type LWPInput struct {
	NoMesin    string `json:"noMesin"`
	Tanggal    string `json:"tanggal"`
	Shift      string `json:"shift"`
	Nik        string `json:"nik"`
	PartName   string `json:"partName"`
	KodePart   string `json:"kodePart"`
	JamMulai   string `json:"jamMulai"`
	JamSelesai string `json:"jamSelesai"`
	HasilOk    int    `json:"hasilOk"`
	Ng         int    `json:"ng"`
}

// POST: /api/lwp (Buat Record Baru)
func CreateLWP(c *gin.Context) {
	var input LWPInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil username dari token untuk ngisi NIK jika kosong
	if input.Nik == "" {
        // Coba ambil username dari token
        val, exists := c.Get("username")
        if !exists || val == nil {
            // Jika token tidak punya username, gunakan fallback atau error
            input.Nik = "UNKNOWN_OPERATOR"
            // Atau bisa return error 401 jika wajib tahu siapa usernya:
            // c.JSON(http.StatusUnauthorized, gin.H{"error": "User identity not found"})
            // return
        } else {
            // Konversi aman
            if str, ok := val.(string); ok {
                input.Nik = str
            } else {
                input.Nik = "INVALID_USER_TYPE"
            }
        }
    }

	// Simpan ke database (di sini kita pakai model PerCycle sebagai contoh sederhana)
	// Idealnya ada tabel LWP khusus, tapi kita gunakan PerCycle biar cepat integrasi
	cycle := models.PerCycle{
		NoMC:         input.NoMesin,
		Item:         input.PartName,
		NoLot:        "NEW-LOT-" + time.Now().Format("150405"), // Generate NoLot unik
		StatusMesin:  "produksi",
		NamaOperator: input.Nik,
	}

	if err := database.DB.Create(&cycle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan LWP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "LWP Created", "data": cycle})
}