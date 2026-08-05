package mongo

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{
		config.ViewNodeKeeperOverview, config.ViewNodeKeeperConfig,
		config.ManageNodeKeeperConfigUpdate, config.ManageNodeKeeperFailover,
	}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for mongo", feature)
		}
	}

	excluded := []config.Feature{
		config.ManageNodeKeeperSwitchover, config.ManageNodeKeeperReinitialize,
		config.ManageNodeKeeperRestart, config.ManageNodeKeeperReload, config.ManageNodeKeeperActivation,
	}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for mongo", feature)
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
	if _, ok := spec.Defaults[keeper.VarDbUser]; ok {
		t.Errorf("expected no credentials required by default, got %+v", spec.Defaults)
	}
	if _, ok := spec.Defaults[keeper.VarKeeperPort]; ok {
		t.Errorf("expected no separate keeper port, got %+v", spec.Defaults)
	}
	if spec.Defaults[keeper.VarDbPort] != "27017" {
		t.Errorf("expected db port default 27017, got %+v", spec.Defaults)
	}
	if spec.EntryScriptReplicasOnly {
		t.Error("expected EntryScriptReplicasOnly false, every node needs --replSet at startup")
	}
	if spec.EntryScript == "" || !strings.Contains(spec.EntryScript, "--replSet") {
		t.Errorf("expected an entry script starting mongod with --replSet, got %q", spec.EntryScript)
	}
	if spec.PostScript == "" || !strings.Contains(spec.PostScript, "rs.initiate") {
		t.Errorf("expected a post-deploy script calling rs.initiate, got %q", spec.PostScript)
	}
	if !strings.Contains(spec.PostScript, string(keeper.VarClusterHosts)) {
		t.Errorf("expected post script to reference %s, got %q", keeper.VarClusterHosts, spec.PostScript)
	}
	if len(spec.Fields) != 1 || spec.Fields[0].Name != keeper.VarClusterHosts {
		t.Errorf("expected a single clusterHosts field, got %+v", spec.Fields)
	}
	if unknown := spec.UnknownVariables(); len(unknown) != 0 {
		t.Errorf("spec references unknown variables: %v", unknown)
	}
}
