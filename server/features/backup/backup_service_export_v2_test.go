package backup

import (
	"ivory/core/config"
	"ivory/features/cluster"
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

// TestExportV2CarriesUsersWithoutPasswords holds the one rule a user export has
// to keep: a backup says who somebody is and how they sign in, and never what
// they sign in with.
func TestExportV2CarriesUsersWithoutPasswords(t *testing.T) {
	s := createTestBackupService(t)
	if _, err := s.userService.CreateOutright("root", "password123", []user.AuthType{user.AuthBasic, user.AuthLdap}, true); err != nil {
		t.Fatalf("failed to seed root: %v", err)
	}

	backupModel, err := s.exportV2()
	if err != nil {
		t.Fatalf("exportV2() error = %v", err)
	}
	if len(backupModel.Users) != 1 {
		t.Fatalf("expected the seeded user, got %+v", backupModel.Users)
	}

	got := backupModel.Users[0]
	if got.Username != "root" || !got.Superuser {
		t.Errorf("got %+v, want the superuser root", got)
	}
	if len(got.AuthTypes) != 2 || got.AuthTypes[0] != string(user.AuthBasic) {
		t.Errorf("got %v, want both ways of signing in", got.AuthTypes)
	}

	marshalled, errMarshal := s.Export()
	if errMarshal != nil {
		t.Fatalf("Export() error = %v", errMarshal)
	}
	if strings.Contains(string(marshalled), "password123") {
		t.Error("the exported file carries a password")
	}
}

func testTemplateRequest(name string) deployment.TemplateRequest {
	return deployment.TemplateRequest{
		Name:        name,
		Description: "three etcd members",
		Keeper:      keeper.NATIVE_ETCD,
		Platform:    platform.Docker,
		Defaults:    deployment.TemplateDefaults{DbUser: "root", Dcs: "10.0.0.1:2379"},
		Commands: []deployment.TemplateCommand{
			{
				Command:  "docker run -d --name {{name}} etcd",
				Defaults: deployment.CommandDefaults{Name: "etcd1", Host: "localhost", KeeperPort: 2379, DbPort: 2379},
			},
			{
				Command:     "docker run -d --name {{name}} etcd",
				PostScripts: []string{"etcdctl auth enable"},
				Defaults:    deployment.CommandDefaults{Name: "etcd2", KeeperPort: 2381, DbPort: 2381},
			},
		},
	}
}

func TestExportV2(t *testing.T) {
	s := createTestBackupService(t)
	if _, err := s.deploymentService.Create(testTemplateRequest("mine")); err != nil {
		t.Fatalf("failed to seed template: %v", err)
	}

	backupModel, err := s.exportV2()
	if err != nil {
		t.Fatalf("exportV2() error = %v", err)
	}
	if len(backupModel.Deployments) != 1 {
		t.Fatalf("expected exactly the one manual template, got %+v", backupModel.Deployments)
	}

	got := backupModel.Deployments[0]
	if got.Name != "mine" || got.Description != "three etcd members" {
		t.Errorf("got %+v, want the stored name and description", got)
	}
	if got.Keeper != string(keeper.NATIVE_ETCD) || got.Platform != string(platform.Docker) {
		t.Errorf("got keeper %q platform %q, want the stored pair", got.Keeper, got.Platform)
	}
	if len(got.Commands) != 2 || len(got.Commands[1].PostScripts) != 1 || got.Commands[1].PostScripts[0] != "etcdctl auth enable" {
		t.Errorf("got commands %+v, want both, post script included", got.Commands)
	}
	// NOTE: without these a single-host template restores with every node on
	// one port, which is the collision the defaults exist to prevent
	if got.Commands[1].Defaults.Name != "etcd2" || got.Commands[1].Defaults.DbPort != 2381 {
		t.Errorf("got defaults %+v, want the second node's own name and port", got.Commands[1].Defaults)
	}
	if got.Commands[0].Defaults.Host != "localhost" {
		t.Errorf("got host %q, want the machine a single-host template's commands land on", got.Commands[0].Defaults.Host)
	}
	// NOTE: without these the restored template opens its deploy screen with
	// the credentials switched off and the DCS address blank
	if got.Defaults.DbUser != "root" || got.Defaults.Dcs != "10.0.0.1:2379" {
		t.Errorf("got template defaults %+v, want the stored user and coordination store", got.Defaults)
	}
}

// NOTE: the shipped templates are computed from the keeper plugins on every
// request, so exporting one would import it back as an editable copy sitting
// next to the original
func TestExportV2SkipsSystemTemplates(t *testing.T) {
	s := createTestBackupService(t)
	list, err := s.deploymentService.List(deployment.ListRequest{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(list) == 0 {
		t.Fatal("expected the etcd defaults to be available")
	}

	backupModel, errExport := s.exportV2()
	if errExport != nil {
		t.Fatalf("exportV2() error = %v", errExport)
	}
	if len(backupModel.Deployments) != 0 {
		t.Fatalf("expected no template exported, got %+v", backupModel.Deployments)
	}
}

// NOTE: everything a cluster is - the pair it is managed by, whether it speaks
// TLS, and each node's own name and ports - V1 could carry none of it
func TestExportV2CarriesWholeCluster(t *testing.T) {
	s := createTestBackupService(t)
	keeperPort, dbPort, sshPort := 2379, 2379, 22
	if _, err := s.clusterService.Update(cluster.Request{
		Name: "etcd-cluster",
		Options: cluster.Options{
			Plugins: cluster.Plugins{Keeper: keeper.NATIVE_ETCD, Database: database.ETCD},
			Tls:     cluster.Tls{Keeper: true},
			Tags:    []string{"prod"},
		},
		Nodes: []cluster.NodeConfig{{
			Name:       "etcd-1",
			Host:       "10.0.0.1",
			SshPort:    &sshPort,
			KeeperPort: &keeperPort,
			DbPort:     &dbPort,
		}},
	}); err != nil {
		t.Fatalf("failed to seed cluster: %v", err)
	}

	backupModel, err := s.exportV2()
	if err != nil {
		t.Fatalf("exportV2() error = %v", err)
	}
	if len(backupModel.Clusters) != 1 {
		t.Fatalf("expected the seeded cluster, got %+v", backupModel.Clusters)
	}

	got := backupModel.Clusters[0]
	if got.Keeper != string(keeper.NATIVE_ETCD) || got.Database != string(database.ETCD) {
		t.Errorf("got keeper %q database %q, want the stored pair", got.Keeper, got.Database)
	}
	if !got.Tls.Keeper || got.Tls.Database {
		t.Errorf("got tls %+v, want the keeper half only", got.Tls)
	}
	if len(got.Nodes) != 1 {
		t.Fatalf("expected the one node, got %+v", got.Nodes)
	}
	node := got.Nodes[0]
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
}

// NOTE: V1 skipped every non-postgres query on export, since importing one
// would have silently turned it into a postgres query
func TestExportV2CarriesQueriesOfEveryPlugin(t *testing.T) {
	s := createTestBackupService(t)
	queryType := query.OTHER
	if _, _, err := s.queryService.Create(query.Manual, query.Request{
		Name:      "keys",
		Type:      &queryType,
		Plugin:    database.REDIS,
		Query:     "KEYS *",
		Varieties: []query.VarietyType{query.MasterOnly},
	}); err != nil {
		t.Fatalf("failed to seed query: %v", err)
	}

	backupModel, errExport := s.exportV2()
	if errExport != nil {
		t.Fatalf("exportV2() error = %v", errExport)
	}
	if len(backupModel.Queries) != 1 {
		t.Fatalf("expected the redis query to be exported, got %+v", backupModel.Queries)
	}

	got := backupModel.Queries[0]
	if got.Plugin != string(database.REDIS) {
		t.Errorf("got plugin %q, want redis", got.Plugin)
	}
	// NOTE: spelled out rather than stored as the domain enum's ordinal, which
	// only means anything next to the list it was written against
	if got.Type != queryTypeOtherV2 {
		t.Errorf("got type %q, want %q", got.Type, queryTypeOtherV2)
	}
	if len(got.Varieties) != 1 || got.Varieties[0] != queryVarietyMasterOnlyV2 {
		t.Errorf("got varieties %+v, want %q", got.Varieties, queryVarietyMasterOnlyV2)
	}
	if got.Query != "KEYS *" {
		t.Errorf("got query %q, want the text as it stands now", got.Query)
	}
}

func TestExportV2CarriesPermissions(t *testing.T) {
	s := createTestBackupService(t)
	perms := permission.PermissionMap{config.ViewClusterOverview: permission.GRANTED}
	if err := s.permissionService.UpdateUserPermissions("user1", perms); err != nil {
		t.Fatalf("failed to seed permissions: %v", err)
	}

	backupModel, err := s.exportV2()
	if err != nil {
		t.Fatalf("exportV2() error = %v", err)
	}
	if len(backupModel.Permissions) != 1 {
		t.Fatalf("expected the seeded user, got %+v", backupModel.Permissions)
	}

	got := backupModel.Permissions[0]
	if got.Username != "user1" {
		t.Errorf("got username %q, want user1", got.Username)
	}
	if got.Permissions[string(config.ViewClusterOverview)] != permissionGrantedV2 {
		t.Errorf("got %+v, want the feature granted", got.Permissions)
	}
}
