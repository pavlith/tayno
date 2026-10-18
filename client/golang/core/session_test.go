package core

import (
	"testing"
	"time"
)

func TestLogoutRevokesSession(t *testing.T) {
	app := &application{sessions: map[string]time.Time{"first": time.Now(), "second": time.Now()}, key: []byte("key")}
	if err := app.Logout([]byte("first")); err != nil {
		t.Fatal(err)
	}
	if app.sessionActive([]byte("first")) {
		t.Fatal("logged out session is still active")
	}
	if app.key == nil {
		t.Fatal("another session still needs the encryption key")
	}
	if err := app.HealthCheck([]byte("first")); err != ErrAuthenticationRequired {
		t.Fatalf("revoked session heartbeat: %v", err)
	}
	if app.sessionActive([]byte("first")) {
		t.Fatal("heartbeat restored revoked session")
	}
	if err := app.Logout([]byte("second")); err != nil {
		t.Fatal(err)
	}
	if app.key != nil {
		t.Fatal("key retained after last logout")
	}
	if err := app.Logout([]byte("second")); err != nil {
		t.Fatal("logout must be idempotent", err)
	}
}
