package clickhouse

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ViewNodeKeeperOverview, config.ViewNodeKeeperConfig, config.ManageNodeKeeperReload}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for clickhouse", feature)
		}
	}

	excluded := []config.Feature{
		config.ManageNodeKeeperConfigUpdate, config.ManageNodeKeeperSwitchover, config.ManageNodeKeeperReinitialize,
		config.ManageNodeKeeperRestart, config.ManageNodeKeeperFailover, config.ManageNodeKeeperActivation,
	}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for clickhouse", feature)
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
	if spec.Defaults[keeper.VarDbPort] != "9000" {
		t.Errorf("expected db port default 9000, got %+v", spec.Defaults)
	}
	if len(spec.Fields) != 2 {
		t.Fatalf("expected the dcs and clusterHosts fields, got %+v", spec.Fields)
	}
	if spec.Fields[0].Name != keeper.VarDcs || spec.Fields[0].Type != keeper.FieldText {
		t.Errorf("expected a dcs text field, got %+v", spec.Fields[0])
	}
	if spec.Fields[1].Name != keeper.VarClusterHosts || spec.Fields[1].Type != keeper.FieldText {
		t.Errorf("expected a clusterHosts text field, got %+v", spec.Fields[1])
	}
	if !strings.Contains(spec.Fields[1].Template, string(keeper.VarHost)) || !strings.Contains(spec.Fields[1].Template, string(keeper.VarDbPort)) {
		t.Errorf("expected clusterHosts template to reference host and db port, got %q", spec.Fields[1].Template)
	}
	if spec.PostScript != "" {
		t.Errorf("expected no post-deploy script, got %q", spec.PostScript)
	}
	if spec.EntryScript == "" {
		t.Error("expected an entry script to generate the cluster config on every node")
	}
	if spec.EntryScriptReplicasOnly {
		t.Error("expected EntryScriptReplicasOnly false (default), clickhouse has no primary/replica asymmetry to skip node 0 for")
	}
	if !strings.Contains(spec.EntryScript, string(keeper.VarDcs)) || !strings.Contains(spec.EntryScript, string(keeper.VarClusterHosts)) {
		t.Errorf("expected the entry script to reference both dcs and clusterHosts, got %q", spec.EntryScript)
	}
	if unknown := spec.UnknownVariables(); len(unknown) != 0 {
		t.Errorf("spec references unknown variables: %v", unknown)
	}
}
