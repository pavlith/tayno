package models

import (
	"time"

	"myst/pkgs/uuid"
)

type Entry struct {
	Id        string
	Website   string
	Username  string
	Password  string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEntry(opts ...EntryOption) Entry {
	e := Entry{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&e)
		}
	}

	return e
}

type EntryOption func(*Entry)

func WithWebsite(website string) EntryOption {
	return func(e *Entry) {
		e.Website = website
	}
}

func WithUsername(username string) EntryOption {
	return func(e *Entry) {
		e.Username = username
	}
}

func WithPassword(password string) EntryOption {
	return func(e *Entry) {
		e.Password = password
	}
}

func WithNotes(notes string) EntryOption {
	return func(e *Entry) {
		e.Notes = notes
	}
}
