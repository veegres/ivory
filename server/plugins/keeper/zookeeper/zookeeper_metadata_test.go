package zookeeper

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ViewNodeKeeperOverview, config.ViewNodeKeeperConfig}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for zookeeper", feature)
		}
	}

	excluded := []config.Feature{
		config.ManageNodeKeeperConfigUpdate, config.ManageNodeKeeperSwitchover, config.ManageNodeKeeperReinitialize,
		config.ManageNodeKeeperRestart, config.ManageNodeKeeperReload, config.ManageNodeKeeperFailover, config.ManageNodeKeeperActivation,
	}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for zookeeper", feature)
		}
	}
}

func TestDeploymentSpec(t *testing.T) {
	spec := NewAdapter().DeploymentSpec()

	if spec.DefaultImage == "" {
		t.Error("expected a default image")
	}
	if len(spec.Ports) != 3 {
		t.Errorf("expected 3 ports (client, peer, leader election), got %d", len(spec.Ports))
	}
	if len(spec.Volumes) != 2 {
		t.Errorf("expected 2 volumes (data, datalog), got %d", len(spec.Volumes))
	}
	if len(spec.Env) == 0 {
		t.Error("expected at least one env var")
	}
	if _, ok := spec.Defaults[keeper.VarKeeperPort]; ok {
		t.Errorf("expected no separate keeper port (client port is the db port), got %+v", spec.Defaults)
	}
	if spec.Defaults[keeper.VarDbPort] != "2181" {
		t.Errorf("expected db port default 2181, got %+v", spec.Defaults)
	}
	if _, ok := spec.Defaults[keeper.VarDbUser]; ok {
		t.Errorf("expected no credentials, zookeeper's admin protocol has no auth concept, got %+v", spec.Defaults)
	}
	if len(spec.Fields) != 3 {
		t.Fatalf("expected peerPort, leaderElectionPort and clusterHosts fields, got %+v", spec.Fields)
	}
	if spec.Fields[0].Name != keeper.VarPeerPort || spec.Fields[0].Type != keeper.FieldPort || spec.Fields[0].Default != "2888" {
		t.Errorf("expected a peerPort port field with default 2888, got %+v", spec.Fields[0])
	}
	if spec.Fields[1].Name != varLeaderElectionPort || spec.Fields[1].Type != keeper.FieldPort || spec.Fields[1].Default != "3888" {
		t.Errorf("expected a leaderElectionPort port field with default 3888, got %+v", spec.Fields[1])
	}
	if spec.Fields[2].Name != keeper.VarClusterHosts || spec.Fields[2].Type != keeper.FieldText {
		t.Errorf("expected a clusterHosts text field, got %+v", spec.Fields[2])
	}
	expectedTemplate := "server.{{index}}={{host}}:{{peerPort}}:{{leaderElectionPort}};{{dbPort}}"
	if spec.Fields[2].Template != expectedTemplate || spec.Fields[2].Separator != " " {
		t.Errorf("unexpected clusterHosts template %q with separator %q", spec.Fields[2].Template, spec.Fields[2].Separator)
	}
	if spec.PostScript != "" {
		t.Errorf("expected no post-deploy script, got %q", spec.PostScript)
	}
	if spec.EntryScript != "" {
		t.Errorf("expected no entry script (bootstrap is pure Env), got %q", spec.EntryScript)
	}
	for _, e := range spec.Env {
		if e.Name == "ZOO_4LW_COMMANDS_WHITELIST" {
			if !strings.Contains(e.Value, "mntr") || !strings.Contains(e.Value, "conf") {
				t.Errorf("expected the 4lw whitelist to include mntr and conf, got %q", e.Value)
			}
		}
	}
	if unknown := spec.UnknownVariables(); len(unknown) != 0 {
		t.Errorf("spec references unknown variables: %v", unknown)
	}
}
