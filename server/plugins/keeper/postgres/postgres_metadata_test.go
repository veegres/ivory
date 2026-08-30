package postgres

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ViewNodeKeeperOverview, config.ViewNodeKeeperConfig, config.ManageNodeKeeperReload, config.ManageNodeKeeperFailover}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for native postgres", feature)
		}
	}

	excluded := []config.Feature{
		config.ManageNodeKeeperConfigUpdate, config.ManageNodeKeeperSwitchover, config.ManageNodeKeeperReinitialize,
		config.ManageNodeKeeperRestart, config.ManageNodeKeeperActivation,
	}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for native postgres", feature)
		}
	}
}

// TestDefaultTemplateDefaults covers what replaced keeper.Requirements: the
// deploy screen's credential fields are filled in by the template that creates
// the deployment, because the keeper endpoint is postgres itself, as the
// account POSTGRES_USER creates.
func TestDefaultTemplateDefaults(t *testing.T) {
	for _, template := range NewAdapter().DefaultTemplates() {
		t.Run(template.Name, func(t *testing.T) {
			if template.Defaults.KeeperUser != "postgres" {
				t.Errorf("expected keeper user %q, got %q", "postgres", template.Defaults.KeeperUser)
			}
			if template.Defaults.DbUser != "postgres" {
				t.Errorf("expected database user %q, got %q", "postgres", template.Defaults.DbUser)
			}
		})
	}
}

// TestDefaultTemplates covers the leader/replica asymmetry: it is expressed as
// a different command at index 0, not as a flag. A replica must rebase from
// the leader before postgres ever starts, since streaming replication cannot
// build the initial database.
func TestDefaultTemplates(t *testing.T) {
	templates := NewAdapter().DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			if strings.Contains(template.Commands[0].Command, "pg_basebackup") {
				t.Error("the leader must not run the replica bootstrap")
			}
			// NOTE: on one VM the leader is reached at {{host}} - host
			// networking joins no docker network, so its container name
			// resolves to nothing
			leader := "-h postgres1"
			if strings.Contains(template.Name, "Single Host") {
				leader = "-h {{host}}"
			}
			for i, command := range template.Commands[1:] {
				if !strings.Contains(command.Command, "pg_basebackup") {
					t.Errorf("replica %d does not rebase from the leader", i+1)
				}
				if !strings.Contains(command.Command, leader) {
					t.Errorf("replica %d has no leader host to bootstrap from", i+1)
				}
				// NOTE: the bootstrap script is parsed again by the
				// container's own shell, so it reads the credentials from the
				// env this same command sets rather than interpolating them -
				// a password holding a `$` or a backtick would otherwise be
				// expanded or executed there
				script := command.Command[strings.Index(command.Command, "sh -c '"):]
				for _, v := range []keeper.Var{keeper.VarDbUser, keeper.VarDbPass} {
					if strings.Contains(script, string(v)) {
						t.Errorf("replica %d interpolates %s into a script the container parses again", i+1, v)
					}
				}
				if !strings.Contains(script, `PGPASSWORD="$POSTGRES_PASSWORD"`) {
					t.Errorf("replica %d does not take its password from the env the command sets", i+1)
				}
			}
		})
	}
}
