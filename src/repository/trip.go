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
	err := r.db.First(&trip, id).Error
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