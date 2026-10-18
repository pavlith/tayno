package core

import (
	"myst/client/golang/models"
)

type Enclave interface {
	Initialize(password string) (key []byte, err error)
	IsInitialized() (bool, error)
	Authenticate(password string) (key []byte, err error)
	CreateKeystore(key []byte, k models.Keystore) (models.Keystore, error)
	Keystore(key []byte, id string) (models.Keystore, error)
	UpdateKeystore(key []byte, k models.Keystore) (models.Keystore, error)
	Keystores(key []byte) (map[string]models.Keystore, error)
	DeleteKeystore(key []byte, id string) error
	UpdateCredentials(key []byte, creds models.Credentials) (models.Credentials, error)
	Credentials(key []byte) (models.Credentials, error)
}

type Remote interface {
	VerifyCredentials(username, password string) (models.User, error)
	Address() string
	Authenticate(username, password string) error
	Register(username, password string, publicKey []byte) (models.User, error)
	Authenticated() bool
	CurrentUser() *models.User
	CreateKeystore(k models.Keystore) (models.Keystore, error)
	UpdateKeystore(k models.Keystore) (models.Keystore, error)
	Keystores(privateKey []byte) (map[string]models.Keystore, error)
	DeleteKeystore(id string) error
	CreateInvitation(keystoreRemoteId, inviteeUsername string) (models.Invitation, error)
	Invitation(id string) (models.Invitation, error)
	AcceptInvitation(id string) (models.Invitation, error)
	DeleteInvitation(id string) (models.Invitation, error)
	FinalizeInvitation(invitationId string, privateKey, inviteePublicKey []byte,
		k models.Keystore) (models.Invitation, error)
	Invitations() (map[string]models.Invitation, error)
	SharedSecret(privateKey []byte, userId string) ([]byte, error)
}
