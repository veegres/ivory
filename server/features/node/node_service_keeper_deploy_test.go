package node

import (
	"errors"
	"fmt"
	"ivory/clients/console/ssh"
	"ivory/core/utils"
	"ivory/plugins/keeper"
	"ivory/plugins/keeper/clickhouse"
	"ivory/plugins/keeper/etcd"
	"ivory/plugins/keeper/patroni"
	"ivory/plugins/keeper/postgres"
	"ivory/plugins/keeper/redis"
	"ivory/plugins/keeper/zookeeper"
	"ivory/plugins/platform"
	"ivory/plugins/platform/linux"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// newDeployTestService builds a service with the real platform adapter and
// all keeper plugins registered, enough for the pure deploy computations.
func newDeployTestService() *Service {
	platformRegistry := utils.NewRegistry[platform.Plugin, platform.Adapter]()
	platformRegistry.Register(platform.Linux, linux.NewAdapter(ssh.NewClient()))
	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	keeperMetadataRegistry.Register(keeper.PATRONI_POSTGRES, patroni.NewAdapter(nil))
	keeperMetadataRegistry.Register(keeper.NATIVE_POSTGRES, postgres.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_ETCD, etcd.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_REDIS, redis.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_CLICKHOUSE, clickhouse.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_ZOOKEEPER, zookeeper.NewAdapter())
	return &Service{
		platformRegistry:       platformRegistry,
		keeperMetadataRegistry: keeperMetadataRegistry,
	}
}

func intPtr(v int) *int {
	return &v
}

func TestService_KeeperDeploySpec(t *testing.T) {
	s := newDeployTestService()

	tests := []struct {
		name     string
		plugin   KeeperPlugin
		expected KeeperDeploySpecResponse
	}{
		{
			name:   "patroni",
			plugin: keeper.PATRONI_POSTGRES,
			expected: KeeperDeploySpecResponse{
				Uri: "ghcr.io/zalando/spilo-18:4.1-p2",
				Fields: DeployFieldsResponse{
					Defaults: map[string]string{
						string(keeper.VarKeeperPort): "8008",
						string(keeper.VarDbPort):     "5432",
						string(keeper.VarDbUser):     "postgres",
					},
					Fields: []DeployFieldResponse{
						{Name: "{{dcs}}", Label: "DCS (etcd, zookeper, etc)", Example: "etcd1:2379, etcd2:2379, etcd3:2379", Type: "text"},
					},
				},
			},
		},
		{
			name:   "postgres",
			plugin: keeper.NATIVE_POSTGRES,
			expected: KeeperDeploySpecResponse{
				Uri: "postgres:18",
				Fields: DeployFieldsResponse{
					Defaults: map[string]string{string(keeper.VarDbPort): "5432", string(keeper.VarDbUser): ""},
					Fields:   []DeployFieldResponse{},
				},
			},
		},
		{
			name:   "etcd",
			plugin: keeper.NATIVE_ETCD,
			expected: KeeperDeploySpecResponse{
				Uri: "quay.io/coreos/etcd:v3.6.5",
				Fields: DeployFieldsResponse{
					Defaults: map[string]string{string(keeper.VarDbPort): "2379", string(keeper.VarDbUser): "root"},
					Fields: []DeployFieldResponse{
						{Name: "{{peerPort}}", Label: "Peer Port", Type: "port", Default: "2380"},
						{Name: "{{clusterHosts}}", Label: "Initial Cluster", Type: "text", Derived: true},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.KeeperDeploySpec(KeeperDeploySpecRequest{Plugin: tt.plugin})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got.Uri != tt.expected.Uri {
				t.Errorf("Uri = %q, want %q", got.Uri, tt.expected.Uri)
			}
			if !reflect.DeepEqual(got.Fields, tt.expected.Fields) {
				t.Errorf("Fields = %+v, want %+v", got.Fields, tt.expected.Fields)
			}
		})
	}
}

// TestService_KeeperDeployPlanTemplates pins the planned default
// options byte-for-byte to the templates that previously lived hardcoded in
// the web app (DatabaseImageOptions) and were later served via deploy-options,
// so moving them into the plan changes nothing for existing users.
func TestService_KeeperDeployPlanTemplates(t *testing.T) {
	s := newDeployTestService()

	patroniTemplate := "--name {{host}}\n" +
		"--hostname {{host}}\n" +
		"--restart unless-stopped\n" +
		"-p {{keeperPort}}:{{keeperPort}}\n" +
		"-p {{dbPort}}:{{dbPort}}\n" +
		"-v /data/postgres:/home/postgres/pgdata\n" +
		"-e SCOPE=\"{{cluster}}\"\n" +
		"-e PATRONI_NAME=\"{{host}}\"\n" +
		"-e ETCD3_HOSTS=\"{{dcs}}\"\n" +
		"-e PGPORT={{dbPort}}\n" +
		"-e APIPORT={{keeperPort}}\n" +
		"-e PGPASSWORD_SUPERUSER=\"{{dbPass}}\"\n" +
		"-e RESTAPI_CONNECT_ADDRESS=\"{{host}}:{{keeperPort}}\"\n" +
		`-e SPILO_CONFIGURATION='{"postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'`
	patroniTemplateSingleHost := "--name {{host}}\n" +
		"--hostname {{host}}\n" +
		"--network host\n" +
		"-e SCOPE=\"{{cluster}}\"\n" +
		"-e PATRONI_NAME=\"{{host}}\"\n" +
		"-e ETCD3_HOSTS=\"{{dcs}}\"\n" +
		"-e PGPORT={{dbPort}}\n" +
		"-e APIPORT={{keeperPort}}\n" +
		"-e PGPASSWORD_SUPERUSER=\"{{dbPass}}\"\n" +
		"-e RESTAPI_CONNECT_ADDRESS=\"{{host}}:{{keeperPort}}\"\n" +
		`-e SPILO_CONFIGURATION='{"postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'`
	postgresTemplate := "--name {{host}}\n" +
		"--hostname {{host}}\n" +
		"--restart unless-stopped\n" +
		"-p {{dbPort}}:{{dbPort}}\n" +
		"-v /data/postgres:/var/lib/postgresql/data\n" +
		"-e PGPORT=\"{{dbPort}}\"\n" +
		"-e POSTGRES_USER=\"{{dbUser}}\"\n" +
		"-e POSTGRES_PASSWORD=\"{{dbPass}}\""
	postgresTemplateSingleHost := "--name {{host}}\n" +
		"--hostname {{host}}\n" +
		"--network host\n" +
		"-e PGPORT=\"{{dbPort}}\"\n" +
		"-e POSTGRES_USER=\"{{dbUser}}\"\n" +
		"-e POSTGRES_PASSWORD=\"{{dbPass}}\""
	etcdTemplate := "--name {{host}}\n" +
		"--hostname {{host}}\n" +
		"--restart unless-stopped\n" +
		"-p {{dbPort}}:{{dbPort}}\n" +
		"-p {{peerPort}}:{{peerPort}}\n" +
		"-v /data/etcd:/data/etcd\n" +
		"-e ETCD_NAME=\"{{host}}\"\n" +
		"-e ETCD_DATA_DIR=\"/data/etcd\"\n" +
		"-e ETCD_INITIAL_CLUSTER=\"{{clusterHosts}}\"\n" +
		"-e ETCD_INITIAL_CLUSTER_STATE=\"new\"\n" +
		"-e ETCD_INITIAL_CLUSTER_TOKEN=\"{{cluster}}\"\n" +
		"-e ETCD_LISTEN_CLIENT_URLS=\"http://0.0.0.0:{{dbPort}}\"\n" +
		"-e ETCD_ADVERTISE_CLIENT_URLS=\"http://{{host}}:{{dbPort}}\"\n" +
		"-e ETCD_LISTEN_PEER_URLS=\"http://0.0.0.0:{{peerPort}}\"\n" +
		"-e ETCD_INITIAL_ADVERTISE_PEER_URLS=\"http://{{host}}:{{peerPort}}\""
	etcdTemplateSingleHost := "--name {{host}}\n" +
		"--hostname {{host}}\n" +
		"--network host\n" +
		"-e ETCD_NAME=\"{{host}}\"\n" +
		"-e ETCD_DATA_DIR=\"/data/etcd\"\n" +
		"-e ETCD_INITIAL_CLUSTER=\"{{clusterHosts}}\"\n" +
		"-e ETCD_INITIAL_CLUSTER_STATE=\"new\"\n" +
		"-e ETCD_INITIAL_CLUSTER_TOKEN=\"{{cluster}}\"\n" +
		"-e ETCD_LISTEN_CLIENT_URLS=\"http://0.0.0.0:{{dbPort}}\"\n" +
		"-e ETCD_ADVERTISE_CLIENT_URLS=\"http://{{host}}:{{dbPort}}\"\n" +
		"-e ETCD_LISTEN_PEER_URLS=\"http://0.0.0.0:{{peerPort}}\"\n" +
		"-e ETCD_INITIAL_ADVERTISE_PEER_URLS=\"http://{{host}}:{{peerPort}}\""

	tests := []struct {
		name       string
		plugin     KeeperPlugin
		singleHost bool
		expected   string
	}{
		{name: "patroni", plugin: keeper.PATRONI_POSTGRES, expected: patroniTemplate},
		{name: "patroni single host", plugin: keeper.PATRONI_POSTGRES, singleHost: true, expected: patroniTemplateSingleHost},
		{name: "postgres", plugin: keeper.NATIVE_POSTGRES, expected: postgresTemplate},
		{name: "postgres single host", plugin: keeper.NATIVE_POSTGRES, singleHost: true, expected: postgresTemplateSingleHost},
		{name: "etcd", plugin: keeper.NATIVE_ETCD, expected: etcdTemplate},
		{name: "etcd single host", plugin: keeper.NATIVE_ETCD, singleHost: true, expected: etcdTemplateSingleHost},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
				Plugin:     tt.plugin,
				Cluster:    "main",
				SingleHost: tt.singleHost,
				Nodes:      []KeeperDeployPlanNodeRequest{{Host: "node1"}},
			})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if plan.Nodes[0].Options != tt.expected {
				t.Errorf("Options mismatch\ngot:\n%s\nwant:\n%s", plan.Nodes[0].Options, tt.expected)
			}
		})
	}
}

func TestService_KeeperDeployPlan(t *testing.T) {
	s := newDeployTestService()

	t.Run("patroni fills defaults and keeps credential variables visible", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.PATRONI_POSTGRES,
			Cluster: "main",
			Values:  map[string]string{"{{dcs}}": "etcd1:2379"},
			Nodes:   []KeeperDeployPlanNodeRequest{{Host: "db1"}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if plan.Image != "ghcr.io/zalando/spilo-18:4.1-p2" {
			t.Errorf("Image = %q, want spilo default", plan.Image)
		}
		node := plan.Nodes[0]
		if node.SshPort != 22 || node.KeeperPort != 8008 || node.DbPort != 5432 {
			t.Errorf("ports = %d/%d/%d, want 22/8008/5432", node.SshPort, node.KeeperPort, node.DbPort)
		}
		if !strings.Contains(node.OptionsPreview, `-e ETCD3_HOSTS="etcd1:2379"`) {
			t.Errorf("preview misses interpolated dcs:\n%s", node.OptionsPreview)
		}
		if !strings.Contains(node.OptionsPreview, `-e PGPASSWORD_SUPERUSER="{{dbPass}}"`) {
			t.Errorf("preview must keep the credential variable visible, not fake a value:\n%s", node.OptionsPreview)
		}
		if !strings.Contains(node.OptionsPreview, "-e APIPORT=8008") || !strings.Contains(node.OptionsPreview, "-e PGPORT=5432") {
			t.Errorf("preview misses interpolated ports:\n%s", node.OptionsPreview)
		}
		if len(plan.Warnings) != 0 {
			t.Errorf("expected no warnings, got %v", plan.Warnings)
		}
	})

	t.Run("patroni warns when the manual dcs is missing", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.PATRONI_POSTGRES,
			Cluster: "main",
			Nodes:   []KeeperDeployPlanNodeRequest{{Host: "db1"}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !slices.Contains(plan.Warnings, "missing value for placeholder {{dcs}}") {
			t.Errorf("expected a dcs warning, got %v", plan.Warnings)
		}
	})

	t.Run("patroni single host offsets the default db and keeper ports per node", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:     keeper.PATRONI_POSTGRES,
			Cluster:    "main",
			SingleHost: true,
			Values:     map[string]string{"{{dcs}}": "etcd1:2379"},
			Nodes: []KeeperDeployPlanNodeRequest{
				{Host: "n1"},
				{Host: "n2"},
				{Host: "n3"},
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for i, node := range plan.Nodes {
			if node.DbPort != 5432+i {
				t.Errorf("node %d dbPort = %d, want %d", i, node.DbPort, 5432+i)
			}
			if node.KeeperPort != 8008+i {
				t.Errorf("node %d keeperPort = %d, want %d", i, node.KeeperPort, 8008+i)
			}
		}
	})

	t.Run("patroni single host honors an explicit per-node db port without offsetting it", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:     keeper.PATRONI_POSTGRES,
			Cluster:    "main",
			SingleHost: true,
			Values:     map[string]string{"{{dcs}}": "etcd1:2379"},
			Nodes: []KeeperDeployPlanNodeRequest{
				{Host: "n1", DbPort: intPtr(6000), KeeperPort: intPtr(6100)},
				{Host: "n2"},
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if plan.Nodes[0].DbPort != 6000 || plan.Nodes[0].KeeperPort != 6100 {
			t.Errorf("node 0 ports = %d/%d, want 6000/6100", plan.Nodes[0].DbPort, plan.Nodes[0].KeeperPort)
		}
		if plan.Nodes[1].DbPort != 5433 || plan.Nodes[1].KeeperPort != 8009 {
			t.Errorf("node 1 ports = %d/%d, want 5433/8009", plan.Nodes[1].DbPort, plan.Nodes[1].KeeperPort)
		}
	})

	t.Run("etcd single host derives unique peer ports and the dcs", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:     keeper.NATIVE_ETCD,
			Cluster:    "dcs",
			SingleHost: true,
			Nodes: []KeeperDeployPlanNodeRequest{
				{Host: "n1", DbPort: intPtr(2381)},
				{Host: "n2", DbPort: intPtr(2382)},
				{Host: "n3", DbPort: intPtr(2383)},
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expectedMembers := "n1=http://n1:2380,n2=http://n2:2381,n3=http://n3:2382"
		if plan.Values["{{clusterHosts}}"] != expectedMembers {
			t.Errorf("Values[{{clusterHosts}}] = %q, want %q", plan.Values["{{clusterHosts}}"], expectedMembers)
		}
		if plan.PostScript == "" {
			t.Error("expected the etcd auth post-deploy script")
		}
		for i, node := range plan.Nodes {
			if node.Ports["{{peerPort}}"] != 2380+i {
				t.Errorf("node %d peerPort = %d, want %d", i, node.Ports["{{peerPort}}"], 2380+i)
			}
			if node.KeeperPort != node.DbPort {
				t.Errorf("node %d keeperPort = %d, want dbPort %d", i, node.KeeperPort, node.DbPort)
			}
			if !strings.Contains(node.OptionsPreview, fmt.Sprintf("-e ETCD_INITIAL_CLUSTER=%q", expectedMembers)) {
				t.Errorf("node %d preview misses the derived member list:\n%s", i, node.OptionsPreview)
			}
		}
		if len(plan.Warnings) != 0 {
			t.Errorf("expected no warnings, got %v", plan.Warnings)
		}
	})

	t.Run("etcd multi host keeps the default peer port", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.NATIVE_ETCD,
			Cluster: "dcs",
			Nodes:   []KeeperDeployPlanNodeRequest{{Host: "n1"}, {Host: "n2"}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if plan.Values["{{clusterHosts}}"] != "n1=http://n1:2380,n2=http://n2:2380" {
			t.Errorf("Values[{{clusterHosts}}] = %q", plan.Values["{{clusterHosts}}"])
		}
		if plan.Nodes[0].DbPort != 2379 || plan.Nodes[1].Ports["{{peerPort}}"] != 2380 {
			t.Errorf("unexpected ports %+v", plan.Nodes)
		}
	})

	t.Run("etcd warns when the given keeper port differs from the db port", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.NATIVE_ETCD,
			Cluster: "dcs",
			Nodes:   []KeeperDeployPlanNodeRequest{{Host: "n1", KeeperPort: intPtr(1234)}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if plan.Nodes[0].KeeperPort != 2379 {
			t.Errorf("keeperPort = %d, want the db port 2379", plan.Nodes[0].KeeperPort)
		}
		if len(plan.Warnings) != 1 || !strings.Contains(plan.Warnings[0], "keeper port 1234 is ignored") {
			t.Errorf("expected an ignored keeper port warning, got %v", plan.Warnings)
		}
	})

	t.Run("postgres has no dcs and keeps both credential variables visible", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.NATIVE_POSTGRES,
			Cluster: "pg",
			Nodes:   []KeeperDeployPlanNodeRequest{{Host: "pg1"}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if plan.Nodes[0].EntryScript != "" {
			t.Errorf("expected the primary (only node) to get an empty EntryScript, got:\n%s", plan.Nodes[0].EntryScript)
		}
		if !strings.Contains(plan.Nodes[0].OptionsPreview, `-e POSTGRES_USER="{{dbUser}}"`) ||
			!strings.Contains(plan.Nodes[0].OptionsPreview, `-e POSTGRES_PASSWORD="{{dbPass}}"`) {
			t.Errorf("preview must keep the credential variables visible, not fake values:\n%s", plan.Nodes[0].OptionsPreview)
		}
		if len(plan.Warnings) != 0 {
			t.Errorf("expected no warnings, got %v", plan.Warnings)
		}
	})

	t.Run("postgres rebases every node but the first from the primary's real host", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.NATIVE_POSTGRES,
			Cluster: "pg",
			Nodes:   []KeeperDeployPlanNodeRequest{{Host: "pg1"}, {Host: "pg2"}, {Host: "pg3"}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if plan.Nodes[0].EntryScript != "" {
			t.Errorf("expected the primary (pg1) to get an empty EntryScript, got:\n%s", plan.Nodes[0].EntryScript)
		}
		for i, host := range []string{"pg2", "pg3"} {
			node := plan.Nodes[i+1]
			if node.EntryScript == "" {
				t.Fatalf("node %q: expected a non-empty EntryScript", host)
			}
			// NOTE: the primary's host is resolved via primaryHostMarker
			// directly into EntryScript at plan time - {{host}} itself
			// stays unresolved until EntryScriptPreview.
			if !strings.Contains(node.EntryScript, `host=pg1`) {
				t.Errorf("node %q: expected EntryScript to already reference the primary's real host pg1, got:\n%s", host, node.EntryScript)
			}
			if !strings.Contains(node.EntryScriptPreview, "application_name="+host) {
				t.Errorf("node %q: expected EntryScriptPreview to resolve its own host into application_name, got:\n%s", host, node.EntryScriptPreview)
			}
			if !strings.Contains(node.EntryScriptPreview, "pg_basebackup") || !strings.Contains(node.EntryScriptPreview, "docker-entrypoint.sh postgres") {
				t.Errorf("node %q: expected EntryScriptPreview to still contain the rebase and the real entrypoint, got:\n%s", host, node.EntryScriptPreview)
			}
		}
		if len(plan.Warnings) != 0 {
			t.Errorf("expected no warnings, got %v", plan.Warnings)
		}
	})

	t.Run("clickhouse gives every node including the primary an entry script", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.NATIVE_CLICKHOUSE,
			Cluster: "ch",
			Values:  map[string]string{"{{dcs}}": "keeper1:9181,keeper2:9181"},
			Nodes:   []KeeperDeployPlanNodeRequest{{Host: "ch1"}, {Host: "ch2"}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expectedClusterHosts := "<replica><host>ch1</host><port>9000</port></replica><replica><host>ch2</host><port>9000</port></replica>"
		if plan.Values["{{clusterHosts}}"] != expectedClusterHosts {
			t.Errorf("Values[{{clusterHosts}}] = %q, want %q", plan.Values["{{clusterHosts}}"], expectedClusterHosts)
		}
		for i, host := range []string{"ch1", "ch2"} {
			node := plan.Nodes[i]
			if node.EntryScript == "" {
				t.Fatalf("node %q: expected a non-empty EntryScript even for the primary, clickhouse leaves EntryScriptReplicasOnly false", host)
			}
			if !strings.Contains(node.EntryScriptPreview, expectedClusterHosts) {
				t.Errorf("node %q: expected EntryScriptPreview to embed the resolved cluster replica list, got:\n%s", host, node.EntryScriptPreview)
			}
			if !strings.Contains(node.EntryScriptPreview, `dcs="keeper1:9181,keeper2:9181"`) {
				t.Errorf("node %q: expected EntryScriptPreview to resolve the dcs field, got:\n%s", host, node.EntryScriptPreview)
			}
		}
		if len(plan.Warnings) != 0 {
			t.Errorf("expected no warnings, got %v", plan.Warnings)
		}
	})

	t.Run("zookeeper assigns each node a 1-based index for ZOO_MY_ID and the ensemble list", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.NATIVE_ZOOKEEPER,
			Cluster: "zk",
			Nodes:   []KeeperDeployPlanNodeRequest{{Host: "zk1"}, {Host: "zk2"}, {Host: "zk3"}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expectedServers := "server.1=zk1:2888:3888;2181 server.2=zk2:2888:3888;2181 server.3=zk3:2888:3888;2181"
		if plan.Values["{{clusterHosts}}"] != expectedServers {
			t.Errorf("Values[{{clusterHosts}}] = %q, want %q", plan.Values["{{clusterHosts}}"], expectedServers)
		}
		for i, host := range []string{"zk1", "zk2", "zk3"} {
			node := plan.Nodes[i]
			if node.Index != i+1 {
				t.Errorf("node %q: Index = %d, want %d", host, node.Index, i+1)
			}
			if !strings.Contains(node.OptionsPreview, fmt.Sprintf(`-e ZOO_MY_ID="%d"`, i+1)) {
				t.Errorf("node %q: expected OptionsPreview to resolve its own ZOO_MY_ID, got:\n%s", host, node.OptionsPreview)
			}
			if !strings.Contains(node.OptionsPreview, expectedServers) {
				t.Errorf("node %q: expected OptionsPreview to embed the resolved ensemble server list, got:\n%s", host, node.OptionsPreview)
			}
			if node.EntryScript != "" {
				t.Errorf("node %q: expected no EntryScript, zookeeper bootstraps purely via Env, got:\n%s", host, node.EntryScript)
			}
		}
		if len(plan.Warnings) != 0 {
			t.Errorf("expected no warnings, got %v", plan.Warnings)
		}
	})

	t.Run("etcd user-edited peer port becomes the base", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:     keeper.NATIVE_ETCD,
			Cluster:    "dcs",
			SingleHost: true,
			Values:     map[string]string{"{{peerPort}}": "3000"},
			Nodes:      []KeeperDeployPlanNodeRequest{{Host: "n1"}, {Host: "n2"}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if plan.Values["{{peerPort}}"] != "3000" {
			t.Errorf("Values[{{peerPort}}] = %q, want the edited base 3000", plan.Values["{{peerPort}}"])
		}
		if plan.Nodes[0].Ports["{{peerPort}}"] != 3000 || plan.Nodes[1].Ports["{{peerPort}}"] != 3001 {
			t.Errorf("unexpected per-node peer ports %+v", plan.Nodes)
		}
		if plan.Values["{{clusterHosts}}"] != "n1=http://n1:3000,n2=http://n2:3001" {
			t.Errorf("Values[{{clusterHosts}}] = %q", plan.Values["{{clusterHosts}}"])
		}
	})

	t.Run("etcd user-edited member list wins over the derived one", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.NATIVE_ETCD,
			Cluster: "dcs",
			Values:  map[string]string{"{{clusterHosts}}": "custom=http://elsewhere:2380"},
			Nodes:   []KeeperDeployPlanNodeRequest{{Host: "n1"}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if plan.Values["{{clusterHosts}}"] != "custom=http://elsewhere:2380" {
			t.Errorf("Values[{{clusterHosts}}] = %q, want the edited value", plan.Values["{{clusterHosts}}"])
		}
		if !strings.Contains(plan.Nodes[0].OptionsPreview, `-e ETCD_INITIAL_CLUSTER="custom=http://elsewhere:2380"`) {
			t.Errorf("preview misses the edited member list:\n%s", plan.Nodes[0].OptionsPreview)
		}
	})

	t.Run("invalid port field value falls back to the default with a warning", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.NATIVE_ETCD,
			Cluster: "dcs",
			Values:  map[string]string{"{{peerPort}}": "oops"},
			Nodes:   []KeeperDeployPlanNodeRequest{{Host: "n1"}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if plan.Nodes[0].Ports["{{peerPort}}"] != 2380 {
			t.Errorf("peerPort = %d, want the default 2380", plan.Nodes[0].Ports["{{peerPort}}"])
		}
		if len(plan.Warnings) != 1 || !strings.Contains(plan.Warnings[0], "invalid port") {
			t.Errorf("expected an invalid port warning, got %v", plan.Warnings)
		}
	})

	t.Run("node options override replaces the rendered template", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.NATIVE_POSTGRES,
			Cluster: "pg",
			Nodes: []KeeperDeployPlanNodeRequest{
				{Host: "pg1", Options: "--name {{host}} -p {{dbPort}}:{{dbPort}}"},
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if plan.Nodes[0].Options != "--name {{host}} -p {{dbPort}}:{{dbPort}}" {
			t.Errorf("Options = %q, want the override", plan.Nodes[0].Options)
		}
		if plan.Nodes[0].OptionsPreview != "--name pg1 -p 5432:5432" {
			t.Errorf("Preview = %q", plan.Nodes[0].OptionsPreview)
		}
	})

	t.Run("image override wins over the spec default", func(t *testing.T) {
		plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
			Plugin:  keeper.NATIVE_POSTGRES,
			Cluster: "pg",
			Image:   "custom:1",
			Nodes:   []KeeperDeployPlanNodeRequest{{Host: "pg1"}},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if plan.Image != "custom:1" {
			t.Errorf("Image = %q, want custom:1", plan.Image)
		}
	})

	t.Run("unknown plugin fails", func(t *testing.T) {
		_, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{Plugin: "unknown"})
		if err == nil {
			t.Fatal("expected error for unknown plugin, got nil")
		}
	})
}

func TestService_KeeperDeployPlanWarningsBlockDeploy(t *testing.T) {
	s := newDeployTestService()

	// NOTE: no {{dcs}} value is provided, so the plan carries patroni's usual
	// "missing value for placeholder {{dcs}}" warning; KeeperDeploy must
	// refuse before ever touching Connection/Vaults (left zero-valued here).
	_, err := s.KeeperDeploy(KeeperDeployRequest{
		Plugin:  keeper.PATRONI_POSTGRES,
		Cluster: "main",
		Node:    KeeperDeployPlanNodeRequest{Host: "db1"},
	})
	if !errors.Is(err, ErrKeeperDeployPlanHasWarnings) {
		t.Fatalf("expected ErrKeeperDeployPlanHasWarnings, got %v", err)
	}
	if !strings.Contains(err.Error(), "{{dcs}}") {
		t.Fatalf("expected the error to name the missing placeholder, got %v", err)
	}
}

func TestService_KeeperDeploySpecUnknownPlugin(t *testing.T) {
	s := &Service{
		platformRegistry:       utils.NewRegistry[platform.Plugin, platform.Adapter](),
		keeperMetadataRegistry: utils.NewRegistry[keeper.Plugin, keeper.Metadata](),
	}

	_, err := s.KeeperDeploySpec(KeeperDeploySpecRequest{Plugin: "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown plugin, got nil")
	}
}

func TestService_getInterpolatedStringDeployKeys(t *testing.T) {
	values := map[string]string{
		"{{cluster}}":    "main",
		"{{dcs}}":        "etcd1:2379",
		"{{host}}":       "db1",
		"{{dbUser}}":     "postgres",
		"{{dbPass}}":     "secret",
		"{{keeperPort}}": "8008",
		"{{dbPort}}":     "5432",
	}

	t.Run("deploy keys", func(t *testing.T) {
		got := keeper.Interpolate(
			"{{cluster}} {{dcs}} {{host}} {{keeperPort}} {{dbPort}} {{dbUser}} {{dbPass}}",
			values,
		)
		want := "main etcd1:2379 db1 8008 5432 postgres secret"
		if got != want {
			t.Errorf("Interpolate() = %q, want %q", got, want)
		}
	})

	t.Run("aux ports", func(t *testing.T) {
		got := keeper.Interpolate(
			"{{host}}:{{peerPort}} {{dbPort}}",
			map[string]string{"{{host}}": "db1", "{{peerPort}}": "2380", "{{dbPort}}": "5432"},
		)
		want := "db1:2380 5432"
		if got != want {
			t.Errorf("Interpolate() = %q, want %q", got, want)
		}
	})

	t.Run("missing and empty values keep placeholders unresolved", func(t *testing.T) {
		got := keeper.Interpolate(
			"{{cluster}} {{dcs}} {{peerPort}}",
			map[string]string{"{{cluster}}": "main", "{{dcs}}": ""},
		)
		want := "main {{dcs}} {{peerPort}}"
		if got != want {
			t.Errorf("Interpolate() = %q, want %q", got, want)
		}
	})
}

func TestBuildNodeValues(t *testing.T) {
	s := &Service{}
	pn := KeeperDeployPlanNode{
		Host:       "db1",
		Index:      1,
		KeeperPort: 8008,
		DbPort:     5432,
		Ports:      map[string]int{"{{peerPort}}": 2380},
	}

	t.Run("strips credentials and fills built-ins", func(t *testing.T) {
		got := s.buildNodeValues("main", map[string]string{string(keeper.VarDbUser): "postgres", string(keeper.VarDbPass): "secret", "{{dcs}}": "etcd1:2379"}, map[string]string{"{{dcs}}": "etcd1:2379"}, pn)
		want := map[string]string{
			string(keeper.VarCluster):    "main",
			string(keeper.VarIndex):      "1",
			string(keeper.VarKeeperPort): "8008",
			string(keeper.VarDbPort):     "5432",
			"{{dcs}}":                    "etcd1:2379",
			"{{peerPort}}":               "2380",
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("buildNodeValues() = %+v, want %+v", got, want)
		}
	})

	t.Run("plan values win over raw request values", func(t *testing.T) {
		got := s.buildNodeValues("main", map[string]string{"{{dcs}}": "stale"}, map[string]string{"{{dcs}}": "fresh"}, pn)
		if got["{{dcs}}"] != "fresh" {
			t.Errorf("{{dcs}} = %q, want the plan value to win", got["{{dcs}}"])
		}
	})
}
