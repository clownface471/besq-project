package models

import (
	"time"
)

type AuditLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Action    string    `json:"action"`   // misal: "START_MACHINE_CUTTING"
	Endpoint  string    `json:"endpoint"` // misal: "/production/cutting/start"
	CreatedAt time.Time `json:"created_at"`
}