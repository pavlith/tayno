package storage

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"myst/server/core"
	"myst/server/models"
)

func (r *Repository) CreateUser(value models.User) (models.User, error) {
	err := r.db.Create(&value).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return models.User{}, core.ErrUsernameTaken
	}
	return value, err
}

func (r *Repository) User(id string) (models.User, error) {
	var value models.User
	err := r.db.Where("id = ?", id).First(&value).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return value, core.ErrUserNotFound
	}
	return value, err
}
func (r *Repository) Users() ([]models.User, error) {
	values := []models.User{}
	err := r.db.Order("created_at ASC, id ASC").Find(&values).Error
	return values, err
}
func (r *Repository) UpdateUser(value models.User) (models.User, error) {
	value.UpdatedAt = time.Now()
	result := r.db.Model(&models.User{}).Where("id = ?", value.Id).
		Select("Username", "PasswordHash", "PublicKey", "UpdatedAt").Updates(&value)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return models.User{}, core.ErrUsernameTaken
	}
	if result.Error != nil {
		return models.User{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.User{}, core.ErrUserNotFound
	}
	return value, nil
}

func (r *Repository) UserByUsername(username string) (models.User, error) {
	var value models.User
	err := r.db.Where("username = ?", username).First(&value).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return value, core.ErrUserNotFound
	}
	return value, err
}
