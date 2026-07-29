package query

import (
	"ivory/clients/console/ssh"
	"ivory/clients/storage"
	"ivory/core/config"
	"ivory/core/service/cert"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/plugins/database"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

// fakeDatabaseAdapter lets tests control every database.Adapter method
// individually; unset funcs return harmless zero values.
type fakeDatabaseAdapter struct {
	getManyFn        func(ctx database.Context, query string, params []any) ([]string, error)
	getOneFn         func(ctx database.Context, query string) (any, error)
	getFieldsFn      func(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error)
	listDatabasesFn  func(ctx database.Context, name string) ([]string, error)
	listSchemasFn    func(ctx database.Context, name string) ([]string, error)
	listTablesFn     func(ctx database.Context, schema string, name string) ([]string, error)
	cancelFn         func(ctx database.Context, pid int) error
	terminateFn      func(ctx database.Context, pid int) error
	activeQueriesFn  func(ctx database.Context, options *database.QueryOptions) (*database.QueryFields, error)
	supportedFeature map[config.Feature]bool
	systemRequests   []database.SystemRequest
	systemCharts     map[database.SystemChartType]string
}

func (f *fakeDatabaseAdapter) GetMany(ctx database.Context, query string, params []any) ([]string, error) {
	if f.getManyFn != nil {
		return f.getManyFn(ctx, query, params)
	}
	return nil, nil
}

func (f *fakeDatabaseAdapter) GetOne(ctx database.Context, query string) (any, error) {
	if f.getOneFn != nil {
		return f.getOneFn(ctx, query)
	}
	return nil, nil
}

func (f *fakeDatabaseAdapter) GetFields(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error) {
	if f.getFieldsFn != nil {
		return f.getFieldsFn(ctx, query, options)
	}
	return &database.QueryFields{}, nil
}

func (f *fakeDatabaseAdapter) ListDatabases(ctx database.Context, name string) ([]string, error) {
	if f.listDatabasesFn != nil {
		return f.listDatabasesFn(ctx, name)
	}
	return nil, nil
}

func (f *fakeDatabaseAdapter) ListSchemas(ctx database.Context, name string) ([]string, error) {
	if f.listSchemasFn != nil {
		return f.listSchemasFn(ctx, name)
	}
	return nil, nil
}

func (f *fakeDatabaseAdapter) ListTables(ctx database.Context, schema string, name string) ([]string, error) {
	if f.listTablesFn != nil {
		return f.listTablesFn(ctx, schema, name)
	}
	return nil, nil
}

func (f *fakeDatabaseAdapter) Cancel(ctx database.Context, pid int) error {
	if f.cancelFn != nil {
		return f.cancelFn(ctx, pid)
	}
	return nil
}

func (f *fakeDatabaseAdapter) Terminate(ctx database.Context, pid int) error {
	if f.terminateFn != nil {
		return f.terminateFn(ctx, pid)
	}
	return nil
}

func (f *fakeDatabaseAdapter) ActiveQueries(ctx database.Context, options *database.QueryOptions) (*database.QueryFields, error) {
	if f.activeQueriesFn != nil {
		return f.activeQueriesFn(ctx, options)
	}
	return &database.QueryFields{}, nil
}

func (f *fakeDatabaseAdapter) SupportedFeatures() map[config.Feature]bool {
	return f.supportedFeature
}

func (f *fakeDatabaseAdapter) SystemRequests() []database.SystemRequest {
	return f.systemRequests
}

func (f *fakeDatabaseAdapter) SystemCharts() map[database.SystemChartType]string {
	return f.systemCharts
}

type testQueryEnv struct {
	service      *Service
	vaultService *vault.Service
	certService  *cert.Service
	adapter      *fakeDatabaseAdapter
}

func createTestQueryService(t *testing.T, adapter *fakeDatabaseAdapter) *testQueryEnv {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "query-service-test-*")
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

	registry := utils.NewRegistry[database.Plugin, database.Adapter]()
	if adapter != nil {
		registry.Register(database.POSTGRES, adapter)
	}

	repository := NewRepository(storage.NewDbBucket[Response](db, "Query"), storage.NewFileStorage("query-logs", ".log"))
	service := NewService(repository, registry, vaultService, certService, "ivory")

	return &testQueryEnv{service: service, vaultService: vaultService, certService: certService, adapter: adapter}
}

func TestServiceGetApplicationName(t *testing.T) {
	env := createTestQueryService(t, nil)
	got := env.service.GetApplicationName("abcdefghij")
	if got != "ivory [abcdefg]" {
		t.Fatalf("expected 'ivory [abcdefg]', got %q", got)
	}
}

func TestServiceSupportedFeatures(t *testing.T) {
	env := createTestQueryService(t, &fakeDatabaseAdapter{supportedFeature: map[config.Feature]bool{config.ViewClusterList: true}})

	t.Run("known plugin returns its feature map", func(t *testing.T) {
		features := env.service.SupportedFeatures(database.POSTGRES)
		if !features[config.ViewClusterList] {
			t.Fatalf("expected ViewClusterList true, got %v", features)
		}
	})

	t.Run("unknown plugin returns an empty map", func(t *testing.T) {
		features := env.service.SupportedFeatures(database.ETCD)
		if len(features) != 0 {
			t.Fatalf("expected an empty map, got %v", features)
		}
	})
}

func TestServiceInitializeSystemQueriesSeedsOncePerPlugin(t *testing.T) {
	adapter := &fakeDatabaseAdapter{
		systemRequests: []database.SystemRequest{
			{Name: "sys-query", Type: database.BLOAT, Query: "select 1", Description: "d"},
		},
	}
	env := createTestQueryService(t, adapter)

	list, err := env.service.GetList(nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	found := 0
	for _, q := range list {
		if q.Name == "sys-query" {
			found++
		}
	}
	if found != 1 {
		t.Fatalf("expected the system query to be seeded exactly once, got %d", found)
	}

	// Re-running initializeSystemQueries (as NewService would on restart)
	// must not duplicate it.
	if err := env.service.initializeSystemQueries(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	list, err = env.service.GetList(nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	found = 0
	for _, q := range list {
		if q.Name == "sys-query" {
			found++
		}
	}
	if found != 1 {
		t.Fatalf("expected the system query to still appear exactly once, got %d", found)
	}
}

func TestServiceGetDatabaseAdapter(t *testing.T) {
	env := createTestQueryService(t, &fakeDatabaseAdapter{})

	t.Run("known plugin resolves the adapter", func(t *testing.T) {
		adapter, ctx, err := env.service.getDatabaseAdapter(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES, Host: "h", Port: 5432}}})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if adapter == nil {
			t.Fatalf("expected a non-nil adapter")
		}
		if ctx.Connection.Config.Host != "h" {
			t.Fatalf("unexpected context: %+v", ctx)
		}
	})

	t.Run("unknown plugin fails", func(t *testing.T) {
		_, _, err := env.service.getDatabaseAdapter(Context{Connection: Connection{Db: DbConfig{Plugin: database.ETCD}}})
		if err == nil {
			t.Fatalf("expected an error for an unregistered plugin")
		}
	})

	t.Run("bad vault id fails", func(t *testing.T) {
		badId := uuid.New()
		_, _, err := env.service.getDatabaseAdapter(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}, VaultId: &badId}})
		if err != ErrVaultProblems {
			t.Fatalf("expected ErrVaultProblems, got %v", err)
		}
	})
}
