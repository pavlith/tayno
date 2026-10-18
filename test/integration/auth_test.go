package integration_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/test/integration/suite"
)

func TestAuthentication(t *testing.T) {
	s := suite.New(t)

	sessionId, err := s.Client1.App.Initialize(s.Client1.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client1.App.Register(sessionId, s.Client1.Username, s.Client1.Password)
	assert.NilError(t, err)

	u, err := s.Client1.App.CurrentUser(sessionId)
	assert.NilError(t, err)
	assert.Equal(t, u.Username, s.Client1.Username)
}

func TestAccountLoginThenUnlock(t *testing.T) {
	s := suite.New(t)
	app := s.Client1.App
	token, err := app.SignupAccount(s.Client1.Username, s.Client1.Password, s.Client1.MasterPassword)
	assert.NilError(t, err)
	username, err := app.AccountStatus(token)
	assert.NilError(t, err)
	assert.Equal(t, username, s.Client1.Username)
	_, err = app.UnlockAccount("missing", s.Client1.MasterPassword)
	assert.Assert(t, err != nil)
	sid, err := app.UnlockAccount(token, s.Client1.MasterPassword)
	assert.NilError(t, err)
	_, err = app.CurrentUser(sid)
	assert.NilError(t, err)
	assert.NilError(t, app.LogoutAccount(token))
	_, err = app.AccountStatus(token)
	assert.Assert(t, err != nil)
	_, err = app.CurrentUser(sid)
	assert.Assert(t, err != nil)
	assert.Assert(t, app.HealthCheck(sid) != nil)
	_, err = app.LoginAccount(s.Client1.Username, "wrong-password")
	assert.Assert(t, err != nil)
	token, err = app.LoginAccount(s.Client1.Username, s.Client1.Password)
	assert.NilError(t, err)
	_, err = app.UnlockAccount(token, s.Client1.MasterPassword)
	assert.NilError(t, err)
	// A valid account from the same server cannot unlock another account's vault.
	other, err := s.Client2.App.SignupAccount(s.Client2.Username, s.Client2.Password, s.Client2.MasterPassword)
	assert.NilError(t, err)
	assert.Assert(t, other != "")
	other, err = app.LoginAccount(s.Client2.Username, s.Client2.Password)
	assert.NilError(t, err)
	_, err = app.UnlockAccount(other, s.Client1.MasterPassword)
	assert.Assert(t, err != nil)
}
