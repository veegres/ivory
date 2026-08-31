package mongo

import (
	"ivory/core/config"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewPlugin().SupportedFeatures()

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

// TestDefaultTemplateDefaults covers what replaced keeper.Requirements: the
// deploy screen's credential fields are filled in by the template that creates
// the deployment, because enabling client auth would take a keyfile the
// deployment does not ship, so there is no account to name.
func TestDefaultTemplateDefaults(t *testing.T) {
	for _, template := range NewPlugin().DefaultTemplates() {
		t.Run(template.Name, func(t *testing.T) {
			if template.Defaults.KeeperUser != "" {
				t.Errorf("expected keeper user %q, got %q", "", template.Defaults.KeeperUser)
			}
			if template.Defaults.DbUser != "" {
				t.Errorf("expected database user %q, got %q", "", template.Defaults.DbUser)
			}
		})
	}
}

// TestDefaultTemplates covers the replica set bootstrap: every member starts
// as a plain mongod, and the last command initiates the set once they are all
// up - which is why rs.initiate() sits there rather than on the first.
func TestDefaultTemplates(t *testing.T) {
	templates := NewPlugin().DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			last := template.Commands[len(template.Commands)-1]
			if len(last.PostScripts) != 1 || !strings.Contains(last.PostScripts[0], "rs.initiate") {
				t.Fatal("expected the last command to initiate the replica set")
			}
			// NOTE: on one VM the members differ by port and are addressed by
			// {{host}}; across VMs they are example addresses. Neither is a
			// container name - no docker network spans either layout, so a
			// name resolves to nothing in both
			member := `10.0.0.1:27017`
			if strings.Contains(template.Name, "Single Host") {
				member = `{{host}}:27017`
			}
			if !strings.Contains(last.PostScripts[0], member) {
				t.Error("expected the members list rs.initiate is given")
			}
			for i, command := range template.Commands[:len(template.Commands)-1] {
				if len(command.PostScripts) > 0 {
					t.Errorf("command %d must not initiate the set before every member is up", i)
				}
			}
			for i, command := range template.Commands {
				if !strings.Contains(command.Command, "--replSet") {
					t.Errorf("command %d does not start with a replica set name", i)
				}
			}
		})
	}
}
