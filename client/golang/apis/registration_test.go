package apis

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"myst/client/golang/core"
	"myst/client/golang/models"
)

type registrationApp struct {
	core.Application
	user *models.User
	err  error
}

func (a registrationApp) CurrentUser([]byte) (*models.User, error) { return a.user, a.err }
func (a registrationApp) Keystores([]byte) (map[string]models.Keystore, error) {
	return map[string]models.Keystore{}, nil
}

func TestVaultRequiresRegisteredUser(t *testing.T) {
	for _, tc := range []struct {
		name   string
		user   *models.User
		err    error
		status int
	}{
		{"unregistered", nil, core.ErrCredentialsNotFound, http.StatusForbidden},
		{"signed out", nil, nil, http.StatusForbidden},
		{"invalid session", nil, core.ErrForbidden, http.StatusForbidden},
		{"registered", &models.User{Id: "user"}, nil, http.StatusOK},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := NewServer(registrationApp{user: tc.user, err: tc.err}, fstest.MapFS{"static/index.html": &fstest.MapFile{Data: []byte("test")}})
			res := httptest.NewRecorder()
			s.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/keystores", nil))
			if res.Code != tc.status {
				t.Fatalf("got %d, want %d: %s", res.Code, tc.status, res.Body.String())
			}
		})
	}
}
