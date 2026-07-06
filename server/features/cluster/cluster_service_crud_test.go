package cluster

import (
	"ivory/clients/storage"
	"ivory/features/node"
	"ivory/features/query"
	"ivory/features/tag"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

func createTestService(t *testing.T) *Service {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "cluster-service-test-*")
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

	clusterRepository := NewRepository(storage.NewDbBucket[Response](db, "Cluster"))
	tagRepository := tag.NewRepository(storage.NewDbBucket[[]string](db, "Tag"))
	tagService := tag.NewService(tagRepository)

	return &Service{clusterRepository: clusterRepository, tagService: tagService}
}

func TestServiceSearch(t *testing.T) {
	s := createTestService(t)

	_, errPg := s.clusterRepository.Create(Request{
		Name:    "patroni-pg",
		Options: Options{Plugins: Plugins{Keeper: keeper.PATRONI_POSTGRES, Database: database.POSTGRES}, Tags: []string{"prod"}},
	})
	if errPg != nil {
		t.Fatalf("failed to create cluster: %v", errPg)
	}
	_, errEtcd := s.clusterRepository.Create(Request{
		Name:    "native-etcd",
		Options: Options{Plugins: Plugins{Keeper: keeper.NATIVE_ETCD, Database: database.ETCD}, Tags: []string{"staging"}},
	})
	if errEtcd != nil {
		t.Fatalf("failed to create cluster: %v", errEtcd)
	}
	if _, err := s.tagService.UpdateCluster("patroni-pg", []string{"prod"}); err != nil {
		t.Fatalf("failed to tag cluster: %v", err)
	}
	if _, err := s.tagService.UpdateCluster("native-etcd", []string{"staging"}); err != nil {
		t.Fatalf("failed to tag cluster: %v", err)
	}

	postgres := query.DbPlugin(database.POSTGRES)
	patroniPostgres := node.KeeperPlugin(keeper.PATRONI_POSTGRES)

	t.Run("no filters returns everything", func(t *testing.T) {
		list, err := s.Search(SearchRequest{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(list) != 2 {
			t.Fatalf("expected 2 clusters, got %d", len(list))
		}
	})

	t.Run("tag narrows to matching clusters", func(t *testing.T) {
		list, err := s.Search(SearchRequest{Tags: []string{"prod"}})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(list) != 1 || list[0].Name != "patroni-pg" {
			t.Fatalf("expected only patroni-pg, got %v", list)
		}
	})

	t.Run("unknown tag returns no results instead of falling back to all", func(t *testing.T) {
		list, err := s.Search(SearchRequest{Tags: []string{"unknown"}})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(list) != 0 {
			t.Fatalf("expected no clusters, got %v", list)
		}
	})

	t.Run("plugin filters combine with tag filter", func(t *testing.T) {
		list, err := s.Search(SearchRequest{Tags: []string{"prod"}, Database: &postgres})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(list) != 1 || list[0].Name != "patroni-pg" {
			t.Fatalf("expected only patroni-pg, got %v", list)
		}
	})

	t.Run("plugin filter alone narrows without tags", func(t *testing.T) {
		list, err := s.Search(SearchRequest{Keeper: &patroniPostgres})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(list) != 1 || list[0].Name != "patroni-pg" {
			t.Fatalf("expected only patroni-pg, got %v", list)
		}
	})
}
