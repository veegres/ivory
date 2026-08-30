package backup

import (
	"encoding/json"
	"ivory/core/config"
	"ivory/features/deployment"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"testing"
)

func testBackupDeployment(name string) backupDeploymentV2 {
	return backupDeploymentV2{
		Name:        name,
		Description: "three etcd members",
		Keeper:      string(keeper.NATIVE_ETCD),
		Platform:    string(platform.Docker),
		Commands: []backupDeploymentCommandV2{
			{
				Command:  "docker run -d --name {{name}} etcd",
				Defaults: backupDeploymentCommandDefaultsV2{Name: "etcd1", KeeperPort: 2379, DbPort: 2379},
			},
			{
				Command:     "docker run -d --name {{name}} etcd",
				PostScripts: []string{"etcdctl auth enable"},
				Defaults:    backupDeploymentCommandDefaultsV2{Name: "etcd2", KeeperPort: 2381, DbPort: 2381},
			},
		},
	}
}

func TestImportV2(t *testing.T) {
	s := createTestBackupService(t)
	keeperPort, dbPort, sshPort := 2379, 2379, 22

	data, errMarshal := json.Marshal(BackupV2{
		Clusters: []backupClusterV2{{
			Name:     "cluster1",
			Keeper:   string(keeper.NATIVE_ETCD),
			Database: string(database.ETCD),
			Tls:      backupTlsV2{Keeper: true},
			Tags:     []string{"prod"},
			Nodes: []backupNodeV2{{
				Name:       "etcd-1",
				Host:       "10.0.0.1",
				SshPort:    &sshPort,
				KeeperPort: &keeperPort,
				DbPort:     &dbPort,
			}},
		}},
		Deployments: []backupDeploymentV2{testBackupDeployment("mine")},
	})
	if errMarshal != nil {
		t.Fatalf("failed to marshal backup: %v", errMarshal)
	}

	if err := s.importV2(data); err != nil {
		t.Fatalf("importV2() error = %v", err)
	}

	clusters, _ := s.clusterService.List()
	if len(clusters) != 1 {
		t.Fatalf("expected the cluster to be restored too, got %v", clusters)
	}
	restoredCluster := clusters[0]
	if restoredCluster.Plugins.Keeper != keeper.NATIVE_ETCD || restoredCluster.Plugins.Database != database.ETCD {
		t.Errorf("got plugins %+v, want the pair the backup named", restoredCluster.Plugins)
	}
	if !restoredCluster.Tls.Keeper {
		t.Errorf("got tls %+v, want the keeper half restored", restoredCluster.Tls)
	}
	if len(restoredCluster.Nodes) != 1 {
		t.Fatalf("expected the one node, got %+v", restoredCluster.Nodes)
	}
	node := restoredCluster.Nodes[0]
	if node.Name != "etcd-1" || node.Host != "10.0.0.1" {
		t.Errorf("got %+v, want the node's own name beside its host", node)
	}
	if node.KeeperPort == nil || *node.KeeperPort != keeperPort {
		t.Errorf("got keeper port %v, want %d", node.KeeperPort, keeperPort)
	}
	if node.DbPort == nil || *node.DbPort != dbPort {
		t.Errorf("got db port %v, want %d", node.DbPort, dbPort)
	}
	if node.SshPort == nil || *node.SshPort != sshPort {
		t.Errorf("got ssh port %v, want %d", node.SshPort, sshPort)
	}

	templates, err := s.deploymentService.List(deployment.ListRequest{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	restored := 0
	for _, tpl := range templates {
		if tpl.Creation == deployment.Manual {
			restored++
			if tpl.Name != "mine" || len(tpl.Commands) != 2 {
				t.Errorf("got %+v, want the exported template", tpl)
			}
			// NOTE: a restore that dropped these would put every node of a
			// single-host template back on one port
			if tpl.Commands[1].Defaults.Name != "etcd2" || tpl.Commands[1].Defaults.KeeperPort != 2381 {
				t.Errorf("got defaults %+v, want the second node's own name and port", tpl.Commands[1].Defaults)
			}
		}
	}
	if restored != 1 {
		t.Fatalf("expected exactly one restored template, got %d", restored)
	}
}

// NOTE: a template restores as the user's own, never as a shipped one - the
// shipped ones are computed and cannot be written at all
func TestImportV2AggregatesTemplateErrorsButContinues(t *testing.T) {
	s := createTestBackupService(t)

	data, errMarshal := json.Marshal(BackupV2{
		Deployments: []backupDeploymentV2{
			{Name: "", Commands: []backupDeploymentCommandV2{{Command: "docker run"}}},
			testBackupDeployment("good"),
		},
	})
	if errMarshal != nil {
		t.Fatalf("failed to marshal backup: %v", errMarshal)
	}

	if err := s.importV2(data); err == nil {
		t.Fatal("expected the nameless template to be reported")
	}

	templates, _ := s.deploymentService.List(deployment.ListRequest{Keeper: nil, Platform: nil})
	found := false
	for _, tpl := range templates {
		if tpl.Name == "good" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected the valid template to be restored despite the broken one")
	}
}

// NOTE: a query is restored for the engine it was written for - V1 could only
// ever describe postgres, and reading a V2 file that way would point a redis
// query at a postgres console
func TestImportV2RestoresQueryPlugin(t *testing.T) {
	s := createTestBackupService(t)

	data, errMarshal := json.Marshal(BackupV2{
		Queries: []backupQueryV2{{
			Name:      "keys",
			Type:      queryTypeOtherV2,
			Plugin:    string(database.REDIS),
			Query:     "KEYS *",
			Varieties: []backupQueryVarietyV2{queryVarietyMasterOnlyV2},
		}},
	})
	if errMarshal != nil {
		t.Fatalf("failed to marshal backup: %v", errMarshal)
	}

	if err := s.importV2(data); err != nil {
		t.Fatalf("importV2() error = %v", err)
	}

	queries, errList := s.queryService.GetList(nil, nil)
	if errList != nil {
		t.Fatalf("GetList() error = %v", errList)
	}
	var restored *query.Response
	for i, q := range queries {
		if q.Name == "keys" {
			restored = &queries[i]
		}
	}
	if restored == nil {
		t.Fatal("expected the query to be restored")
	}
	if restored.Plugin != database.REDIS {
		t.Errorf("got %q, want the redis query restored as redis", restored.Plugin)
	}
	if restored.Type != query.OTHER {
		t.Errorf("got type %v, want the type the backup spelled out", restored.Type)
	}
	if len(restored.Varieties) != 1 || restored.Varieties[0] != query.MasterOnly {
		t.Errorf("got varieties %+v, want master only", restored.Varieties)
	}
	if restored.Custom != "KEYS *" {
		t.Errorf("got query %q, want the stored text", restored.Custom)
	}
}

func TestImportV2RestoresPermissions(t *testing.T) {
	s := createTestBackupService(t)

	data, errMarshal := json.Marshal(BackupV2{
		Permissions: []backupPermissionsV2{{
			Username: "user1",
			Permissions: map[string]backupPermissionStatusV2{
				string(config.ViewClusterOverview): permissionGrantedV2,
				"nothing.like.a.feature":           permissionGrantedV2,
			},
		}},
	})
	if errMarshal != nil {
		t.Fatalf("failed to marshal backup: %v", errMarshal)
	}

	if err := s.importV2(data); err != nil {
		t.Fatalf("importV2() error = %v", err)
	}

	all, errGet := s.permissionService.GetAllUserPermissions()
	if errGet != nil {
		t.Fatalf("GetAllUserPermissions() error = %v", errGet)
	}
	if len(all) != 1 || all[0].Username != "user1" {
		t.Fatalf("expected the one restored user, got %+v", all)
	}
	perms := all[0].Permissions
	if perms[config.ViewClusterOverview] != permission.GRANTED {
		t.Errorf("got %+v, want the feature granted", perms)
	}
	// NOTE: a key that is no feature at all is dropped rather than restored as
	// a permission nothing can ever check
	if _, ok := perms[config.Feature("nothing.like.a.feature")]; ok {
		t.Errorf("got %+v, want the unknown feature dropped", perms)
	}
}

func TestImportV2RejectsMalformedJSON(t *testing.T) {
	s := createTestBackupService(t)
	if err := s.importV2([]byte("not-json")); err == nil {
		t.Fatal("expected an error for malformed backup data")
	}
}
