package core

import (
	"time"

	"github.com/pkg/errors"

	"myst/pkgs/rand"
)

func (app *application) newSession() (sessionId []byte, err error) {
	sessionKey, err := rand.Bytes(16)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate session key")
	}

	app.sessions[string(sessionKey)] = time.Now()

	return sessionKey, nil
}

func (app *application) sessionActive(sessionId []byte) bool {
	if sessionId == nil {
		return false
	}
	_, ok := app.sessions[string(sessionId)]
	return ok
}

func (app *application) Logout(sessionId []byte) error {
	app.mux.Lock()
	defer app.mux.Unlock()
	delete(app.sessions, string(sessionId))
	if len(app.sessions) == 0 {
		app.key = nil
	}
	return nil
}

func (app *application) startHealthCheck() {
	// TODO @rdnt: make this configurable
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		app.mux.Lock()

		for token, account := range app.accounts {
			if time.Now().After(account.expires) {
				app.logoutAccountLocked(token)
			}
		}

		for sessionId, lastHealthCheck := range app.sessions {
			elapsed := time.Since(lastHealthCheck)

			if elapsed < time.Minute {
				continue
			}

			delete(app.sessions, sessionId)

			if len(app.sessions) == 0 && app.key != nil {
				app.key = nil
			}
		}

		app.mux.Unlock()
	}
}
