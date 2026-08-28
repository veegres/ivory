package platform

import "testing"

func TestPluginCurrent(t *testing.T) {
	tests := []struct {
		name     string
		stored   PluginType
		expected PluginType
	}{
		{name: "a renamed platform resolves to its current name", stored: "linux", expected: Docker},
		{name: "a current platform is left untouched", stored: Docker, expected: Docker},
		{name: "an unknown platform is left untouched", stored: "kubernetes", expected: "kubernetes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stored.Current(); got != tt.expected {
				t.Errorf("Current() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestRenamedPluginsResolveToCurrentNames keeps the rename table honest: a
// target that is itself renamed would leave stored data one hop short of the
// name the registries actually answer to.
func TestRenamedPluginsResolveToCurrentNames(t *testing.T) {
	for stored, current := range renamedPlugins {
		if _, ok := renamedPlugins[current]; ok {
			t.Errorf("platform %q renames to %q, which is itself renamed", stored, current)
		}
	}
}
