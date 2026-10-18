package core

import "myst/server/models"

type UpdateKeystoreOptions struct {
	Name    *string
	Payload *[]byte
}

type UserInvitationsOptions struct {
	Status *models.Status
}

type Application interface {
	CreateUser(username, password string, publicKey []byte) (models.User, error)
	AuthenticateUser(username, password string) (models.User, error)
	User(id string) (models.User, error)
	UserByUsername(username string) (models.User, error)

	CreateKeystore(name, ownerId string, payload []byte) (models.Keystore, error)
	Keystore(id string) (models.Keystore, error)
	UserKeystore(userId, keystoreId string) (models.Keystore, error)
	Keystores() ([]models.Keystore, error)
	UpdateKeystore(userId, keystoreId string, opts UpdateKeystoreOptions) (models.Keystore, error)
	DeleteKeystore(userId, keystoreId string) error
	UserKeystores(userId string) ([]models.Keystore, error)

	CreateInvitation(keystoreId, inviterId, inviteeUsername string) (models.Invitation, error)
	AcceptInvitation(userId, invitationId string) (models.Invitation, error)
	DeleteInvitation(userId, invitationId string) (models.Invitation, error)
	FinalizeInvitation(usrId, invId string, encKeystoreKey []byte) (models.Invitation, error)
	UserInvitation(userId, invitationId string) (models.Invitation, error)
	UserInvitations(userId string, opts UserInvitationsOptions) ([]models.Invitation, error)
}

type application struct {
	invitations InvitationRepository
	users       UserRepository
	keystores   KeystoreRepository
}

func New(opts ...Option) Application {
	app := &application{}

	for _, opt := range opts {
		if opt != nil {
			opt(app)
		}
	}

	return app
}

type Option func(app *application)

func WithKeystoreRepository(repo KeystoreRepository) Option {
	return func(app *application) {
		app.keystores = repo
	}
}

func WithUserRepository(repo UserRepository) Option {
	return func(app *application) {
		app.users = repo
	}
}

func WithInvitationRepository(repo InvitationRepository) Option {
	return func(app *application) {
		app.invitations = repo
	}
}
