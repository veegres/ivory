package query

import (
	"ivory/clients/storage"
	"ivory/plugins/database"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

func createTestRepository(t *testing.T) *Repository {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "query-test-*")
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

	return NewRepository(storage.NewDbBucket[Response](db, "Query"), nil)
}

func createQuery(t *testing.T, r *Repository, name string, queryType Type, plugin DbPlugin, creation CreationType) {
	t.Helper()
	_, _, err := r.Create(Response{Name: name, Type: queryType, Plugin: plugin, Creation: creation, Default: "q", Custom: "q"})
	if err != nil {
		t.Fatalf("failed to create query %q: %v", name, err)
	}
}

func TestRepositoryListByFilter(t *testing.T) {
	repository := createTestRepository(t)
	createQuery(t, repository, "pg-bloat", BLOAT, database.POSTGRES, System)
	createQuery(t, repository, "pg-legacy", ACTIVITY, "", System) // record from before the plugin field existed
	createQuery(t, repository, "etcd-activity", ACTIVITY, database.ETCD, System)
	createQuery(t, repository, "etcd-manual", ACTIVITY, database.ETCD, Manual)

	postgres := DbPlugin(database.POSTGRES)
	etcd := DbPlugin(database.ETCD)
	activity := ACTIVITY

	tests := []struct {
		name      string
		queryType *Type
		plugin    *DbPlugin
		expected  []string
	}{
		{name: "no filters", expected: []string{"pg-bloat", "pg-legacy", "etcd-activity", "etcd-manual"}},
		{name: "postgres includes legacy empty plugin", plugin: &postgres, expected: []string{"pg-bloat", "pg-legacy"}},
		{name: "etcd only", plugin: &etcd, expected: []string{"etcd-activity", "etcd-manual"}},
		{name: "type only", queryType: &activity, expected: []string{"pg-legacy", "etcd-activity", "etcd-manual"}},
		{name: "type and plugin", queryType: &activity, plugin: &postgres, expected: []string{"pg-legacy"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := repository.ListByFilter(tt.queryType, tt.plugin)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			names := make([]string, 0, len(list))
			for _, q := range list {
				names = append(names, q.Name)
			}
			if len(names) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, names)
			}
			for i, name := range tt.expected {
				if names[i] != name {
					t.Fatalf("expected %v, got %v", tt.expected, names)
				}
			}
		})
	}
}

func TestRepositoryHasSystemQueriesForPlugin(t *testing.T) {
	repository := createTestRepository(t)
	createQuery(t, repository, "pg-legacy", ACTIVITY, "", System)
	createQuery(t, repository, "etcd-manual", ACTIVITY, database.ETCD, Manual)

	t.Run("legacy empty plugin counts as postgres", func(t *testing.T) {
		exists, err := repository.HasSystemQueriesForPlugin(database.POSTGRES)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !exists {
			t.Error("expected postgres system queries to exist")
		}
	})

	t.Run("manual queries do not count", func(t *testing.T) {
		exists, err := repository.HasSystemQueriesForPlugin(database.ETCD)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if exists {
			t.Error("expected no etcd system queries")
		}
	})
}
