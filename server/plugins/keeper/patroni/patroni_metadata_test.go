package patroni

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"testing"
)

func TestSupportedFeaturesAllSupported(t *testing.T) {
	features := NewAdapter(nil).SupportedFeatures()

	all := []config.Feature{
		config.ViewNodeKeeperOverview, config.ViewNodeKeeperConfig, config.ManageNodeKeeperConfigUpdate,
		config.ManageNodeKeeperSwitchover, config.ManageNodeKeeperReinitialize, config.ManageNodeKeeperRestart,
		config.ManageNodeKeeperReload, config.ManageNodeKeeperFailover, config.ManageNodeKeeperActivation,
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
	if user, ok := spec.Defaults[keeper.VarDbUser]; !ok || user != "postgres" {
		t.Errorf("expected credentials with the spilo-required username postgres, got %+v", spec.Defaults)
	}
	if spec.Defaults[keeper.VarKeeperPort] != "8008" {
		t.Errorf("expected keeper port default 8008, got %+v", spec.Defaults)
	}
	if spec.Defaults[keeper.VarDbPort] != "5432" {
		t.Errorf("expected db port default 5432, got %+v", spec.Defaults)
	}
	if len(spec.Fields) != 1 || spec.Fields[0].Name != keeper.VarDcs || spec.Fields[0].Type != keeper.FieldText {
		t.Errorf("expected a dcs text field (external DCS address), got %+v", spec.Fields)
	}
	if len(spec.PostDeploy) != 0 {
		t.Errorf("expected no post-deploy commands, got %+v", spec.PostDeploy)
	}
	if unknown := spec.UnknownVariables(); len(unknown) != 0 {
		t.Errorf("spec references unknown variables: %v", unknown)
	}
}
