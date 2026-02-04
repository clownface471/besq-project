package models

import (
	"time"
)

type WorkOrder struct {
	ID           uint      `gorm:"primaryKey"`
	OrderNumber  string    `gorm:"unique;not null"` // Contoh: WO-2026-001
	PartName     string    `gorm:"not null"`        // Nama spare part
	Quantity     int       `gorm:"not null"`        // Jumlah yang harus dibuat
	Status       string    `gorm:"default:'PENDING'"` // PENDING, CUTTING, PRESSING, DONE
	OperatorID   uint      // Siapa yang sedang mengerjakan
	CreatedAt    time.Time
	UpdatedAt    time.Time
}