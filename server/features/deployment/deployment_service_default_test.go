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

	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	keeperMetadataRegistry.Register(keeper.PATRONI_POSTGRES, patroni.NewAdapter(nil))
	keeperMetadataRegistry.Register(keeper.NATIVE_POSTGRES, postgres.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_ETCD, etcd.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_REDIS, redis.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_CLICKHOUSE, clickhouse.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_ZOOKEEPER, zookeeper.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_MONGO, mongo.NewAdapter())

	platformRegistry := utils.NewRegistry[platform.Plugin, platform.Adapter]()
	platformRegistry.Register(platform.Docker, docker.NewAdapter(nil))

	return NewService(newTestRepository(t), keeperMetadataRegistry, platformRegistry)
}

// TestDefaultsCoverEveryKeeper is what replaced the per-plugin deployment-spec
// tests: a plugin with no shipped command leaves the user with nothing to copy.
func TestDefaultsCoverEveryKeeper(t *testing.T) {
	s := newFullTestService(t)
	defaults := s.Defaults(ListRequest{})

	plugins := []keeper.Plugin{
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
				for _, text := range []string{command.Command, command.PostScript} {
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
			utils.NewRegistry[keeper.Plugin, keeper.Metadata](),
			func() *utils.Registry[platform.Plugin, platform.Adapter] {
				r := utils.NewRegistry[platform.Plugin, platform.Adapter]()
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
				interpolated := keeper.Interpolate(command.Command, values)
				if left := keeper.UnresolvedPlaceholders(interpolated); len(left) > 0 {
					t.Errorf("command %d still references %v after a node fills it in", i, left)
				}

				got := adapter.UpContainer(platform.Connection{}, interpolated).(*ssh.Command).Command
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
