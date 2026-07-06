package models

import (
	"time"
	"gorm.io/gorm"
)

type BusinessTrip struct {
	gorm.Model
	UserID         uint      `json:"user_id"`
	User           User      `json:"user" gorm:"foreignKey:UserID"` // Relasi Belong To ke tabel User
	Destination    string    `gorm:"type:varchar(255);not null" json:"destination"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Status         string    `gorm:"type:varchar(20);default:'PENDING'" json:"status"` // PENDING, APPROVED, REJECTED, ON_DUTY, COMPLETED
	AttachmentPath string    `gorm:"type:text" json:"attachment_path,omitempty"`      // Lokasi path file laporan lokal
}