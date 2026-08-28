package deployment

import (
	"errors"
	"ivory/clients/storage"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

func newTestRepository(t *testing.T) *Repository {
	t.Helper()

	db, err := bolt.Open(filepath.Join(t.TempDir(), "test.db"), 0600, nil)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return NewRepository(storage.NewDbBucket[Template](db, "DeploymentTemplate"))
}

const testKeeper KeeperPlugin = "native_etcd"
const testPlatform PlatformPlugin = "docker"

func testTemplate(name string) Template {
	return Template{
		Name:     name,
		Keeper:   testKeeper,
		Platform: testPlatform,
		Commands: []TemplateCommand{{Command: "docker run -d --name {{name}} etcd"}},
	}
}

func TestRepository_CreateAndGet(t *testing.T) {
	r := newTestRepository(t)

	created, err := r.Create(uuid.New(), testTemplate("mine"))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.Id == "" {
		t.Fatal("expected the created template to carry an id")
	}

	id, err := uuid.Parse(created.Id)
	if err != nil {
		t.Fatalf("expected a uuid id, got %q", created.Id)
	}

	got, err := r.Get(id)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Name != "mine" || len(got.Commands) != 1 {
		t.Errorf("Get() = %+v, want the created template", got)
	}
}

func TestRepository_GetByName(t *testing.T) {
	r := newTestRepository(t)
	if _, err := r.Create(uuid.New(), testTemplate("mine")); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	t.Run("existing name", func(t *testing.T) {
		got, err := r.GetByName("mine", testKeeper, testPlatform)
		if err != nil {
			t.Fatalf("GetByName() error = %v", err)
		}
		if got == nil || got.Name != "mine" {
			t.Errorf("GetByName() = %+v, want the stored template", got)
		}
	})

	// NOTE: a name is only taken within the list it was picked in, so the same
	// one under another keeper is free
	t.Run("same name under another keeper is free", func(t *testing.T) {
		got, err := r.GetByName("mine", "native_redis", testPlatform)
		if err != nil {
			t.Fatalf("GetByName() error = %v", err)
		}
		if got != nil {
			t.Errorf("GetByName() = %+v, want nil", got)
		}
	})

	t.Run("free name returns no template and no error", func(t *testing.T) {
		got, err := r.GetByName("other", testKeeper, testPlatform)
		if err != nil {
			t.Fatalf("GetByName() error = %v", err)
		}
		if got != nil {
			t.Errorf("GetByName() = %+v, want nil", got)
		}
	})
}

func TestRepository_ListFiltersByPlugin(t *testing.T) {
	r := newTestRepository(t)

	etcdTemplate := testTemplate("etcd-one")
	redisTemplate := testTemplate("redis-one")
	redisTemplate.Keeper = "native_redis"
	if _, err := r.Create(uuid.New(), etcdTemplate); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := r.Create(uuid.New(), redisTemplate); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	tests := []struct {
		name     string
		criteria ListRequest
		expected int
	}{
		{name: "no filter", criteria: ListRequest{}, expected: 2},
		{name: "by keeper", criteria: ListRequest{Keeper: keeperPtr("native_etcd")}, expected: 1},
		{name: "by platform", criteria: ListRequest{Platform: platformPtr("docker")}, expected: 2},
		{name: "by unknown platform", criteria: ListRequest{Platform: platformPtr("k8s")}, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := r.List(tt.criteria)
			if err != nil {
				t.Fatalf("List() error = %v", err)
			}
			if len(list) != tt.expected {
				t.Errorf("List() returned %d templates, want %d", len(list), tt.expected)
			}
		})
	}
}

func TestRepository_UpdateAndDelete(t *testing.T) {
	r := newTestRepository(t)
	created, err := r.Create(uuid.New(), testTemplate("mine"))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	id := uuid.MustParse(created.Id)

	t.Run("update keeps the id", func(t *testing.T) {
		updated := testTemplate("renamed")
		got, err := r.Update(id, updated)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}
		if got.Id != created.Id || got.Name != "renamed" {
			t.Errorf("Update() = %+v, want the renamed template under the same id", got)
		}
	})

	t.Run("delete removes it", func(t *testing.T) {
		if err := r.Delete(id); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}
		if _, err := r.Get(id); !errors.Is(err, storage.ErrNotFound) {
			t.Errorf("expected ErrNotFound after delete, got %v", err)
		}
	})
}

func keeperPtr(v KeeperPlugin) *KeeperPlugin       { return &v }
func platformPtr(v PlatformPlugin) *PlatformPlugin { return &v }

// TestRepository_BackfillsCreation covers a template stored before the field
// existed: only the shipped templates are ever system, and those are computed
// rather than stored, so anything in the bucket is manual by definition.
func TestRepository_BackfillsCreation(t *testing.T) {
	r := newTestRepository(t)

	stored := testTemplate("legacy")
	stored.Creation = ""
	created, err := r.Create(uuid.New(), stored)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	t.Run("get", func(t *testing.T) {
		got, err := r.Get(uuid.MustParse(created.Id))
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if got.Creation != Manual {
			t.Errorf("Get() creation = %q, want manual", got.Creation)
		}
	})

	t.Run("list", func(t *testing.T) {
		list, err := r.List(ListRequest{})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if len(list) != 1 || list[0].Creation != Manual {
			t.Errorf("List() = %+v, want one manual template", list)
		}
	})
}

// TestRepository_RenamesStoredPlatform covers a template stored before the
// linux platform was renamed to docker: it has to keep listing and reading back
// under the current name, filter included, or it disappears from the UI.
func TestRepository_RenamesStoredPlatform(t *testing.T) {
	r := newTestRepository(t)

	stored := testTemplate("legacy")
	stored.Platform = "linux"
	created, err := r.Create(uuid.New(), stored)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	t.Run("get", func(t *testing.T) {
		got, err := r.Get(uuid.MustParse(created.Id))
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if got.Platform != testPlatform {
			t.Errorf("Get() platform = %q, want %q", got.Platform, testPlatform)
		}
	})

	t.Run("list filtered by the current platform", func(t *testing.T) {
		list, err := r.List(ListRequest{Platform: platformPtr(testPlatform)})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if len(list) != 1 {
			t.Fatalf("List() = %d templates, want 1", len(list))
		}
		if list[0].Platform != testPlatform {
			t.Errorf("List() platform = %q, want %q", list[0].Platform, testPlatform)
		}
	})

	t.Run("get by name within the current platform", func(t *testing.T) {
		got, err := r.GetByName("legacy", testKeeper, testPlatform)
		if err != nil {
			t.Fatalf("GetByName() error = %v", err)
		}
		if got == nil {
			t.Fatalf("GetByName() = nil, want the stored template")
		}
		if got.Platform != testPlatform {
			t.Errorf("GetByName() platform = %q, want %q", got.Platform, testPlatform)
		}
	})
}
