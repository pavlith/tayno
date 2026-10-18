package core

import (
	"time"

	"github.com/pkg/errors"

	"myst/server/models"
)

// CreateInvitation creates an invitation for the given keystoreId, from the
// given inviter to the invitee. The inviter should be the owner of the
// keystore.
// Errors returned:
//   - ErrKeystoreNotFound if the keystore is not found.
//   - ErrInviterNotFound if the inviter is not found.
//   - ErrInviteeNotFound if the invitee is not found.
//   - ErrForbidden if the inviter is not the owner of the keystore, or an
//     invitation is already active for the given keystore, inviter, and invitee
//   - ErrInvalidInvitee if the inviter and invitee are the same user.
//
// TODO: return not found if the user does not have read access to the keystore
func (app *application) CreateInvitation(
	keystoreId, inviterId, inviteeUsername string) (models.Invitation, error) {
	inviter, err := app.users.User(inviterId)
	if errors.Is(err, ErrUserNotFound) {
		return models.Invitation{}, ErrInviterNotFound
	} else if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get inviter by id")
	}

	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get keystore by id")
	}

	if k.OwnerId != inviter.Id {
		return models.Invitation{}, ErrForbidden
	}

	invitee, err := app.users.UserByUsername(inviteeUsername)
	if errors.Is(err, ErrUserNotFound) {
		return models.Invitation{}, ErrInviteeNotFound
	} else if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get invitee by username")
	}

	if inviter.Id == invitee.Id {
		return models.Invitation{}, ErrInvalidInvitee
	}

	invs, err := app.invitations.Invitations()
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get invitations")
	}

	for _, inv := range invs {
		if inv.InviterId == inviter.Id && inv.InviteeId == invitee.Id && inv.KeystoreId == k.Id &&
			(inv.Status != models.Cancelled && inv.Status != models.Deleted && inv.Status != models.
				Declined) {
			return models.Invitation{}, ErrForbidden
		}
	}

	inv := models.NewInvitation(
		models.WithKeystoreId(k.Id),
		models.WithInviterId(inviter.Id),
		models.WithInviteeId(invitee.Id),
	)

	inv, err = app.invitations.CreateInvitation(inv)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to create invitation")
	}

	return inv, nil
}

// AcceptInvitation allows a user to accept an models.
// If the user doesn't exist, ErrUserNotFound is returned.
// ErrInvitationNotFound is returned if the invitation doesn't exist, or if the
// user does not have access to the invitation (e.g. they are the invitee and
// the invitation is marked as deleted).
// A user can only accept an invitation if they are the invitee and the
// invitation's status is pending. If either is not true then ErrForbidden is
// returned.
func (app *application) AcceptInvitation(
	userId string, invitationId string) (models.Invitation, error) {
	inv, err := app.UserInvitation(userId, invitationId)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get user invitation")
	}

	if userId != inv.InviteeId {
		return models.Invitation{}, ErrForbidden
	}

	if inv.Status != models.Pending {
		return models.Invitation{}, ErrForbidden
	}

	inv.Status = models.Accepted
	inv.AcceptedAt = time.Now()

	inv, err = app.invitations.UpdateInvitation(inv)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to update invitation")
	}

	return inv, nil
}

// DeleteInvitation deletes, declines or cancels an invitation, depending on
// the type of association of the user against the invitation (inviter:
// deletes/cancels, invitee: declines) and the invitation's current status.
// ErrUserNotFound is returned if the user doesn't exist.
// ErrInvitationNotFound is returned if the invitation doesn't exist, or if the
// user does not have access to the invitation (e.g. they are the invitee and
// the invitation is marked as deleted).
// ErrForbidden is returned if the invitation is in a state where the operation
// cannot be performed. See deleteInvitation, cancelInvitation and
// declineInvitation for more details on the respective conditions.
func (app *application) DeleteInvitation(
	userId, invitationId string) (models.Invitation, error) {
	inv, err := app.invitations.Invitation(invitationId)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	if userId != inv.InviterId && userId != inv.InviteeId {
		return models.Invitation{}, ErrInvitationNotFound
	}

	if userId == inv.InviterId {
		inv, err = app.deleteOrCancelInvitation(inv)
		if err != nil {
			return models.Invitation{}, errors.WithMessage(err,
				"failed to delete or cancel invitation")
		}
	} else if userId == inv.InviteeId {
		inv, err = app.declineInvitation(inv)
		if err != nil {
			return models.Invitation{}, errors.WithMessage(err,
				"failed to decline invitation")
		}
	}

	return inv, nil
}

func (app *application) deleteOrCancelInvitation(inv models.Invitation,
) (models.Invitation, error) {
	if inv.Pending() {
		inv.Status = models.Deleted
		inv.DeletedAt = time.Now()
	} else if inv.Accepted() {
		inv.Status = models.Cancelled
		inv.CancelledAt = time.Now()
	} else {
		return models.Invitation{}, ErrForbidden
	}

	var err error
	inv, err = app.invitations.UpdateInvitation(inv)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to update invitation")
	}

	return inv, nil
}

func (app *application) declineInvitation(inv models.Invitation) (models.Invitation, error) {
	if !inv.Pending() {
		return models.Invitation{}, ErrForbidden
	}

	inv.Status = models.Declined
	inv.DeclinedAt = time.Now()

	var err error
	inv, err = app.invitations.UpdateInvitation(inv)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to update invitation")
	}

	return inv, nil
}

// FinalizeInvitation finalizes an invitation, setting the encrypted keystore
// key for the keystore associated with the invitation, allowing decryption of
// its payload by the invitee. If the user doesn't have access to the
// invitation, ErrInvitationNotFound is returned. ErrForbidden is returned if
// the user is not the inviter or the invitation is not in the accepted state.
func (app *application) FinalizeInvitation(
	userId, invitationId string, encryptedKeystoreKey []byte) (models.Invitation, error) {
	inv, err := app.invitations.Invitation(invitationId)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	if userId != inv.InviterId && userId != inv.InviteeId {
		return models.Invitation{}, ErrInvitationNotFound
	}

	if userId != inv.InviterId {
		return models.Invitation{}, ErrForbidden
	}

	if !inv.Accepted() {
		return models.Invitation{}, ErrForbidden
	}

	inv.EncryptedKeystoreKey = encryptedKeystoreKey
	inv.Status = models.Finalized
	inv.FinalizedAt = time.Now()

	inv, err = app.invitations.UpdateInvitation(inv)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to update invitation")
	}

	return inv, nil
}
