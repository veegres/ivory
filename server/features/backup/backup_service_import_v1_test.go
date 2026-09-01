package backup

import (
	"encoding/json"
	"ivory/core/config"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
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

// TestImportV1FromFrozenFile is the compatibility test proper: a file V1 wrote,
// checked in and never regenerated, imported through the same entry point the
// router calls - so the filename dispatch is exercised too - and then read back
// out of the current models field by field.
//
// It is what a struct literal built inside a test cannot be. That literal is
// written against today's types and changes with them, so it keeps passing
// while the mapping underneath rots; this file cannot move, so the day a root
// model drifts out from under importV1, this is what says so.
func TestImportV1FromFrozenFile(t *testing.T) {
	s := createTestBackupService(t)
	if err := s.Import(createMultipartFile(t, "ivory.v1.bak", readGolden(t, "ivory.v1.bak"))); err != nil {
		t.Fatalf("Import() error = %v", err)
	}

	// NOTE: V1 predates plugins entirely, so every cluster in one restores as
	// patroni over postgres - the only pair Ivory had when V1 was written
	t.Run("clusters restore as the pair V1 implied", func(t *testing.T) {
		clusters, err := s.clusterService.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if len(clusters) != 2 {
			t.Fatalf("expected both clusters, got %d", len(clusters))
		}
		c, errGet := s.clusterService.Get("prod-patroni")
		if errGet != nil {
			t.Fatalf("Get() error = %v", errGet)
		}
		if c.Plugins.Keeper != keeper.PATRONI_POSTGRES || c.Plugins.Database != database.POSTGRES {
			t.Errorf("got %+v, want patroni over postgres", c.Plugins)
		}
		if len(c.Tags) != 2 || c.Tags[0] != "prod" {
			t.Errorf("got tags %v, want the file's own", c.Tags)
		}
	})

	// NOTE: a V1 sidecar is a host and a keeper port and nothing else, so a
	// node's name is its host - and two nodes on one VM would collide, which
	// is what the -2 suffix exists to prevent
	t.Run("a repeated host still yields unique node names", func(t *testing.T) {
		c, err := s.clusterService.Get("prod-patroni")
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if len(c.Nodes) != 3 {
			t.Fatalf("expected all three sidecars, got %+v", c.Nodes)
		}
		names := []string{c.Nodes[0].Name, c.Nodes[1].Name, c.Nodes[2].Name}
		want := []string{"10.0.0.1", "10.0.0.2", "10.0.0.1-2"}
		for i := range want {
			if names[i] != want[i] {
				t.Errorf("node %d is named %q, want %q", i, names[i], want[i])
			}
		}
		if c.Nodes[2].KeeperPort == nil || *c.Nodes[2].KeeperPort != 8009 {
			t.Errorf("got keeper port %v, want the sidecar's own 8009", c.Nodes[2].KeeperPort)
		}
		// NOTE: the format never carried these, so their absence is V1, not a
		// dropped field - the distinction this whole file exists to keep
		if c.Nodes[0].DbPort != nil || c.Nodes[0].SshPort != nil {
			t.Errorf("got %+v, want no db/ssh port: V1 had neither", c.Nodes[0])
		}
	})

	// NOTE: V1 stored these as ordinals, which only mean anything beside the
	// list they were written against - so every one of them is checked
	t.Run("ordinal enums map across their whole range", func(t *testing.T) {
		queries, err := s.queryService.GetList(nil, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		byName := make(map[string]query.Response, len(queries))
		for _, q := range queries {
			byName[q.Name] = q
		}
		want := map[string]query.Type{
			"table bloat":     query.BLOAT,
			"running queries": query.ACTIVITY,
			"replication lag": query.REPLICATION,
			"table sizes":     query.STATISTIC,
			"scratch":         query.OTHER,
		}
		for name, queryType := range want {
			got, ok := byName[name]
			if !ok {
				t.Errorf("query %q was not restored", name)
				continue
			}
			if got.Type != queryType {
				t.Errorf("query %q restored as type %v, want %v", name, got.Type, queryType)
			}
		}
		bloat := byName["table bloat"]
		if len(bloat.Varieties) != 2 || bloat.Varieties[0] != query.DatabaseSensitive || bloat.Varieties[1] != query.ReplicaRecommended {
			t.Errorf("got varieties %v, want both the file named", bloat.Varieties)
		}
	})

	// NOTE: a permission key is stored as its string, so a renamed feature has
	// to be resolved on the way in or every grant written under the old name
	// silently resets
	t.Run("a renamed feature key resolves and an unknown one is dropped", func(t *testing.T) {
		permissions, err := s.permissionService.GetAllUserPermissions()
		if err != nil {
			t.Fatalf("GetAllUserPermissions() error = %v", err)
		}
		var alice permission.PermissionMap
		for _, up := range permissions {
			if up.Username == "basic:alice" {
				alice = up.Permissions
			}
		}
		if alice == nil {
			t.Fatalf("expected the file's user, got %v", permissions)
		}
		if alice[config.ViewNodeSystem] != permission.GRANTED {
			t.Errorf("view.node.platform did not resolve to %s: got %v", config.ViewNodeSystem, alice[config.ViewNodeSystem])
		}
		if alice[config.ViewClusterList] != permission.PENDING {
			t.Errorf("got %v for view.cluster.list, want pending", alice[config.ViewClusterList])
		}
		if _, exists := alice[config.Feature("gone.feature.removed")]; exists {
			t.Error("a feature that no longer exists was restored rather than dropped")
		}
	})
}
