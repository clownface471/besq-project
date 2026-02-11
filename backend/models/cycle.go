// FILE: backend/models/cycle.go
package models

import "gorm.io/gorm"

type PerCycle struct {
	gorm.Model
	NoMC         string `json:"no_mc" binding:"required"`         // Nomor Mesin
	Item         string `json:"item" binding:"required"`          // Nama Item
	NoLot        string `json:"no_lot" binding:"required"`        // Nomor Lot
	StatusMesin  string `json:"status_mesin" binding:"required"`  // produksi, mati, rusak, reparasi
	NamaOperator string `json:"nama_operator"`                    // Nama Operator
}