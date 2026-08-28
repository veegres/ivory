package docker

import (
	"ivory/core/config"
	"testing"
)

func TestSupportedFeatures(t *testing.T) {
	features := NewAdapter(nil).SupportedFeatures()

	supported := []config.Feature{config.ViewNodeSystem, config.ManageNodeSystem}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for docker", feature)
		}
	}

	if len(features) != len(supported) {
		t.Errorf("expected only the system features to be declared, got %v", features)
	}
}
