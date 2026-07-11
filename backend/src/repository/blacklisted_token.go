package repository

import (
	"go-perjalanan-dinas/models"
	"time"

	"gorm.io/gorm"
)

type BlacklistedTokenRepository interface {
	Add(token string, expiresAt time.Time) error
	IsBlacklisted(token string) (bool, error)
}

type blacklistedTokenRepository struct {
	db *gorm.DB
}

func NewBlacklistedTokenRepository(db *gorm.DB) BlacklistedTokenRepository {
	return &blacklistedTokenRepository{db: db}
}

func (r *blacklistedTokenRepository) Add(token string, expiresAt time.Time) error {
	if r.db == nil {
		return nil
	}

	record := models.BlacklistedToken{Token: token, ExpiresAt: expiresAt}
	return r.db.Where(models.BlacklistedToken{Token: token}).Assign(models.BlacklistedToken{ExpiresAt: expiresAt}).FirstOrCreate(&record).Error
}

func (r *blacklistedTokenRepository) IsBlacklisted(token string) (bool, error) {
	if r.db == nil || token == "" {
		return false, nil
	}

	var count int64
	err := r.db.Model(&models.BlacklistedToken{}).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
