package models

import (
	"encoding/hex"
	"fmt"
	"time"

	"myst/pkgs/uuid"
)

type User struct {
	Id           string `gorm:"column:id;primaryKey"`
	Username     string `gorm:"not null;uniqueIndex"`
	PasswordHash string `gorm:"not null"`
	PublicKey    []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u User) String() string {
	return fmt.Sprintln(u.Id, u.Username, "****", hex.EncodeToString(u.PublicKey))
}

func NewUser(opts ...UserOption) User {
	u := User{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&u)
		}
	}

	return u
}
