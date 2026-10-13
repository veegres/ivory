package cluster

import (
	"ivory/core/service/cert"
	"ivory/core/utils"
	"ivory/features/node"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// TestWithNodeNamesDisambiguatesSharedHosts covers a cluster stored before
// node names existed whose nodes all sit on one VM. Defaulting every name to
// its host made the three identical, and a cluster whose names collide is
// rejected by validateNodeNames - so it could not be updated or edited in the
// list at all.
func TestWithNodeNamesDisambiguatesSharedHosts(t *testing.T) {
	tests := []struct {
		name     string
		nodes    []NodeConfig
		expected []string
	}{
		{
			name:     "several nodes on one host",
			nodes:    []NodeConfig{{Host: "10.0.0.1"}, {Host: "10.0.0.1"}, {Host: "10.0.0.1"}},
			expected: []string{"10.0.0.1", "10.0.0.1-2", "10.0.0.1-3"},
		},
		{
			name:     "distinct hosts keep their own name",
			nodes:    []NodeConfig{{Host: "10.0.0.1"}, {Host: "10.0.0.2"}},
			expected: []string{"10.0.0.1", "10.0.0.2"},
		},
		{
			name:     "an existing name is never overwritten and never reused",
			nodes:    []NodeConfig{{Name: "10.0.0.1", Host: "10.0.0.1"}, {Host: "10.0.0.1"}},
			expected: []string{"10.0.0.1", "10.0.0.1-2"},
		},
		{
			name:     "a suffix already taken is skipped",
			nodes:    []NodeConfig{{Name: "10.0.0.1-2", Host: "10.0.0.1"}, {Host: "10.0.0.1"}, {Host: "10.0.0.1"}},
			expected: []string{"10.0.0.1-2", "10.0.0.1", "10.0.0.1-3"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Response{Nodes: test.nodes}.withNodeNames()

			names := make([]string, 0, len(got.Nodes))
			for _, n := range got.Nodes {
				names = append(names, n.Name)
			}
			for i, name := range test.expected {
				if names[i] != name {
					t.Errorf("node %d: expected name %q, got %q", i, name, names[i])
				}
			}

			s := &Service{}
			if err := s.validateNodeNames(got.Nodes); err != nil {
				t.Errorf("the defaulted names must survive the cluster's own validation: %v", err)
			}
		})
	}
}

// TestMapKeeperResponseListSkipsResponsesWithoutAHost covers what Detect and Fix
// must not write into a cluster: a response that carries no endpoint describes
// an attribute of a node (native postgres' sync state, read from the primary) or
// a member the keeper cannot yet address, not a node. Storing one would put a
// node with no host into the cluster - and mapping one derefs a host that is not
// there.
func TestMapKeeperResponseListSkipsResponsesWithoutAHost(t *testing.T) {
	host := "db1"
	port := 5432
	name := "postgres2"

	got := mapKeeperResponseList([]node.KeeperOneResponse{
		{DiscoveredHost: &host, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port},
		{Sync: true, DiscoveredName: &name},
	})

	if len(got) != 1 {
		t.Fatalf("expected only the response describing a node, got %d: %v", len(got), got)
	}
	if got[0].Host != host {
		t.Errorf("expected the node at %q, got %q", host, got[0].Host)
	}
}

// TestMapKeeperResponseDisambiguatesSharedHosts covers a Detect or a Fix of a
// cluster whose members live on one VM. A keeper that identifies its members by
// endpoint reports no name of its own, so every node falls back to the same
// host - and a cluster whose names collide is rejected by validateNodeNames on
// every later write, leaving it impossible to edit.
func TestMapKeeperResponseDisambiguatesSharedHosts(t *testing.T) {
	host := "localhost"
	first, second, third := 27017, 27018, 27019
	named := "mongo2"

	tests := []struct {
		name      string
		responses []node.KeeperOneResponse
		expected  []string
	}{
		{
			name: "nameless members on one host",
			responses: []node.KeeperOneResponse{
				{DiscoveredHost: &host, DiscoveredKeeperPort: &first},
				{DiscoveredHost: &host, DiscoveredKeeperPort: &second},
				{DiscoveredHost: &host, DiscoveredKeeperPort: &third},
			},
			expected: []string{"localhost", "localhost-2", "localhost-3"},
		},
		{
			name: "a reported name is kept and never reused",
			responses: []node.KeeperOneResponse{
				{DiscoveredHost: &host, DiscoveredKeeperPort: &first, DiscoveredName: &named},
				{DiscoveredHost: &host, DiscoveredKeeperPort: &second},
			},
			expected: []string{"mongo2", "localhost"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := mapKeeperResponseList(test.responses)

			if len(got) != len(test.expected) {
				t.Fatalf("expected %d nodes, got %d: %v", len(test.expected), len(got), got)
			}
			for i, name := range test.expected {
				if got[i].Name != name {
					t.Errorf("node %d: expected name %q, got %q", i, name, got[i].Name)
				}
			}

			s := &Service{}
			if err := s.validateNodeNames(got); err != nil {
				t.Errorf("a discovered cluster must survive the cluster's own validation: %v", err)
			}
		})
	}
}

// TestMapKeeperResponseMapIsOrdered covers the map-keyed half of discovery:
// which node keeps the bare name and which one is suffixed must not depend on
// map iteration order, or two Detects of the same cluster would name its nodes
// differently.
func TestMapKeeperResponseMapIsOrdered(t *testing.T) {
	host := "localhost"
	first, second, third := 27017, 27018, 27019
	responses := map[string]node.KeeperOneResponse{
		"localhost:27017": {DiscoveredHost: &host, DiscoveredKeeperPort: &first},
		"localhost:27018": {DiscoveredHost: &host, DiscoveredKeeperPort: &second},
		"localhost:27019": {DiscoveredHost: &host, DiscoveredKeeperPort: &third},
	}
	expected := map[int]string{27017: "localhost", 27018: "localhost-2", 27019: "localhost-3"}

	for range 10 {
		got := mapKeeperResponseMap(responses)

		if len(got) != len(expected) {
			t.Fatalf("expected %d nodes, got %d: %v", len(expected), len(got), got)
		}
		for _, n := range got {
			if name := expected[*n.KeeperPort]; n.Name != name {
				t.Errorf("port %d: expected name %q, got %q", *n.KeeperPort, name, n.Name)
			}
		}
	}
}

// TestRequestBinding covers the shape rules the request types state in their
// binding tags, which every route and tool runs before a service sees them.
func TestRequestBinding(t *testing.T) {
	zero, outOfRange := 0, 70000

	validRequest := func() Request {
		return Request{
			Name:    "prod",
			Nodes:   []NodeConfig{{Name: "pg1", Host: "10.0.0.1"}},
			Options: Options{Plugins: Plugins{Keeper: "patroni_postgres", Database: "postgres"}},
		}
	}

	tests := []struct {
		name   string
		modify func(r *Request)
		want   string
	}{
		{name: "a complete request passes", modify: func(r *Request) {}},
		{name: "ports may be left out", modify: func(r *Request) { r.Nodes[0].SshPort = nil }},
		{name: "cluster name is required", modify: func(r *Request) { r.Name = " " }, want: "'Request.Name' Error:Field validation for 'Name' failed on the 'notblank' tag"},
		{name: "at least one node", modify: func(r *Request) { r.Nodes = []NodeConfig{} }, want: "'Request.Nodes' Error:Field validation for 'Nodes' failed on the 'min' tag"},
		{name: "a node names itself", modify: func(r *Request) { r.Nodes[0].Name = "" }, want: "'Request.Nodes[0].Name' Error:Field validation for 'Name' failed on the 'required' tag"},
		{name: "a node names its host", modify: func(r *Request) { r.Nodes[0].Host = "" }, want: "'Request.Nodes[0].Host' Error:Field validation for 'Host' failed on the 'required' tag"},
		{name: "a stated port must be a port", modify: func(r *Request) { r.Nodes[0].DbPort = &outOfRange }, want: "'Request.Nodes[0].DbPort' Error:Field validation for 'DbPort' failed on the 'max' tag"},
		{name: "a stated zero port is refused", modify: func(r *Request) { r.Nodes[0].KeeperPort = &zero }, want: "'Request.Nodes[0].KeeperPort' Error:Field validation for 'KeeperPort' failed on the 'min' tag"},
		{name: "unknown keeper plugin", modify: func(r *Request) { r.Plugins.Keeper = "pgpool" }, want: "'Request.Options.Plugins.Keeper' Error:Field validation for 'Keeper' failed on the 'valid' tag"},
		{name: "database plugin is required", modify: func(r *Request) { r.Plugins.Database = "" }, want: "'Request.Options.Plugins.Database' Error:Field validation for 'Database' failed on the 'required' tag"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := validRequest()
			test.modify(&request)
			assertBinding(t, &request, test.want)
		})
	}
}

// TestDeployRequestBinding holds what Deploy used to check by hand before it
// started: a cluster name, at least one node, and a platform it knows.
func TestDeployRequestBinding(t *testing.T) {
	tests := []struct {
		name   string
		modify func(r *DeployRequest)
		want   string
	}{
		{name: "a complete request passes", modify: func(r *DeployRequest) {}},
		{name: "cluster name not provided", modify: func(r *DeployRequest) { r.CommonConfig.Cluster = "" }, want: "'DeployRequest.CommonConfig.Cluster' Error:Field validation for 'Cluster' failed on the 'required' tag"},
		{name: "no nodes provided", modify: func(r *DeployRequest) { r.Nodes = nil }, want: "'DeployRequest.Nodes' Error:Field validation for 'Nodes' failed on the 'required' tag"},
		{name: "platform not provided", modify: func(r *DeployRequest) { r.Platform = "" }, want: "'DeployRequest.Platform' Error:Field validation for 'Platform' failed on the 'required' tag"},
		{name: "unknown platform", modify: func(r *DeployRequest) { r.Platform = "k8s" }, want: "'DeployRequest.Platform' Error:Field validation for 'Platform' failed on the 'valid' tag"},
		{name: "blank command", modify: func(r *DeployRequest) { r.Nodes[0].Command = "  " }, want: "'DeployRequest.Nodes[0].Command' Error:Field validation for 'Command' failed on the 'notblank' tag"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := validDeployRequest()
			test.modify(&request)
			assertBinding(t, &request, test.want)
		})
	}
}

func assertBinding(t *testing.T, request any, want string) {
	t.Helper()
	utils.RegisterValidators()
	err := binding.Validator.ValidateStruct(request)
	if want == "" {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("expected an error containing %q, got %v", want, err)
	}
}

// TestOptionsNodeConnections holds the same tls rule discovery applies -
// certificates only on the side whose flag is on - and that a half the node
// cannot be reached on is left out rather than filled with a guess.
func TestOptionsNodeConnections(t *testing.T) {
	port := func(value int) *int { return &value }
	keeperVault, dbVault, sshVault, ca := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	options := Options{
		Plugins: Plugins{Keeper: keeper.PATRONI_POSTGRES, Database: database.POSTGRES},
		Tls:     Tls{Keeper: true},
		Certs:   cert.Certs{ClientCAId: &ca},
		Vaults:  Vaults{KeeperId: &keeperVault, DatabaseId: &dbVault, SshKeyId: &sshVault},
	}
	full := NodeConfig{Name: "pg1", Host: "10.0.0.1", SshPort: port(22), KeeperPort: port(8008), DbPort: port(5432)}

	t.Run("every half of a fully configured node", func(t *testing.T) {
		c := options.nodeConnections("prod", full)
		if c.Keeper == nil || c.Keeper.Host != "10.0.0.1" || c.Keeper.Port != 8008 || c.Keeper.Cluster != "prod" ||
			c.Keeper.Plugin != keeper.PATRONI_POSTGRES || c.Keeper.VaultId != &keeperVault {
			t.Fatalf("keeper = %+v", c.Keeper)
		}
		if c.Keeper.Certs == nil || c.Keeper.Certs.ClientCAId != &ca {
			t.Errorf("keeper tls is on, so its certs should travel: %+v", c.Keeper.Certs)
		}
		if c.Platform == nil || c.Platform.Port != 22 || c.Platform.VaultId != sshVault || c.Platform.Platform != node.DefaultPlatform {
			t.Errorf("platform = %+v", c.Platform)
		}
		if c.Database == nil || c.Database.Db.Port != 5432 || c.Database.Db.Plugin != database.POSTGRES || c.Database.VaultId != &dbVault {
			t.Fatalf("database = %+v", c.Database)
		}
		if c.Database.Certs != nil {
			t.Errorf("database tls is off, so no certs should travel: %+v", c.Database.Certs)
		}
	})

	t.Run("a missing port or ssh vault leaves that half out", func(t *testing.T) {
		bare := options
		bare.Vaults.SshKeyId = nil
		c := bare.nodeConnections("prod", NodeConfig{Name: "pg1", Host: "10.0.0.1", SshPort: port(22)})
		if c.Keeper != nil || c.Platform != nil || c.Database != nil {
			t.Errorf("expected nothing resolvable, got %+v", c)
		}
	})
}

func TestRouterOverviewRequestBinding(t *testing.T) {
	tests := []struct {
		name    string
		request RouterOverviewRequest
		want    string
	}{
		{name: "a name alone asks every node", request: RouterOverviewRequest{Name: "prod"}},
		{name: "a host and port ask one endpoint", request: RouterOverviewRequest{Name: "prod", Host: "10.0.0.1", Port: 8008}},
		{name: "a host without its port", request: RouterOverviewRequest{Name: "prod", Host: "10.0.0.1"}, want: "'RouterOverviewRequest.Port' Error:Field validation for 'Port' failed on the 'required_with' tag"},
		{name: "a port without its host", request: RouterOverviewRequest{Name: "prod", Port: 8008}, want: "'RouterOverviewRequest.Host' Error:Field validation for 'Host' failed on the 'required_with' tag"},
		{name: "a port out of range", request: RouterOverviewRequest{Name: "prod", Host: "10.0.0.1", Port: 70000}, want: "'RouterOverviewRequest.Port' Error:Field validation for 'Port' failed on the 'max' tag"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertBinding(t, &test.request, test.want)
		})
	}
}
