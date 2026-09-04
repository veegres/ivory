package user

import (
	"errors"
	"ivory/clients/storage"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

func newTestRepository(t *testing.T) *Repository {
	t.Helper()

	db, err := bolt.Open(filepath.Join(t.TempDir(), "test.db"), 0600, nil)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return NewRepository(storage.NewDbBucket[User](db, "User"))
}

func TestRepositoryCreateAndGet(t *testing.T) {
	r := newTestRepository(t)

	created, err := r.Create(User{Username: "alice", Password: "encrypted", Superuser: true})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.Username != "alice" {
		t.Fatalf("expected username 'alice', got %q", created.Username)
	}

	got, errGet := r.Get("alice")
	if errGet != nil {
		t.Fatalf("expected no error, got %v", errGet)
	}
	if got.Password != "encrypted" || !got.Superuser {
		t.Fatalf("expected the stored user back, got %+v", got)
	}
}

func TestRepositoryCreateRejectsDuplicates(t *testing.T) {
	r := newTestRepository(t)

	if _, err := r.Create(User{Username: "alice", Password: "one"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := r.Create(User{Username: "alice", Password: "two"}); !errors.Is(err, storage.ErrAlreadyExists) {
		t.Fatalf("expected storage.ErrAlreadyExists, got %v", err)
	}
}

func TestRepositoryGetMissing(t *testing.T) {
	r := newTestRepository(t)

	if _, err := r.Get("nobody"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected storage.ErrNotFound, got %v", err)
	}
}

func TestRepositoryList(t *testing.T) {
	r := newTestRepository(t)

	for _, name := range []string{"carol", "alice", "bob"} {
		if _, err := r.Create(User{Username: name}); err != nil {
			t.Fatalf("failed to seed %s: %v", name, err)
		}
	}

	users, err := r.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}
	for i, expected := range []string{"alice", "bob", "carol"} {
		if users[i].Username != expected {
			t.Fatalf("expected %q at %d, got %q", expected, i, users[i].Username)
		}
	}
}

func TestRepositoryUpdate(t *testing.T) {
	r := newTestRepository(t)

	if _, err := r.Create(User{Username: "alice", Password: "old"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := r.Update(User{Username: "alice", Password: "new", Superuser: true}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, errGet := r.Get("alice")
	if errGet != nil {
		t.Fatalf("expected no error, got %v", errGet)
	}
	if got.Password != "new" || !got.Superuser {
		t.Fatalf("expected the updated user, got %+v", got)
	}
}

func TestRepositoryDelete(t *testing.T) {
	r := newTestRepository(t)

	if _, err := r.Create(User{Username: "alice"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := r.Delete("alice"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := r.Get("alice"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected storage.ErrNotFound, got %v", err)
	}
}

func TestRepositoryDeleteAll(t *testing.T) {
	r := newTestRepository(t)

	for _, name := range []string{"alice", "bob"} {
		if _, err := r.Create(User{Username: name}); err != nil {
			t.Fatalf("failed to seed %s: %v", name, err)
		}
	}
	if err := r.DeleteAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	users, err := r.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected no users, got %d", len(users))
	}
}

func TestRepositorySuperusers(t *testing.T) {
	r := newTestRepository(t)

	seed := []User{
		{Username: "carol", Superuser: true},
		{Username: "alice", Superuser: true},
		{Username: "bob"},
	}
	for _, u := range seed {
		if _, err := r.Create(u); err != nil {
			t.Fatalf("failed to seed %s: %v", u.Username, err)
		}
	}

	superusers, err := r.Superusers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(superusers) != 2 || superusers[0] != "alice" || superusers[1] != "carol" {
		t.Fatalf("expected sorted [alice carol], got %v", superusers)
	}
}
