package backup

import (
	"encoding/json"
	"ivory/clients/storage"
	"ivory/core/utils"
	"ivory/features/cluster"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/features/tag"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
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
	queryService := query.NewService(queryRepository, utils.NewRegistry[database.Plugin, database.Adapter](), nil, nil, "ivory")

	permissionRepository := permission.NewRepository(storage.NewDbBucket[permission.PermissionMap](db, "Permission"))
	permissionService := permission.NewService(permissionRepository)

	return NewService(clusterService, queryService, permissionService)
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
	if got := s.GetFileName(); got != "ivory.v1.bak" {
		t.Fatalf("expected 'ivory.v1.bak', got %q", got)
	}
}

func TestServiceExportProducesValidV1JSON(t *testing.T) {
	s := createTestBackupService(t)

	port := 5432
	if _, err := s.clusterService.Update(cluster.Request{
		Name: "cluster1",
		Options: cluster.Options{
			Plugins: cluster.Plugins{Keeper: keeper.PATRONI_POSTGRES, Database: database.POSTGRES},
			Tags:    []string{"prod"},
		},
		Nodes: []cluster.NodeConfig{{Host: "host1", KeeperPort: &port}},
	}); err != nil {
		t.Fatalf("failed to seed cluster: %v", err)
	}

	data, errExport := s.Export()
	if errExport != nil {
		t.Fatalf("expected no error, got %v", errExport)
	}

	var backupModel BackupV1
	if err := json.Unmarshal(data, &backupModel); err != nil {
		t.Fatalf("expected exported data to be valid BackupV1 JSON: %v", err)
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

	t.Run("unsupported version in filename is rejected", func(t *testing.T) {
		s := createTestBackupService(t)
		if err := s.Import(createMultipartFile(t, "ivory.v2.bak", data)); err == nil {
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
