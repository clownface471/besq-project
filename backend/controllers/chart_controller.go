package controllers

import (
	"factory-api/database"
	"factory-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Helper: Cek koneksi DB
func isDBConnected(c *gin.Context) bool {
	if database.MySQL == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Database Statistik (MySQL) tidak terhubung. Pastikan VPN aktif atau konfigurasi DB benar.",
			"data":  []models.ChartSeries{}, 
		})
		return false
	}
	return true
}

// --- LEVEL 1: MANAGER VIEW (Overview Per Proses) ---
func GetManagerOverview(c *gin.Context) {
	if !isDBConnected(c) { return }

	tanggal := c.Query("tanggal") 
	var results []models.ChartSeries

	// Rumus: Target = (Durasi Jam) * (Target Qty/Jam)
	query := `
		SELECT 
			t.proses AS label,
			COALESCE(SUM((TIME_TO_SEC(TIMEDIFF(t.SELESAI, t.MULAI)) / 3600.0) * s.tgtQtyPJam), 0) AS target,
			COALESCE(SUM(t.Total), 0) AS actual,
			COALESCE(SUM(t.NG), 0) AS actual_ng
		FROM vtrx_lwp_prs t
		LEFT JOIN v_stdlot s ON t.itemCode = s.itemCode COLLATE utf8mb4_unicode_ci
		WHERE t.tanggal = ?
		GROUP BY t.proses
	`

	if err := database.MySQL.Raw(query, tanggal).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data manager: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// --- LEVEL 2: LEADER VIEW (Overview Per Mesin) ---
func GetLeaderProcessView(c *gin.Context) {
	if !isDBConnected(c) { return }

	tanggal := c.Query("tanggal")
	proses := c.Query("proses") 

	var results []models.ChartSeries

	// Rumus Konsisten: Target = (Durasi Jam) * (Target Qty/Jam)
	query := `
		SELECT 
			t.noMC AS label,
			COALESCE(SUM((TIME_TO_SEC(TIMEDIFF(t.SELESAI, t.MULAI)) / 3600.0) * s.tgtQtyPJam), 0) AS target,
			COALESCE(SUM(t.Total), 0) AS actual,
			COALESCE(SUM(t.NG), 0) AS actual_ng
		FROM vtrx_lwp_prs t
		LEFT JOIN v_stdlot s ON t.itemCode = s.itemCode COLLATE utf8mb4_unicode_ci
		WHERE t.tanggal = ? AND t.proses = ?
		GROUP BY t.noMC
		ORDER BY t.noMC
	`

	if err := database.MySQL.Raw(query, tanggal, proses).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data leader: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// --- LEVEL 3: MACHINE DETAIL (Per Jam) WITH SHIFT FILTER ---
func GetMachineDetail(c *gin.Context) {
	if !isDBConnected(c) { return }

	tanggal := c.Query("tanggal")
	noMC := c.Query("no_mc")
	shift := c.Query("shift") // Get shift parameter

	var results []models.ChartSeries
	
	// Define shift hour ranges
	var startHour, endHour int
	switch shift {
	case "1":
		startHour = 0
		endHour = 8
	case "2":
		startHour = 8
		endHour = 16
	case "3":
		startHour = 16
		endHour = 24
	default:
		// If no shift specified, show all 24 hours
		startHour = 0
		endHour = 24
	}

	// Updated query with shift filtering
	query := `
WITH RECURSIVE 
jam_master AS (
    SELECT ? AS jam_angka
    UNION ALL
    SELECT jam_angka + 1 FROM jam_master WHERE jam_angka < ?
),
raw_data AS (
    SELECT 
        t.tanggal, t.nama AS operator, t.MULAI, t.SELESAI,
        s.itemName, s.tgtQtyPJam, j.jam_angka,
        
        -- Hitung detik overlap di jam ini
        (GREATEST(0, TIME_TO_SEC(TIMEDIFF(LEAST(t.SELESAI, MAKETIME(j.jam_angka + 1, 0, 0)), GREATEST(t.MULAI, MAKETIME(j.jam_angka, 0, 0)))))) AS seconds_in_slot,
        
        -- Total durasi transaksi asli (untuk pembagi rasio)
        NULLIF(TIME_TO_SEC(TIMEDIFF(t.SELESAI, t.MULAI)), 0) AS total_duration_sec,

        t.Total, t.OK, t.NG
    FROM vtrx_lwp_prs t
    -- JOIN SUDAH DIPERBARUI KE itemCode
    LEFT JOIN v_stdlot s ON t.itemCode = s.itemCode COLLATE utf8mb4_unicode_ci
    CROSS JOIN jam_master j
    WHERE 
        t.tanggal = ? 
        AND t.noMC = ? 
        AND t.MULAI < MAKETIME(j.jam_angka + 1, 0, 0) 
        AND t.SELESAI > MAKETIME(j.jam_angka, 0, 0)
)
SELECT 
    CONCAT(LPAD(jam_angka, 2, '0'), ':00') AS label,
    
    -- Target = (Detik di Slot Ini / 3600) * Target Per Jam
    COALESCE(SUM((seconds_in_slot / 3600.0) * tgtQtyPJam), 0) AS target,
    
    -- Actual & NG dialokasikan berdasarkan rasio waktu
    COALESCE(SUM(ROUND(Total * (seconds_in_slot / total_duration_sec))), 0) AS actual,
    COALESCE(SUM(ROUND(NG * (seconds_in_slot / total_duration_sec))), 0) AS actual_ng,
    
    -- FITUR YANG HILANG KARENA CONFLICT (JANGAN SAMPAI KETINGGALAN)
    COALESCE(GROUP_CONCAT(DISTINCT itemName SEPARATOR ', '), '-') AS item_name,
    
    CONCAT(COALESCE(MAX(itemName), '-'), ' (', COALESCE(MAX(operator), '-'), ')') AS extra_info
FROM raw_data
GROUP BY jam_angka
ORDER BY jam_angka ASC;
	`

	if err := database.MySQL.Raw(query, startHour, endHour-1, tanggal, noMC).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil detail mesin: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}