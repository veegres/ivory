package patroni

import (
	"ivory/core/config"
	"testing"
)

func TestSupportedFeaturesAllSupported(t *testing.T) {
	features := NewAdapter(nil).SupportedFeatures()

	all := []env.Feature{
		env.ViewNodeKeeperOverview, env.ViewNodeKeeperConfig, env.ManageNodeKeeperConfigUpdate,
		env.ManageNodeKeeperSwitchover, env.ManageNodeKeeperReinitialize, env.ManageNodeKeeperRestart,
		env.ManageNodeKeeperReload, env.ManageNodeKeeperFailover, env.ManageNodeKeeperActivation,
	}
	for _, feature := range all {
		if !features[feature] {
			t.Errorf("feature %v must be supported for patroni", feature)
		}
	}
}

func TestDeploymentSpec(t *testing.T) {
	spec := NewAdapter(nil).DeploymentSpec()

	if spec.DefaultImage == "" {
		t.Error("expected a default image")
	}
	if len(spec.Ports) != 2 {
		t.Errorf("expected 2 ports (keeper, db), got %d", len(spec.Ports))
	}
	if len(spec.Volumes) == 0 {
		t.Error("expected at least one volume")
	}
	if len(spec.Env) == 0 {
		t.Error("expected at least one env var")
	}
	if spec.DefaultValues["username"] != "postgres" {
		t.Errorf("expected default username postgres, got %q", spec.DefaultValues["username"])
	}
}
