package core

import (
	"encoding/base64"
	"strings"
	"time"

	"myst/pkgs/rand"
)

type accountSession struct {
	username string
	expires  time.Time
	sessions []string
}

func (app *application) accountLocked(token string) (accountSession, error) {
	account, ok := app.accounts[token]
	if !ok || time.Now().After(account.expires) {
		app.logoutAccountLocked(token)
		return accountSession{}, ErrAuthenticationRequired
	}
	return account, nil
}

func (app *application) newAccountLocked(username string) (string, error) {
	bytes, err := rand.Bytes(32)
	if err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(bytes)
	if app.accounts == nil {
		app.accounts = make(map[string]accountSession)
	}
	app.accounts[token] = accountSession{username: username, expires: time.Now().Add(12 * time.Hour)}
	return token, nil
}

func (app *application) LoginAccount(username, password string) (string, error) {
	app.mux.Lock()
	defer app.mux.Unlock()
	user, err := app.remote.VerifyCredentials(strings.TrimSpace(username), password)
	if err != nil {
		return "", err
	}
	return app.newAccountLocked(user.Username)
}

func (app *application) AccountStatus(token string) (string, error) {
	app.mux.Lock()
	defer app.mux.Unlock()
	account, err := app.accountLocked(token)
	return account.username, err
}

func (app *application) LogoutAccount(token string) error {
	app.mux.Lock()
	defer app.mux.Unlock()
	app.logoutAccountLocked(token)
	return nil
}

func (app *application) logoutAccountLocked(token string) {
	for _, sid := range app.accounts[token].sessions {
		delete(app.sessions, sid)
	}
	delete(app.accounts, token)
	if len(app.sessions) == 0 {
		app.key = nil
	}
}

func (app *application) UnlockAccount(token, password string) ([]byte, error) {
	app.mux.Lock()
	defer app.mux.Unlock()
	account, err := app.accountLocked(token)
	if err != nil {
		return nil, err
	}
	sid, err := app.authenticateLocked(password, account.username)
	if err != nil {
		return nil, err
	}
	account.sessions = append(account.sessions, string(sid))
	app.accounts[token] = account
	return sid, nil
}

// SignupAccount preserves an existing unregistered vault and registers its key pair.
// A failed registration can be retried with the same master password.
func (app *application) SignupAccount(username, password, masterPassword string) (string, error) {
	app.mux.Lock()
	defer app.mux.Unlock()
	username = strings.TrimSpace(username)
	if username == "" {
		return "", ErrorInvalidUsername
	}
	if len(password) < 8 || len(masterPassword) < 8 {
		return "", ErrInvalidPassword
	}
	initialized, err := app.enclave.IsInitialized()
	if err != nil {
		return "", err
	}
	var sid []byte
	if initialized {
		key, err := app.enclave.Authenticate(masterPassword)
		if err != nil {
			return "", err
		}
		creds, err := app.enclave.Credentials(key)
		if err != nil {
			return "", err
		}
		if creds.Address != "" {
			return "", ErrEnclaveExists
		}
		app.key = key
		sid, err = app.newSession()
		if err != nil {
			return "", err
		}
	} else {
		sid, err = app.initializeLocked(masterPassword)
		if err != nil {
			return "", err
		}
	}
	defer func() {
		delete(app.sessions, string(sid))
		if len(app.sessions) == 0 {
			app.key = nil
		}
	}()
	if _, err := app.registerLocked(sid, username, password); err != nil {
		return "", err
	}
	return app.newAccountLocked(username)
}
