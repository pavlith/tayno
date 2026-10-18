package remote

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"

	"myst/client/golang/core"
	"myst/client/golang/models"
	"myst/pkgs/crypto"
	"myst/server/apis/generated"
)

func (r *Remote) Invitation(id string) (models.Invitation, error) {
	if !r.Authenticated() {
		return models.Invitation{}, core.ErrAuthenticationRequired
	}

	res, err := r.client.GetInvitationWithResponse(context.Background(), id)
	if err != nil {
		return models.Invitation{}, errors.Wrap(err, "failed to get invitation")
	}

	if res.StatusCode() == http.StatusNotFound {
		return models.Invitation{}, core.ErrInvitationNotFound
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return models.Invitation{}, errors.New("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *Remote) CreateInvitation(keystoreRemoteId, inviteeUsername string,
) (models.Invitation, error) {
	if !r.Authenticated() {
		return models.Invitation{}, core.ErrAuthenticationRequired
	}

	res, err := r.client.CreateInvitationWithResponse(
		context.Background(), keystoreRemoteId, generated.CreateInvitationJSONRequestBody{
			Invitee: inviteeUsername,
		},
	)
	if err != nil {
		return models.Invitation{}, errors.Wrap(err, "failed to create invitation")
	}

	if res.StatusCode() != http.StatusCreated || res.JSON201 == nil {
		return models.Invitation{}, errors.New("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON201)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *Remote) Invitations() (map[string]models.Invitation, error) {
	if !r.Authenticated() {
		return nil, core.ErrAuthenticationRequired
	}

	res, err := r.client.InvitationsWithResponse(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get invitations")
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return nil, errors.New("invalid response")
	}

	invs := map[string]models.Invitation{}
	for _, inv := range *res.JSON200 {
		inv, err := InvitationFromJSON(inv)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse invitation")
		}

		invs[inv.Id] = inv
	}

	return invs, nil
}

func (r *Remote) AcceptInvitation(invitationId string) (models.Invitation, error) {
	res, err := r.client.AcceptInvitationWithResponse(
		context.Background(), invitationId,
	)
	if err != nil {
		return models.Invitation{}, errors.Wrap(err, "failed to accept invitation")
	}

	if res.StatusCode() == http.StatusNotFound {
		return models.Invitation{}, core.ErrInvitationNotFound
	} else if res.StatusCode() == http.StatusForbidden {
		return models.Invitation{}, core.ErrForbidden
	} else if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return models.Invitation{}, errors.New("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *Remote) DeleteInvitation(id string) (models.Invitation, error) {
	if !r.Authenticated() {
		return models.Invitation{}, core.ErrAuthenticationRequired
	}

	res, err := r.client.DeleteInvitationWithResponse(
		context.Background(), id,
	)
	if err != nil {
		return models.Invitation{}, errors.Wrap(err, "failed to delete invitation")
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return models.Invitation{}, errors.New("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *Remote) FinalizeInvitation(
	invitationId string, privateKey []byte, inviteePublicKey []byte, k models.Keystore,
) (models.Invitation, error) {
	// derive shared secret using the user's private key and the invitee's public key
	sharedSecret, err := curve25519.X25519(privateKey, inviteePublicKey)
	if err != nil {
		return models.Invitation{}, errors.Wrap(err, "failed to derive shared secret")
	}

	// encrypt the keystore key with the asymmetric key
	encryptedKeystoreKey, err := crypto.AES256CBC_Encrypt(sharedSecret, k.Key)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to encrypt keystore key")
	}

	res, err := r.client.FinalizeInvitationWithResponse(context.Background(),
		invitationId, generated.FinalizeInvitationJSONRequestBody{
			KeystoreKey: encryptedKeystoreKey,
		},
	)
	if err != nil {
		return models.Invitation{}, errors.Wrap(err, "failed to finalize invitation")
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return models.Invitation{}, errors.New("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}
