package remote

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"

	"myst/client/golang/core"
	"myst/client/golang/models"
	"myst/server/apis/generated"
)

// Reauthenticate authenticates with the previously provided username and
// password, if the Authenticate method has been called in the past.
func (r *Remote) Reauthenticate() error {
	if r.username == "" || r.password == "" {
		return nil
	}

	err := r.authenticate()
	if err != nil {
		return errors.WithMessage(err, "failed to reauthenticate")
	}

	return nil
}

// VerifyCredentials checks an account without changing the active vault's remote identity.
func (r *Remote) VerifyCredentials(username, password string) (models.User, error) {
	probe, err := New(WithAddress(r.address))
	if err != nil {
		return models.User{}, err
	}
	if err := probe.Authenticate(username, password); err != nil {
		return models.User{}, err
	}
	return *probe.CurrentUser(), nil
}

func (r *Remote) Authenticate(username, password string) error {
	// store for future re-authentication
	r.username = username
	r.password = password

	err := r.authenticate()
	if err != nil {
		return errors.WithMessage(err, "failed to authenticate")
	}

	return nil
}

func (r *Remote) authenticate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client.LoginWithResponse(
		ctx, generated.LoginJSONRequestBody{
			Username: r.username,
			Password: r.password,
		},
	)
	if err != nil {
		return errors.Wrap(err, "login failed")
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return core.ErrAuthenticationFailed
	}

	if res.JSON200.Token == "" {
		return errors.New("invalid token")
	}

	u := models.User{
		Id:        res.JSON200.User.Id,
		Username:  res.JSON200.User.Username,
		PublicKey: res.JSON200.User.PublicKey,
	}

	r.user = &u
	r.bearerToken = res.JSON200.Token

	return nil
}

func (r *Remote) Authenticated() bool {
	return r.bearerToken != ""
}

func (r *Remote) CurrentUser() *models.User {
	return r.user
}

func (r *Remote) Register(username, password string, publicKey []byte) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client.RegisterWithResponse(
		ctx, generated.RegisterJSONRequestBody{
			Username:  username,
			Password:  password,
			PublicKey: publicKey,
		},
	)
	if err != nil {
		return models.User{}, errors.Wrap(err, "failed to register")
	}

	if res.StatusCode() != http.StatusCreated || res.JSON201 == nil {
		return models.User{}, errors.New("invalid response")
	}

	if res.JSON201.Token == "" {
		return models.User{}, errors.New("invalid token")
	}

	u := models.User{
		Id:        res.JSON201.User.Id,
		Username:  res.JSON201.User.Username,
		PublicKey: publicKey,
	}

	r.user = &u
	r.bearerToken = res.JSON201.Token

	return u, nil
}

func (r *Remote) SharedSecret(privateKey []byte, userId string) ([]byte, error) {
	invs, err := r.Invitations()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query invitations")
	}

	for _, inv := range invs {
		if inv.Inviter.Id != userId && inv.Invitee.Id != userId {
			continue
		}

		var pub []byte
		if inv.Inviter.Id == r.CurrentUser().Id {
			pub = inv.Invitee.PublicKey
		} else {
			pub = inv.Inviter.PublicKey
		}

		sharedSecret, err := curve25519.X25519(privateKey, pub)
		if err != nil {
			return nil, errors.Wrap(err, "failed to derive shared secret")
		}

		return sharedSecret, nil
	}

	return nil, core.ErrSharedSecretNotFound
}
