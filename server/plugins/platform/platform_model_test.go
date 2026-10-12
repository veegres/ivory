package platform

import "testing"

func TestPluginTypeValid(t *testing.T) {
	if !Docker.Valid() {
		t.Error("docker should be valid")
	}
	for _, plugin := range []PluginType{"", "kubernetes"} {
		if plugin.Valid() {
			t.Errorf("%q should not be valid", plugin)
		}
	}
}
