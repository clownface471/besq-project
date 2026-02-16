package models

// Kita tidak menggunakan gorm.Model karena tabel ini tabel lama (legacy)
// yang mungkin tidak punya created_at, updated_at, dll.
type TmEmploy struct {
	IDEmploy int    `gorm:"column:IDEmploy;primaryKey" json:"id_employ"`
	NIK      string `gorm:"column:nik;primaryKey" json:"nik"`
	Nama     string `gorm:"column:nama" json:"nama"`
	Passw    string `gorm:"column:passw" json:"-"` // Password tidak dikirim di JSON response
	EmailAddr string `gorm:"column:emailAddr" json:"email_addr"` // Tambahkan kolom email
    
    // Opsional: Jika di tabel ada kolom jabatan/dept, bisa ditambahkan untuk mapping Role
	// Dept  string `gorm:"column:dept"` 
}

// Paksa GORM menggunakan nama tabel yang spesifik
func (TmEmploy) TableName() string {
	return "tm_employ"
}