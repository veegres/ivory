package zookeeper

import (
	"ivory/core/config"
	"strconv"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ViewNodeKeeperOverview, config.ViewNodeKeeperConfig}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for zookeeper", feature)
		}
	}

	excluded := []config.Feature{
		config.ManageNodeKeeperConfigUpdate, config.ManageNodeKeeperSwitchover, config.ManageNodeKeeperReinitialize,
		config.ManageNodeKeeperRestart, config.ManageNodeKeeperReload, config.ManageNodeKeeperFailover, config.ManageNodeKeeperActivation,
	}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for zookeeper", feature)
		}
	}
}

func TestRequirements(t *testing.T) {
	req := NewAdapter().Requirements()

	if req.KeeperPort != 2181 {
		t.Errorf("expected the keeper endpoint to be declared as 2181, got %d", req.KeeperPort)
	}
	if req.KeeperCredentials {
		t.Error("expected no keeper credentials: zookeeper ships without auth")
	}
	if req.DbCredentials {
		t.Error("expected no credentials: zookeeper ships without auth")
	}
}

// TestDefaultTemplates covers the one thing no host-derived value can provide:
// a unique small integer per member, written literally into each command.
func TestDefaultTemplates(t *testing.T) {
	templates := NewAdapter().DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			for i, command := range template.Commands {
				id := strconv.Itoa(i + 1)
				if !strings.Contains(command.Command, `ZOO_MY_ID="`+id+`"`) {
					t.Errorf("command %d does not carry myid %s", i, id)
				}
				if !strings.Contains(command.Command, "ZOO_4LW_COMMANDS_WHITELIST") {
					t.Errorf("command %d does not whitelist the four-letter commands List/Config need", i)
				}
				if getServers(command.Command) != getServers(template.Commands[0].Command) {
					t.Errorf("command %d has a different server list, which would split the ensemble", i)
				}
			}
		})
	}
}

// getServers pulls the ZOO_SERVERS value out of a command, so the test can
// compare the ensemble's view without depending on the rest of the text.
func getServers(command string) string {
	const prefix = `ZOO_SERVERS="`
	start := strings.Index(command, prefix)
	if start < 0 {
		return ""
	}
	rest := command[start+len(prefix):]
	end := strings.Index(rest, `"`)
	if end < 0 {
		return ""
	}
	return rest[:end]
}
