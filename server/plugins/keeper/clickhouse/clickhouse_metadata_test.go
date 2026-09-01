package clickhouse

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strconv"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewPlugin().SupportedFeatures()

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

// TestDefaultTemplateDefaults covers what replaced keeper.Requirements: the
// deploy screen's credential fields are filled in by the template that creates
// the deployment, because clickhouse answers keeper and database questions on
// the same native endpoint, as the account CLICKHOUSE_USER creates. The DCS
// default is there for the same reason - {{dcs}} carries the text clickhouse's
// own config expects, so the template that wrote the command is the only thing
// that can say what shape that is.
func TestDefaultTemplateDefaults(t *testing.T) {
	for _, template := range NewPlugin().DefaultTemplates() {
		t.Run(template.Name, func(t *testing.T) {
			if template.Defaults.KeeperUser != "default" {
				t.Errorf("expected keeper user %q, got %q", "default", template.Defaults.KeeperUser)
			}
			if template.Defaults.DbUser != "default" {
				t.Errorf("expected database user %q, got %q", "default", template.Defaults.DbUser)
			}
			// NOTE: three nodes on one VM share a loopback but not a port, and
			// under --network host that loopback is where the shipped zookeeper
			// single-host template listens. A multi-host ensemble runs on
			// machines a template cannot know, so it names none at all.
			if !strings.Contains(template.Name, "Single Host") {
				if template.Defaults.Dcs != "" {
					t.Errorf("expected no DCS default for a multi-host ensemble, got %q", template.Defaults.Dcs)
				}
				return
			}
			nodes := []string{
				"<node><host>localhost</host><port>2181</port></node>",
				"<node><host>localhost</host><port>2182</port></node>",
				"<node><host>localhost</host><port>2183</port></node>",
			}
			for _, node := range nodes {
				if !strings.Contains(template.Defaults.Dcs, node) {
					t.Errorf("the DCS default is missing coordinator %s, got %q", node, template.Defaults.Dcs)
				}
			}
		})
	}
}

// TestDefaultTemplates covers clickhouse's lack of leader/replica asymmetry -
// every node generates the same cluster config file, so the multi-host commands
// are identical - and the one thing that does separate the single-host nodes:
// three servers on one port namespace need three sets of listening ports.
func TestDefaultTemplates(t *testing.T) {
	templates := NewPlugin().DefaultTemplates()

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
				// NOTE: the image's entrypoint ends in exec "$@", so a script
				// passed the ordinary way is the last thing it runs and the
				// config file lands after a server is already up on the
				// untouched one. --entrypoint is what puts it first
				if !strings.Contains(command.Command, "--entrypoint sh") {
					t.Errorf("command %d writes its config after the entrypoint has already started a server", i)
				}
				// NOTE: CLICKHOUSE_DB is the sole trigger for an init pass
				// whose client is hardcoded to 127.0.0.1 with no --port, so on
				// a shared VM it talks to another node's server and kills this
				// one
				if strings.Contains(command.Command, "CLICKHOUSE_DB") {
					t.Errorf("command %d runs the entrypoint's port-blind init pass", i)
				}
				// NOTE: a replica registers its fetch endpoint under this name;
				// left to default to the machine's hostname it does not
				// resolve, and the replicas silently never catch up
				if !strings.Contains(command.Command, "<interserver_http_host>{{host}}</interserver_http_host>") {
					t.Errorf("command %d leaves its fetch endpoint on an unresolvable hostname", i)
				}
				// NOTE: the ensemble is one answer for the whole deployment, so
				// it is the template's default rather than text in the command
				if !strings.Contains(command.Command, string(keeper.VarDcs)) {
					t.Errorf("command %d does not point its coordinator list at {{dcs}}", i)
				}
			}
			if singleHost {
				assertSingleHostPorts(t, template)
				return
			}
			for i, command := range template.Commands {
				if !strings.Contains(command.Command, "<replica><host>10.0.0.1</host>") {
					t.Errorf("command %d is missing the shard's replica list", i)
				}
			}
		})
	}
}

// assertSingleHostPorts checks the ports three replicas on one VM cannot share:
// the native port Ivory connects on plus every other port the image binds on
// its own. mysql and postgresql are in the list despite nothing using them -
// the image binds 9004 and 9005 whether or not they are wanted, so they
// collide like any other.
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
		for _, tag := range []string{"http_port", "interserver_http_port", "mysql_port", "postgresql_port"} {
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

// TestHasLeaderIsFalse pins the one thing that separates clickhouse from every
// other keeper: it has no single-primary model, so nothing elects a leader and
// the overview must not warn that none was found. It cannot be inferred from
// SupportedFeatures - zookeeper declares neither switchover nor failover and
// still reports one.
func TestHasLeaderIsFalse(t *testing.T) {
	if NewPlugin().HasLeader() {
		t.Error("clickhouse coordinates through ClickHouse Keeper and elects no leader")
	}
}
