package cluster

import (
	"ivory/features/node"
	"testing"
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
