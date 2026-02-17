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

// GET: /api/pressing/weekly-stats
func GetPressingWeeklyStats(c *gin.Context) {
	namaOperator := c.Query("nama")
	if namaOperator == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'nama' operator wajib diisi"})
		return
	}

	// Hitung tanggal 7 hari yang lalu
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -6) // 6 hari ke belakang + hari ini = 7 hari

	type DailyStats struct {
		Tanggal    time.Time `json:"tanggal"`
		Total      int       `json:"total"`
		OK         int       `json:"ok"`
		NG         int       `json:"ng"`
		Efficiency float64   `json:"efficiency"`
	}

	// HAPUS BARIS INI - tidak digunakan:
	// var dailyStats []DailyStats

	// Query untuk aggregate data per hari
	query := `
		SELECT 
			DATE(CONCAT(thn, '-', LPAD(bln, 2, '0'), '-', LPAD(tgl, 2, '0'))) as tanggal,
			COALESCE(SUM(TOTAL), 0) as total,
			COALESCE(SUM(OK), 0) as ok,
			COALESCE(SUM(NG), 0) as ng
		FROM vtrx_lwp_prs
		WHERE nama = ?
			AND DATE(CONCAT(thn, '-', LPAD(bln, 2, '0'), '-', LPAD(tgl, 2, '0'))) BETWEEN ? AND ?
		GROUP BY DATE(CONCAT(thn, '-', LPAD(bln, 2, '0'), '-', LPAD(tgl, 2, '0')))
		ORDER BY tanggal ASC
	`

	rows, err := database.MySQL.Raw(query, namaOperator, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal query data mingguan: " + err.Error()})
		return
	}
	defer rows.Close()

	// Map untuk menyimpan data per tanggal
	dataMap := make(map[string]DailyStats)
	for rows.Next() {
		var stat DailyStats
		if err := rows.Scan(&stat.Tanggal, &stat.Total, &stat.OK, &stat.NG); err != nil {
			continue
		}
		
		// Hitung efficiency
		if stat.Total > 0 {
			stat.Efficiency = (float64(stat.OK) / float64(stat.Total)) * 100
		}
		
		dataMap[stat.Tanggal.Format("2006-01-02")] = stat
	}

	// Format data untuk frontend (pastikan 7 hari lengkap)
	dayNames := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	dayShorts := []string{"Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"}
	
	var weeklyData []gin.H
	
	for i := 0; i < 7; i++ {
		currentDate := startDate.AddDate(0, 0, i)
		dateKey := currentDate.Format("2006-01-02")
		dayOfWeek := int(currentDate.Weekday())
		
		stat, exists := dataMap[dateKey]
		if !exists {
			stat = DailyStats{
				Tanggal:    currentDate,
				Total:      0,
				OK:         0,
				NG:         0,
				Efficiency: 0,
			}
		}
		
		weeklyData = append(weeklyData, gin.H{
			"day":        dayNames[dayOfWeek],
			"short":      dayShorts[dayOfWeek],
			"date":       dateKey,
			"total":      stat.Total,
			"ok":         stat.OK,
			"ng":         stat.NG,
			"efficiency": stat.Efficiency,
		})
	}

	// Hitung summary untuk hari ini
	today := endDate.Format("2006-01-02")
	todayStats, _ := dataMap[today]

	c.JSON(http.StatusOK, gin.H{
		"weeklyData": weeklyData,
		"summary": gin.H{
			"todayCompleted": todayStats.OK,
			"todayTarget":    60, // Bisa disesuaikan dengan target real
			"efficiency":     int(todayStats.Efficiency),
		},
	})
}

// GetPressingLWPData - Ambil data LWP dari database untuk operator tertentu
func GetPressingLWPData(c *gin.Context) {
	namaOperator := c.Query("nik")
	tanggal := c.Query("tanggal") // Format: YYYY-MM-DD
	
	if namaOperator == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'nik' wajib diisi"})
		return
	}

	// Jika tanggal kosong, gunakan hari ini
	if tanggal == "" {
		tanggal = time.Now().Format("2006-01-02")
	}

	// Parse tanggal untuk mendapatkan thn, bln, tgl
	parsedDate, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal salah (gunakan YYYY-MM-DD)"})
		return
	}

	thn := parsedDate.Year()
	bln := int(parsedDate.Month())
	tgl := parsedDate.Day()

	type LWPRecord struct {
		NoMC              string    `json:"noMesin"`
		Tanggal           string    `json:"tanggal"`
		Shift             string    `json:"shift"`
		NPK               string    `json:"nik"`
		ItemCode          string    `json:"kodePart"`
		ItemName          string    `json:"partName"`
		LotNo             string    `json:"noLot"`
		JamMulai          string    `json:"jamMulai"`
		JamSelesai        string    `json:"jamSelesai"`
		OK                int       `json:"hasilOk"`
		NG                int       `json:"ng"`
		Bintik            int       `json:"bintik"`
		Kotor             int       `json:"kotor"`
		Sompel            int       `json:"sompel"`
		Belang            int       `json:"belang"`
		Scratch           int       `json:"scratch"`
		Tipis             int       `json:"tipis"`
		Gelombang         int       `json:"gelombang"`
		SilverMark        int       `json:"silverMark"`
		FlowMark          int       `json:"flowMark"`
		Gompal            int       `json:"gompal"`
		Bending           int       `json:"bending"`
		ShortShot         int       `json:"shortShot"`
		Bushing           int       `json:"bushing"`
		Nilon             int       `json:"nilon"`
		Sink              int       `json:"sink"`
		HitamPutih        int       `json:"hitamPutih"`
		Lain2             int       `json:"lain2"`
	}

	var records []LWPRecord

	// Query ke database dengan JOIN ke v_stdlot untuk mendapatkan itemName
	query := `
		SELECT 
			v.noMC,
			CONCAT(v.thn, '-', LPAD(v.bln, 2, '0'), '-', LPAD(v.tgl, 2, '0')) as tanggal,
			v.shift,
			v.NPK,
			v.itemcode,
			COALESCE(s.itemName, '-') as itemName,
			v.lotno,
			TIME_FORMAT(v.MULAI, '%H:%i') as jamMulai,
			TIME_FORMAT(v.SELESAI, '%H:%i') as jamSelesai,
			v.OK,
			v.NG,
			COALESCE(v.bintik, 0) as bintik,
			COALESCE(v.kotor, 0) as kotor,
			COALESCE(v.sompel, 0) as sompel,
			COALESCE(v.belang, 0) as belang,
			COALESCE(v.scratch, 0) as scratch,
			COALESCE(v.tipis, 0) as tipis,
			COALESCE(v.gelombang, 0) as gelombang,
			COALESCE(v.silvermark, 0) as silverMark,
			COALESCE(v.flowmark, 0) as flowMark,
			COALESCE(v.gompal, 0) as gompal,
			COALESCE(v.bending, 0) as bending,
			COALESCE(v.shortshoot, 0) as shortShot,
			COALESCE(v.bushing, 0) as bushing,
			COALESCE(v.nilon, 0) as nilon,
			COALESCE(v.sink, 0) as sink,
			COALESCE(v.hitamputih, 0) as hitamPutih,
			COALESCE(v.lain2, 0) as lain2
		FROM vtrx_lwp_prs v
		LEFT JOIN v_stdlot s ON v.itemcode = s.itemCode COLLATE utf8mb4_unicode_ci
		WHERE v.NPK = ? 
			AND v.thn = ? 
			AND v.bln = ? 
			AND v.tgl = ?
		ORDER BY v.MULAI DESC
	`

	if err := database.MySQL.Raw(query, namaOperator, thn, bln, tgl).Scan(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data LWP: " + err.Error()})
		return
	}

	// Group data berdasarkan mesin
	type MachineGroup struct {
		NoMesin   string                   `json:"noMesin"`
		Tanggal   string                   `json:"tanggal"`
		Shift     string                   `json:"shift"`
		NIK       string                   `json:"nik"`
		PartName  string                   `json:"partName"`
		KodePart  string                   `json:"kodePart"`
		Details   []map[string]interface{} `json:"details"`
	}

	machineMap := make(map[string]*MachineGroup)

	for _, record := range records {
		key := record.NoMC + "_" + record.Shift

		if _, exists := machineMap[key]; !exists {
			machineMap[key] = &MachineGroup{
				NoMesin:  record.NoMC,
				Tanggal:  record.Tanggal,
				Shift:    record.Shift,
				NIK:      record.NPK,
				PartName: record.ItemName,
				KodePart: record.ItemCode,
				Details:  []map[string]interface{}{},
			}
		}

		// Build klasifikasi reject array
		klasifikasiReject := []map[string]interface{}{}
		rejectTypes := map[string]int{
			"Bintik":       record.Bintik,
			"Kotor":        record.Kotor,
			"Sompel":       record.Sompel,
			"Belang":       record.Belang,
			"Scratch":      record.Scratch,
			"Tipis":        record.Tipis,
			"Gelombang":    record.Gelombang,
			"Silver Mark":  record.SilverMark,
			"Flow Mark":    record.FlowMark,
			"Gompal":       record.Gompal,
			"Bending":      record.Bending,
			"Short Shot":   record.ShortShot,
			"Bushing":      record.Bushing,
			"Nilon":        record.Nilon,
			"Sink":         record.Sink,
			"Hitam Putih":  record.HitamPutih,
			"Lain-lain":    record.Lain2,
		}

		for jenis, qty := range rejectTypes {
			if qty > 0 {
				klasifikasiReject = append(klasifikasiReject, map[string]interface{}{
					"jenis": jenis,
					"qty":   qty,
				})
			}
		}

		detail := map[string]interface{}{
			"noLot":             record.LotNo,
			"jamMulai":          record.JamMulai,
			"jamSelesai":        record.JamSelesai,
			"hasilOk":           record.OK,
			"ng":                record.NG,
			"klasifikasiReject": klasifikasiReject,
		}

		machineMap[key].Details = append(machineMap[key].Details, detail)
	}

	// Convert map to array
	result := []MachineGroup{}
	for _, group := range machineMap {
		result = append(result, *group)
	}

	c.JSON(http.StatusOK, gin.H{
		"lwpRecords": result,
		"total":      len(result),
	})
}