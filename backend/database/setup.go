package database

import (
	"factory-api/models"
	"github.com/glebarez/sqlite" // GANTI KE DRIVER PURE GO INI
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Menggunakan driver glebarez/sqlite agar tidak butuh CGO/GCC di Windows
	database, err := gorm.Open(sqlite.Open("factory.db"), &gorm.Config{})
	if err != nil {
		// Kita cetak error aslinya supaya tahu kenapa gagal
		panic("Gagal koneksi ke database: " + err.Error())
	}

	database.AutoMigrate(&models.User{}, &models.AuditLog{}, &models.WorkOrder{}, &models.PerCycle{})
	DB = database
}

func RecordActivity(userID uint, username, action, endpoint string) {
	log := models.AuditLog{
		UserID:    userID,
		Username:  username,
		Action:    action,
		Endpoint:  endpoint,
		CreatedAt: time.Now(),
	}
	DB.Create(&log)
}