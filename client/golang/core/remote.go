package core

import (
	"github.com/pkg/errors"

	"myst/client/golang/models"
)

func (app *application) Register(sessionId []byte, username, password string) (models.User, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	return app.registerLocked(sessionId, username, password)
}

func (app *application) registerLocked(sessionId []byte, username, password string) (models.User, error) {
	if !app.sessionActive(sessionId) {
		return models.User{}, ErrForbidden
	}

	var mustInit bool

	creds, err := app.enclave.Credentials(app.key)
	if errors.Is(err, ErrCredentialsNotFound) {
		mustInit = true
	} else if err != nil {
		return models.User{}, errors.WithMessage(err, "failed to query credentials")
	}

	if !mustInit && creds.Address != "" && creds.Address != app.remote.Address() {
		return models.User{}, ErrRemoteAddressMismatch
	}

	u, err := app.remote.Register(username, password, creds.PublicKey)
	if err != nil {
		return models.User{}, errors.WithMessage(err, "failed to register user")
	}

	creds.Address = app.remote.Address()
	creds.Username = username
	creds.Password = password

	_, err = app.enclave.UpdateCredentials(app.key, creds)
	if err != nil {
		return models.User{}, errors.WithMessage(err, "failed to update credentials")
	}

	return u, nil
}

func (app *application) CurrentUser(sessionId []byte) (*models.User, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return nil, ErrForbidden
	}

	rem, err := app.enclave.Credentials(app.key)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query credentials")
	}

	if rem.Address == "" {
		return nil, ErrCredentialsNotFound
	}

	u := app.remote.CurrentUser()
	if u == nil {
		return nil, nil
	}

	u.PublicKey = rem.PublicKey

	return u, nil
}

func (app *application) SharedSecret(sessionId []byte, userId string) ([]byte, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return nil, ErrForbidden
	}

	creds, err := app.enclave.Credentials(app.key)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query credentials")
	}

	return app.remote.SharedSecret(creds.PrivateKey, userId)
}
