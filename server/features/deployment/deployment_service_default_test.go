package deployment

import (
	"ivory/clients/console/ssh"
	"ivory/core/utils"
	"ivory/plugins/keeper"
	"ivory/plugins/keeper/clickhouse"
	"ivory/plugins/keeper/etcd"
	"ivory/plugins/keeper/mongo"
	"ivory/plugins/keeper/patroni"
	"ivory/plugins/keeper/postgres"
	"ivory/plugins/keeper/redis"
	"ivory/plugins/keeper/zookeeper"
	"ivory/plugins/platform"
	"ivory/plugins/platform/docker"
	"strings"
	"testing"
)

// newFullTestService registers every keeper plugin, so the shipped catalog is
// visible in full.
func newFullTestService(t *testing.T) *Service {
	t.Helper()

	keeperRegistry := utils.NewRegistry[keeper.PluginType, keeper.Plugin]()
	keeperRegistry.Register(keeper.PATRONI_POSTGRES, patroni.NewAdapter(nil))
	keeperRegistry.Register(keeper.NATIVE_POSTGRES, postgres.NewAdapter())
	keeperRegistry.Register(keeper.NATIVE_ETCD, etcd.NewAdapter())
	keeperRegistry.Register(keeper.NATIVE_REDIS, redis.NewAdapter())
	keeperRegistry.Register(keeper.NATIVE_CLICKHOUSE, clickhouse.NewAdapter())
	keeperRegistry.Register(keeper.NATIVE_ZOOKEEPER, zookeeper.NewAdapter())
	keeperRegistry.Register(keeper.NATIVE_MONGO, mongo.NewAdapter())

	platformRegistry := utils.NewRegistry[platform.PluginType, platform.Plugin]()
	platformRegistry.Register(platform.Docker, docker.NewAdapter(nil))

	return NewService(newTestRepository(t), keeperRegistry, platformRegistry)
}

// TestDefaultsCoverEveryKeeper is what replaced the per-plugin deployment-spec
// tests: a plugin with no shipped command leaves the user with nothing to copy.
func TestDefaultsCoverEveryKeeper(t *testing.T) {
	s := newFullTestService(t)
	defaults := s.Defaults(ListRequest{})

	plugins := []keeper.PluginType{
		keeper.PATRONI_POSTGRES, keeper.NATIVE_POSTGRES, keeper.NATIVE_ETCD,
		keeper.NATIVE_REDIS, keeper.NATIVE_CLICKHOUSE, keeper.NATIVE_ZOOKEEPER, keeper.NATIVE_MONGO,
	}

	for _, plugin := range plugins {
		t.Run(string(plugin), func(t *testing.T) {
			var multiHost, singleHost int
			for _, template := range defaults {
				if template.Keeper != plugin {
					continue
				}
				if strings.Contains(template.Name, "Single Host") {
					singleHost++
				} else {
					multiHost++
				}
			}
			if multiHost != 1 || singleHost != 1 {
				t.Errorf("expected exactly one multi-host and one single-host template, got %d/%d", multiHost, singleHost)
			}
		})
	}
}

func TestDefaultsAreWellFormed(t *testing.T) {
	s := newFullTestService(t)

	for _, template := range s.Defaults(ListRequest{}) {
		t.Run(template.Name, func(t *testing.T) {
			if template.Creation != System {
				t.Errorf("expected a system template, got %q", template.Creation)
			}
			if template.Description == "" {
				t.Error("expected a description explaining what the template deploys")
			}
			if len(template.Commands) == 0 {
				t.Fatal("expected at least one command")
			}
			for i, command := range template.Commands {
				if command.Command == "" {
					t.Fatalf("command %d is empty", i)
				}
				// NOTE: the catalog is hand-written, so a typo would otherwise
				// only surface as an unresolved placeholder at deploy time
				for _, text := range append([]string{command.Command}, command.PostScripts...) {
					if unknown := keeper.UnknownPlaceholders(text); len(unknown) > 0 {
						t.Errorf("command %d references unknown variables: %v", i, unknown)
					}
				}
				if !strings.Contains(command.Command, string(keeper.VarName)) {
					t.Errorf("command %d does not name its deployment with {{name}}", i)
				}
			}
		})
	}
}

// TestDefaultsKeepCredentialsOutOfNestedScripts holds the one boundary
// interpolation cannot cross. A value is filled into an argument after the
// command has been split, so nothing parses it - except a shell the command
// itself starts. Where a template's argument is "sh -c '...'", the container
// parses that script a second time and a `$` or a backtick in a password is
// expansion again. Such a script reads what it needs from env its own command
// sets, the way postgres reads $POSTGRES_PASSWORD.
//
// A post script step is a plain command with no shell, so credentials there
// are fine - which is exactly why etcd's auth steps can name {{dbPass}}.
func TestDefaultsKeepCredentialsOutOfNestedScripts(t *testing.T) {
	s := newFullTestService(t)
	credentials := []keeper.Var{keeper.VarKeeperUser, keeper.VarKeeperPass, keeper.VarDbUser, keeper.VarDbPass}

	for _, template := range s.Defaults(ListRequest{}) {
		t.Run(template.Name, func(t *testing.T) {
			for i, command := range template.Commands {
				texts := append([]string{command.Command}, command.PostScripts...)
				for _, text := range texts {
					script := nestedScript(text)
					for _, v := range credentials {
						if strings.Contains(script, string(v)) {
							t.Errorf("command %d interpolates %s into a script the container parses again", i, v)
						}
					}
				}
			}
		})
	}
}

// nestedScript returns the part of a command a second shell parses - the tail
// after "sh -c '" - or empty when the command runs no script of its own.
func nestedScript(command string) string {
	i := strings.Index(command, "sh -c '")
	if i < 0 {
		return ""
	}
	return command[i:]
}

// TestSingleHostDefaultsAvoidCollisions covers what the deleted singleHost
// flag used to compute: several nodes on one VM need host networking, no
// published ports, no shared volume, and a distinct peer port per member.
func TestSingleHostDefaultsAvoidCollisions(t *testing.T) {
	s := newFullTestService(t)

	for _, template := range s.Defaults(ListRequest{}) {
		if !strings.Contains(template.Name, "Single Host") {
			continue
		}
		t.Run(template.Name, func(t *testing.T) {
			for i, command := range template.Commands {
				if !strings.Contains(command.Command, "--network host") {
					t.Errorf("command %d must use host networking", i)
				}
				if strings.Contains(command.Command, "\n  -p ") {
					t.Errorf("command %d publishes a port, which collides on one VM", i)
				}
				if strings.Contains(command.Command, "\n  -v ") {
					t.Errorf("command %d mounts a volume, which collides on one VM", i)
				}
				// NOTE: docker refuses --hostname and --network host together,
				// so a command carrying both never starts at all
				if strings.Contains(command.Command, "--hostname") {
					t.Errorf("command %d sets --hostname, which docker rejects alongside --network host", i)
				}
			}
		})
	}
}

// TestDefaultsStateNodeDefaults holds the promise the deploy form depends on:
// a shipped template fills its own node cards in, so deploying one asks for
// nothing but the host and its ssh port. There is no fallback behind this - a
// command that states no port leaves the field empty for the user to fill in,
// so a shipped template that forgot one would not deploy at all.
func TestDefaultsStateNodeDefaults(t *testing.T) {
	s := newFullTestService(t)

	for _, template := range s.Defaults(ListRequest{}) {
		t.Run(template.Name, func(t *testing.T) {
			names := map[string]bool{}
			for i, command := range template.Commands {
				if command.Defaults.Name == "" {
					t.Errorf("command %d states no default node name", i)
				}
				if names[command.Defaults.Name] {
					t.Errorf("command %d reuses the node name %q, which a cluster rejects", i, command.Defaults.Name)
				}
				names[command.Defaults.Name] = true
				if command.Defaults.KeeperPort == 0 {
					t.Errorf("command %d states no default keeper port", i)
				}
				if command.Defaults.DbPort == 0 {
					t.Errorf("command %d states no default database port", i)
				}
			}
		})
	}
}

// TestDefaultsStateCredentialDefaults holds the rule that replaced
// keeper.Requirements: the deploy screen opens a credential pair filled in when
// the template names a username for it, and switched off when it does not. So a
// template whose commands reference {{dbUser}}/{{dbPass}} has to name the
// database user it creates - otherwise the screen would open with the pair off
// and the deploy would fail on a placeholder nothing filled in.
//
// The reverse is deliberately not asserted. A template names a username for a
// pair its commands never mention whenever Ivory needs the account to reach the
// engine afterwards rather than to start it: etcd's root, clickhouse's keeper
// endpoint being the database itself.
func TestDefaultsStateCredentialDefaults(t *testing.T) {
	s := newFullTestService(t)

	for _, template := range s.Defaults(ListRequest{}) {
		t.Run(template.Name, func(t *testing.T) {
			pairs := []struct {
				label string
				vars  []keeper.Var
				user  string
			}{
				{label: "keeper", vars: []keeper.Var{keeper.VarKeeperUser, keeper.VarKeeperPass}, user: template.Defaults.KeeperUser},
				{label: "database", vars: []keeper.Var{keeper.VarDbUser, keeper.VarDbPass}, user: template.Defaults.DbUser},
			}
			for _, pair := range pairs {
				if !referencesAny(template, pair.vars) {
					continue
				}
				if pair.user == "" {
					t.Errorf("the commands reference the %s credentials but the template names no %s user", pair.label, pair.label)
				}
			}
		})
	}
}

// referencesAny reports whether any of a template's commands or post scripts
// names one of the given variables.
func referencesAny(template Template, vars []keeper.Var) bool {
	for _, command := range template.Commands {
		for _, text := range append([]string{command.Command}, command.PostScripts...) {
			for _, v := range vars {
				if strings.Contains(text, string(v)) {
					return true
				}
			}
		}
	}
	return false
}

// TestSingleHostDefaultsGiveEachNodeItsOwnPorts is the collision the shipped
// single-host templates existed to demonstrate and could not survive: three
// nodes sharing one port namespace need three sets of ports, and the deploy
// form takes them from the commands themselves.
func TestSingleHostDefaultsGiveEachNodeItsOwnPorts(t *testing.T) {
	s := newFullTestService(t)

	for _, template := range s.Defaults(ListRequest{}) {
		if !strings.Contains(template.Name, "Single Host") {
			continue
		}
		t.Run(template.Name, func(t *testing.T) {
			keeperPorts := map[int]bool{}
			dbPorts := map[int]bool{}
			for i, command := range template.Commands {
				if keeperPorts[command.Defaults.KeeperPort] {
					t.Errorf("command %d reuses keeper port %d on the same VM", i, command.Defaults.KeeperPort)
				}
				if dbPorts[command.Defaults.DbPort] {
					t.Errorf("command %d reuses database port %d on the same VM", i, command.Defaults.DbPort)
				}
				keeperPorts[command.Defaults.KeeperPort] = true
				dbPorts[command.Defaults.DbPort] = true
			}
		})
	}
}

func TestMultiHostDefaultsPublishPortsAndVolumes(t *testing.T) {
	s := newFullTestService(t)

	for _, template := range s.Defaults(ListRequest{}) {
		if strings.Contains(template.Name, "Single Host") {
			continue
		}
		t.Run(template.Name, func(t *testing.T) {
			command := template.Commands[0].Command
			if !strings.Contains(command, "--restart unless-stopped") {
				t.Error("expected a restart policy on a one-node-per-VM deployment")
			}
			if !strings.Contains(command, "\n  -p ") {
				t.Error("expected at least one published port")
			}
			if !strings.Contains(command, "\n  -v ") {
				t.Error("expected a data volume")
			}
			if strings.Contains(command, "--network host") {
				t.Error("host networking belongs to the single-host template only")
			}
		})
	}
}

func TestDefaultsFilterByPlugin(t *testing.T) {
	s := newFullTestService(t)

	t.Run("by keeper", func(t *testing.T) {
		for _, template := range s.Defaults(ListRequest{Keeper: keeperPtr(keeper.NATIVE_ETCD)}) {
			if template.Keeper != keeper.NATIVE_ETCD {
				t.Errorf("expected only etcd templates, got %q", template.Keeper)
			}
		}
	})

	t.Run("an unregistered platform has no templates", func(t *testing.T) {
		if got := s.Defaults(ListRequest{Platform: platformPtr("k8s")}); len(got) != 0 {
			t.Errorf("expected no templates for an unregistered platform, got %d", len(got))
		}
	})

	// NOTE: a keeper plugin that is compiled out must not leave its shipped
	// templates dangling in the list
	t.Run("an unregistered keeper has no templates", func(t *testing.T) {
		empty := NewService(
			newTestRepository(t),
			utils.NewRegistry[keeper.PluginType, keeper.Plugin](),
			func() *utils.Registry[platform.PluginType, platform.Plugin] {
				r := utils.NewRegistry[platform.PluginType, platform.Plugin]()
				r.Register(platform.Docker, docker.NewAdapter(nil))
				return r
			}(),
		)
		if got := empty.Defaults(ListRequest{}); len(got) != 0 {
			t.Errorf("expected no templates without keeper plugins, got %d", len(got))
		}
	})
}

// TestDefaultsProduceRunnableCommands walks every shipped command all the way
// to the shell command that would actually run: interpolate it with a node's
// values, then hand it to the platform adapter. It is the one test that covers
// the whole chain, and the only place a broken template would surface before a
// user hits it.
func TestDefaultsProduceRunnableCommands(t *testing.T) {
	s := newFullTestService(t)
	adapter := docker.NewAdapter(ssh.NewClient())

	values := keeper.Values{
		Cluster: "main", Name: "node-1", Host: "10.0.0.1",
		SshPort: "22", KeeperPort: "8008", DbPort: "5432",
		DbUser: "ivory", DbPass: "s3cr3t",
	}

	for _, template := range s.Defaults(ListRequest{}) {
		t.Run(template.Name, func(t *testing.T) {
			for i, command := range template.Commands {
				if left := keeper.UnresolvedPlaceholders(keeper.Interpolate(command.Command, values)); len(left) > 0 {
					t.Errorf("command %d still references %v after a node fills it in", i, left)
				}

				// NOTE: split first, fill each argument in, then hand the
				// adapter the finished command - the real execution path
				args := platform.SplitCommand(command.Command)
				for j, argument := range args {
					args[j] = keeper.Interpolate(argument, values)
				}
				got := adapter.UpContainer(platform.Connection{}, args).(*ssh.Command).Command
				if !strings.HasPrefix(got, "docker run -d") {
					t.Errorf("command %d does not run a container: %q", i, got)
				}
				if strings.Contains(got, "{{") {
					t.Errorf("command %d leaked a placeholder into the shell", i)
				}
				// NOTE: readability newlines between flags must collapse, but
				// the ones inside a quoted startup script have to survive -
				// they are real statement separators
				if strings.Contains(got, "\n") && !strings.Contains(got, "sh -c") {
					t.Errorf("command %d has a newline outside a quoted script", i)
				}
			}
		})
	}
}
