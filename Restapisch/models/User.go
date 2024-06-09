package models

type User struct {
	User_id  uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Level    int    `gorm:"default:2"` // Default level untuk pengguna baru
}
