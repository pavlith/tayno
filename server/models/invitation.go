package models

import (
	"fmt"
	"time"

	"myst/pkgs/uuid"
)

type Invitation struct {
	Keystore             *Keystore `gorm:"foreignKey:KeystoreId;references:Id;constraint:OnDelete:CASCADE" json:"-"`
	Inviter              *User     `gorm:"foreignKey:InviterId;references:Id;constraint:OnDelete:RESTRICT" json:"-"`
	Invitee              *User     `gorm:"foreignKey:InviteeId;references:Id;constraint:OnDelete:RESTRICT" json:"-"`
	Id                   string    `gorm:"column:id;primaryKey"`
	KeystoreId           string    `gorm:"column:keystore_id;not null;index"`
	InviterId            string    `gorm:"column:inviter_id;not null;index"`
	InviteeId            string    `gorm:"column:invitee_id;not null;index"`
	EncryptedKeystoreKey []byte
	Status               Status `gorm:"not null;index;check:status IN ('pending','deleted','declined','accepted','cancelled','finalized')"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
	DeclinedAt           time.Time
	AcceptedAt           time.Time
	CancelledAt          time.Time
	FinalizedAt          time.Time
}

func NewInvitation(opts ...InvitationOption) Invitation {
	k := Invitation{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    Pending,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&k)
		}
	}

	return k
}

type InvitationOption func(i *Invitation)

func WithKeystoreId(id string) InvitationOption {
	return func(i *Invitation) {
		i.KeystoreId = id
	}
}

func WithInviterId(id string) InvitationOption {
	return func(i *Invitation) {
		i.InviterId = id
	}
}

func WithInviteeId(id string) InvitationOption {
	return func(i *Invitation) {
		i.InviteeId = id
	}
}

func (i Invitation) String() string {
	return fmt.Sprintf(
		"Invitation{Id: %s, KeystoreId: %s, InviterId: %s, InviteeId: %s, EncryptedKeystoreKey: %s, Status: %s}",
		i.Id,
		i.KeystoreId, i.InviterId, i.InviteeId,
		i.EncryptedKeystoreKey,
		i.Status,
	)
}

func (i Invitation) Pending() bool {
	return i.Status == Pending
}

func (i Invitation) Accepted() bool {
	return i.Status == Accepted
}

func (i Invitation) Deleted() bool {
	return i.Status == Deleted
}

func (i Invitation) Declined() bool {
	return i.Status == Declined
}

func (i Invitation) Finalized() bool {
	return i.Status == Finalized
}
