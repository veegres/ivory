package backup

import (
	"ivory/core/config"
	"testing"
)

func TestSyncPermissionV1(t *testing.T) {
	tests := []struct {
		name     string
		stored   string
		expected config.Feature
		err      error
	}{
		{
			name:     "a current key resolves to itself",
			stored:   "view.cluster.list",
			expected: config.ViewClusterList,
		},
		{
			// NOTE: a backup written before the rename is the one live source of
			// the old keys, so importing one must not silently drop the grant
			name:     "a renamed key resolves to its current feature",
			stored:   "view.node.platform",
			expected: config.ViewNodeSystem,
		},
		{
			name:     "another renamed key resolves to its current feature",
			stored:   "manage.node.platform",
			expected: config.ManageNodeSystem,
		},
		{
			name:   "an unknown key is rejected",
			stored: "view.node.nonsense",
			err:    ErrInvalidFeature,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := syncPermissionV1(tt.stored)
			if tt.err != nil {
				if err != tt.err {
					t.Fatalf("expected %v, got %v", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got != tt.expected {
				t.Errorf("syncPermissionV1(%q) = %q, want %q", tt.stored, got, tt.expected)
			}
		})
	}
}

// TestToClusterDisambiguatesSingleHostNodes covers a V1 file holding a cluster
// whose nodes all sit on one VM. V1 knew no node names, so every node's name
// defaults to its host; leaving three of them identical made the restore fail
// the cluster's uniqueness check and put back nothing at all.
func TestToClusterDisambiguatesSingleHostNodes(t *testing.T) {
	bc := backupClusterV1{
		Name: "single-host",
		Sidecars: []backupSidecarV1{
			{Host: "10.0.0.1", Port: 8008},
			{Host: "10.0.0.1", Port: 8009},
			{Host: "10.0.0.1", Port: 8010},
			{Host: "10.0.0.2", Port: 8008},
		},
	}

	nodes := bc.toCluster().Nodes

	expected := []string{"10.0.0.1", "10.0.0.1-2", "10.0.0.1-3", "10.0.0.2"}
	if len(nodes) != len(expected) {
		t.Fatalf("expected %d nodes, got %d", len(expected), len(nodes))
	}
	for i, name := range expected {
		if nodes[i].Name != name {
			t.Errorf("node %d: expected name %q, got %q", i, name, nodes[i].Name)
		}
		if nodes[i].Host != bc.Sidecars[i].Host {
			t.Errorf("node %d: host must stay the connection identity, got %q", i, nodes[i].Host)
		}
		if nodes[i].KeeperPort == nil || *nodes[i].KeeperPort != bc.Sidecars[i].Port {
			t.Errorf("node %d: expected keeper port %d", i, bc.Sidecars[i].Port)
		}
	}
}

// TestBackupV1ShapeIsFrozen holds the sacred rule for V1: it shipped long ago,
// so its types can never change again. Only importV1 moves, adjusting to
// whatever the root models have become.
func TestBackupV1ShapeIsFrozen(t *testing.T) {
	assertShapeIsFrozen(t, "ivory.v1.bak", &BackupV1{})
}
