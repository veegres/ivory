package redis

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ViewNodeKeeperOverview, config.ViewNodeKeeperConfig, config.ManageNodeKeeperFailover}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for redis", feature)
		}
	}

	excluded := []config.Feature{
		config.ManageNodeKeeperConfigUpdate, config.ManageNodeKeeperSwitchover, config.ManageNodeKeeperReinitialize,
		config.ManageNodeKeeperRestart, config.ManageNodeKeeperReload, config.ManageNodeKeeperActivation,
	}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for redis", feature)
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
	if user, ok := spec.Defaults[keeper.VarDbUser]; !ok || user != "default" {
		t.Errorf("expected credentials with the redis default username, got %+v", spec.Defaults)
	}
	if _, ok := spec.Defaults[keeper.VarKeeperPort]; ok {
		t.Errorf("expected no separate keeper port, got %+v", spec.Defaults)
	}
	if spec.Defaults[keeper.VarDbPort] != "6379" {
		t.Errorf("expected db port default 6379, got %+v", spec.Defaults)
	}
	if len(spec.Fields) != 0 {
		t.Errorf("expected no fields, got %+v", spec.Fields)
	}
	if spec.PostScript != "" {
		t.Errorf("expected no post-deploy script, got %q", spec.PostScript)
	}
	if !spec.EntryScriptReplicasOnly {
		t.Error("expected EntryScriptReplicasOnly true, the primary must not replicate from itself")
	}
	if spec.EntryScript == "" {
		t.Error("expected an entry script to start replicas via REDIS_REPLICATION_MODE")
	}
	if !strings.Contains(spec.EntryScript, string(keeper.VarLeaderHost)) {
		t.Errorf("expected the entry script to reference %s, got %q", keeper.VarLeaderHost, spec.EntryScript)
	}
	if !strings.Contains(spec.EntryScript, "REDIS_REPLICATION_MODE") {
		t.Errorf("expected the entry script to set REDIS_REPLICATION_MODE, got %q", spec.EntryScript)
	}
	if unknown := spec.UnknownVariables(); len(unknown) != 0 {
		t.Errorf("spec references unknown variables: %v", unknown)
	}
}
