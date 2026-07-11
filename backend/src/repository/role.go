package repository

import (
	"go-perjalanan-dinas/models"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role models.Role) (models.Role, error)
	FindAll() ([]models.Role, error)
	FindByID(id uint) (models.Role, error)
	FindByName(name string) (models.Role, error)
	Update(role models.Role) (models.Role, error)
	Delete(id uint) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) Create(role models.Role) (models.Role, error) {
	err := r.db.Create(&role).Error
	return role, err
}

func (r *roleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *roleRepository) FindByID(id uint) (models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	return role, err
}

func (r *roleRepository) FindByName(name string) (models.Role, error) {
	var role models.Role
	err := r.db.Where("UPPER(name) = UPPER(?)", name).First(&role).Error
	return role, err
}

func (r *roleRepository) Update(role models.Role) (models.Role, error) {
	err := r.db.Save(&role).Error
	return role, err
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}
