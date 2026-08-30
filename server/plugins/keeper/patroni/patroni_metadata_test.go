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

// TestDefaultTemplateDefaults covers what replaced keeper.Requirements: the
// deploy screen's credential fields are filled in by the template that creates
// the deployment, because the shipped spilo deployment leaves the rest api
// unauthenticated and names its superuser postgres.
func TestDefaultTemplateDefaults(t *testing.T) {
	for _, template := range NewAdapter(nil).DefaultTemplates() {
		t.Run(template.Name, func(t *testing.T) {
			if template.Defaults.KeeperUser != "" {
				t.Errorf("expected keeper user %q, got %q", "", template.Defaults.KeeperUser)
			}
			if template.Defaults.DbUser != "postgres" {
				t.Errorf("expected database user %q, got %q", "postgres", template.Defaults.DbUser)
			}
		})
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
			singleHost := strings.Contains(template.Name, "Single Host")
			// NOTE: the DCS is literal text - Ivory never deploys the
			// coordinator, so a shipped template carries an address to edit.
			// A container name is not one: it resolves on neither a plain
			// docker run across VMs nor host networking.
			dcs := `ETCD3_HOSTS="10.0.0.1:2379`
			if singleHost {
				dcs = `ETCD3_HOSTS="{{host}}:2379,{{host}}:2381,{{host}}:2383"`
			}
			for i, command := range template.Commands {
				if !strings.Contains(command.Command, dcs) {
					t.Errorf("command %d has no external DCS address to edit", i)
				}
				if !strings.Contains(command.Command, string(keeper.VarKeeperPort)) {
					t.Errorf("command %d does not expose patroni's own rest api port", i)
				}
				// NOTE: configure_spilo.py never reads PATRONI_NAME, so a node
				// carrying only that registers under the container's hostname
				// and all three collide
				if strings.Contains(command.Command, "PATRONI_NAME") {
					t.Errorf("command %d names itself through PATRONI_NAME, which spilo ignores", i)
				}
				if !strings.Contains(command.Command, `"name":"{{name}}"`) {
					t.Errorf("command %d does not give spilo its own node name", i)
				}
				// NOTE: spilo resolves getaddrinfo(gethostname()) on its first
				// line, before it reads any config; without --hostname there
				// has to be an --add-host or the container dies at once
				if singleHost && !strings.Contains(command.Command, "--add-host") {
					t.Errorf("command %d leaves the VM hostname unresolvable, which kills spilo at startup", i)
				}
			}
		})
	}
}
