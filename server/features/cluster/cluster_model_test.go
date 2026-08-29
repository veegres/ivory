package cluster

import "testing"

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
