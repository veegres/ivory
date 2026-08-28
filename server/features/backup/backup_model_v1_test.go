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
