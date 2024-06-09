package models

import "time"

type Pendaftar_sekolah struct {
	Id_pendaftaran uint      `gorm:"primaryKey:autoIncrement"`
	User_id        int       `gorm:"foreignKey:user_id"`
	Nama_lengkap   string    `gorm:"not null"`
	Asal_sekolah   string    `gorm:"not null"`
	Tp_lahir       string    `gorm:"not null"`
	Tgl_lahir      time.Time `gorm:"not null"`
	JK             string    `gorm:"type:ENUM('Laki-laki', 'Perempuan');default:'Laki-laki'"`
	Foto           string    `gorm:"not null"`
}
