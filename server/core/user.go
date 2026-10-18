package core

import (
	"github.com/pkg/errors"

	"myst/pkgs/crypto"
	"myst/server/models"
)

// CreateUser creates a user with the given username, password and public key.
// If the username or password are invalid, ErrInvalidUsername or
// ErrInvalidPassword will be returned respectively. If a user with the given
// username already exists, ErrUsernameTaken will be returned.
func (app *application) CreateUser(username, password string, publicKey []byte) (models.User, error) {
	if username == "" {
		return models.User{}, ErrInvalidUsername
	}

	if password == "" {
		return models.User{}, ErrInvalidPassword
	}

	hash, err := crypto.HashPassword(password)
	if err != nil {
		return models.User{}, errors.WithMessage(err, "failed to hash password")
	}

	u, err := app.users.CreateUser(models.NewUser(
		models.WithUsername(username),
		models.WithPasswordHash(hash),
		models.WithPublicKey(publicKey),
	))
	if err != nil {
		return models.User{}, errors.WithMessage(err, "failed to create user")
	}

	return u, nil
}

// AuthenticateUser authenticates a user with the given username and password.
// If the user doesn't exist, or the password does not match,
// ErrAuthenticationFailed is returned.
func (app *application) AuthenticateUser(username, password string) (models.User, error) {
	u, err := app.users.UserByUsername(username)
	if errors.Is(err, ErrUserNotFound) {
		return models.User{}, ErrAuthenticationFailed
	} else if err != nil {
		return models.User{}, errors.WithMessage(err, "failed to get user by username")
	}

	ok, err := crypto.VerifyPassword(password, u.PasswordHash)
	if err != nil {
		return models.User{}, errors.WithMessage(err, "failed to verify password")
	}

	if !ok {
		return models.User{}, ErrAuthenticationFailed
	}

	return u, nil
}

// User returns a models.
func (app *application) User(userId string) (models.User, error) {
	u, err := app.users.User(userId)
	if err != nil {
		return models.User{}, errors.WithMessage(err, "failed to get user")
	}

	return u, nil
}

// UserByUsername returns a user by username.
func (app *application) UserByUsername(username string) (models.User, error) {
	u, err := app.users.UserByUsername(username)
	if err != nil {
		return models.User{}, errors.WithMessage(err, "failed to get user by username")
	}

	return u, nil
}

// UserInvitations returns a list of invitations for the given models.
// UserInvitationsOptions can be passed to filter out invitations by status.
// If the invitation is deleted, it is only included in the response if the user
// in question is the inviter.
func (app *application) UserInvitations(userId string, opts UserInvitationsOptions) ([]models.Invitation, error) {
	u, err := app.users.User(userId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user")
	}

	invs, err := app.invitations.Invitations()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get invitations")
	}

	invitations := []models.Invitation{}
	for _, inv := range invs {
		// skip invitations irrelevant to the user
		if inv.InviterId != u.Id && inv.InviteeId != u.Id {
			continue
		}

		// do not show deleted invitations if user is the invitee
		if inv.Deleted() && inv.InviteeId == u.Id {
			continue
		}

		// skip invitations not matching the status filter, if provided
		if opts.Status != nil && *opts.Status != inv.Status {
			continue
		}

		invitations = append(invitations, inv)
	}

	return invitations, nil
}

// UserInvitation returns an models.
// ErrUserNotFound will be returned if the user is not found.
// ErrInvitationNotFound will be returned if the invitation doesn't exist, if
// the user is not associated with it, or if the user is the invitee and
// the invitation is marked as deleted.
func (app *application) UserInvitation(userId, invitationId string) (models.Invitation, error) {
	u, err := app.users.User(userId)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get user")
	}

	i, err := app.invitations.Invitation(invitationId)
	if err != nil {
		return models.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	// do not show invitations irrelevant to the user
	if i.InviterId != u.Id && i.InviteeId != u.Id {
		return models.Invitation{}, ErrInvitationNotFound
	}

	// do not show deleted invitations if user is the invitee
	if i.Deleted() && i.InviteeId == u.Id {
		return models.Invitation{}, ErrInvitationNotFound
	}

	return i, nil
}
