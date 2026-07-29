package backup

import (
	"ivory/features/cluster"
	"ivory/features/query"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"testing"
)

func TestExportV1(t *testing.T) {
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

	bloatType := query.BLOAT
	if _, _, err := s.queryService.Create(query.Manual, query.Request{
		Name:  "custom-postgres-query",
		Type:  &bloatType,
		Query: "select 1",
	}); err != nil {
		t.Fatalf("failed to seed postgres query: %v", err)
	}

	etcdPlugin := query.DbPlugin(database.ETCD)
	if _, _, err := s.queryService.Create(query.Manual, query.Request{
		Name:   "custom-etcd-query",
		Type:   &bloatType,
		Query:  "get /",
		Plugin: etcdPlugin,
	}); err != nil {
		t.Fatalf("failed to seed etcd query: %v", err)
	}

	if _, err := s.permissionService.CreateUserPermissions("basic", "alice"); err != nil {
		t.Fatalf("failed to seed permissions: %v", err)
	}

	backupModel, err := s.exportV1()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("clusters are mapped with their sidecars", func(t *testing.T) {
		if len(backupModel.Clusters) != 1 {
			t.Fatalf("expected 1 cluster, got %v", backupModel.Clusters)
		}
		bc := backupModel.Clusters[0]
		if bc.Name != "cluster1" {
			t.Fatalf("expected name 'cluster1', got %q", bc.Name)
		}
		if len(bc.Sidecars) != 1 || bc.Sidecars[0].Host != "host1" || bc.Sidecars[0].Port != 5432 {
			t.Fatalf("expected sidecar host1:5432, got %v", bc.Sidecars)
		}
	})

	t.Run("non-postgres queries are skipped since v1 has no plugin field", func(t *testing.T) {
		if len(backupModel.Queries) != 1 {
			t.Fatalf("expected only the postgres query to be exported, got %v", backupModel.Queries)
		}
		if backupModel.Queries[0].Name != "custom-postgres-query" {
			t.Fatalf("expected 'custom-postgres-query', got %q", backupModel.Queries[0].Name)
		}
		if backupModel.Queries[0].Type != BLOAT_V1 {
			t.Fatalf("expected BLOAT_V1, got %v", backupModel.Queries[0].Type)
		}
	})

	t.Run("permissions are mapped for every stored user", func(t *testing.T) {
		found := false
		for _, p := range backupModel.Permissions {
			if p.Username == "basic:alice" {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected permissions for basic:alice, got %v", backupModel.Permissions)
		}
	})
}

func TestExportV1SkipsSystemQueries(t *testing.T) {
	s := createTestBackupService(t)

	bloatType := query.BLOAT
	if _, _, err := s.queryService.Create(query.System, query.Request{
		Name:  "system-query",
		Type:  &bloatType,
		Query: "select 1",
	}); err != nil {
		t.Fatalf("failed to seed system query: %v", err)
	}

	backupModel, err := s.exportV1()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	for _, q := range backupModel.Queries {
		if q.Name == "system-query" {
			t.Fatalf("expected system queries to be excluded from the export, got %v", backupModel.Queries)
		}
	}
}
