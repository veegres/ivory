package backup

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestImportV1(t *testing.T) {
	s := createTestBackupService(t)

	backupModel := BackupV1{
		Clusters: []backupClusterV1{
			{Name: "cluster1", Tags: []string{"prod"}, Sidecars: []backupSidecarV1{{Host: "host1", Port: 5432}}},
		},
		Queries: []backupQueryV1{
			{Name: "restored-query", Type: BLOAT_V1, Default: "select 1", Custom: "select 1"},
		},
		Permissions: []backupPermissionsV1{
			{Username: "basic:alice", Permissions: map[string]backupPermissionTypeV1{"view.cluster.list": GRANTED_V1}},
		},
	}
	data, errMarshal := json.Marshal(backupModel)
	if errMarshal != nil {
		t.Fatalf("failed to marshal backup model: %v", errMarshal)
	}

	if err := s.importV1(data); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("cluster is restored", func(t *testing.T) {
		clusters, err := s.clusterService.List()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(clusters) != 1 || clusters[0].Name != "cluster1" {
			t.Fatalf("expected restored cluster1, got %v", clusters)
		}
	})

	t.Run("query is restored", func(t *testing.T) {
		queries, err := s.queryService.GetList(nil, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		found := false
		for _, q := range queries {
			if q.Name == "restored-query" {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected restored-query to be restored, got %v", queries)
		}
	})

	t.Run("permissions are restored", func(t *testing.T) {
		perms, err := s.permissionService.GetAllUserPermissions()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		found := false
		for _, p := range perms {
			if p.Username == "basic:alice" {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected restored permissions for basic:alice, got %v", perms)
		}
	})
}

func TestImportV1RejectsMalformedJSON(t *testing.T) {
	s := createTestBackupService(t)
	if err := s.importV1([]byte("not-json")); err == nil {
		t.Fatalf("expected an error for malformed JSON")
	}
}

func TestImportV1AggregatesClusterErrorsButContinues(t *testing.T) {
	s := createTestBackupService(t)

	backupModel := BackupV1{
		Clusters: []backupClusterV1{
			{Name: "", Sidecars: []backupSidecarV1{{Host: "host1", Port: 5432}}}, // invalid: empty name
			{Name: "valid-cluster", Sidecars: []backupSidecarV1{{Host: "host2", Port: 5433}}},
		},
	}
	data, errMarshal := json.Marshal(backupModel)
	if errMarshal != nil {
		t.Fatalf("failed to marshal backup model: %v", errMarshal)
	}

	err := s.importV1(data)
	if err == nil {
		t.Fatalf("expected an aggregated error for the invalid cluster")
	}
	if !strings.Contains(err.Error(), "cluster[0]") {
		t.Fatalf("expected the error to reference cluster[0], got %v", err)
	}

	clusters, errList := s.clusterService.List()
	if errList != nil {
		t.Fatalf("expected no error, got %v", errList)
	}
	if len(clusters) != 1 || clusters[0].Name != "valid-cluster" {
		t.Fatalf("expected the valid cluster to still be imported, got %v", clusters)
	}
}

func TestImportV1SkipsQueriesWithInvalidType(t *testing.T) {
	s := createTestBackupService(t)

	backupModel := BackupV1{
		Queries: []backupQueryV1{
			{Name: "bad-query", Type: backupQueryTypeV1(99), Default: "select 1"},
		},
	}
	data, errMarshal := json.Marshal(backupModel)
	if errMarshal != nil {
		t.Fatalf("failed to marshal backup model: %v", errMarshal)
	}

	if err := s.importV1(data); err != nil {
		t.Fatalf("expected no error since an unmappable query is silently skipped, got %v", err)
	}

	queries, errList := s.queryService.GetList(nil, nil)
	if errList != nil {
		t.Fatalf("expected no error, got %v", errList)
	}
	for _, q := range queries {
		if q.Name == "bad-query" {
			t.Fatalf("expected the invalid query to be skipped, got %v", queries)
		}
	}
}
