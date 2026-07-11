package repository

import (
	"go-perjalanan-dinas/models"

	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(department models.Department) (models.Department, error)
	FindAll() ([]models.Department, error)
	FindByID(id uint) (models.Department, error)
	Update(department models.Department) (models.Department, error)
	Delete(id uint) error
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db}
}

func (r *departmentRepository) Create(department models.Department) (models.Department, error) {
	err := r.db.Create(&department).Error
	return department, err
}

func (r *departmentRepository) FindAll() ([]models.Department, error) {
	var departments []models.Department
	err := r.db.Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) FindByID(id uint) (models.Department, error) {
	var department models.Department
	err := r.db.First(&department, id).Error
	return department, err
}

func (r *departmentRepository) Update(department models.Department) (models.Department, error) {
	err := r.db.Save(&department).Error
	return department, err
}

func (r *departmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Department{}, id).Error
}
