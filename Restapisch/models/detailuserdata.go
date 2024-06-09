package models

type Detail_pendaftar struct {
	User_id         int    `gorm:"foreignKey:user_id"`
	Nama_ibu        string `gorm:"not null"`
	Nama_ayah       string `gorm:"not null"`
	Alamat_orangtua string `gorm:"not null"`
	Pekerjaan_ibu   string `gorm:"not null"`
	Pekerjaan_ayah  string `gorm:"not null"`
	Alamat_siswa    string `gorm:"not null"`
	Tlp_orangtua    string `gorm:"not null"`
	Status          string `gorm:"type:ENUM('terdaftar','belum selesai');default:'terdaftar'"`
}
