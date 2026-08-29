package postgres

import (
	"ivory/core/config"
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

func TestRequirements(t *testing.T) {
	req := NewAdapter().Requirements()

	if req.DbPort != 5432 {
		t.Errorf("expected db port 5432, got %d", req.DbPort)
	}
	if req.KeeperPort != 5432 {
		t.Errorf("expected the keeper endpoint to be declared as 5432 (plain postgres has no management api), got %d", req.KeeperPort)
	}
	if !req.KeeperCredentials || req.KeeperUser != "" {
		t.Errorf("expected keeper credentials with a username of the user's own choice, got %v/%q", req.KeeperCredentials, req.KeeperUser)
	}
	if !req.DbCredentials {
		t.Error("expected postgres to consume database credentials")
	}
	if req.DbUser != "" {
		t.Errorf("expected a free choice of username, got the locked %q", req.DbUser)
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
			leader := "host=postgres1"
			if strings.Contains(template.Name, "Single Host") {
				leader = "host={{host}}"
			}
			for i, command := range template.Commands[1:] {
				if !strings.Contains(command.Command, "pg_basebackup") {
					t.Errorf("replica %d does not rebase from the leader", i+1)
				}
				if !strings.Contains(command.Command, leader) {
					t.Errorf("replica %d has no leader host to bootstrap from", i+1)
				}
			}
		})
	}
}
