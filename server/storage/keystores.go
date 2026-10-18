package storage

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"myst/server/core"
	"myst/server/models"
)

func (r *Repository) CreateKeystore(value models.Keystore) (models.Keystore, error) {
	err := r.db.Create(&value).Error
	return value, err
}

func (r *Repository) Keystore(id string) (models.Keystore, error) {
	var value models.Keystore
	err := r.db.Where("id = ?", id).First(&value).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return value, core.ErrKeystoreNotFound
	}
	return value, err
}
func (r *Repository) Keystores() ([]models.Keystore, error) {
	values := []models.Keystore{}
	err := r.db.Order("created_at ASC, id ASC").Find(&values).Error
	return values, err
}
func (r *Repository) UpdateKeystore(value models.Keystore) (models.Keystore, error) {
	value.UpdatedAt = time.Now()
	result := r.db.Model(&models.Keystore{}).Where("id = ?", value.Id).
		Select("Name", "Payload", "OwnerId", "UpdatedAt").Updates(&value)
	if result.Error != nil {
		return models.Keystore{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.Keystore{}, core.ErrKeystoreNotFound
	}
	return value, nil
}

func (r *Repository) DeleteKeystore(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("keystore_id = ?", id).Delete(&models.Invitation{}).Error; err != nil {
			return err
		}
		result := tx.Where("id = ?", id).Delete(&models.Keystore{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return core.ErrKeystoreNotFound
		}
		return nil
	})
}
