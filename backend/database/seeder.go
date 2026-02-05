package database

import (
	"factory-api/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"gorm.io/gorm"
)

func SeedUsers() {
	// Daftar akun yang ingin dibuat otomatis
	users := []models.User{
		{
			Username: "superadmin",
			Password: "admin123",
			Role:     "ADMIN",
		},
		{
			Username: "op_cutting",
			Password: "user123",
			Role:     "OPERATOR_CUTTING",
		},
		{
			Username: "op_pressing",
			Password: "user123",
			Role:     "OPERATOR_PRESSING",
		},
	}

	for _, u := range users {
		var existingUser models.User
		// Cek apakah username sudah ada
		err := DB.Where("username = ?", u.Username).First(&existingUser).Error
		
		if err == gorm.ErrRecordNotFound {
			// Hash password
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			
			newUser := models.User{
				Username: u.Username,
				Password: string(hashedPassword),
				Role:     u.Role,
			}
			
			DB.Create(&newUser)
			log.Printf("Seeder: User created -> %s (%s)\n", u.Username, u.Role)
		} else {
			log.Printf("Seeder: User %s already exists, skipping...\n", u.Username)
		}
	}
}