package node

import (
	"ivory/plugins/platform"
	"testing"
)

func TestPlatformVaultConnectionPlatformOrDefault(t *testing.T) {
	tests := []struct {
		name       string
		connection PlatformVaultConnection
		expected   PlatformPlugin
	}{
		{
			// NOTE: clusters stored before platforms were selectable have no
			// platform at all, and must keep resolving to the only adapter
			// that existed then
			name:       "an empty platform falls back to docker",
			connection: PlatformVaultConnection{Host: "db1"},
			expected:   platform.Docker,
		},
		{
			name:       "an explicit platform is kept",
			connection: PlatformVaultConnection{Host: "db1", Platform: platform.Docker},
			expected:   platform.Docker,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.connection.PlatformOrDefault(); got != tt.expected {
				t.Errorf("PlatformOrDefault() = %q, want %q", got, tt.expected)
			}
		})
	}
}
