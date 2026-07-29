package etcd

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ViewNodeKeeperOverview, config.ManageNodeKeeperSwitchover}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for etcd", feature)
		}
	}

	excluded := []config.Feature{
		config.ViewNodeKeeperConfig, config.ManageNodeKeeperConfigUpdate, config.ManageNodeKeeperReinitialize,
		config.ManageNodeKeeperRestart, config.ManageNodeKeeperReload, config.ManageNodeKeeperFailover, config.ManageNodeKeeperActivation,
	}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for etcd", feature)
		}
	}
}

func TestDeploymentSpec(t *testing.T) {
	spec := NewAdapter().DeploymentSpec()

	if spec.DefaultImage == "" {
		t.Error("expected a default image")
	}
	if len(spec.Ports) != 2 {
		t.Errorf("expected 2 ports (client, peer), got %d", len(spec.Ports))
	}
	if len(spec.Volumes) == 0 {
		t.Error("expected at least one volume")
	}
	if len(spec.Env) == 0 {
		t.Error("expected at least one env var")
	}
	if _, ok := spec.Defaults[keeper.VarKeeperPort]; ok {
		t.Errorf("expected no separate keeper port (client port is the db port), got %+v", spec.Defaults)
	}
	if spec.Defaults[keeper.VarDbPort] != "2379" {
		t.Errorf("expected db port default 2379, got %+v", spec.Defaults)
	}
	if len(spec.Fields) != 2 {
		t.Fatalf("expected the peerPort and initialCluster fields, got %+v", spec.Fields)
	}
	if spec.Fields[0].Name != keeper.VarPeerPort || spec.Fields[0].Type != keeper.FieldPort || spec.Fields[0].Default != "2380" {
		t.Errorf("expected a peerPort port field with default 2380, got %+v", spec.Fields[0])
	}
	if spec.Fields[1].Name != keeper.VarInitialCluster || spec.Fields[1].Type != keeper.FieldText {
		t.Errorf("expected an initialCluster text field, got %+v", spec.Fields[1])
	}
	if spec.Fields[1].Template != "{{host}}=http://{{host}}:{{peerPort}}" || spec.Fields[1].Separator != "," {
		t.Errorf("unexpected initialCluster template %q with separator %q", spec.Fields[1].Template, spec.Fields[1].Separator)
	}
	if user, ok := spec.Defaults[keeper.VarDbUser]; !ok || user != "root" {
		t.Errorf("expected credentials with the etcd-required username root (auth enable needs it), got %+v", spec.Defaults)
	}
	if len(spec.PostDeploy) != 3 {
		t.Fatalf("expected 3 post-deploy auth commands, got %+v", spec.PostDeploy)
	}
	if !strings.Contains(spec.PostDeploy[0], "user add") || !strings.Contains(spec.PostDeploy[0], "{{dbUser}}:{{dbPass}}") {
		t.Errorf("expected the first post-deploy command to create the credentials user, got %q", spec.PostDeploy[0])
	}
	if !strings.Contains(spec.PostDeploy[1], "user grant-role {{dbUser}} root") {
		t.Errorf("expected the second post-deploy command to grant the root role, got %q", spec.PostDeploy[1])
	}
	if !strings.Contains(spec.PostDeploy[2], "auth enable") {
		t.Errorf("expected the last post-deploy command to enable authentication, got %q", spec.PostDeploy[2])
	}
	for _, port := range spec.Ports {
		if port == "2380" {
			t.Error("peer port must be the {{peerPort}} placeholder, not a literal 2380")
		}
	}
	for _, e := range spec.Env {
		if strings.Contains(e.Value, "2380") {
			t.Errorf("env %s must use the {{peerPort}} placeholder, not a literal 2380", e.Name)
		}
	}
	if unknown := spec.UnknownVariables(); len(unknown) != 0 {
		t.Errorf("spec references unknown variables: %v", unknown)
	}
}
