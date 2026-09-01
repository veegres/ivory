package backup

import (
	"encoding/json"
	"ivory/clients/storage"
	"ivory/core/utils"
	"ivory/features/cluster"
	"ivory/features/deployment"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/features/tag"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"ivory/plugins/keeper/etcd"
	"ivory/plugins/keeper/patroni"
	"ivory/plugins/platform"
	"ivory/plugins/platform/docker"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"bytes"

	"github.com/boltdb/bolt"
)

func createTestBackupService(t *testing.T) *Service {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "backup-service-test-*")
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

	clusterRepository := cluster.NewRepository(storage.NewDbBucket[cluster.Response](db, "Cluster"))
	tagRepository := tag.NewRepository(storage.NewDbBucket[[]string](db, "Tag"))
	tagService := tag.NewService(tagRepository)
	clusterService := cluster.NewService(clusterRepository, nil, tagService, nil, nil, nil)

	queryRepository := query.NewRepository(storage.NewDbBucket[query.Response](db, "Query"), nil)
	queryService := query.NewService(queryRepository, utils.NewRegistry[database.PluginType, database.Adapter](), nil, nil, "ivory")

	permissionRepository := permission.NewRepository(storage.NewDbBucket[permission.PermissionMap](db, "Permission"))
	permissionService := permission.NewService(permissionRepository)

	deploymentRepository := deployment.NewRepository(storage.NewDbBucket[deployment.Template](db, "DeploymentTemplate"))
	keeperRegistry := utils.NewRegistry[keeper.PluginType, keeper.Plugin]()
	keeperRegistry.Register(keeper.NATIVE_ETCD, etcd.NewPlugin())
	// NOTE: the frozen backup files carry a patroni template, which the
	// deployment service refuses to create unless its keeper is registered
	keeperRegistry.Register(keeper.PATRONI_POSTGRES, patroni.NewPlugin(nil))
	platformRegistry := utils.NewRegistry[platform.PluginType, platform.Plugin]()
	platformRegistry.Register(platform.Docker, docker.NewPlugin(nil))
	deploymentService := deployment.NewService(deploymentRepository, keeperRegistry, platformRegistry)

	return NewService(clusterService, queryService, permissionService, deploymentService)
}

// readGolden returns one of the frozen backup files under testdata. They are
// the artifact this package's compatibility tests are written against: a file
// some older Ivory wrote, checked in so it can never drift with the code the
// way a struct literal built in a test does. Generated once, then frozen -
// a test that fails against one is the code moving, not the file.
func readGolden(t *testing.T, name string) []byte {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("failed to read the frozen backup file %s: %v", name, err)
	}
	return data
}

// assertShapeIsFrozen re-marshals a frozen file through the version struct it
// belongs to and requires the bytes back unchanged. It is the sacred rule as a
// test: a released version's types never change again, so a renamed json tag
// drops a key, a removed field drops it too, and an added one appears with its
// zero value - every one of those shows up here as a diff. It says nothing
// about whether the values are right; that is what the import tests are for.
func assertShapeIsFrozen(t *testing.T, name string, model any) {
	t.Helper()

	frozen := readGolden(t, name)
	if err := json.Unmarshal(frozen, model); err != nil {
		t.Fatalf("the frozen file no longer parses as its own version: %v", err)
	}
	got, errMarshal := json.MarshalIndent(model, "", "  ")
	if errMarshal != nil {
		t.Fatalf("failed to re-marshal: %v", errMarshal)
	}
	if string(append(got, '\n')) != string(frozen) {
		t.Errorf("the wire shape of %s changed - a released version is frozen, so this means a new version, not an edit\n--- frozen file\n%s\n--- re-marshalled\n%s",
			name, frozen, got)
	}
}

func createMultipartFile(t *testing.T, fieldFilename string, content []byte) *multipart.FileHeader {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, errPart := writer.CreateFormFile("file", fieldFilename)
	if errPart != nil {
		t.Fatalf("failed to create form file: %v", errPart)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("failed to write content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(32 << 20); err != nil {
		t.Fatalf("failed to parse multipart form: %v", err)
	}
	files := req.MultipartForm.File["file"]
	if len(files) != 1 {
		t.Fatalf("expected exactly one file, got %d", len(files))
	}
	return files[0]
}

func TestServiceGetFileName(t *testing.T) {
	s := createTestBackupService(t)
	if got := s.GetFileName(); got != "ivory.v2.bak" {
		t.Fatalf("expected 'ivory.v2.bak', got %q", got)
	}
}

func TestServiceExportProducesValidLatestJSON(t *testing.T) {
	s := createTestBackupService(t)

	port := 5432
	if _, err := s.clusterService.Update(cluster.Request{
		Name: "cluster1",
		Options: cluster.Options{
			Plugins: cluster.Plugins{Keeper: keeper.PATRONI_POSTGRES, Database: database.POSTGRES},
			Tags:    []string{"prod"},
		},
		Nodes: []cluster.NodeConfig{{Name: "host1", Host: "host1", KeeperPort: &port}},
	}); err != nil {
		t.Fatalf("failed to seed cluster: %v", err)
	}

	data, errExport := s.Export()
	if errExport != nil {
		t.Fatalf("expected no error, got %v", errExport)
	}

	var backupModel BackupV2
	if err := json.Unmarshal(data, &backupModel); err != nil {
		t.Fatalf("expected exported data to be valid BackupV2 JSON: %v", err)
	}
	if len(backupModel.Clusters) != 1 || backupModel.Clusters[0].Name != "cluster1" {
		t.Fatalf("expected exported cluster1, got %v", backupModel.Clusters)
	}
}

func TestServiceImportDispatchesByFilename(t *testing.T) {
	backupModel := BackupV1{Clusters: []backupClusterV1{{Name: "dispatched-cluster", Sidecars: []backupSidecarV1{{Host: "h1", Port: 1}}}}}
	data, errMarshal := json.Marshal(backupModel)
	if errMarshal != nil {
		t.Fatalf("failed to marshal backup model: %v", errMarshal)
	}

	t.Run("filename with .v1. is imported as v1", func(t *testing.T) {
		s := createTestBackupService(t)
		if err := s.Import(createMultipartFile(t, "ivory.v1.bak", data)); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		clusters, _ := s.clusterService.List()
		if len(clusters) != 1 {
			t.Fatalf("expected 1 cluster imported, got %v", clusters)
		}
	})

	t.Run("filename without any version segment falls back to v1", func(t *testing.T) {
		s := createTestBackupService(t)
		if err := s.Import(createMultipartFile(t, "ivory.bak", data)); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		clusters, _ := s.clusterService.List()
		if len(clusters) != 1 {
			t.Fatalf("expected 1 cluster imported, got %v", clusters)
		}
	})

	t.Run("filename with .v2. is imported as v2", func(t *testing.T) {
		s := createTestBackupService(t)
		port := 1
		v2, errMarshal := json.Marshal(BackupV2{Clusters: []backupClusterV2{{
			Name:     "dispatched-cluster",
			Keeper:   string(keeper.NATIVE_ETCD),
			Database: string(database.ETCD),
			Nodes:    []backupNodeV2{{Name: "n1", Host: "h1", KeeperPort: &port}},
		}}})
		if errMarshal != nil {
			t.Fatalf("failed to marshal backup model: %v", errMarshal)
		}
		if err := s.Import(createMultipartFile(t, "ivory.v2.bak", v2)); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		clusters, _ := s.clusterService.List()
		if len(clusters) != 1 {
			t.Fatalf("expected 1 cluster imported, got %v", clusters)
		}
	})

	t.Run("unsupported version in filename is rejected", func(t *testing.T) {
		s := createTestBackupService(t)
		if err := s.Import(createMultipartFile(t, "ivory.v9.bak", data)); err == nil {
			t.Fatalf("expected an error for an unsupported backup version")
		}
	})

	t.Run("malformed data is rejected", func(t *testing.T) {
		s := createTestBackupService(t)
		if err := s.Import(createMultipartFile(t, "ivory.v1.bak", []byte("not-json"))); err == nil {
			t.Fatalf("expected an error for malformed backup data")
		}
	})
}
