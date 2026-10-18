package core

import (
	"github.com/pkg/errors"

	"myst/client/golang/models"
)

func (app *application) CreateInvitation(
	sessionId []byte,
	keystoreId string, inviteeUsername string) (models.Invitation, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return models.Invitation{}, ErrForbidden
	}

	k, err := app.enclave.Keystore(app.key, keystoreId)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get keystore")
	}

	// if keystore's remoteId is empty, then upload the keystore and save remoteId
	if k.RemoteId == "" {
		k, err = app.remote.CreateKeystore(k)
		if err != nil {
			return models.Invitation{}, errors.WithMessage(err, "failed to create keystore")
		}

		k, err = app.enclave.UpdateKeystore(app.key, k)
		if err != nil {
			return models.Invitation{}, errors.WithMessage(err, "failed to update keystore")
		}
	}

	inv, err := app.remote.CreateInvitation(k.RemoteId, inviteeUsername)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to create invitation")
	}

	return inv, nil
}

func (app *application) AcceptInvitation(sessionId []byte, id string) (models.Invitation, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return models.Invitation{}, ErrForbidden
	}

	inv, err := app.remote.AcceptInvitation(id)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to accept invitation")
	}

	return inv, nil
}

func (app *application) DeleteInvitation(sessionId []byte, id string) (models.Invitation, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return models.Invitation{}, ErrForbidden
	}

	inv, err := app.remote.DeleteInvitation(id)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to delete invitation")
	}

	return inv, nil
}

// FinalizeInvitation is the explicit call that the inviter makes to
// mark an invitation as finalized, calculate the shared encryption key
// using their private key and the invitee's public key, store the encryption
// key for future use (silent finalization) and finally encrypt the keystore
// key and attach it to the models.
func (app *application) FinalizeInvitation(sessionId []byte, invitationId, remoteKeystoreId string,
	inviteePublicKey []byte) (models.Invitation, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return models.Invitation{}, ErrForbidden
	}

	k, err := app.keystoreByRemoteId(remoteKeystoreId)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get keystore by remoteId")
	}

	creds, err := app.enclave.Credentials(app.key)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get credentials")
	}

	inv, err := app.remote.FinalizeInvitation(invitationId, creds.PrivateKey, inviteePublicKey, k)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to finalize invitation")
	}

	return inv, nil
}

func (app *application) Invitation(sessionId []byte, id string) (models.Invitation, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return models.Invitation{}, ErrForbidden
	}

	inv, err := app.remote.Invitation(id)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	return inv, nil
}

func (app *application) Invitations(sessionId []byte) (map[string]models.Invitation, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return nil, ErrForbidden
	}

	invs, err := app.remote.Invitations()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get invitations")
	}

	return invs, nil
}
