package models

import (
	"time"

	"gorm.io/gorm"
)

type Reimbursement struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	TripID          uint           `gorm:"not null" json:"trip_id"`
	Trip            *BusinessTrip  `gorm:"foreignKey:TripID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"trip,omitempty"`
	Title           string         `gorm:"type:varchar(100);not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	Amount          float64        `gorm:"type:numeric(15,2);not null" json:"amount"`
	Status          string         `gorm:"type:varchar(30);default:'PENDING'" json:"status"`
	RejectedReason  string         `gorm:"type:varchar(255)" json:"rejected_reason,omitempty"`
	TransactionDate time.Time      `gorm:"type:date;not null" json:"transaction_date"` // Hanya menyimpan tanggal (YYYY-MM-DD)
	ReviewedAt      *time.Time     `json:"reviewed_at,omitempty"`
	AttachmentPath  string         `gorm:"type:text" json:"attachment_path,omitempty"`
	AttachmentPaths string         `gorm:"type:text" json:"attachment_paths,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
