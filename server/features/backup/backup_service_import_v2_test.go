package backup

import (
	"encoding/json"
	"ivory/core/config"
	"ivory/features/deployment"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/features/user"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"strings"
	"testing"
)

func testBackupDeployment(name string) backupDeploymentV2 {
	return backupDeploymentV2{
		Name:        name,
		Description: "three etcd members",
		Keeper:      string(keeper.NATIVE_ETCD),
		Platform:    string(platform.Docker),
		Defaults:    backupDeploymentDefaultsV2{DbUser: "root", Dcs: "10.0.0.1:2379"},
		Commands: []backupDeploymentCommandV2{
			{
				Command:  "docker run -d --name {{name}} etcd",
				Defaults: backupDeploymentCommandDefaultsV2{Name: "etcd1", Host: "localhost", KeeperPort: 2379, DbPort: 2379},
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
			if tpl.Commands[0].Defaults.Host != "localhost" {
				t.Errorf("got host %q, want the machine a single-host template's commands land on", tpl.Commands[0].Defaults.Host)
			}
			// NOTE: a restore that dropped these would open the deploy screen
			// with its credentials switched off and no DCS to dial
			if tpl.Defaults.DbUser != "root" || tpl.Defaults.Dcs != "10.0.0.1:2379" {
				t.Errorf("got template defaults %+v, want the exported user and coordination store", tpl.Defaults)
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

func TestImportV2AggregatesQueryErrorsButContinues(t *testing.T) {
	s := createTestBackupService(t)

	data, errMarshal := json.Marshal(BackupV2{
		Queries: []backupQueryV2{
			{Name: "broken", Type: backupQueryTypeV2("not-a-type"), Plugin: string(database.REDIS), Query: "KEYS *"},
			{Name: "good", Type: queryTypeOtherV2, Plugin: string(database.REDIS), Query: "KEYS *"},
		},
	})
	if errMarshal != nil {
		t.Fatalf("failed to marshal backup: %v", errMarshal)
	}

	if err := s.importV2(data); err == nil {
		t.Fatal("expected the unmappable query to be reported")
	}

	queries, errList := s.queryService.GetList(nil, nil)
	if errList != nil {
		t.Fatalf("GetList() error = %v", errList)
	}
	found := false
	for _, q := range queries {
		if q.Name == "good" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected the valid query to be restored despite the broken one")
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

// TestImportV2FromFrozenFile is the V1 test's counterpart for the current
// format: the same frozen-file discipline, applied to the version Export
// writes today. It is worth having before V2 is frozen rather than after -
// the file is what will catch a root model drifting out from under importV2
// once nobody is looking at this code any more.
func TestImportV2FromFrozenFile(t *testing.T) {
	s := createTestBackupService(t)
	if err := s.Import(createMultipartFile(t, "ivory.v2.bak", readGolden(t, "ivory.v2.bak"))); err != nil {
		t.Fatalf("Import() error = %v", err)
	}

	// NOTE: a backup carries who somebody is and never their password, so a
	// restored user who signs in with one comes back waiting for a registration
	t.Run("users restore without a password", func(t *testing.T) {
		users, err := s.userService.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if len(users) != 3 {
			t.Fatalf("expected the file's three users, got %+v", users)
		}

		byName := make(map[string]user.Response, len(users))
		for _, u := range users {
			byName[u.Username] = u
		}

		root, ok := byName["root"]
		if !ok {
			t.Fatalf("expected root to come back, got %+v", users)
		}
		if !root.Superuser || len(root.AuthTypes) != 3 {
			t.Errorf("got %+v, want a superuser with every way of signing in", root)
		}
		if root.Registration == nil || root.Registration.Status != user.RegistrationMissing {
			t.Errorf("got %+v, want a registration still to issue", root.Registration)
		}

		bob, okBob := byName["bob"]
		if !okBob {
			t.Fatalf("expected bob to come back, got %+v", users)
		}
		if bob.Superuser || len(bob.AuthTypes) != 1 || bob.AuthTypes[0] != user.AuthLdap {
			t.Errorf("got %+v, want an ldap-only regular user", bob)
		}
		if bob.Registration != nil {
			t.Errorf("got %+v, want nothing to register for somebody who signs in elsewhere", bob.Registration)
		}
	})

	// NOTE: a record was keyed by the authority that vouched for a person once,
	// and a file written then still names them that way
	t.Run("a permission key written under an authority loses its prefix", func(t *testing.T) {
		permissions, err := s.permissionService.GetAllUserPermissions()
		if err != nil {
			t.Fatalf("GetAllUserPermissions() error = %v", err)
		}
		for _, up := range permissions {
			if strings.Contains(up.Username, ":") {
				t.Errorf("got %q, want a bare username", up.Username)
			}
		}
	})

	// NOTE: everything V1 could not say about a cluster - which pair manages
	// it, whether it speaks TLS, and each node's own name and three ports
	t.Run("a cluster restores whole", func(t *testing.T) {
		c, err := s.clusterService.Get("prod-patroni")
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if c.Plugins.Keeper != keeper.PATRONI_POSTGRES || c.Plugins.Database != database.POSTGRES {
			t.Errorf("got %+v, want the pair the file named", c.Plugins)
		}
		if !c.Tls.Keeper || c.Tls.Database {
			t.Errorf("got tls %+v, want only the keeper half", c.Tls)
		}
		if len(c.Nodes) != 2 {
			t.Fatalf("expected both nodes, got %+v", c.Nodes)
		}
		n := c.Nodes[0]
		if n.Name != "patroni1" || n.Host != "10.0.0.1" {
			t.Errorf("got %+v, want the node's own name beside its host", n)
		}
		if n.SshPort == nil || *n.SshPort != 22 || n.KeeperPort == nil || *n.KeeperPort != 8008 || n.DbPort == nil || *n.DbPort != 5432 {
			t.Errorf("got ports %v/%v/%v, want all three", n.SshPort, n.KeeperPort, n.DbPort)
		}
	})

	// NOTE: a second cluster on a different engine is the point - V1 mapped
	// every cluster to patroni over postgres, and this is what says V2 does not
	t.Run("clusters keep their own engines", func(t *testing.T) {
		c, err := s.clusterService.Get("coordination")
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if c.Plugins.Keeper != keeper.NATIVE_ETCD || c.Plugins.Database != database.ETCD {
			t.Errorf("got %+v, want etcd, not the first cluster's pair", c.Plugins)
		}
	})

	// NOTE: stored as words rather than the ordinals V1 used, so they still
	// mean what they said when the list they index grows
	t.Run("queries keep their engine and their spelled-out enums", func(t *testing.T) {
		queries, err := s.queryService.GetList(nil, nil)
		if err != nil {
			t.Fatalf("GetList() error = %v", err)
		}
		byName := make(map[string]query.Response, len(queries))
		for _, q := range queries {
			byName[q.Name] = q
		}
		member, ok := byName["member list"]
		if !ok {
			t.Fatalf("the non-postgres query was not restored: %v", byName)
		}
		if member.Plugin != database.ETCD {
			t.Errorf("got plugin %v, want etcd: V1 could not carry this at all", member.Plugin)
		}
		if member.Type != query.OTHER {
			t.Errorf("got type %v, want other", member.Type)
		}
		if bloat := byName["table bloat"]; bloat.Type != query.BLOAT || len(bloat.Varieties) != 2 {
			t.Errorf("got %+v, want bloat with both varieties", bloat)
		}
	})

	// NOTE: a template restored without these opens its deploy screen with the
	// credentials switched off, nothing to coordinate through, and every node
	// of a single-host template back on one port
	t.Run("a deployment template keeps its defaults", func(t *testing.T) {
		templates, err := s.deploymentService.List(deployment.ListRequest{})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		var restored *deployment.Template
		for i, tpl := range templates {
			if tpl.Creation == deployment.Manual {
				restored = &templates[i]
			}
		}
		if restored == nil {
			t.Fatalf("the manual template was not restored: %v", templates)
		}
		if restored.Defaults.DbUser != "postgres" || restored.Defaults.Dcs != "localhost:2479,localhost:2481" {
			t.Errorf("got template defaults %+v, want the file's user and coordination store", restored.Defaults)
		}
		if len(restored.Commands) != 2 {
			t.Fatalf("expected both commands, got %d", len(restored.Commands))
		}
		first, second := restored.Commands[0].Defaults, restored.Commands[1].Defaults
		if first.Host != "localhost" || second.Host != "localhost" {
			t.Errorf("got hosts %q/%q, want the machine both commands land on", first.Host, second.Host)
		}
		if first.KeeperPort == second.KeeperPort || first.DbPort == second.DbPort {
			t.Errorf("got %+v and %+v, want each node's own ports", first, second)
		}
		if len(restored.Commands[1].PostScripts) != 1 {
			t.Errorf("got post scripts %v, want the one the file named", restored.Commands[1].PostScripts)
		}
	})

	t.Run("a renamed feature key resolves and an unknown one is dropped", func(t *testing.T) {
		permissions, err := s.permissionService.GetAllUserPermissions()
		if err != nil {
			t.Fatalf("GetAllUserPermissions() error = %v", err)
		}
		var alice permission.PermissionMap
		for _, up := range permissions {
			if up.Username == "alice" {
				alice = up.Permissions
			}
		}
		if alice == nil {
			t.Fatalf("expected the file's user, got %v", permissions)
		}
		if alice[config.ViewNodeSystem] != permission.GRANTED {
			t.Errorf("view.node.platform did not resolve to %s: got %v", config.ViewNodeSystem, alice[config.ViewNodeSystem])
		}
		if _, exists := alice[config.Feature("gone.feature.removed")]; exists {
			t.Error("a feature that no longer exists was restored rather than dropped")
		}
	})
}
