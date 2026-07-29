package query

import (
	"ivory/core/service/cert"
	"ivory/plugins/database"
	"testing"

	"github.com/google/uuid"
)

func TestServiceConsoleQuery(t *testing.T) {
	adapter := &fakeDatabaseAdapter{
		getFieldsFn: func(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error) {
			return &database.QueryFields{Rows: [][]any{{query}}}, nil
		},
	}
	env := createTestQueryService(t, adapter)

	t.Run("success maps the adapter response", func(t *testing.T) {
		resp, err := env.service.ConsoleQuery(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}}}, "select 1", nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Rows) != 1 || resp.Rows[0][0] != "select 1" {
			t.Fatalf("unexpected response: %+v", resp)
		}
	})

	t.Run("unknown plugin fails", func(t *testing.T) {
		_, err := env.service.ConsoleQuery(Context{Connection: Connection{Db: DbConfig{Plugin: database.ETCD}}}, "select 1", nil)
		if err == nil {
			t.Fatalf("expected an error for an unregistered plugin")
		}
	})
}

func TestServiceTemplateQuery(t *testing.T) {
	adapter := &fakeDatabaseAdapter{
		getFieldsFn: func(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error) {
			return &database.QueryFields{Rows: [][]any{{"row1"}}}, nil
		},
	}
	env := createTestQueryService(t, adapter)

	bloatType := BLOAT
	id, _, err := env.service.Create(Manual, Request{Name: "q1", Type: &bloatType, Query: "select 1"})
	if err != nil {
		t.Fatalf("failed to seed query: %v", err)
	}

	t.Run("runs the stored custom query and logs the result", func(t *testing.T) {
		resp, errRun := env.service.TemplateQuery(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}}}, *id, nil)
		if errRun != nil {
			t.Fatalf("expected no error, got %v", errRun)
		}
		if len(resp.Rows) != 1 {
			t.Fatalf("unexpected response: %+v", resp)
		}
		log, errLog := env.service.GetLog(*id)
		if errLog != nil {
			t.Fatalf("expected no error, got %v", errLog)
		}
		if len(log) != 1 {
			t.Fatalf("expected the result to be logged, got %v", log)
		}
	})

	t.Run("unknown id fails", func(t *testing.T) {
		if _, err := env.service.TemplateQuery(Context{}, uuid.New(), nil); err == nil {
			t.Fatalf("expected an error for an unknown query id")
		}
	})

	t.Run("empty custom query is rejected", func(t *testing.T) {
		emptyId, _, errCreate := env.service.repository.Create(Response{Name: "empty", Type: BLOAT, Plugin: database.POSTGRES})
		if errCreate != nil {
			t.Fatalf("failed to seed empty query: %v", errCreate)
		}
		if _, err := env.service.TemplateQuery(Context{}, *emptyId, nil); err != ErrQueryEmpty {
			t.Fatalf("expected ErrQueryEmpty, got %v", err)
		}
	})
}

func TestServiceCancelAndTerminateQuery(t *testing.T) {
	cancelled := false
	terminated := false
	adapter := &fakeDatabaseAdapter{
		cancelFn:    func(ctx database.Context, pid int) error { cancelled = true; return nil },
		terminateFn: func(ctx database.Context, pid int) error { terminated = true; return nil },
	}
	env := createTestQueryService(t, adapter)
	ctx := Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}}}

	if err := env.service.CancelQuery(ctx, 42); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !cancelled {
		t.Fatalf("expected Cancel to be called")
	}

	if err := env.service.TerminateQuery(ctx, 42); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !terminated {
		t.Fatalf("expected Terminate to be called")
	}

	t.Run("unknown plugin fails for both", func(t *testing.T) {
		badCtx := Context{Connection: Connection{Db: DbConfig{Plugin: database.ETCD}}}
		if err := env.service.CancelQuery(badCtx, 1); err == nil {
			t.Fatalf("expected an error for cancel")
		}
		if err := env.service.TerminateQuery(badCtx, 1); err == nil {
			t.Fatalf("expected an error for terminate")
		}
	})
}

func TestServiceRunningQueriesByApplicationName(t *testing.T) {
	var gotParams []any
	adapter := &fakeDatabaseAdapter{
		activeQueriesFn: func(ctx database.Context, options *database.QueryOptions) (*database.QueryFields, error) {
			gotParams = options.Params
			return &database.QueryFields{}, nil
		},
	}
	env := createTestQueryService(t, adapter)

	_, err := env.service.RunningQueriesByApplicationName(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}}, Session: "abcdefghij"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(gotParams) != 1 || gotParams[0] != "ivory [abcdefg]" {
		t.Fatalf("expected the application name to be passed as a param, got %v", gotParams)
	}
}

func TestServiceDatabasesQuery(t *testing.T) {
	adapter := &fakeDatabaseAdapter{
		listDatabasesFn: func(ctx database.Context, name string) ([]string, error) { return []string{"db1"}, nil },
	}
	env := createTestQueryService(t, adapter)

	dbs, err := env.service.DatabasesQuery(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}}}, "db")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(dbs) != 1 || dbs[0] != "db1" {
		t.Fatalf("unexpected result: %v", dbs)
	}
}

func TestServiceSchemasQuery(t *testing.T) {
	called := false
	adapter := &fakeDatabaseAdapter{
		listSchemasFn: func(ctx database.Context, name string) ([]string, error) {
			called = true
			return []string{"public"}, nil
		},
	}
	env := createTestQueryService(t, adapter)

	t.Run("no database name skips the adapter call", func(t *testing.T) {
		called = false
		schemas, err := env.service.SchemasQuery(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}}}, "s")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(schemas) != 0 {
			t.Fatalf("expected an empty result, got %v", schemas)
		}
		if called {
			t.Fatalf("expected the adapter to not be called")
		}
	})

	t.Run("database name set calls the adapter", func(t *testing.T) {
		called = false
		dbName := "mydb"
		schemas, err := env.service.SchemasQuery(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES, Name: &dbName}}}, "s")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(schemas) != 1 || schemas[0] != "public" {
			t.Fatalf("unexpected result: %v", schemas)
		}
		if !called {
			t.Fatalf("expected the adapter to be called")
		}
	})
}

func TestServiceTablesQuery(t *testing.T) {
	called := false
	adapter := &fakeDatabaseAdapter{
		listTablesFn: func(ctx database.Context, schema string, name string) ([]string, error) {
			called = true
			return []string{"t1"}, nil
		},
	}
	env := createTestQueryService(t, adapter)
	dbName := "mydb"

	t.Run("no schema skips the adapter call", func(t *testing.T) {
		called = false
		tables, err := env.service.TablesQuery(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES, Name: &dbName}}}, "", "t")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tables) != 0 || called {
			t.Fatalf("expected no adapter call and an empty result, got %v called=%v", tables, called)
		}
	})

	t.Run("schema and database name set calls the adapter", func(t *testing.T) {
		called = false
		tables, err := env.service.TablesQuery(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES, Name: &dbName}}}, "public", "t")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tables) != 1 || !called {
			t.Fatalf("expected the adapter to be called, got %v called=%v", tables, called)
		}
	})
}

func TestServiceChartQuery(t *testing.T) {
	adapter := &fakeDatabaseAdapter{
		systemCharts: map[database.SystemChartType]string{database.Databases: "select count(*) from pg_database"},
		getOneFn:     func(ctx database.Context, query string) (any, error) { return 3, nil },
	}
	env := createTestQueryService(t, adapter)

	t.Run("known chart for a known plugin succeeds", func(t *testing.T) {
		chart, err := env.service.ChartQuery(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}}}, database.Databases)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if chart.Value != 3 {
			t.Fatalf("expected value 3, got %v", chart.Value)
		}
	})

	t.Run("unsupported plugin fails", func(t *testing.T) {
		_, err := env.service.ChartQuery(Context{Connection: Connection{Db: DbConfig{Plugin: database.ETCD}}}, database.Databases)
		if err == nil {
			t.Fatalf("expected an error for an unsupported plugin")
		}
	})

	t.Run("unsupported chart type for a known plugin fails", func(t *testing.T) {
		_, err := env.service.ChartQuery(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}}}, database.TotalSize)
		if err == nil {
			t.Fatalf("expected an error for an unsupported chart type")
		}
	})
}

func TestServiceMapContextCerts(t *testing.T) {
	env := createTestQueryService(t, &fakeDatabaseAdapter{})

	t.Run("certs enrich the tls config", func(t *testing.T) {
		_, ctx, err := env.service.getDatabaseAdapter(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}, Certs: &cert.Certs{}}})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if ctx.Connection.TlsConfig == nil {
			t.Fatalf("expected a non-nil tls config")
		}
	})

	t.Run("bad ca cert id propagates the error", func(t *testing.T) {
		badId := uuid.New()
		_, _, err := env.service.getDatabaseAdapter(Context{Connection: Connection{Db: DbConfig{Plugin: database.POSTGRES}, Certs: &cert.Certs{ClientCAId: &badId}}})
		if err == nil {
			t.Fatalf("expected an error for an unknown ca cert")
		}
	})
}
