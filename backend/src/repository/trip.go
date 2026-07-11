package repository

import (
	"go-perjalanan-dinas/models"

	"gorm.io/gorm"
)

type TripRepository interface {
	Save(trip models.BusinessTrip) (models.BusinessTrip, error)
	FindByID(id uint) (models.BusinessTrip, error)
	Update(trip models.BusinessTrip) (models.BusinessTrip, error)
	FindAll() ([]models.BusinessTrip, error)
	FindByUserID(userID uint) ([]models.BusinessTrip, error)
	FindByDepartmentID(departmentID uint) ([]models.BusinessTrip, error)
	Delete(id uint) error
	FindUsersByDepartmentID(departmentID uint) ([]models.User, error)
	FindAllUsers() ([]models.User, error)
}

type tripRepository struct {
	db *gorm.DB
}

func NewTripRepository(db *gorm.DB) TripRepository {
	return &tripRepository{db}
}

func (r *tripRepository) Save(trip models.BusinessTrip) (models.BusinessTrip, error) {
	err := r.db.Create(&trip).Error
	return trip, err
}

func (r *tripRepository) FindByID(id uint) (models.BusinessTrip, error) {
	var trip models.BusinessTrip
	err := r.db.Preload("User").First(&trip, id).Error
	return trip, err
}

func (r *tripRepository) Update(trip models.BusinessTrip) (models.BusinessTrip, error) {
	err := r.db.Save(&trip).Error
	if err != nil {
		return trip, err
	}

	// Preload "User" agar data nama karyawan yang dinas otomatis ikut terbawa
	var updatedTrip models.BusinessTrip
	err = r.db.Preload("User").First(&updatedTrip, trip.ID).Error
	return updatedTrip, err
}

func (r *tripRepository) FindAll() ([]models.BusinessTrip, error) {
	var trips []models.BusinessTrip
	// Preload "User" agar data nama karyawan yang dinas otomatis ikut terbawa
	err := r.db.Preload("User").Find(&trips).Error
	return trips, err
}

func (r *tripRepository) FindByUserID(userID uint) ([]models.BusinessTrip, error) {
	var trips []models.BusinessTrip
	err := r.db.Preload("User").Where("user_id = ?", userID).Order("created_at desc").Find(&trips).Error
	return trips, err
}

func (r *tripRepository) Delete(id uint) error {
	return r.db.Delete(&models.BusinessTrip{}, id).Error
}

func (r *tripRepository) FindByDepartmentID(departmentID uint) ([]models.BusinessTrip, error) {
	var trips []models.BusinessTrip
	err := r.db.Preload("User").
		Joins("JOIN users ON users.id = business_trips.user_id").
		Where("users.department_id = ?", departmentID).
		Order("business_trips.created_at desc").
		Find(&trips).Error
	return trips, err
}

func (r *tripRepository) FindUsersByDepartmentID(departmentID uint) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("department_id = ? AND role = ?", departmentID, "EMPLOYEE").Find(&users).Error
	return users, err
}

func (r *tripRepository) FindAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("role = ?", "EMPLOYEE").Find(&users).Error
	return users, err
}
