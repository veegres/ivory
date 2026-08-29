package mongo

import (
	"ivory/core/config"
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

func TestRequirements(t *testing.T) {
	req := NewAdapter().Requirements()

	if req.DbPort != 27017 {
		t.Errorf("expected db port 27017, got %d", req.DbPort)
	}
	if req.KeeperPort != 27017 {
		t.Errorf("expected the keeper endpoint to be declared as 27017, got %d", req.KeeperPort)
	}
	if req.KeeperCredentials {
		t.Error("expected no keeper credentials: the deployed replica set runs unauthenticated")
	}
	if req.DbCredentials {
		t.Error("expected no credentials: replica set auth also needs an internal keyfile")
	}
}

// TestDefaultTemplates covers the replica set bootstrap: every member starts
// as a plain mongod, and the last command initiates the set once they are all
// up - which is why rs.initiate() sits there rather than on the first.
func TestDefaultTemplates(t *testing.T) {
	templates := NewAdapter().DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			last := template.Commands[len(template.Commands)-1]
			if !strings.Contains(last.PostScript, "rs.initiate") {
				t.Fatal("expected the last command to initiate the replica set")
			}
			if !strings.Contains(last.PostScript, `mongo-1:27017`) {
				t.Error("expected the members list rs.initiate is given")
			}
			for i, command := range template.Commands[:len(template.Commands)-1] {
				if command.PostScript != "" {
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
