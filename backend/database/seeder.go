package database

import (
	"factory-api/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func SeedAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "ADMIN").Count(&count)

	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.User{
			Username: "superadmin",
			Password: string(hashedPassword),
			Role:     "ADMIN",
		}
		DB.Create(&admin)
		log.Println("Seeder: Admin user created (superadmin:admin123)")
	} else {
		log.Println("Seeder: Admin already exists, skipping...")
	}
}