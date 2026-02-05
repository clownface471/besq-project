package controllers

import (
	"factory-api/database"
	"net/http"
	"github.com/gin-gonic/gin"
)

// Struct untuk input filter dari Frontend
type ChartFilter struct {
	Tanggal string `json:"tanggal"`  // YYYY-MM-DD
	Mesin   string `json:"mesin"`
	Shift   int    `json:"shift"`
}

// Struct untuk output ke Frontend
type ChartData struct {
	JamLabel   string  `json:"jam_label"`
	NilaiTotal float64 `json:"nilai_total"`
	NilaiNG    float64 `json:"nilai_ng"`
}

func GetProductionChart(c *gin.Context) {
	var filter ChartFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default jam kerja shift
	jamMulai, jamSelesai := 7, 15
	if filter.Shift == 2 {
		jamMulai, jamSelesai = 15, 23
	} else if filter.Shift == 3 {
		jamMulai, jamSelesai = 23, 7
	}

	// Query SQL Kompleks (MySQL)
	// Menggunakan parameter binding (?) untuk keamanan
	query := `
	WITH RECURSIVE jam_master AS (
		SELECT 0 AS jam_angka
		UNION ALL
		SELECT jam_angka + 1 FROM jam_master WHERE jam_angka < 23
	),
	raw_data AS (
		SELECT 
			h.Tanggal, h.noMC, j.jam_angka,
			ROUND(SUM(d.OK * (GREATEST(0, TIME_TO_SEC(TIMEDIFF(LEAST(d.SELESAI, MAKETIME(j.jam_angka + 1, 0, 0)), GREATEST(d.MULAI, MAKETIME(j.jam_angka, 0, 0)))) ) / NULLIF(TIME_TO_SEC(TIMEDIFF(d.SELESAI, d.MULAI)), 0)))) AS qty_ok_allocated,
			ROUND(SUM(d.NG * (GREATEST(0, TIME_TO_SEC(TIMEDIFF(LEAST(d.SELESAI, MAKETIME(j.jam_angka + 1, 0, 0)), GREATEST(d.MULAI, MAKETIME(j.jam_angka, 0, 0)))) ) / NULLIF(TIME_TO_SEC(TIMEDIFF(d.SELESAI, d.MULAI)), 0)))) AS qty_ng_allocated
		FROM trh_lwp_prs h
		JOIN tr_lwp_prs d ON d.IDhprs = h.IDhprs 
		CROSS JOIN jam_master j
		WHERE h.Tanggal = ? AND d.MULAI < MAKETIME(j.jam_angka + 1, 0, 0) AND d.SELESAI > MAKETIME(j.jam_angka, 0, 0)
		GROUP BY h.Tanggal, h.noMC, j.jam_angka
	)
	SELECT 
		CONCAT(LPAD(jam_angka, 2, '0'), ':00') AS jam_label,
		COALESCE(SUM(qty_ok_allocated + qty_ng_allocated), 0) AS nilai_total,
		COALESCE(SUM(qty_ng_allocated), 0) AS nilai_ng
	FROM raw_data
	WHERE noMC = ? AND jam_angka BETWEEN ? AND ?
	GROUP BY jam_angka
	ORDER BY jam_angka ASC;
	`

	var results []ChartData
	// Eksekusi Raw Query
	// Urutan parameter harus sesuai dengan tanda tanya (?) di query: Tanggal, NoMC, JamMulai, JamSelesai
	if err := database.DB.Raw(query, filter.Tanggal, filter.Mesin, jamMulai, jamSelesai).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data grafik: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}