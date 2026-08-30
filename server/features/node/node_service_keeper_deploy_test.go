package node

import (
	"errors"
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
	"reflect"
	"strings"
	"testing"
)

// newDeployTestService registers every keeper plugin, enough for the pure
// deploy computations that touch no connection.
func newDeployTestService() *Service {
	keeperRegistry := utils.NewRegistry[keeper.PluginType, keeper.Plugin]()
	keeperRegistry.Register(keeper.PATRONI_POSTGRES, patroni.NewAdapter(nil))
	keeperRegistry.Register(keeper.NATIVE_POSTGRES, postgres.NewAdapter())
	keeperRegistry.Register(keeper.NATIVE_ETCD, etcd.NewAdapter())
	keeperRegistry.Register(keeper.NATIVE_REDIS, redis.NewAdapter())
	keeperRegistry.Register(keeper.NATIVE_CLICKHOUSE, clickhouse.NewAdapter())
	keeperRegistry.Register(keeper.NATIVE_ZOOKEEPER, zookeeper.NewAdapter())
	keeperRegistry.Register(keeper.NATIVE_MONGO, mongo.NewAdapter())
	return &Service{
		platformRegistry: utils.NewRegistry[platform.PluginType, platform.Plugin](),
		keeperRegistry:   keeperRegistry,
	}
}

// TestService_ValidateKeeperPlugin covers what is left of the old deploy spec:
// the plugin a deploy names has to be one Ivory actually has, because the
// cluster stores it and every later request resolves its adapter through it.
// What that plugin's deployment consumes is no longer asked of it at all - a
// template states the usernames its own commands create, which is asserted in
// each plugin's own DefaultTemplates test.
func TestService_ValidateKeeperPlugin(t *testing.T) {
	s := newDeployTestService()

	plugins := []KeeperPlugin{
		keeper.PATRONI_POSTGRES, keeper.NATIVE_ETCD, keeper.NATIVE_POSTGRES,
		keeper.NATIVE_REDIS, keeper.NATIVE_CLICKHOUSE, keeper.NATIVE_ZOOKEEPER, keeper.NATIVE_MONGO,
	}
	for _, plugin := range plugins {
		t.Run(string(plugin), func(t *testing.T) {
			if err := s.ValidateKeeperPlugin(plugin); err != nil {
				t.Errorf("ValidateKeeperPlugin(%q) error = %v", plugin, err)
			}
		})
	}

	t.Run("unknown plugin", func(t *testing.T) {
		empty := &Service{keeperRegistry: utils.NewRegistry[keeper.PluginType, keeper.Plugin]()}
		if err := empty.ValidateKeeperPlugin("unknown"); err == nil {
			t.Fatal("expected error for unknown plugin, got nil")
		}
	})
}

func TestService_getValues(t *testing.T) {
	s := &Service{}
	request := KeeperDeployRequest{
		Cluster:    "main",
		Name:       "etcd-1",
		KeeperPort: 8008,
		DbPort:     5432,
		Connection: PlatformVaultConnection{Host: "db1", Port: 2222},
	}

	t.Run("builds the scope from the request alone", func(t *testing.T) {
		got := s.getValues(request)
		want := keeper.Values{
			Cluster:    "main",
			Name:       "etcd-1",
			Host:       "db1",
			SshPort:    "2222",
			KeeperPort: "8008",
			DbPort:     "5432",
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("getValues() = %+v, want %+v", got, want)
		}
	})

	t.Run("credentials are left for the vault at execution time", func(t *testing.T) {
		got := s.getValues(request)
		if got.DbUser != "" || got.DbPass != "" {
			t.Errorf("expected credentials to stay empty until execution, got %q/%q", got.DbUser, got.DbPass)
		}
	})

	// NOTE: this is the whole point of scoping values per command - two nodes
	// of one cluster must not be able to see each other's values
	t.Run("one node's values never reach another's command", func(t *testing.T) {
		other := request
		other.Name = "etcd-2"
		other.DbPort = 2381
		other.Connection.Host = "db2"

		first := keeper.Interpolate("{{name}}@{{host}}:{{dbPort}}", s.getValues(request))
		second := keeper.Interpolate("{{name}}@{{host}}:{{dbPort}}", s.getValues(other))
		if first != "etcd-1@db1:5432" {
			t.Errorf("first node rendered %q", first)
		}
		if second != "etcd-2@db2:2381" {
			t.Errorf("second node rendered %q", second)
		}
	})
}

func TestService_KeeperDeployUpRejectsUnknownVariables(t *testing.T) {
	s := newDeployTestService()

	_, err := s.KeeperDeployUp(KeeperDeployRequest{
		Name:       "db1",
		KeeperPort: 2181,
		DbPort:     2181,
		Command:    `docker run -d --name {{name}} -e ZOO_MY_ID={{index}} zookeeper:3.9`,
		Connection: PlatformVaultConnection{Host: "db1", Port: 22},
	})
	if err == nil {
		t.Fatal("expected an unknown-variable error, got nil")
	}
	if !strings.Contains(err.Error(), "{{index}}") {
		t.Fatalf("expected the error to name the unknown variable, got %v", err)
	}
}

func TestService_KeeperDeployUpRequiresHost(t *testing.T) {
	s := newDeployTestService()

	_, err := s.KeeperDeployUp(KeeperDeployRequest{Name: "db1", Command: "docker run -d redis:7"})
	if err == nil || !strings.Contains(err.Error(), "host") {
		t.Fatalf("expected a missing-host error, got %v", err)
	}
}

// TestService_KeeperDeployUpRequiresPorts covers the deploy contract that
// replaced the port magic: nothing falls back to a database port or to ssh 22.
func TestService_KeeperDeployUpRequiresPorts(t *testing.T) {
	s := newDeployTestService()
	valid := KeeperDeployRequest{
		Name:       "db1",
		KeeperPort: 8008,
		DbPort:     5432,
		Command:    "docker run -d spilo",
		Connection: PlatformVaultConnection{Host: "db1", Port: 22},
	}

	tests := map[string]func(r *KeeperDeployRequest){
		"keeper port unset":   func(r *KeeperDeployRequest) { r.KeeperPort = 0 },
		"database port unset": func(r *KeeperDeployRequest) { r.DbPort = 0 },
		"ssh port unset":      func(r *KeeperDeployRequest) { r.Connection.Port = 0 },
	}
	for name, unset := range tests {
		t.Run(name, func(t *testing.T) {
			r := valid
			unset(&r)
			if _, err := s.KeeperDeployUp(r); !errors.Is(err, ErrKeeperDeployPortsRequired) {
				t.Fatalf("expected ErrKeeperDeployPortsRequired, got %v", err)
			}
		})
	}
}

// TestService_KeeperDeployJudgesTheRequestOnly covers what replaced the
// plugin's Requirements: whether a deployment has keeper or database
// credentials at all is the user's answer on the deploy screen, so a request
// carrying neither vault is judged on its own ports rather than refused here. A
// command that really needs one still fails visibly, on the unresolved
// placeholder it names.
func TestService_KeeperDeployJudgesTheRequestOnly(t *testing.T) {
	s := newDeployTestService()

	for _, plugin := range []KeeperPlugin{keeper.PATRONI_POSTGRES, keeper.NATIVE_ETCD} {
		t.Run(string(plugin), func(t *testing.T) {
			_, err := s.KeeperDeploy(KeeperDeployRequest{
				Plugin:     plugin,
				Name:       "db1",
				Command:    "docker run -d image",
				Connection: PlatformVaultConnection{Host: "db1", Port: 22},
			})
			if !errors.Is(err, ErrKeeperDeployPortsRequired) {
				t.Fatalf("expected the request to be judged on its own ports, got %v", err)
			}
		})
	}

	t.Run("unknown plugin", func(t *testing.T) {
		_, err := s.KeeperDeploy(KeeperDeployRequest{Plugin: "unknown", Name: "db1", Command: "docker run -d redis:7"})
		if err == nil {
			t.Fatal("expected error for unknown plugin, got nil")
		}
	})
}

func TestService_KeeperPostDeployRejectsUnknownVariables(t *testing.T) {
	s := newDeployTestService()

	logs, err := s.KeeperPostDeploy(KeeperDeployRequest{
		Name:        "db1",
		PostScripts: []string{`etcdctl user add {{dbUesr}}`},
		Connection:  PlatformVaultConnection{Host: "db1", Port: 22},
	})
	if len(logs) != 1 || !strings.Contains(logs[0], "{{dbUesr}}") {
		t.Fatalf("expected one log line naming the unknown variable, got %v", logs)
	}
	if err == nil || !strings.Contains(err.Error(), "{{dbUesr}}") {
		t.Fatalf("expected the failure to be returned as an error naming the variable, got %v", err)
	}
}

func TestService_KeeperPostDeploySkipsEmptyScript(t *testing.T) {
	s := newDeployTestService()

	logs, err := s.KeeperPostDeploy(KeeperDeployRequest{Name: "db1"})
	if logs != nil {
		t.Fatalf("expected no logs without a post script, got %v", logs)
	}
	if err != nil {
		t.Fatalf("expected no error without a post script, got %v", err)
	}
}
