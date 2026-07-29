package vault

import (
	"ivory/clients/console/ssh"
	"ivory/clients/storage"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

func createTestVaultService(t *testing.T) *Service {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "vault-service-test-*")
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

	secretRepository := secret.NewRepository(storage.NewDbBucket[string](db, "Secret"))
	secretService := secret.NewService(secretRepository, encryption.NewService())
	if err := secretService.SetDefault(); err != nil {
		t.Fatalf("failed to set default secret: %v", err)
	}

	vaultRepository := NewRepository(storage.NewDbBucket[Vault](db, "Vault"))
	return NewService(vaultRepository, ssh.NewClient(), secretService, encryption.NewService())
}

func TestServiceCreatePasswordVault(t *testing.T) {
	s := createTestVaultService(t)

	t.Run("missing username is rejected", func(t *testing.T) {
		_, _, err := s.Create(Vault{Type: DATABASE_PASSWORD, Secret: "pw"})
		if err != ErrVaultUsernameIsEmpty {
			t.Fatalf("expected ErrVaultUsernameIsEmpty, got %v", err)
		}
	})

	t.Run("missing secret is rejected", func(t *testing.T) {
		_, _, err := s.Create(Vault{Type: DATABASE_PASSWORD, Username: "user"})
		if err != ErrVaultSecretIsEmpty {
			t.Fatalf("expected ErrVaultSecretIsEmpty, got %v", err)
		}
	})

	t.Run("valid password vault is created with the secret hidden and encrypted", func(t *testing.T) {
		key, created, err := s.Create(Vault{Type: DATABASE_PASSWORD, Username: "user", Secret: "pw"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if created.Secret != "********" {
			t.Fatalf("expected the returned secret to be hidden, got %q", created.Secret)
		}

		stored, errGet := s.vaultRepository.Get(*key)
		if errGet != nil {
			t.Fatalf("expected no error, got %v", errGet)
		}
		if stored.Secret == "pw" {
			t.Fatalf("expected the stored secret to be encrypted, not stored in plaintext")
		}

		decrypted, errDec := s.GetDecrypted(*key)
		if errDec != nil {
			t.Fatalf("expected no error, got %v", errDec)
		}
		if decrypted.Secret != "pw" {
			t.Fatalf("expected decrypted secret 'pw', got %q", decrypted.Secret)
		}
	})
}

func TestServiceCreateSSHKeyVault(t *testing.T) {
	s := createTestVaultService(t)

	key, created, err := s.Create(Vault{Type: SSH_KEY, Username: "user"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.Metadata == nil || *created.Metadata == "" {
		t.Fatalf("expected the public key to be stored as metadata")
	}
	if created.Secret != "********" {
		t.Fatalf("expected the returned secret to be hidden, got %q", created.Secret)
	}

	decrypted, errDec := s.GetDecrypted(*key)
	if errDec != nil {
		t.Fatalf("expected no error, got %v", errDec)
	}
	if decrypted.Secret == "" {
		t.Fatalf("expected a decrypted private key")
	}
}

func TestServiceUpdate(t *testing.T) {
	s := createTestVaultService(t)

	key, _, err := s.Create(Vault{Type: DATABASE_PASSWORD, Username: "user", Secret: "pw"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedKey, updated, errUpdate := s.Update(*key, Vault{Type: DATABASE_PASSWORD, Username: "user2", Secret: "pw2"})
	if errUpdate != nil {
		t.Fatalf("expected no error, got %v", errUpdate)
	}
	if *updatedKey != *key {
		t.Fatalf("expected the same key to be returned")
	}
	if updated.Username != "user2" {
		t.Fatalf("expected the username to be updated, got %q", updated.Username)
	}

	decrypted, errDec := s.GetDecrypted(*key)
	if errDec != nil {
		t.Fatalf("expected no error, got %v", errDec)
	}
	if decrypted.Secret != "pw2" {
		t.Fatalf("expected the updated secret 'pw2', got %q", decrypted.Secret)
	}
}

func TestServiceGetHidesSecret(t *testing.T) {
	s := createTestVaultService(t)

	key, _, err := s.Create(Vault{Type: DATABASE_PASSWORD, Username: "user", Secret: "pw"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, errGet := s.Get(*key)
	if errGet != nil {
		t.Fatalf("expected no error, got %v", errGet)
	}
	if got.Secret != "********" {
		t.Fatalf("expected hidden secret, got %q", got.Secret)
	}

	unknown := uuid.New()
	if _, err := s.Get(unknown); err == nil {
		t.Fatalf("expected an error getting an unknown vault")
	}
}

func TestServiceMap(t *testing.T) {
	s := createTestVaultService(t)

	if _, _, err := s.Create(Vault{Type: DATABASE_PASSWORD, Username: "db-user", Secret: "pw"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, _, err := s.Create(Vault{Type: SSH_PASSWORD, Username: "ssh-user", Secret: "pw"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("no filter returns all vaults with secrets hidden", func(t *testing.T) {
		all, err := s.Map(nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(all) != 2 {
			t.Fatalf("expected 2 vaults, got %d", len(all))
		}
		for _, v := range all {
			if v.Secret != "********" {
				t.Fatalf("expected hidden secret, got %q", v.Secret)
			}
		}
	})

	t.Run("filter by type narrows the result", func(t *testing.T) {
		sshType := SSH_PASSWORD
		filtered, err := s.Map(&sshType)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(filtered) != 1 {
			t.Fatalf("expected 1 vault, got %d", len(filtered))
		}
	})
}

func TestServiceDeleteAndDeleteAll(t *testing.T) {
	s := createTestVaultService(t)

	key, _, err := s.Create(Vault{Type: DATABASE_PASSWORD, Username: "user", Secret: "pw"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("delete removes the vault", func(t *testing.T) {
		if err := s.Delete(*key); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, err := s.Get(*key); err == nil {
			t.Fatalf("expected an error getting a deleted vault")
		}
	})

	t.Run("delete all clears every vault", func(t *testing.T) {
		if _, _, err := s.Create(Vault{Type: DATABASE_PASSWORD, Username: "user", Secret: "pw"}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := s.DeleteAll(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		all, errMap := s.Map(nil)
		if errMap != nil {
			t.Fatalf("expected no error, got %v", errMap)
		}
		if len(all) != 0 {
			t.Fatalf("expected no vaults, got %v", all)
		}
	})
}

func TestServiceReencrypt(t *testing.T) {
	s := createTestVaultService(t)

	key, _, err := s.Create(Vault{Type: DATABASE_PASSWORD, Username: "user", Secret: "pw"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	oldSecret := s.secretService.Get()
	newSecret := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	if err := s.Reencrypt(oldSecret, newSecret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	stored, errGet := s.vaultRepository.Get(*key)
	if errGet != nil {
		t.Fatalf("expected no error, got %v", errGet)
	}
	decrypted, errDec := s.encryption.Decrypt(stored.Secret, newSecret)
	if errDec != nil {
		t.Fatalf("expected the secret to be decryptable with the new key, got %v", errDec)
	}
	if decrypted != "pw" {
		t.Fatalf("expected decrypted secret 'pw', got %q", decrypted)
	}
}

func TestSecretService(t *testing.T) {
	s := createTestVaultService(t)
	if s.SecretService() == nil {
		t.Fatalf("expected a non-nil secret service")
	}
}
