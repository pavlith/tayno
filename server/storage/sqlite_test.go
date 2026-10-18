package storage_test

import (
	"bytes"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"myst/server/core"
	"myst/server/models"
	"myst/server/storage"
)

func openRepo(t *testing.T) *storage.Repository {
	t.Helper()
	r, err := storage.Open(filepath.Join(t.TempDir(), "nested", "test.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := r.Close(); err != nil {
			t.Error(err)
		}
	})
	return r
}

func fixture(t *testing.T, r *storage.Repository) (models.User, models.User, models.Keystore, models.Invitation) {
	t.Helper()
	owner, err := r.CreateUser(models.NewUser(models.WithUsername("owner"), models.WithPasswordHash("hash"), models.WithPublicKey([]byte{0, 1, 255})))
	if err != nil {
		t.Fatal(err)
	}
	invitee, err := r.CreateUser(models.NewUser(models.WithUsername("invitee")))
	if err != nil {
		t.Fatal(err)
	}
	k, err := r.CreateKeystore(models.NewKeystore(models.WithOwnerId(owner.Id), models.WithName("secret"), models.WithPayload([]byte{255, 0, 128})))
	if err != nil {
		t.Fatal(err)
	}
	inv, err := r.CreateInvitation(models.NewInvitation(models.WithKeystoreId(k.Id), models.WithInviterId(owner.Id), models.WithInviteeId(invitee.Id)))
	if err != nil {
		t.Fatal(err)
	}
	return owner, invitee, k, inv
}

func TestPersistenceAndMigrationOnReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "server.sqlite")
	r, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	owner, _, k, inv := fixture(t, r)
	inv.Status = models.Deleted
	inv.DeletedAt = time.Now().UTC().Truncate(time.Microsecond)
	inv.EncryptedKeystoreKey = []byte{0, 255, 8}
	if _, err = r.UpdateInvitation(inv); err != nil {
		t.Fatal(err)
	}
	if err = r.Close(); err != nil {
		t.Fatal(err)
	}
	r, err = storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	u, err := r.UserByUsername("owner")
	if err != nil || u.Id != owner.Id || !bytes.Equal(u.PublicKey, owner.PublicKey) {
		t.Fatalf("user round trip: %+v, %v", u, err)
	}
	got, err := r.Keystore(k.Id)
	if err != nil || !bytes.Equal(got.Payload, k.Payload) {
		t.Fatalf("payload round trip: %+v, %v", got, err)
	}
	loaded, err := r.Invitation(inv.Id)
	if err != nil || loaded.Status != models.Deleted || !loaded.DeletedAt.Equal(inv.DeletedAt) || !bytes.Equal(loaded.EncryptedKeystoreKey, inv.EncryptedKeystoreKey) {
		t.Fatalf("invitation round trip: %+v, %v", loaded, err)
	}
	all, err := r.Invitations()
	if err != nil || len(all) != 1 {
		t.Fatalf("deleted lifecycle invitation must remain queryable: %v %v", all, err)
	}
}

func TestUpdatesClearValuesAndNeverInsertMissingRows(t *testing.T) {
	r := openRepo(t)
	owner, _, k, inv := fixture(t, r)
	owner.PublicKey = nil
	if _, err := r.UpdateUser(owner); err != nil {
		t.Fatal(err)
	}
	u, err := r.User(owner.Id)
	if err != nil || len(u.PublicKey) != 0 {
		t.Fatalf("clear public key: %v %v", u, err)
	}
	k.Name = ""
	k.Payload = []byte{}
	if _, err := r.UpdateKeystore(k); err != nil {
		t.Fatal(err)
	}
	loaded, err := r.Keystore(k.Id)
	if err != nil || loaded.Name != "" || len(loaded.Payload) != 0 {
		t.Fatalf("clear keystore: %v %v", loaded, err)
	}
	inv.AcceptedAt = time.Now()
	inv.EncryptedKeystoreKey = []byte{1}
	if _, err := r.UpdateInvitation(inv); err != nil {
		t.Fatal(err)
	}
	inv.AcceptedAt = time.Time{}
	inv.EncryptedKeystoreKey = nil
	if _, err := r.UpdateInvitation(inv); err != nil {
		t.Fatal(err)
	}
	got, err := r.Invitation(inv.Id)
	if err != nil || !got.AcceptedAt.IsZero() || len(got.EncryptedKeystoreKey) != 0 {
		t.Fatalf("clear invitation: %v %v", got, err)
	}
	owner.Id = "missing"
	k.Id = "missing"
	inv.Id = "missing"
	if _, err := r.UpdateUser(owner); !errors.Is(err, core.ErrUserNotFound) {
		t.Fatal(err)
	}
	if _, err := r.UpdateKeystore(k); !errors.Is(err, core.ErrKeystoreNotFound) {
		t.Fatal(err)
	}
	if _, err := r.UpdateInvitation(inv); !errors.Is(err, core.ErrInvitationNotFound) {
		t.Fatal(err)
	}
	if _, err := r.User("missing"); !errors.Is(err, core.ErrUserNotFound) {
		t.Fatal(err)
	}
	if err := r.DeleteInvitation("missing"); !errors.Is(err, core.ErrInvitationNotFound) {
		t.Fatal(err)
	}
	if err := r.DeleteKeystore("missing"); !errors.Is(err, core.ErrKeystoreNotFound) {
		t.Fatal(err)
	}
}

func TestConstraintsAndKeystoreDeletion(t *testing.T) {
	r := openRepo(t)
	_, _, k, inv := fixture(t, r)
	invalid := models.NewKeystore(models.WithOwnerId("missing"))
	if _, err := r.CreateKeystore(invalid); err == nil {
		t.Fatal("orphan keystore accepted")
	}
	duplicate := inv
	duplicate.Id = "duplicate"
	if _, err := r.CreateInvitation(duplicate); !errors.Is(err, core.ErrForbidden) {
		t.Fatalf("active invitation unique constraint: %v", err)
	}
	inv.Status = models.Declined
	if _, err := r.UpdateInvitation(inv); err != nil {
		t.Fatal(err)
	}
	if _, err := r.CreateInvitation(duplicate); err != nil {
		t.Fatalf("reinvite after decline: %v", err)
	}
	if err := r.DeleteKeystore(k.Id); err != nil {
		t.Fatal(err)
	}
	invs, err := r.Invitations()
	if err != nil || len(invs) != 0 {
		t.Fatalf("associated invitations not deleted: %v, %v", invs, err)
	}
	if _, err := r.Keystore(k.Id); !errors.Is(err, core.ErrKeystoreNotFound) {
		t.Fatal(err)
	}
}

func TestConcurrentUsernameUniqueness(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shared.sqlite")
	r1, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer r1.Close()
	r2, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer r2.Close()
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, repo := range []*storage.Repository{r1, r2} {
		wg.Add(1)
		go func(r *storage.Repository) {
			defer wg.Done()
			_, err := r.CreateUser(models.NewUser(models.WithUsername("same")))
			results <- err
		}(repo)
	}
	wg.Wait()
	close(results)
	var success, duplicate int
	for err := range results {
		if err == nil {
			success++
		} else if errors.Is(err, core.ErrUsernameTaken) {
			duplicate++
		} else {
			t.Fatal(err)
		}
	}
	if success != 1 || duplicate != 1 {
		t.Fatalf("success=%d duplicate=%d", success, duplicate)
	}
}
