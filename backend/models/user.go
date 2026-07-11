package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           uint       `gorm:"primaryKey; autoIncrement" json:"id"`
	Name         string     `gorm:"type:varchar(100);not null" json:"name"`
	Email        string     `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password     string     `gorm:"type:varchar(255);not null" json:"-"`             // tag json:"-" agar password tidak ikut bocor saat di-render ke JSON
	Role         string     `gorm:"type:varchar(20);default:'EMPLOYEE'" json:"role"` // Denormalized untuk JWT & middleware
	RoleID       *uint      `gorm:"index" json:"role_id"`                            // Foreign key ke Role
	RoleData     Role       `gorm:"foreignKey:RoleID" json:"-"`                      // Relasi Belong To
	DepartmentID *uint      `json:"department_id"`
	Department   Department `gorm:"foreignKey:DepartmentID" json:"department"` // Relasi Belong To ke tabel Department
}
