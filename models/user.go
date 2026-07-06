package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	Email    string `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"` // tag json:"-" agar password tidak ikut bocor saat di-render ke JSON
	Role     string `gorm:"type:varchar(20);default:'EMPLOYEE'" json:"role"` // Nilai: EMPLOYEE, MANAGER, atau HRD
}