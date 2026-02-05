package database

import (
	"factory-api/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite" // Pastikan import ini ada
	"gorm.io/gorm"
)

// DB Utama (SQLite) untuk User, Login, Cycle Button
var DB *gorm.DB

// DB Statistik (MySQL) untuk Grafik Laporan
var MySQL *gorm.DB

func ConnectDatabase() {
	// 1. KONEKSI SQLITE (Database Utama Aplikasi)
	// Ini mengembalikan fungsi user management ke kondisi semula
	sqliteDB, err := gorm.Open(sqlite.Open("besq.db"), &gorm.Config{})
	if err != nil {
		panic("Gagal koneksi ke SQLite: " + err.Error())
	}

	// Migrasi tabel aplikasi ke SQLite
	sqliteDB.AutoMigrate(&models.User{}, &models.AuditLog{}, &models.PerCycle{})
	DB = sqliteDB
	fmt.Println("✅ SQLite Connected (besq.db)")

	// 2. KONEKSI MYSQL (Database Statistik/Pabrik)
	// Sesuaikan nama database dengan yang ada di HeidiSQL kamu!
	// Jika nama DB kamu bukan 'besq_factory', ganti bagian ini.
	dsn := "root:@tcp(127.0.0.1:3306)/besq_factory?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true"
	
	mysqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// Kita gunakan Println, bukan Panic.
		// Jadi kalau MySQL mati/salah nama DB, aplikasi tetap bisa jalan untuk login & cycle.
		fmt.Println("⚠️  PERINGATAN: Gagal koneksi ke MySQL untuk Statistik. Fitur Grafik tidak akan jalan.")
		fmt.Println("   Error:", err.Error())
	} else {
		MySQL = mysqlDB
		fmt.Println("✅ MySQL Connected (Data Statistik Siap)")
	}
}

// Fungsi Helper RecordActivity (Masuk ke SQLite)
func RecordActivity(userID uint, username string, action string, details string) {
	fmt.Printf("[AUDIT] User: %s (ID: %d) | Action: %s | Details: %s\n", username, userID, action, details)

	audit := models.AuditLog{
		Username: username, 
		Action:   action + " - " + details,
	}
	
	go func() {
		// Simpan log ke SQLite
		DB.Create(&audit)
	}()
}