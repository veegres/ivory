package postgres

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewPlugin().SupportedFeatures()

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
	for _, template := range NewPlugin().DefaultTemplates() {
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
	templates := NewPlugin().DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			if strings.Contains(template.Commands[0].Command, "pg_basebackup") {
				t.Error("the leader must not run the replica bootstrap")
			}
			// NOTE: the image's own pg_hba.conf grants replication to loopback
			// only - its catch-all "host all all all" excludes the replication
			// pseudo-database - so without this hook a replica connecting from
			// anywhere else matches no line at all
			if !strings.Contains(template.Commands[0].Command, "host replication all all scram-sha-256") {
				t.Error("the leader accepts no replication connection from its replicas")
			}
			// NOTE: on one VM the leader is reached at {{host}}; across VMs it
			// is an example address. Neither is a container name - no docker
			// network spans either layout, so a name resolves to nothing
			leader := "-h 10.0.0.1"
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
				// NOTE: without set -e a failed rebase falls through to the
				// entrypoint, which initdb's a brand-new standalone primary -
				// three of them read back as a successfully deployed cluster
				if !strings.Contains(command.Command, "set -e") {
					t.Errorf("replica %d turns a failed rebase into a standalone primary", i+1)
				}
				// NOTE: pg_basebackup runs as root, so creating $PGDATA also
				// creates its parent root-owned and mode 700, which the
				// entrypoint's gosu phase cannot even traverse
				if !strings.Contains(command.Command, "chown -R postgres:postgres /var/lib/postgresql") {
					t.Errorf("replica %d leaves the data directory's parent root-owned", i+1)
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
