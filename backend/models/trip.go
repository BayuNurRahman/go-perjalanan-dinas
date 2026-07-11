package models

import (
	"time"

	"gorm.io/gorm"
)

type BusinessTrip struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	UserID          uint           `json:"user_id"`
	User            User           `json:"user" gorm:"foreignKey:UserID"` // Relasi Belong To ke tabel User
	Description     string         `gorm:"type:text;not null" json:"description"`
	Destination     string         `gorm:"type:varchar(255);not null" json:"destination"`
	StartDate       time.Time      `json:"start_date"`
	EndDate         time.Time      `json:"end_date"`
	Initiator       string         `gorm:"type:varchar(100);not null" json:"initiator"` // Nama orang yang mengajukan perjalanan dinas
	Summary         string         `gorm:"type:text" json:"summary,omitempty"`           // Ringkasan perjalanan dinas
	Nomor_Surat     string         `gorm:"type:varchar(50);not null" json:"nomor_surat,omitempty"` // Nomor surat perjalanan dinas
	Status          string         `gorm:"type:varchar(20);default:'PENDING'" json:"status"` // PENDING, APPROVED, REJECTED, ON_DUTY, COMPLETED
	AttachmentPath  string         `gorm:"type:text" json:"attachment_path,omitempty"`       // Lokasi path file laporan lokal
	AttachmentPaths string         `gorm:"type:text" json:"attachment_paths,omitempty"`      // Daftar path file yang diunggah untuk perjalanan dinas
	Notes           string         `gorm:"type:text" json:"notes,omitempty"`                 // Catatan tambahan untuk claim atau perubahan
}
