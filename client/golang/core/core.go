package core

import (
	"sync"
	"time"

	"myst/client/golang/models"
	"myst/pkgs/logger"
)

var log = logger.New("app", logger.Blue)

type UpdateEntryOptions struct {
	Password *string
	Notes    *string
}

type Application interface {
	LoginAccount(username, password string) (string, error)
	SignupAccount(username, password, masterPassword string) (string, error)
	AccountStatus(token string) (string, error)
	LogoutAccount(token string) error
	UnlockAccount(token, password string) ([]byte, error)
	Initialize(password string) (sessionId []byte, err error)
	IsInitialized(sessionId []byte) (bool, error)
	Authenticate(password string) (sessionId []byte, err error)
	Logout(sessionId []byte) error
	HealthCheck(sessionId []byte) error

	CreateKeystore(sessionId []byte, name string) (models.Keystore, error)
	DeleteKeystore(sessionId []byte, id string) error
	Keystore(sessionId []byte, id string) (models.Keystore, error)
	Keystores(sessionId []byte) (map[string]models.Keystore, error)

	CreateEntry(sessionId []byte,
		keystoreId string, website, username, password, notes string) (models.Entry, error)
	UpdateEntry(sessionId []byte,
		keystoreId string, entryId string, opts UpdateEntryOptions) (models.Entry, error)
	DeleteEntry(sessionId []byte,
		keystoreId, entryId string) error

	Register(sessionId []byte, username, password string) (models.User, error)
	CurrentUser(sessionId []byte) (*models.User, error)
	SharedSecret(sessionId []byte, userId string) ([]byte, error)

	CreateInvitation(sessionId []byte,
		keystoreId string, inviteeUsername string) (models.Invitation, error)
	AcceptInvitation(sessionId []byte, id string) (models.Invitation, error)
	DeleteInvitation(sessionId []byte, id string) (models.Invitation, error)
	FinalizeInvitation(sessionId []byte,
		invitationId, remoteKeystoreId string, inviteePublicKey []byte) (models.Invitation, error)
	Invitations(sessionId []byte) (map[string]models.Invitation, error)
	Invitation(sessionId []byte, id string) (models.Invitation, error)

	Sync() error
}

type application struct {
	enclave  Enclave
	remote   Remote
	sessions map[string]time.Time
	accounts map[string]accountSession
	key      []byte
	mux      sync.Mutex
}

func New(opts ...Option) Application {
	app := &application{
		sessions: make(map[string]time.Time),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(app)
		}
	}

	go app.startHealthCheck()

	return app
}

type Option func(app *application)

func WithRemote(remote Remote) Option {
	return func(app *application) {
		app.remote = remote
	}
}

func WithEnclave(enclave Enclave) Option {
	return func(app *application) {
		app.enclave = enclave
	}
}
