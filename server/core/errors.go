package core

import "errors"

var (
	ErrInvitationNotFound   = errors.New("invitation not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrUsernameTaken        = errors.New("username taken")
	ErrForbidden            = errors.New("forbidden")
	ErrKeystoreNotFound     = errors.New("keystore not found")
	ErrInvalidUsername      = errors.New("invalid username")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrInvalidInvitee       = errors.New("invalid inviter")
	ErrInviterNotFound      = errors.New("inviter not found")
	ErrInviteeNotFound      = errors.New("invitee not found")
)
