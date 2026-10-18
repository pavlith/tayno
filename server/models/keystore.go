package models

import (
	"fmt"
	"time"

	"myst/pkgs/uuid"
)

type Keystore struct {
	Owner     *User  `gorm:"foreignKey:OwnerId;references:Id;constraint:OnDelete:RESTRICT" json:"-"`
	Id        string `gorm:"column:id;primaryKey"`
	Name      string
	Payload   []byte
	OwnerId   string `gorm:"column:owner_id;not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (k Keystore) String() string {
	return fmt.Sprintln(k.Id, k.Name, k.OwnerId, k.CreatedAt, k.UpdatedAt)
}

func NewKeystore(opts ...KeystoreOption) Keystore {
	k := Keystore{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&k)
		}
	}

	return k
}

type KeystoreOption func(*Keystore)

func WithName(name string) KeystoreOption {
	return func(k *Keystore) {
		k.Name = name
	}
}

func WithPayload(payload []byte) KeystoreOption {
	return func(k *Keystore) {
		k.Payload = payload
	}
}

func WithOwnerId(id string) KeystoreOption {
	return func(k *Keystore) {
		k.OwnerId = id
	}
}
