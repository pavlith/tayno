package suite

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"

	"myst/pkgs/crypto"
	"myst/server/apis"
	"myst/server/apis/generated"
	"myst/server/core"
	"myst/server/storage"
	"path/filepath"
)

type Server struct {
	Address string

	repo   *storage.Repository
	App    core.Application
	Server *apis.Server
	Client *generated.ClientWithResponses
}

func newServer(t *testing.T, address string) *Server {
	repo, err := storage.Open(filepath.Join(t.TempDir(), "server.sqlite"))
	assert.NilError(t, err)

	jwtSigningKey, err := crypto.GenerateRandomBytes(32)
	assert.NilError(t, err)

	app := core.New(
		core.WithKeystoreRepository(repo),
		core.WithUserRepository(repo),
		core.WithInvitationRepository(repo),
	)
	assert.NilError(t, err)

	server := apis.NewServer(app, jwtSigningKey)

	clientAddr := fmt.Sprintf("http://%s/api", address)

	client, err := generated.NewClientWithResponses(clientAddr)
	assert.NilError(t, err)

	return &Server{
		Address: address,
		repo:    repo,
		App:     app,
		Server:  server,
		Client:  client,
	}
}

func (s *Server) start(t *testing.T) {
	err := s.Server.Start(s.Address)
	assert.NilError(t, err)

}

func (s *Server) stop(t *testing.T) {
	_ = s.Server.Stop()

	err := s.repo.Close()
	assert.NilError(t, err)
}
