package management

import (
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/clients/console/ssh"
	"ivory/clients/storage"
	coreConfig "ivory/core/config"
	"ivory/core/service/cert"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/core/service/token"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/features/auth"
	"ivory/features/backup"
	"ivory/features/cluster"
	"ivory/features/config"
	"ivory/features/deployment"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/features/tag"
	"ivory/features/user"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"ivory/tools"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

type testManagementEnv struct {
	service           *Service
	secretService     *secret.Service
	vaultService      *vault.Service
	certService       *cert.Service
	clusterService    *cluster.Service
	tagService        *tag.Service
	queryService      *query.Service
	configService     *config.Service
	permissionService *permission.Service
	basicProvider     *basic.Provider
}

func createTestManagementService(t *testing.T) *testManagementEnv {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "management-service-test-*")
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

	vaultService := vault.NewService(
		vault.NewRepository(storage.NewDbBucket[vault.Vault](db, "Vault")),
		ssh.NewClient(),
		secretService,
		encryption.NewService(),
	)

	certService := cert.NewService(
		cert.NewRepository(storage.NewDbBucket[cert.Cert](db, "Cert"), storage.NewFileStorage("cert", "")),
	)

	tagService := tag.NewService(tag.NewRepository(storage.NewDbBucket[[]string](db, "Tag")))

	queryService := query.NewService(
		query.NewRepository(storage.NewDbBucket[query.Response](db, "Query"), storage.NewFileStorage("query-logs", ".log")),
		utils.NewRegistry[database.PluginType, database.Adapter](),
		vaultService,
		certService,
		"ivory",
	)

	deploymentService := deployment.NewService(
		deployment.NewRepository(storage.NewDbBucket[deployment.Template](db, "DeploymentTemplate")),
		utils.NewRegistry[keeper.PluginType, keeper.Plugin](),
		utils.NewRegistry[platform.PluginType, platform.Plugin](),
	)

	clusterService := cluster.NewService(
		cluster.NewRepository(storage.NewDbBucket[cluster.Response](db, "Cluster")),
		nil,
		tagService,
		queryService,
		vaultService,
		utils.NewRegistry[tools.Tool, tools.Adapter](),
	)

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
	basicProvider := basic.NewProvider(userService)
	ldapProvider := ldap.NewProvider()
	oidcProvider := oidc.NewProvider()

	configService := config.NewService(
		storage.NewFileStorage("config", ""),
		encryption.NewService(),
		secretService,
		nil,
		userService,
		basicProvider,
		ldapProvider,
		oidcProvider,
	)

	authService := auth.NewService(token.NewService(secretService), basicProvider, ldapProvider, oidcProvider, userService)

	backupService := backup.NewService(clusterService, queryService, permissionService, deploymentService, userService)

	toolRegistry := utils.NewRegistry[tools.Tool, tools.Adapter]()

	env := &coreConfig.Environment{Version: coreConfig.Version{Tag: "v-test"}}

	service := NewService(
		env,
		authService,
		vaultService,
		clusterService,
		certService,
		tagService,
		nil,
		queryService,
		deploymentService,
		nil,
		secretService,
		configService,
		permissionService,
		backupService,
		userService,
		toolRegistry,
	)

	return &testManagementEnv{
		service:           service,
		secretService:     secretService,
		vaultService:      vaultService,
		certService:       certService,
		clusterService:    clusterService,
		tagService:        tagService,
		queryService:      queryService,
		configService:     configService,
		permissionService: permissionService,
		basicProvider:     basicProvider,
	}
}

func TestServiceErase(t *testing.T) {
	env := createTestManagementService(t)

	if _, _, err := env.vaultService.Create(vault.Vault{Type: vault.DATABASE_PASSWORD, Username: "u", Secret: "pw"}); err != nil {
		t.Fatalf("failed to seed vault: %v", err)
	}
	port := 5432
	if _, err := env.clusterService.Update(cluster.Request{
		Name:  "c1",
		Nodes: []cluster.NodeConfig{{Name: "h1", Host: "h1", KeeperPort: &port}},
	}); err != nil {
		t.Fatalf("failed to seed cluster: %v", err)
	}
	if _, err := env.permissionService.CreateUserPermissions("alice", permission.NOT_PERMITTED); err != nil {
		t.Fatalf("failed to seed permissions: %v", err)
	}

	if err := env.service.Erase(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	vaults, errVaults := env.vaultService.Map(nil)
	if errVaults != nil {
		t.Fatalf("expected no error, got %v", errVaults)
	}
	if len(vaults) != 0 {
		t.Fatalf("expected all vaults erased, got %v", vaults)
	}

	clusters, errClusters := env.clusterService.List()
	if errClusters != nil {
		t.Fatalf("expected no error, got %v", errClusters)
	}
	if len(clusters) != 0 {
		t.Fatalf("expected all clusters erased, got %v", clusters)
	}

	perms, errPerms := env.permissionService.GetAllUserPermissions()
	if errPerms != nil {
		t.Fatalf("expected no error, got %v", errPerms)
	}
	if len(perms) != 0 {
		t.Fatalf("expected all permissions erased, got %v", perms)
	}

	if !env.secretService.IsRefEmpty() {
		t.Fatalf("expected the secret ref to be cleared")
	}
}

func TestServiceFree(t *testing.T) {
	env := createTestManagementService(t)
	if err := env.service.Free(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestServiceChangeSecret(t *testing.T) {
	env := createTestManagementService(t)

	if err := env.configService.SetAppConfig(config.NewAppConfig{AppConfig: config.AppConfig{Company: "Acme"}}); err != nil {
		t.Fatalf("failed to set app config: %v", err)
	}
	if _, err := env.configService.GetAppConfig(); err != nil {
		t.Fatalf("failed to load app config: %v", err)
	}

	key, _, errCreate := env.vaultService.Create(vault.Vault{Type: vault.DATABASE_PASSWORD, Username: "u", Secret: "pw"})
	if errCreate != nil {
		t.Fatalf("failed to seed vault: %v", errCreate)
	}

	if err := env.service.ChangeSecret("ivory", "newsecret"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !env.secretService.Verify("newsecret") {
		t.Fatalf("expected the secret to be rotated to 'newsecret'")
	}

	decrypted, errDec := env.vaultService.GetDecrypted(*key)
	if errDec != nil {
		t.Fatalf("expected the vault to still be decryptable after rotation, got %v", errDec)
	}
	if decrypted.Secret != "pw" {
		t.Fatalf("expected decrypted secret 'pw', got %q", decrypted.Secret)
	}
}

func TestServiceChangeSecretFailsWithWrongPreviousKey(t *testing.T) {
	env := createTestManagementService(t)
	if err := env.service.ChangeSecret("wrong", "newsecret"); err == nil {
		t.Fatalf("expected an error for a wrong previous key")
	}
}

func TestServiceGetAppInfoWhenNotConfigured(t *testing.T) {
	env := createTestManagementService(t)

	info := env.service.GetAppInfo(false, false, "", "")

	if info.Config.Configured {
		t.Fatalf("expected Configured=false, got %+v", info.Config)
	}
	if info.Config.Company != "Ivory" {
		t.Fatalf("expected the fallback company 'Ivory', got %q", info.Config.Company)
	}
	if info.Config.Error == "" {
		t.Fatalf("expected a non-empty config error")
	}
	if info.Auth.Authorised {
		t.Fatalf("expected Authorised=false")
	}
	if info.Auth.User != nil {
		t.Fatalf("expected a nil user, got %+v", info.Auth.User)
	}
	if len(info.Auth.Supported) != 0 {
		t.Fatalf("expected no supported auth types, got %v", info.Auth.Supported)
	}
}

func TestServiceGetAppInfoWhenConfigured(t *testing.T) {
	env := createTestManagementService(t)

	req := config.NewAppConfig{AppConfig: config.AppConfig{Company: "Acme"}}
	if err := env.configService.SetAppConfig(req); err != nil {
		t.Fatalf("failed to set app config: %v", err)
	}
	if _, err := env.configService.GetAppConfig(); err != nil {
		t.Fatalf("failed to load app config: %v", err)
	}

	t.Run("propagates an existing auth error without looking up permissions", func(t *testing.T) {
		info := env.service.GetAppInfo(false, true, "bob", "token expired")
		if info.Config.Company != "Acme" {
			t.Fatalf("expected company 'Acme', got %q", info.Config.Company)
		}
		if info.Auth.Error != "token expired" {
			t.Fatalf("expected the auth error to propagate, got %q", info.Auth.Error)
		}
		if info.Auth.User != nil {
			t.Fatalf("expected a nil user when there is an auth error, got %+v", info.Auth.User)
		}
	})

	t.Run("resolves permissions for an authorised user when auth is disabled", func(t *testing.T) {
		info := env.service.GetAppInfo(true, false, "bob", "")
		if info.Auth.Error != "" {
			t.Fatalf("expected no auth error, got %q", info.Auth.Error)
		}
		if info.Auth.User == nil || info.Auth.User.Username != "bob" {
			t.Fatalf("expected user 'bob', got %+v", info.Auth.User)
		}
		if !info.Auth.Authorised {
			t.Fatalf("expected Authorised=true")
		}
	})
}

func TestServiceBackupDelegation(t *testing.T) {
	env := createTestManagementService(t)

	if got := env.service.BackupFileName(); got != "ivory.v2.bak" {
		t.Fatalf("expected 'ivory.v2.bak', got %q", got)
	}

	data, err := env.service.BackupExport()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected non-empty backup data")
	}
}
