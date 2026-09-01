package patroni

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
	"testing"
)

func TestSupportedFeaturesAllSupported(t *testing.T) {
	features := NewPlugin(nil).SupportedFeatures()

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
	for _, template := range NewPlugin(nil).DefaultTemplates() {
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
	templates := NewPlugin(nil).DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			singleHost := strings.Contains(template.Name, "Single Host")
			// NOTE: Ivory never deploys the coordinator, so its address comes
			// in through {{dcs}}, which the deploy screen asks for once and
			// every node's command reads. Only the single-host template can
			// name one ahead of time - etcd's own single-host ports, reached
			// through localhost, which under --network host is the VM the whole
			// cluster runs on. A multi-host store lives on machines a template
			// cannot know, exactly like a multi-host node's host, so it states
			// nothing rather than example text the operator has to notice.
			dcs := ""
			if singleHost {
				dcs = "localhost:2479,localhost:2481,localhost:2483"
			}
			if template.Defaults.Dcs != dcs {
				t.Errorf("expected the DCS default %q, got %q", dcs, template.Defaults.Dcs)
			}
			for i, command := range template.Commands {
				if !strings.Contains(command.Command, `ETCD3_HOSTS="`+string(keeper.VarDcs)+`"`) {
					t.Errorf("command %d does not point its DCS at {{dcs}}", i)
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
				// line, before it reads any config; without --hostname the
				// command has to make the VM's own name resolve, or the
				// container dies at once. It maps the name it is actually
				// given rather than one written into the template, which
				// nobody but the operator of that VM could know
				if singleHost && !strings.Contains(command.Command, `echo "127.0.0.1 $(hostname)" >> /etc/hosts`) {
					t.Errorf("command %d leaves the VM hostname unresolvable, which kills spilo at startup", i)
				}
				if singleHost && !strings.Contains(command.Command, "exec /bin/sh /launch.sh init") {
					t.Errorf("command %d does not hand back over to spilo's own startup", i)
				}
			}
		})
	}
}
