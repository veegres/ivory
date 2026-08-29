package clickhouse

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strconv"
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

	if req.KeeperPort != 9000 {
		t.Errorf("expected the keeper endpoint to be declared as 9000, got %d", req.KeeperPort)
	}
	if !req.KeeperCredentials || req.KeeperUser != "" {
		t.Errorf("expected keeper credentials with a username of the user's own choice, got %v/%q", req.KeeperCredentials, req.KeeperUser)
	}
	if !req.DbCredentials {
		t.Error("expected clickhouse to consume database credentials")
	}
	if req.DbUser != "" {
		t.Errorf("expected a free choice of username, got the locked %q", req.DbUser)
	}
}

// TestDefaultTemplates covers clickhouse's lack of leader/replica asymmetry -
// every node generates the same cluster config file, so the multi-host commands
// are identical - and the one thing that does separate the single-host nodes:
// three servers on one port namespace need three sets of listening ports.
func TestDefaultTemplates(t *testing.T) {
	templates := NewAdapter().DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			singleHost := strings.Contains(template.Name, "Single Host")
			for i, command := range template.Commands {
				if !singleHost && command.Command != template.Commands[0].Command {
					t.Errorf("command %d differs, but clickhouse has no leader/replica asymmetry", i)
				}
				if !strings.Contains(command.Command, "ivory-cluster.xml") {
					t.Errorf("command %d does not generate the cluster config", i)
				}
				if !strings.Contains(command.Command, "<node><host>keeper-1</host>") {
					t.Errorf("command %d is missing the coordinator list", i)
				}
			}
			if singleHost {
				assertSingleHostPorts(t, template)
				return
			}
			for i, command := range template.Commands {
				if !strings.Contains(command.Command, "<replica><host>clickhouse1</host>") {
					t.Errorf("command %d is missing the shard's replica list", i)
				}
			}
		})
	}
}

// assertSingleHostPorts checks the ports three replicas on one VM cannot share:
// the native port Ivory connects on plus the http and interserver ports the
// image binds on its own.
func assertSingleHostPorts(t *testing.T, template keeper.DeploymentTemplate) {
	t.Helper()

	seen := map[int]bool{}
	for i, command := range template.Commands {
		if !strings.Contains(command.Command, "<replica><host>{{host}}</host>") {
			t.Errorf("command %d should address its replicas by {{host}}, which is all host networking resolves", i)
		}
		if !strings.Contains(command.Command, "<tcp_port>{{dbPort}}</tcp_port>") {
			t.Errorf("command %d does not override the image's fixed native port", i)
		}
		if strings.Contains(command.Command, "--hostname") {
			t.Errorf("command %d sets --hostname, which docker rejects alongside --network host", i)
		}
		for _, tag := range []string{"http_port", "interserver_http_port"} {
			port := portOf(t, command.Command, tag)
			if seen[port] {
				t.Errorf("command %d reuses %s %d, which collides on one host", i, tag, port)
			}
			seen[port] = true
		}
		if port := template.Commands[i].Defaults.DbPort; port == 0 {
			t.Errorf("command %d states no default database port", i)
		}
	}
}

func portOf(t *testing.T, command string, tag string) int {
	t.Helper()

	open, close := "<"+tag+">", "</"+tag+">"
	from := strings.Index(command, open)
	to := strings.Index(command, close)
	if from < 0 || to < from {
		t.Fatalf("command has no %s", tag)
	}
	port, err := strconv.Atoi(command[from+len(open) : to])
	if err != nil {
		t.Fatalf("unreadable %s: %v", tag, err)
	}
	return port
}
