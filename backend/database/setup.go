package database

import (
	"factory-api/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    // Pastikan DSN sesuai dengan database kamu
	// Format Template:
	// dsn := "USER:PASSWORD@tcp(HOST:PORT)/NAMA_DB?KONFIGURASI"
	dsn := "root:@tcp(127.0.0.1:3306)/besq_factory?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Gagal koneksi ke MySQL: " + err.Error())
	}

	database.AutoMigrate(&models.User{}, &models.AuditLog{}, &models.PerCycle{})

	DB = database
}

// --- UPDATE FUNGSI INI ---
// Sekarang menerima 4 parameter: userID, username, action, details
func RecordActivity(userID uint, username string, action string, details string) {
	
	// Log ke console untuk debugging
	fmt.Printf("[AUDIT] User: %s (ID: %d) | Action: %s | Details: %s\n", username, userID, action, details)

	// Simpan ke Database
	// Catatan: Pastikan model AuditLog kamu punya field 'Details' atau gabungkan jika tidak
	audit := models.AuditLog{
		Username: username, 
		Action:   action + " - " + details, // Kita gabung detail ke action jika kolom terbatas
	}
	
	// Gunakan Goroutine agar tidak memblokir proses utama (Fire and Forget)
	go func() {
		DB.Create(&audit)
	}()
}