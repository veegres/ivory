package tag

import (
	"ivory/clients/storage"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/boltdb/bolt"
)

func createTestTagService(t *testing.T) *Service {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "tag-service-test-*")
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

	return NewService(NewRepository(storage.NewDbBucket[[]string](db, "Tag")))
}

func TestServiceUpdateCluster(t *testing.T) {
	s := createTestTagService(t)

	t.Run("assigns lowercased tags to a cluster", func(t *testing.T) {
		tags, err := s.UpdateCluster("cluster1", []string{"PROD", "Web"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		sort.Strings(tags)
		if len(tags) != 2 || tags[0] != "prod" || tags[1] != "web" {
			t.Fatalf("expected lowercased tags, got %v", tags)
		}

		clusters, errGet := s.Get("prod")
		if errGet != nil {
			t.Fatalf("expected no error, got %v", errGet)
		}
		if len(clusters) != 1 || clusters[0] != "cluster1" {
			t.Fatalf("expected cluster1 tagged as prod, got %v", clusters)
		}
	})

	t.Run("moving a cluster to new tags removes it from old ones", func(t *testing.T) {
		if _, err := s.UpdateCluster("cluster1", []string{"staging"}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if _, err := s.Get("prod"); err == nil {
			t.Fatalf("expected the emptied 'prod' tag to be deleted")
		}
		if _, err := s.Get("web"); err == nil {
			t.Fatalf("expected the emptied 'web' tag to be deleted")
		}

		clusters, err := s.Get("staging")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(clusters) != 1 || clusters[0] != "cluster1" {
			t.Fatalf("expected cluster1 tagged as staging, got %v", clusters)
		}
	})

	t.Run("a tag shared by multiple clusters keeps the other cluster when one moves away", func(t *testing.T) {
		if _, err := s.UpdateCluster("cluster2", []string{"staging"}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, err := s.UpdateCluster("cluster1", []string{"prod"}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		clusters, err := s.Get("staging")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(clusters) != 1 || clusters[0] != "cluster2" {
			t.Fatalf("expected only cluster2 tagged as staging, got %v", clusters)
		}
	})

	t.Run("empty tag list removes the cluster from every tag", func(t *testing.T) {
		if _, err := s.UpdateCluster("cluster1", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, err := s.Get("prod"); err == nil {
			t.Fatalf("expected the 'prod' tag to be deleted once empty")
		}
	})
}

func TestServiceListAndGetMap(t *testing.T) {
	s := createTestTagService(t)
	if _, err := s.UpdateCluster("cluster1", []string{"prod", "web"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := s.UpdateCluster("cluster2", []string{"prod"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	list, err := s.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	sort.Strings(list)
	if len(list) != 2 || list[0] != "prod" || list[1] != "web" {
		t.Fatalf("expected [prod web], got %v", list)
	}

	tagMap, err := s.GetMap()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	sort.Strings(tagMap["prod"])
	if len(tagMap["prod"]) != 2 || tagMap["prod"][0] != "cluster1" || tagMap["prod"][1] != "cluster2" {
		t.Fatalf("expected prod tagged with both clusters, got %v", tagMap["prod"])
	}
}

func TestServiceDelete(t *testing.T) {
	s := createTestTagService(t)
	if _, err := s.UpdateCluster("cluster1", []string{"PROD"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := s.Delete("PROD"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := s.Get("prod"); err == nil {
		t.Fatalf("expected the tag to be deleted regardless of case")
	}
}

func TestServiceDeleteAll(t *testing.T) {
	s := createTestTagService(t)
	if _, err := s.UpdateCluster("cluster1", []string{"prod", "web"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := s.DeleteAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	list, err := s.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected no tags, got %v", list)
	}
}
