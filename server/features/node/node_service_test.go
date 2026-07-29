package node

import (
	"ivory/clients/console/ssh"
	"ivory/clients/storage"
	"ivory/core/config"
	"ivory/core/service/cert"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

type fakePlatformAdapter struct {
	platform.Adapter
}

type fakeKeeperMetadata struct {
	features map[config.Feature]bool
}

func (f *fakeKeeperMetadata) SupportedFeatures() map[config.Feature]bool {
	return f.features
}

func (f *fakeKeeperMetadata) DeploymentSpec() keeper.DeploymentSpec {
	return keeper.DeploymentSpec{}
}

func createTestNodeService(t *testing.T) (*Service, *vault.Service) {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "node-service-test-*")
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

	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(oldWd) })

	secretService := secret.NewService(
		secret.NewRepository(storage.NewDbBucket[string](db, "Secret")),
		encryption.NewService(),
	)
	if err := secretService.SetDefault(); err != nil {
		t.Fatalf("failed to set default secret: %v", err)
	}

	vaultService := vault.NewService(
		vault.NewRepository(storage.NewDbBucket[vault.Vault](db, "Vault")),
		ssh.NewClient(),
		secretService,
		encryption.NewService(),
	)

	certService := cert.NewService(
		cert.NewRepository(storage.NewDbBucket[cert.Cert](db, "Cert"), storage.NewFileStorage("cert", "")),
	)

	platformRegistry := utils.NewRegistry[platform.Plugin, platform.Adapter]()
	platformRegistry.Register(platform.Linux, &fakePlatformAdapter{})

	keeperRegistry := utils.NewRegistry[keeper.Plugin, keeper.Adapter]()
	keeperRegistry.Register("fake", &fakeKeeperAdapter{})

	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	keeperMetadataRegistry.Register("fake", &fakeKeeperMetadata{features: map[config.Feature]bool{config.ViewClusterList: true}})

	s := NewService(platformRegistry, keeperRegistry, keeperMetadataRegistry, vaultService, certService, nil)
	return s, vaultService
}

func TestServiceSupportedFeatures(t *testing.T) {
	s, _ := createTestNodeService(t)

	t.Run("known plugin returns its feature map", func(t *testing.T) {
		features := s.SupportedFeatures("fake")
		if !features[config.ViewClusterList] {
			t.Fatalf("expected ViewClusterList to be supported, got %v", features)
		}
	})

	t.Run("unknown plugin returns an empty map", func(t *testing.T) {
		features := s.SupportedFeatures("unknown")
		if len(features) != 0 {
			t.Fatalf("expected an empty map, got %v", features)
		}
	})
}

func TestServiceGetPlatformAdapter(t *testing.T) {
	s, vaultService := createTestNodeService(t)

	key, _, errCreate := vaultService.Create(vault.Vault{Type: vault.SSH_KEY, Username: "root"})
	if errCreate != nil {
		t.Fatalf("failed to seed ssh key vault: %v", errCreate)
	}

	t.Run("known platform and vault produce a usable connection", func(t *testing.T) {
		adapter, conn, err := s.getPlatformAdapter(PlatformVaultConnection{Host: "host1", Port: 22, VaultId: *key})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if adapter == nil {
			t.Fatalf("expected a non-nil adapter")
		}
		if conn.Host != "host1" || conn.Port != 22 || conn.Username != "root" {
			t.Fatalf("unexpected connection: %+v", conn)
		}
		if len(conn.PrivateKey) == 0 {
			t.Fatalf("expected a decrypted private key")
		}
	})

	t.Run("unknown vault id fails", func(t *testing.T) {
		badId := *key
		badId[0] ^= 0xFF
		_, _, err := s.getPlatformAdapter(PlatformVaultConnection{Host: "host1", Port: 22, VaultId: badId})
		if err == nil {
			t.Fatalf("expected an error for an unknown vault id")
		}
	})
}

func TestServiceGetKeeperAdapter(t *testing.T) {
	s, vaultService := createTestNodeService(t)

	key, _, errCreate := vaultService.Create(vault.Vault{Type: vault.KEEPER_PASSWORD, Username: "keeper-user", Secret: "keeper-pass"})
	if errCreate != nil {
		t.Fatalf("failed to seed keeper vault: %v", errCreate)
	}

	t.Run("no vault id and no certs returns bare adapter", func(t *testing.T) {
		adapter, tlsConfig, cred, err := s.getKeeperAdapter(KeeperOptions{Plugin: "fake"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if adapter == nil {
			t.Fatalf("expected a non-nil adapter")
		}
		if tlsConfig != nil {
			t.Fatalf("expected a nil tls config, got %+v", tlsConfig)
		}
		if cred != nil {
			t.Fatalf("expected nil credentials, got %+v", cred)
		}
	})

	t.Run("vault id resolves decrypted credentials", func(t *testing.T) {
		_, _, cred, err := s.getKeeperAdapter(KeeperOptions{Plugin: "fake", VaultId: key})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if cred == nil || cred.Username != "keeper-user" || cred.Password != "keeper-pass" {
			t.Fatalf("unexpected credentials: %+v", cred)
		}
	})

	t.Run("certs enrich the tls config", func(t *testing.T) {
		_, tlsConfig, _, err := s.getKeeperAdapter(KeeperOptions{Plugin: "fake", Certs: &cert.Certs{}})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tlsConfig == nil {
			t.Fatalf("expected a non-nil tls config")
		}
	})

	t.Run("unknown plugin fails", func(t *testing.T) {
		_, _, _, err := s.getKeeperAdapter(KeeperOptions{Plugin: "unknown"})
		if err == nil {
			t.Fatalf("expected an error for an unknown plugin")
		}
	})

	t.Run("unknown vault id fails", func(t *testing.T) {
		badId := *key
		badId[0] ^= 0xFF
		_, _, _, err := s.getKeeperAdapter(KeeperOptions{Plugin: "fake", VaultId: &badId})
		if err == nil {
			t.Fatalf("expected an error for an unknown vault id")
		}
	})
}

func TestServiceGetPlatformVaultConnection(t *testing.T) {
	s, vaultService := createTestNodeService(t)

	key, _, errCreate := vaultService.Create(vault.Vault{Type: vault.SSH_KEY, Username: "root"})
	if errCreate != nil {
		t.Fatalf("failed to seed ssh key vault: %v", errCreate)
	}

	t.Run("valid vault id resolves the connection", func(t *testing.T) {
		conn, err := s.getPlatformVaultConnection(PlatformVaultConnection{Host: "h", Port: 22, VaultId: *key})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if conn.Username != "root" {
			t.Fatalf("expected username 'root', got %q", conn.Username)
		}
	})

	t.Run("unknown vault id fails", func(t *testing.T) {
		badId := *key
		badId[0] ^= 0xFF
		_, err := s.getPlatformVaultConnection(PlatformVaultConnection{Host: "h", Port: 22, VaultId: badId})
		if err == nil {
			t.Fatalf("expected an error for an unknown vault id")
		}
	})
}

func TestServiceGetPlatformCredConnection(t *testing.T) {
	s, _ := createTestNodeService(t)

	conn := s.getPlatformCredConnection(PlatformCredConnection{Host: "h", Port: 22, Username: "user", Password: "pass"})
	if conn.Host != "h" || conn.Port != 22 || conn.Username != "user" || conn.Password != "pass" {
		t.Fatalf("unexpected connection: %+v", conn)
	}
}
