package clickhouse

import (
	"ivory/core/config"
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

func TestRequirements(t *testing.T) {
	req := NewAdapter().Requirements()

	if req.DbPort != 9000 {
		t.Errorf("expected the native tcp port 9000, got %d", req.DbPort)
	}
	if req.KeeperPort != nil {
		t.Errorf("expected no separate keeper port, got %v", *req.KeeperPort)
	}
	if !req.Credentials {
		t.Error("expected clickhouse to consume database credentials")
	}
	if req.DbUser != "" {
		t.Errorf("expected a free choice of username, got the locked %q", req.DbUser)
	}
}

// TestDefaultTemplates covers clickhouse's lack of asymmetry: every node
// generates the same cluster config file, so there is no special first node.
func TestDefaultTemplates(t *testing.T) {
	templates := NewAdapter().DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			for i, command := range template.Commands {
				if command.Command != template.Commands[0].Command {
					t.Errorf("command %d differs, but clickhouse has no leader/replica asymmetry", i)
				}
				if !strings.Contains(command.Command, "ivory-cluster.xml") {
					t.Errorf("command %d does not generate the cluster config", i)
				}
				if !strings.Contains(command.Command, "<replica><host>clickhouse-1</host>") {
					t.Errorf("command %d is missing the shard's replica list", i)
				}
				if !strings.Contains(command.Command, "<node><host>keeper-1</host>") {
					t.Errorf("command %d is missing the coordinator list", i)
				}
			}
		})
	}
}
