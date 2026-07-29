package secret

import (
	"crypto/md5"
	"errors"
	"ivory/clients/storage"
	"ivory/core/service/encryption"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

func createTestSecretService(t *testing.T) *Service {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "secret-service-test-*")
	if errDir != nil {
		t.Fatalf("failed to create temp dir: %v", errDir)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	db, errOpen := bolt.Open(filepath.Join(tmpDir, "test.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() {
		db.Close()
	})

	repository := NewRepository(storage.NewDbBucket[string](db, "Secret"))
	return NewService(repository, encryption.NewService())
}

func TestServiceSetDefault(t *testing.T) {
	s := createTestSecretService(t)

	if !s.IsRefEmpty() {
		t.Fatalf("expected a fresh service to have an empty ref")
	}
	if status := s.Status(); status.Key || status.Ref {
		t.Fatalf("expected empty status, got %+v", status)
	}

	if err := s.SetDefault(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.IsRefEmpty() {
		t.Fatalf("expected ref to be set")
	}
	if s.Get() != md5.Sum([]byte("ivory")) {
		t.Fatalf("expected key to be md5 of the default ref")
	}

	if err := s.SetDefault(); err != ErrSecretAlreadySet {
		t.Fatalf("expected ErrSecretAlreadySet, got %v", err)
	}
}

func TestServiceSetAndVerify(t *testing.T) {
	s := createTestSecretService(t)

	if err := s.Set(""); err != ErrSecretCannotBeEmpty {
		t.Fatalf("expected ErrSecretCannotBeEmpty, got %v", err)
	}

	if err := s.Set("mysecret"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !s.Verify("mysecret") {
		t.Fatalf("expected verify to succeed with the set secret")
	}
	if s.Verify("wrong") {
		t.Fatalf("expected verify to fail with a wrong secret")
	}
}

func TestServiceRestartRecoversRefWithCorrectSecretOnly(t *testing.T) {
	tmpDir, errDir := os.MkdirTemp("", "secret-service-restart-test-*")
	if errDir != nil {
		t.Fatalf("failed to create temp dir: %v", errDir)
	}
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	db, errOpen := bolt.Open(filepath.Join(tmpDir, "test.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() { db.Close() })

	bucket := storage.NewDbBucket[string](db, "Secret")
	repository := NewRepository(bucket)
	enc := encryption.NewService()

	s := NewService(repository, enc)
	if err := s.Set("mysecret"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	restarted := NewService(NewRepository(bucket), enc)
	if restarted.IsRefEmpty() {
		t.Fatalf("expected ref to be recovered from storage")
	}
	if !restarted.IsSecretEmpty() {
		t.Fatalf("expected the recovered service to have no working key until the correct secret is provided")
	}
	if status := restarted.Status(); !status.Ref || status.Key {
		t.Fatalf("expected ref-only status, got %+v", status)
	}

	if err := restarted.Set("mysecret"); err != nil {
		t.Fatalf("expected no error providing the correct secret, got %v", err)
	}
	if !restarted.Verify("mysecret") {
		t.Fatalf("expected verify to succeed after providing the correct secret")
	}

	if err := restarted.Set("wrong"); !errors.Is(err, ErrWrongSecret) {
		t.Fatalf("expected ErrWrongSecret, got %v", err)
	}
}

func TestServiceUpdate(t *testing.T) {
	s := createTestSecretService(t)

	if _, _, err := s.Update("a", "a"); err != ErrSecretsAreSame {
		t.Fatalf("expected ErrSecretsAreSame, got %v", err)
	}

	if _, _, err := s.Update("a", "b"); err != ErrNoSecret {
		t.Fatalf("expected ErrNoSecret, got %v", err)
	}

	if err := s.SetDefault(); err != nil {
		t.Fatalf("failed to set default: %v", err)
	}

	t.Run("wrong previous secret is rejected", func(t *testing.T) {
		if _, _, err := s.Update("wrong", "newpass"); err != ErrWrongSecret {
			t.Fatalf("expected ErrWrongSecret, got %v", err)
		}
	})

	t.Run("empty previous secret falls back to the default ref", func(t *testing.T) {
		oldKey, newKey, err := s.Update("", "newpass")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if oldKey != md5.Sum([]byte("ivory")) {
			t.Fatalf("expected old key to be md5 of the default ref")
		}
		if newKey != md5.Sum([]byte("newpass")) {
			t.Fatalf("expected new key to be md5 of newpass")
		}
		if !s.Verify("newpass") {
			t.Fatalf("expected verify to succeed with the new secret")
		}
	})

	t.Run("empty new secret falls back to the default ref", func(t *testing.T) {
		_, newKey, err := s.Update("newpass", "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if newKey != md5.Sum([]byte("ivory")) {
			t.Fatalf("expected new key to be md5 of the default ref")
		}
	})
}

func TestServiceClean(t *testing.T) {
	s := createTestSecretService(t)

	if err := s.Set("mysecret"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := s.Clean(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !s.IsRefEmpty() {
		t.Fatalf("expected ref to be cleared")
	}
	if !s.IsSecretEmpty() {
		t.Fatalf("expected key to be cleared")
	}
	ref, err := s.repository.Get()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ref != "" {
		t.Fatalf("expected persisted ref to be cleared, got %q", ref)
	}
}
