package keeper

import "testing"

func TestPluginTypeValid(t *testing.T) {
	for _, plugin := range PluginTypes {
		if !plugin.Valid() {
			t.Errorf("%s should be valid", plugin)
		}
	}
	for _, plugin := range []PluginType{"", "patroni", "pgpool"} {
		if plugin.Valid() {
			t.Errorf("%q should not be valid", plugin)
		}
	}
}
