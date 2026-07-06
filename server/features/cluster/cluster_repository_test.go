package cluster

import (
	"ivory/clients/storage"
	"ivory/features/node"
	"ivory/features/query"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

func createTestRepository(t *testing.T) *Repository {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "cluster-test-*")
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

	return NewRepository(storage.NewDbBucket[Response](db, "Cluster"))
}

func createCluster(t *testing.T, r *Repository, name string, k node.KeeperPlugin, db query.DbPlugin) {
	t.Helper()
	_, err := r.Create(Request{Name: name, Options: Options{Plugins: Plugins{Keeper: k, Database: db}}})
	if err != nil {
		t.Fatalf("failed to create cluster %q: %v", name, err)
	}
}

func TestRepositorySearch(t *testing.T) {
	repository := createTestRepository(t)
	createCluster(t, repository, "patroni-pg", keeper.PATRONI_POSTGRES, database.POSTGRES)
	createCluster(t, repository, "native-pg", keeper.NATIVE_POSTGRES, database.POSTGRES)
	createCluster(t, repository, "native-etcd", keeper.NATIVE_ETCD, database.ETCD)

	patroniPostgres := node.KeeperPlugin(keeper.PATRONI_POSTGRES)
	nativePostgres := node.KeeperPlugin(keeper.NATIVE_POSTGRES)
	postgres := query.DbPlugin(database.POSTGRES)
	etcd := query.DbPlugin(database.ETCD)

	tests := []struct {
		name     string
		criteria SearchCriteria
		expected []string
	}{
		{name: "no filters", criteria: SearchCriteria{}, expected: []string{"patroni-pg", "native-pg", "native-etcd"}},
		{name: "keeper only", criteria: SearchCriteria{Keeper: &nativePostgres}, expected: []string{"native-pg"}},
		{name: "database only", criteria: SearchCriteria{Database: &postgres}, expected: []string{"patroni-pg", "native-pg"}},
		{name: "keeper and database", criteria: SearchCriteria{Keeper: &patroniPostgres, Database: &postgres}, expected: []string{"patroni-pg"}},
		{name: "keeper and database mismatch", criteria: SearchCriteria{Keeper: &patroniPostgres, Database: &etcd}, expected: []string{}},
		{name: "names filter", criteria: SearchCriteria{Names: []string{"native-pg", "native-etcd"}}, expected: []string{"native-pg", "native-etcd"}},
		{name: "empty names filter excludes everything", criteria: SearchCriteria{Names: []string{}}, expected: []string{}},
		{name: "names and database", criteria: SearchCriteria{Names: []string{"native-pg", "native-etcd"}, Database: &etcd}, expected: []string{"native-etcd"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := repository.Search(tt.criteria)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			names := make([]string, 0, len(list))
			for _, c := range list {
				names = append(names, c.Name)
			}
			if len(names) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, names)
			}
			expectedSet := make(map[string]bool, len(tt.expected))
			for _, name := range tt.expected {
				expectedSet[name] = true
			}
			for _, name := range names {
				if !expectedSet[name] {
					t.Fatalf("expected %v, got %v", tt.expected, names)
				}
			}
		})
	}
}
