package config

import (
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/clients/storage"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/core/service/token"
	"ivory/features/permission"
	"ivory/features/user"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

type nopVerifier struct{}

func (nopVerifier) VerifyPassword(_ string, _ string) error {
	return nil
}

type testConfigEnv struct {
	service       *Service
	userService   *user.Service
	secretService *secret.Service
	basicProvider *basic.Provider
	ldapProvider  *ldap.Provider
	oidcProvider  *oidc.Provider
}

func createTestConfigService(t *testing.T) *testConfigEnv {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "config-service-test-*")
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

	oldWd, errWd := os.Getwd()
	if errWd != nil {
		t.Fatalf("failed to get working dir: %v", errWd)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		os.Chdir(oldWd)
	})

	secretService := secret.NewService(
		secret.NewRepository(storage.NewDbBucket[string](db, "Secret")),
		encryption.NewService(),
	)
	if err := secretService.SetDefault(); err != nil {
		t.Fatalf("failed to set default secret: %v", err)
	}

	permissionService := permission.NewService(
		permission.NewRepository(storage.NewDbBucket[permission.PermissionMap](db, "Permission")),
	)
	userService := user.NewService(
		user.NewRepository(storage.NewDbBucket[user.User](db, "User")),
		encryption.NewService(),
		secretService,
		permissionService,
		token.NewService(secretService),
	)

	basicProvider := basic.NewProvider(nopVerifier{})
	ldapProvider := ldap.NewProvider()
	oidcProvider := oidc.NewProvider()

	configFiles := storage.NewFileStorage("config", "")
	service := NewService(configFiles, encryption.NewService(), secretService, nil, userService, basicProvider, ldapProvider, oidcProvider)

	return &testConfigEnv{
		service:       service,
		userService:   userService,
		secretService: secretService,
		basicProvider: basicProvider,
		ldapProvider:  ldapProvider,
		oidcProvider:  oidcProvider,
	}
}

func TestServiceGetIsConfiguredInitiallyFalse(t *testing.T) {
	env := createTestConfigService(t)
	if env.service.GetIsConfigured() {
		t.Fatalf("expected a fresh service to be unconfigured")
	}
}

func TestServiceSetAppConfigValidation(t *testing.T) {
	env := createTestConfigService(t)

	t.Run("empty company name is rejected", func(t *testing.T) {
		err := env.service.SetAppConfig(NewAppConfig{AppConfig: AppConfig{Company: ""}})
		if err != ErrCompanyNameEmpty {
			t.Fatalf("expected ErrCompanyNameEmpty, got %v", err)
		}
	})
}

func TestServiceSetAppConfigSuperuser(t *testing.T) {
	t.Run("switching authentication on without a superuser is refused", func(t *testing.T) {
		env := createTestConfigService(t)

		req := NewAppConfig{AppConfig: AppConfig{Company: "Acme", Auth: AuthConfig{Basic: &basic.Config{}}}}
		if err := env.service.SetAppConfig(req); err != ErrSuperuserRequired {
			t.Fatalf("expected ErrSuperuserRequired, got %v", err)
		}
		if env.service.GetIsConfigured() {
			t.Fatalf("expected nothing to be saved")
		}
	})

	t.Run("the superuser named at setup is created and can sign in", func(t *testing.T) {
		env := createTestConfigService(t)

		req := NewAppConfig{
			User:      &UserSetupRequest{Username: "root", Password: "password123", AuthTypes: []user.AuthType{user.AuthBasic}},
			AppConfig: AppConfig{Company: "Acme", Auth: AuthConfig{Basic: &basic.Config{}}},
		}
		if err := env.service.SetAppConfig(req); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if err := env.userService.VerifyPassword("root", "password123"); err != nil {
			t.Fatalf("expected the superuser to be created, got %v", err)
		}
		hasSuperuser, errSuper := env.userService.HasSuperuser()
		if errSuper != nil || !hasSuperuser {
			t.Fatalf("expected a superuser, got %v (%v)", hasSuperuser, errSuper)
		}
	})

	t.Run("a superuser that cannot be created stops the config being saved", func(t *testing.T) {
		env := createTestConfigService(t)

		req := NewAppConfig{
			User:      &UserSetupRequest{Username: "root", Password: "", AuthTypes: []user.AuthType{user.AuthBasic}},
			AppConfig: AppConfig{Company: "Acme", Auth: AuthConfig{Basic: &basic.Config{}}},
		}
		if err := env.service.SetAppConfig(req); err == nil {
			t.Fatalf("expected an error for a superuser without a password")
		}
		if env.service.GetIsConfigured() {
			t.Fatalf("expected nothing to be saved")
		}
	})

	t.Run("a setup that names no superuser is fine once one exists", func(t *testing.T) {
		env := createTestConfigService(t)
		if _, err := env.userService.CreateOutright("root", "password123", []user.AuthType{user.AuthBasic}, true); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}

		req := NewAppConfig{AppConfig: AppConfig{Company: "Acme", Auth: AuthConfig{Basic: &basic.Config{}}}}
		if err := env.service.SetAppConfig(req); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("no authentication needs no superuser", func(t *testing.T) {
		env := createTestConfigService(t)

		req := NewAppConfig{AppConfig: AppConfig{Company: "Acme"}}
		if err := env.service.SetAppConfig(req); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestServiceSetAndGetAppConfig(t *testing.T) {
	env := createTestConfigService(t)

	req := NewAppConfig{
		User: &UserSetupRequest{Username: "root", Password: "password123", AuthTypes: []user.AuthType{user.AuthBasic}},
		AppConfig: AppConfig{
			Company: "Acme",
			Auth: AuthConfig{
				Basic: &basic.Config{},
			},
		},
	}
	if err := env.service.SetAppConfig(req); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !env.basicProvider.Configured() {
		t.Fatalf("expected the basic provider to be configured")
	}

	appConfig, errGet := env.service.GetAppConfig()
	if errGet != nil {
		t.Fatalf("expected no error, got %v", errGet)
	}
	if appConfig.Company != "Acme" {
		t.Fatalf("expected company 'Acme', got %q", appConfig.Company)
	}
	if appConfig.Auth.Basic == nil {
		t.Fatalf("expected basic auth to stay enabled, got %+v", appConfig.Auth)
	}

	if !env.service.GetIsConfigured() {
		t.Fatalf("expected the service to be configured after loading the saved config")
	}
}

func TestServiceSetAppConfigAlreadySetRejectsWrongSecret(t *testing.T) {
	env := createTestConfigService(t)

	initial := NewAppConfig{AppConfig: AppConfig{Company: "Acme"}}
	if err := env.service.SetAppConfig(initial); err != nil {
		t.Fatalf("failed to set initial config: %v", err)
	}
	if _, err := env.service.GetAppConfig(); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	t.Run("wrong secret is rejected once configured", func(t *testing.T) {
		update := NewAppConfig{Secret: "wrong", AppConfig: AppConfig{Company: "NewCo"}}
		if err := env.service.SetAppConfig(update); err != ErrConfigAlreadySet {
			t.Fatalf("expected ErrConfigAlreadySet, got %v", err)
		}
	})

	t.Run("correct secret allows an update", func(t *testing.T) {
		update := NewAppConfig{Secret: "ivory", AppConfig: AppConfig{Company: "NewCo"}}
		if err := env.service.SetAppConfig(update); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// SetAppConfig persists to file but does not refresh the in-memory
		// cache, so a fresh service is used to verify what was actually saved.
		fresh := NewService(
			storage.NewFileStorage("config", ""),
			encryption.NewService(),
			env.secretService,
			nil,
			env.userService,
			basic.NewProvider(nopVerifier{}),
			ldap.NewProvider(),
			oidc.NewProvider(),
		)
		appConfig, err := fresh.GetAppConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if appConfig.Company != "NewCo" {
			t.Fatalf("expected company 'NewCo', got %q", appConfig.Company)
		}
	})
}

func TestServiceSetAppConfigRollsBackOnAuthValidationError(t *testing.T) {
	env := createTestConfigService(t)

	req := NewAppConfig{
		User: &UserSetupRequest{Username: "root", Password: "password123", AuthTypes: []user.AuthType{user.AuthBasic}},
		AppConfig: AppConfig{
			Company: "Acme",
			Auth: AuthConfig{
				Basic: &basic.Config{},
				Ldap:  &ldap.Config{}, // invalid: all fields empty
			},
		},
	}
	if err := env.service.SetAppConfig(req); err == nil {
		t.Fatalf("expected an error for the invalid ldap config")
	}

	if env.basicProvider.Configured() {
		t.Fatalf("expected the basic provider config to be rolled back after the ldap validation failure")
	}
}

func TestServiceReencrypt(t *testing.T) {
	env := createTestConfigService(t)

	req := NewAppConfig{
		User: &UserSetupRequest{Username: "root", Password: "password123", AuthTypes: []user.AuthType{user.AuthBasic}},
		AppConfig: AppConfig{
			Company: "Acme",
			Auth: AuthConfig{
				Basic: &basic.Config{},
				Ldap: &ldap.Config{
					Url:      "ldap://localhost:389",
					BindDN:   "cn=admin",
					BindPass: "secretpw",
					BaseDN:   "dc=acme",
					Filter:   "(uid=%s)",
				},
			},
		},
	}
	if err := env.service.SetAppConfig(req); err != nil {
		t.Fatalf("failed to set config: %v", err)
	}
	if _, err := env.service.GetAppConfig(); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if _, _, err := env.secretService.Update("ivory", "newsecret"); err != nil {
		t.Fatalf("failed to rotate secret: %v", err)
	}

	if err := env.service.Reencrypt(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// A fresh service sharing the same files and (rotated) secret should decrypt correctly.
	fresh := NewService(
		storage.NewFileStorage("config", ""),
		encryption.NewService(),
		env.secretService,
		nil,
		env.userService,
		basic.NewProvider(nopVerifier{}),
		ldap.NewProvider(),
		oidc.NewProvider(),
	)
	appConfig, err := fresh.GetAppConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if appConfig.Auth.Ldap == nil || appConfig.Auth.Ldap.BindPass != "secretpw" {
		t.Fatalf("expected the bind password to still decrypt to 'secretpw' after reencryption, got %+v", appConfig.Auth.Ldap)
	}
}

func TestServiceDeleteAll(t *testing.T) {
	env := createTestConfigService(t)

	req := NewAppConfig{
		User: &UserSetupRequest{Username: "root", Password: "password123", AuthTypes: []user.AuthType{user.AuthBasic}},
		AppConfig: AppConfig{
			Company: "Acme",
			Auth: AuthConfig{
				Basic: &basic.Config{},
			},
		},
	}
	if err := env.service.SetAppConfig(req); err != nil {
		t.Fatalf("failed to set config: %v", err)
	}
	if _, err := env.service.GetAppConfig(); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if err := env.service.DeleteAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if env.basicProvider.Configured() {
		t.Fatalf("expected the basic provider config to be cleared")
	}
	if _, err := env.service.GetAppConfig(); err != ErrConfigNotSpecified {
		t.Fatalf("expected ErrConfigNotSpecified after delete, got %v", err)
	}
}
