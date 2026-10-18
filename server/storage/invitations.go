package storage

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"myst/server/core"
	"myst/server/models"
)

func (r *Repository) CreateInvitation(value models.Invitation) (models.Invitation, error) {
	err := r.db.Create(&value).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return models.Invitation{}, core.ErrForbidden
	}
	return value, err
}

func (r *Repository) Invitation(id string) (models.Invitation, error) {
	var value models.Invitation
	err := r.db.Where("id = ?", id).First(&value).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return value, core.ErrInvitationNotFound
	}
	return value, err
}
func (r *Repository) Invitations() ([]models.Invitation, error) {
	values := []models.Invitation{}
	err := r.db.Order("created_at ASC, id ASC").Find(&values).Error
	return values, err
}
func (r *Repository) UpdateInvitation(value models.Invitation) (models.Invitation, error) {
	value.UpdatedAt = time.Now()
	result := r.db.Model(&models.Invitation{}).Where("id = ?", value.Id).
		Select("KeystoreId", "InviterId", "InviteeId", "EncryptedKeystoreKey", "Status", "AcceptedAt", "DeclinedAt", "DeletedAt", "CancelledAt", "FinalizedAt", "UpdatedAt").Updates(&value)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return models.Invitation{}, core.ErrForbidden
	}
	if result.Error != nil {
		return models.Invitation{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.Invitation{}, core.ErrInvitationNotFound
	}
	return value, nil
}

func (r *Repository) DeleteInvitation(id string) error {
	result := r.db.Where("id = ?", id).Delete(&models.Invitation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return core.ErrInvitationNotFound
	}
	return nil
}
