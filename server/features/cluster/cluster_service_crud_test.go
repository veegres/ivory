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

func TestServiceUpdate(t *testing.T) {
	s := createTestService(t)
	port := 5432

	t.Run("empty name is rejected", func(t *testing.T) {
		if _, err := s.Update(Request{Nodes: []NodeConfig{{Host: "h1", KeeperPort: &port}}}); err != ErrClusterNameEmpty {
			t.Fatalf("expected ErrClusterNameEmpty, got %v", err)
		}
	})

	t.Run("no nodes is rejected", func(t *testing.T) {
		if _, err := s.Update(Request{Name: "c1"}); err != ErrClusterKeepersEmpty {
			t.Fatalf("expected ErrClusterKeepersEmpty, got %v", err)
		}
	})

	t.Run("creates a new cluster and tags it", func(t *testing.T) {
		created, err := s.Update(Request{
			Name:  "c1",
			Nodes: []NodeConfig{{Host: "h1", KeeperPort: &port}},
			Options: Options{
				Tags: []string{"PROD"},
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(created.Tags) != 1 || created.Tags[0] != "prod" {
			t.Fatalf("expected the tag to be lowercased to 'prod', got %v", created.Tags)
		}

		got, errGet := s.Get("c1")
		if errGet != nil {
			t.Fatalf("expected no error, got %v", errGet)
		}
		if len(got.Nodes) != 1 || got.Nodes[0].Host != "h1" {
			t.Fatalf("expected 1 node with host h1, got %v", got.Nodes)
		}

		clusters, errTag := s.tagService.Get("prod")
		if errTag != nil {
			t.Fatalf("expected no error, got %v", errTag)
		}
		if len(clusters) != 1 || clusters[0] != "c1" {
			t.Fatalf("expected c1 tagged as prod, got %v", clusters)
		}
	})

	t.Run("updating retags the cluster", func(t *testing.T) {
		if _, err := s.Update(Request{
			Name:    "c1",
			Nodes:   []NodeConfig{{Host: "h1", KeeperPort: &port}},
			Options: Options{Tags: []string{"staging"}},
		}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, err := s.tagService.Get("prod"); err == nil {
			t.Fatalf("expected the old 'prod' tag to be gone")
		}
		clusters, err := s.tagService.Get("staging")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(clusters) != 1 || clusters[0] != "c1" {
			t.Fatalf("expected c1 tagged as staging, got %v", clusters)
		}
	})
}

func TestServiceGet(t *testing.T) {
	s := createTestService(t)
	if _, err := s.clusterRepository.Create(Request{Name: "c1"}); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	if _, err := s.Get("c1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := s.Get("unknown"); err == nil {
		t.Fatalf("expected an error for an unknown cluster")
	}
}

func TestServiceList(t *testing.T) {
	s := createTestService(t)
	if _, err := s.clusterRepository.Create(Request{Name: "c1"}); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}
	if _, err := s.clusterRepository.Create(Request{Name: "c2"}); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	list, err := s.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(list))
	}
}

func TestServiceDelete(t *testing.T) {
	s := createTestService(t)
	port := 5432
	if _, err := s.Update(Request{
		Name:    "c1",
		Nodes:   []NodeConfig{{Host: "h1", KeeperPort: &port}},
		Options: Options{Tags: []string{"prod"}},
	}); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	if err := s.Delete("c1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := s.Get("c1"); err == nil {
		t.Fatalf("expected the cluster to be gone")
	}
	if _, err := s.tagService.Get("prod"); err == nil {
		t.Fatalf("expected the cluster's tag association to be cleaned up")
	}
}

func TestServiceDeleteAll(t *testing.T) {
	s := createTestService(t)
	if _, err := s.clusterRepository.Create(Request{Name: "c1"}); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}
	if _, err := s.clusterRepository.Create(Request{Name: "c2"}); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	if err := s.DeleteAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	list, err := s.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected no clusters, got %v", list)
	}
}
