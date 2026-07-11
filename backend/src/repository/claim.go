package repository

import (
	"go-perjalanan-dinas/models"

	"gorm.io/gorm"
)

type ReimbursementRepository interface {
	Create(claim models.Reimbursement) (models.Reimbursement, error)
	FindAll() ([]models.Reimbursement, error)
	FindByID(id uint) (models.Reimbursement, error)
	FindByTripID(tripID uint) ([]models.Reimbursement, error)
	Update(claim models.Reimbursement) (models.Reimbursement, error)
	VerifyOwnership(claimID uint, userID uint) (bool, error)
	Delete(id uint) error
}

type reimbursementRepository struct {
	db *gorm.DB
}

func NewReimbursementRepository(db *gorm.DB) ReimbursementRepository {
	return &reimbursementRepository{db: db}
}

func (r *reimbursementRepository) Create(claim models.Reimbursement) (models.Reimbursement, error) {
	err := r.db.Create(&claim).Error
	return claim, err
}

func (r *reimbursementRepository) FindAll() ([]models.Reimbursement, error) {
	var claims []models.Reimbursement
	err := r.db.Order("created_at desc").Find(&claims).Error
	return claims, err
}

func (r *reimbursementRepository) FindByID(id uint) (models.Reimbursement, error) {
	var claim models.Reimbursement
	err := r.db.First(&claim, id).Error
	return claim, err
}

func (r *reimbursementRepository) FindByTripID(tripID uint) ([]models.Reimbursement, error) {
	var claims []models.Reimbursement
	err := r.db.Where("trip_id = ?", tripID).Order("created_at desc").Find(&claims).Error
	return claims, err
}

func (r *reimbursementRepository) Update(claim models.Reimbursement) (models.Reimbursement, error) {
	err := r.db.Save(&claim).Error
	return claim, err
}

func (r *reimbursementRepository) VerifyOwnership(claimID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Table("reimbursements").
		Joins("JOIN business_trips ON business_trips.id = reimbursements.trip_id").
		Where("reimbursements.id = ? AND business_trips.user_id = ?", claimID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *reimbursementRepository) Delete(id uint) error {
	return r.db.Delete(&models.Reimbursement{}, id).Error
}
