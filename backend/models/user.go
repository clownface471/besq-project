package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // "-" supaya password tidak bocor di JSON
	Role     string `gorm:"not null" json:"role"` // ADMIN, CUTTING, PRESSING
}

// Struct untuk request body saat registrasi/tambah operator
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}