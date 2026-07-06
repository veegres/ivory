package postgres

import (
	"ivory/core/config"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []env.Feature{env.ViewNodeKeeperOverview, env.ViewNodeKeeperConfig, env.ManageNodeKeeperReload, env.ManageNodeKeeperFailover}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for native postgres", feature)
		}
	}

	excluded := []env.Feature{
		env.ManageNodeKeeperConfigUpdate, env.ManageNodeKeeperSwitchover, env.ManageNodeKeeperReinitialize,
		env.ManageNodeKeeperRestart, env.ManageNodeKeeperActivation,
	}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for native postgres", feature)
		}
	}
}

func TestDeploymentSpec(t *testing.T) {
	spec := NewAdapter().DeploymentSpec()

	if spec.DefaultImage == "" {
		t.Error("expected a default image")
	}
	if len(spec.Ports) != 1 {
		t.Errorf("expected 1 port (db), got %d", len(spec.Ports))
	}
	if len(spec.Volumes) == 0 {
		t.Error("expected at least one volume")
	}
	if len(spec.Env) == 0 {
		t.Error("expected at least one env var")
	}
	if spec.DefaultValues["dcs"] != "empty" {
		t.Errorf("expected default dcs empty, got %q", spec.DefaultValues["dcs"])
	}
}
