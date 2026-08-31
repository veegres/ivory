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
