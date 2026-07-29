package query

import (
	"ivory/clients/storage"
	"ivory/plugins/database"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

// createTestRepositoryWithLogs is like createTestRepository but backs
// queryLogFiles with a real FileStorage, needed by anything that touches
// query logs (e.g. Service.Delete calling HasLog).
func createTestRepositoryWithLogs(t *testing.T) *Repository {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "query-service-crud-test-*")
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

	return NewRepository(storage.NewDbBucket[Response](db, "Query"), storage.NewFileStorage("query-logs", ".log"))
}

func TestServiceCreateDefaultsPlugin(t *testing.T) {
	repository := createTestRepository(t)
	s := &Service{repository: repository}

	queryType := ACTIVITY
	_, created, err := s.Create(Manual, Request{Name: "no-plugin", Type: &queryType, Query: "select 1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.Plugin != database.POSTGRES {
		t.Errorf("expected plugin to default to postgres, got %q", created.Plugin)
	}

	_, createdEtcd, errEtcd := s.Create(Manual, Request{Name: "etcd", Type: &queryType, Plugin: database.ETCD, Query: "member list"})
	if errEtcd != nil {
		t.Fatalf("expected no error, got %v", errEtcd)
	}
	if createdEtcd.Plugin != database.ETCD {
		t.Errorf("expected etcd plugin, got %q", createdEtcd.Plugin)
	}
}

func TestServiceCreateRequiresAllFields(t *testing.T) {
	repository := createTestRepository(t)
	s := &Service{repository: repository}
	queryType := ACTIVITY

	tests := []struct {
		name string
		req  Request
	}{
		{name: "missing name", req: Request{Type: &queryType, Query: "select 1"}},
		{name: "missing type", req: Request{Name: "n", Query: "select 1"}},
		{name: "missing query", req: Request{Name: "n", Type: &queryType}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := s.Create(Manual, tt.req); err != ErrAllFieldsRequired {
				t.Fatalf("expected ErrAllFieldsRequired, got %v", err)
			}
		})
	}
}

func TestServiceUpdateKeepsPluginImmutable(t *testing.T) {
	repository := createTestRepository(t)
	s := &Service{repository: repository}

	queryType := ACTIVITY
	id, _, err := s.Create(Manual, Request{Name: "etcd", Type: &queryType, Plugin: database.ETCD, Query: "member list"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, updated, errUpdate := s.Update(*id, Request{Name: "etcd", Type: &queryType, Plugin: database.POSTGRES, Query: "alarm list"})
	if errUpdate != nil {
		t.Fatalf("expected no error, got %v", errUpdate)
	}
	if updated.Plugin != database.ETCD {
		t.Errorf("expected plugin to stay etcd, got %q", updated.Plugin)
	}
}

func TestServiceUpdateManualQueryAllowsChangingNameTypeAndDescription(t *testing.T) {
	repository := createTestRepository(t)
	s := &Service{repository: repository}

	bloatType := BLOAT
	desc := "old"
	id, _, err := s.Create(Manual, Request{Name: "old-name", Type: &bloatType, Description: &desc, Query: "select 1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	activityType := ACTIVITY
	newDesc := "new"
	_, updated, errUpdate := s.Update(*id, Request{Name: "new-name", Type: &activityType, Description: &newDesc, Query: "select 2"})
	if errUpdate != nil {
		t.Fatalf("expected no error, got %v", errUpdate)
	}
	if updated.Name != "new-name" || updated.Type != ACTIVITY || *updated.Description != "new" {
		t.Fatalf("expected the manual query to be fully updatable, got %+v", updated)
	}
	if updated.Custom != "select 2" {
		t.Fatalf("expected custom query to update, got %q", updated.Custom)
	}
}

func TestServiceUpdateSystemQueryRestrictions(t *testing.T) {
	repository := createTestRepository(t)
	s := &Service{repository: repository}

	bloatType := BLOAT
	desc := "system desc"
	id, _, err := s.Create(System, Request{Name: "sys-query", Type: &bloatType, Description: &desc, Query: "select 1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("name change is rejected", func(t *testing.T) {
		if _, _, err := s.Update(*id, Request{Name: "renamed", Type: &bloatType, Description: &desc, Query: "select 2"}); err != ErrNameChangeNotAllowed {
			t.Fatalf("expected ErrNameChangeNotAllowed, got %v", err)
		}
	})

	t.Run("type change is rejected", func(t *testing.T) {
		activityType := ACTIVITY
		if _, _, err := s.Update(*id, Request{Name: "sys-query", Type: &activityType, Description: &desc, Query: "select 2"}); err != ErrTypeChangeNotAllowed {
			t.Fatalf("expected ErrTypeChangeNotAllowed, got %v", err)
		}
	})

	t.Run("description change is rejected", func(t *testing.T) {
		otherDesc := "changed"
		if _, _, err := s.Update(*id, Request{Name: "sys-query", Type: &bloatType, Description: &otherDesc, Query: "select 2"}); err != ErrDescriptionChangeNotAllowed {
			t.Fatalf("expected ErrDescriptionChangeNotAllowed, got %v", err)
		}
	})

	t.Run("custom query text can still be updated", func(t *testing.T) {
		_, updated, errUpdate := s.Update(*id, Request{Name: "sys-query", Type: &bloatType, Description: &desc, Query: "select 2"})
		if errUpdate != nil {
			t.Fatalf("expected no error, got %v", errUpdate)
		}
		if updated.Custom != "select 2" {
			t.Fatalf("expected custom to update to 'select 2', got %q", updated.Custom)
		}
		if updated.Default != "select 1" {
			t.Fatalf("expected default to remain 'select 1', got %q", updated.Default)
		}
	})
}

func TestServiceGetListNormalizesLegacyPlugin(t *testing.T) {
	repository := createTestRepository(t)
	s := &Service{repository: repository}
	createQuery(t, repository, "pg-legacy", ACTIVITY, "", System)

	list, err := s.GetList(nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 query, got %d", len(list))
	}
	if list[0].Plugin != database.POSTGRES {
		t.Errorf("expected normalized postgres plugin, got %q", list[0].Plugin)
	}
}

func TestServiceDelete(t *testing.T) {
	repository := createTestRepositoryWithLogs(t)
	s := &Service{repository: repository}

	bloatType := BLOAT
	manualId, _, err := s.Create(Manual, Request{Name: "manual", Type: &bloatType, Query: "select 1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	systemId, _, err := s.Create(System, Request{Name: "sys", Type: &bloatType, Query: "select 1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("manual query can be deleted", func(t *testing.T) {
		if err := s.Delete(*manualId); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, err := repository.Get(*manualId); err == nil {
			t.Fatalf("expected the query to be gone")
		}
	})

	t.Run("system query deletion is restricted", func(t *testing.T) {
		if err := s.Delete(*systemId); err != ErrDeletionOfSystemQueriesRestricted {
			t.Fatalf("expected ErrDeletionOfSystemQueriesRestricted, got %v", err)
		}
	})

	t.Run("unknown id fails", func(t *testing.T) {
		fakeRepo := createTestRepository(t)
		fakeService := &Service{repository: fakeRepo}
		if err := fakeService.Delete(*manualId); err == nil {
			t.Fatalf("expected an error deleting an unknown query")
		}
	})
}

func TestServiceDeleteAllReseedsSystemQueries(t *testing.T) {
	env := createTestQueryService(t, &fakeDatabaseAdapter{
		systemRequests: []database.SystemRequest{
			{Name: "sys-query", Type: database.BLOAT, Query: "select 1", Description: "d"},
		},
	})

	bloatType := BLOAT
	if _, _, err := env.service.Create(Manual, Request{Name: "manual", Type: &bloatType, Query: "select 1"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := env.service.DeleteAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	list, err := env.service.GetList(nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(list) != 1 || list[0].Name != "sys-query" {
		t.Fatalf("expected only the reseeded system query to remain, got %v", list)
	}
}
