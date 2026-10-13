package plugins

import (
	"ivory/core/config"
	"testing"
)

// TestKeeperScheduleFeaturesAgreeWithTheirOperation holds the same line
// SupportedFeatures already holds against the adapter, one level in: an
// operation's schedule may only be declared where the operation itself
// exists. Declaring a schedule without one puts an input on screen for a
// request that could never be made.
func TestKeeperScheduleFeaturesAgreeWithTheirOperation(t *testing.T) {
	for plugin, adapter := range NewContext(nil, nil).KeeperRegistry.All() {
		t.Run(string(plugin), func(t *testing.T) {
			features := adapter.SupportedFeatures()

			if features[config.ManageNodeKeeperSwitchoverSchedule] && !features[config.ManageNodeKeeperSwitchover] {
				t.Error("scheduled switchover declared without a supported switchover")
			}
			if features[config.ManageNodeKeeperRestartSchedule] && !features[config.ManageNodeKeeperRestart] {
				t.Error("scheduled restart declared without a supported restart")
			}
		})
	}
}

// TestKeeperLagIsReadable holds every keeper to stating what its lag counts.
func TestKeeperLagIsReadable(t *testing.T) {
	for plugin, adapter := range NewContext(nil, nil).KeeperRegistry.All() {
		t.Run(string(plugin), func(t *testing.T) {
			lag := adapter.Lag()
			if lag.Unit == "" {
				t.Error("lag unit is not declared")
			}
			if lag.Threshold < 0 {
				t.Errorf("threshold %d is negative", lag.Threshold)
			}
		})
	}
}
