package core

import "myst/server/models"

type InvitationRepository interface {
	CreateInvitation(models.Invitation) (models.Invitation, error)
	Invitation(id string) (models.Invitation, error)
	UpdateInvitation(models.Invitation) (models.Invitation, error)
	DeleteInvitation(id string) error
	Invitations() ([]models.Invitation, error)
}

type KeystoreRepository interface {
	CreateKeystore(models.Keystore) (models.Keystore, error)
	Keystore(id string) (models.Keystore, error)
	UpdateKeystore(k models.Keystore) (models.Keystore, error)
	// DeleteKeystore atomically removes the keystore and associated invitations.
	DeleteKeystore(id string) error
	Keystores() ([]models.Keystore, error)
}

type UserRepository interface {
	CreateUser(models.User) (models.User, error)
	User(id string) (models.User, error)
	UserByUsername(username string) (models.User, error)
	UpdateUser(models.User) (models.User, error)
	Users() ([]models.User, error)
}
