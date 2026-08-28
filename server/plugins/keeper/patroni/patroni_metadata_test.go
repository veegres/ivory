package patroni

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
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

func TestRequirements(t *testing.T) {
	req := NewAdapter(nil).Requirements()

	if req.DbPort != 5432 {
		t.Errorf("expected db port 5432, got %d", req.DbPort)
	}
	if req.KeeperPort == nil || *req.KeeperPort != 8008 {
		t.Errorf("expected patroni's own rest api port 8008, got %v", req.KeeperPort)
	}
	if !req.Credentials {
		t.Error("expected patroni to consume database credentials")
	}
	if req.DbUser != "postgres" {
		t.Errorf("expected spilo's locked superuser postgres, got %q", req.DbUser)
	}
}

// TestDefaultTemplates covers patroni's external DCS: Ivory never deploys the
// coordinator, the user points at one they already run.
func TestDefaultTemplates(t *testing.T) {
	templates := NewAdapter(nil).DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			for i, command := range template.Commands {
				// NOTE: the DCS is literal text now - Ivory never deploys the
				// coordinator, so a shipped template carries an address to edit
				if !strings.Contains(command.Command, `ETCD3_HOSTS="etcd-1:2379`) {
					t.Errorf("command %d has no external DCS address to edit", i)
				}
				if !strings.Contains(command.Command, string(keeper.VarKeeperPort)) {
					t.Errorf("command %d does not expose patroni's own rest api port", i)
				}
			}
		})
	}
}
