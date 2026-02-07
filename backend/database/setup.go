package database

import (
	"factory-api/models"
	"fmt"
	
	// --- DRIVER MYSQL (Untuk Statistik) ---
	"gorm.io/driver/mysql"
	
	// --- DRIVER SQLITE PURE GO (Untuk User/Login) ---
	// PENTING: Jangan pakai "gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite" 
	
	"gorm.io/gorm"
)

// DB Utama (SQLite) untuk User, Login, Cycle Button
var DB *gorm.DB

// DB Statistik (MySQL) untuk Grafik Laporan dari Server
var MySQL *gorm.DB

func ConnectDatabase() {
	// ==========================================
	// 1. KONEKSI SQLITE (Database Lokal Aplikasi)
	// ==========================================
	// Menggunakan driver glebarez/sqlite yang aman untuk Windows tanpa GCC
	sqliteDB, err := gorm.Open(sqlite.Open("besq.db"), &gorm.Config{})
	if err != nil {
		panic("Gagal koneksi ke SQLite: " + err.Error())
	}

	// Migrasi tabel aplikasi
	sqliteDB.AutoMigrate(&models.User{}, &models.AuditLog{}, &models.PerCycle{})
	DB = sqliteDB
	fmt.Println("✅ SQLite Connected (besq.db) - Pure Go Mode")

	// ==========================================
	// 2. KONEKSI MYSQL (Database Statistik Remote via VPN)
	// ==========================================
	
	// Ganti data di bawah ini sesuai Server Asli (HeidiSQL)
	dbUser := "root"            // Username DB Server
	dbPass := "am2bitnami"                // Password DB Server
	dbHost := "192.168.1.230"     // IP Address Server (JANGAN localhost/127.0.0.1 jika via VPN)
	dbPort := "3306"            // Port
	dbName := "besq_prd"    // Nama Database di Server

	// Format DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true&timeout=5s", 
		dbUser, dbPass, dbHost, dbPort, dbName)
	
	mysqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// Gunakan Println saja agar aplikasi TETAP JALAN walau VPN mati
		fmt.Println("⚠️  PERINGATAN: Gagal koneksi ke MySQL Statistik (Cek VPN/IP).")
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