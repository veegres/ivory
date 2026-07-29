package postgres

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ViewNodeKeeperOverview, config.ViewNodeKeeperConfig, config.ManageNodeKeeperReload, config.ManageNodeKeeperFailover}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for native postgres", feature)
		}
	}

	excluded := []config.Feature{
		config.ManageNodeKeeperConfigUpdate, config.ManageNodeKeeperSwitchover, config.ManageNodeKeeperReinitialize,
		config.ManageNodeKeeperRestart, config.ManageNodeKeeperActivation,
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
	if user, ok := spec.Defaults[keeper.VarDbUser]; !ok || user != "" {
		t.Errorf("expected credentials with a user-chosen username, got %+v", spec.Defaults)
	}
	if _, ok := spec.Defaults[keeper.VarKeeperPort]; ok {
		t.Errorf("expected no separate keeper port, got %+v", spec.Defaults)
	}
	if spec.Defaults[keeper.VarDbPort] != "5432" {
		t.Errorf("expected db port default 5432, got %+v", spec.Defaults)
	}
	if len(spec.Fields) != 0 {
		t.Errorf("expected no fields, got %+v", spec.Fields)
	}
	if len(spec.PostDeploy) != 0 {
		t.Errorf("expected no post-deploy commands, got %+v", spec.PostDeploy)
	}
	if spec.EntryScript == "" {
		t.Error("expected an entry script to rebase replicas via pg_basebackup")
	}
	if !strings.Contains(spec.EntryScript, string(keeper.VarPrimaryHost)) {
		t.Errorf("expected the entry script to reference %s, got %q", keeper.VarPrimaryHost, spec.EntryScript)
	}
	if !strings.Contains(spec.EntryScript, "pg_basebackup") {
		t.Errorf("expected the entry script to run pg_basebackup, got %q", spec.EntryScript)
	}
	if unknown := spec.UnknownVariables(); len(unknown) != 0 {
		t.Errorf("spec references unknown variables: %v", unknown)
	}
}
