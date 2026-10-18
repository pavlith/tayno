package storage_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/client/golang/models"
	"myst/client/golang/storage"
	"myst/pkgs/rand"
)

func TestRepository(t *testing.T) {
	t.Parallel()

	repo := storage.New(t.TempDir())

	pass, err := rand.String(16)
	assert.NilError(t, err)

	keypair, err := repo.Initialize(pass)
	assert.NilError(t, err)

	keypair2, err := repo.Authenticate(pass)
	assert.NilError(t, err)

	assert.DeepEqual(t, keypair, keypair2)

	k, err := repo.CreateKeystore(keypair, models.NewKeystore(models.WithName("test")))
	assert.NilError(t, err)

	k2, err := repo.Keystore(keypair, k.Id)
	assert.NilError(t, err)

	assert.Equal(t, k.Id, k2.Id)
	assert.Equal(t, k.Name, k2.Name)
}
