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
	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	keeperMetadataRegistry.Register(keeper.PATRONI_POSTGRES, patroni.NewAdapter(nil))
	keeperMetadataRegistry.Register(keeper.NATIVE_POSTGRES, postgres.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_ETCD, etcd.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_REDIS, redis.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_CLICKHOUSE, clickhouse.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_ZOOKEEPER, zookeeper.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_MONGO, mongo.NewAdapter())
	return &Service{
		platformRegistry:       utils.NewRegistry[platform.Plugin, platform.Adapter](),
		keeperMetadataRegistry: keeperMetadataRegistry,
	}
}

func intPtr(v int) *int {
	return &v
}

func TestService_KeeperDeploySpec(t *testing.T) {
	s := newDeployTestService()

	tests := []struct {
		name     string
		plugin   KeeperPlugin
		expected KeeperDeploySpecResponse
	}{
		{
			name:     "patroni exposes its own keeper endpoint and a locked superuser",
			plugin:   keeper.PATRONI_POSTGRES,
			expected: KeeperDeploySpecResponse{DbPort: 5432, KeeperPort: intPtr(8008), Credentials: true, DbUser: "postgres"},
		},
		{
			name:     "etcd serves its keeper endpoint on the database port",
			plugin:   keeper.NATIVE_ETCD,
			expected: KeeperDeploySpecResponse{DbPort: 2379, Credentials: true, DbUser: "root"},
		},
		{
			name:     "native postgres consumes credentials but leaves the username free",
			plugin:   keeper.NATIVE_POSTGRES,
			expected: KeeperDeploySpecResponse{DbPort: 5432, Credentials: true},
		},
		{
			name:     "zookeeper consumes no credentials at all",
			plugin:   keeper.NATIVE_ZOOKEEPER,
			expected: KeeperDeploySpecResponse{DbPort: 2181},
		},
		{
			name:     "mongo consumes no credentials at all",
			plugin:   keeper.NATIVE_MONGO,
			expected: KeeperDeploySpecResponse{DbPort: 27017},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.KeeperDeploySpec(KeeperDeploySpecRequest{Plugin: tt.plugin})
			if err != nil {
				t.Fatalf("KeeperDeploySpec() error = %v", err)
			}
			if !reflect.DeepEqual(*got, tt.expected) {
				t.Errorf("KeeperDeploySpec() = %+v, want %+v", *got, tt.expected)
			}
		})
	}

	t.Run("unknown plugin", func(t *testing.T) {
		empty := &Service{keeperMetadataRegistry: utils.NewRegistry[keeper.Plugin, keeper.Metadata]()}
		if _, err := empty.KeeperDeploySpec(KeeperDeploySpecRequest{Plugin: "unknown"}); err == nil {
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

func TestService_KeeperDeployRequiresDatabaseCredentials(t *testing.T) {
	s := newDeployTestService()

	t.Run("a plugin that consumes credentials needs a database vault", func(t *testing.T) {
		_, err := s.KeeperDeploy(KeeperDeployRequest{
			Plugin:     keeper.PATRONI_POSTGRES,
			Name:       "db1",
			Command:    "docker run -d spilo",
			Connection: PlatformVaultConnection{Host: "db1", Port: 22},
		})
		if !errors.Is(err, ErrKeeperDeployDatabaseCredentialsRequired) {
			t.Fatalf("expected ErrKeeperDeployDatabaseCredentialsRequired, got %v", err)
		}
	})

	t.Run("unknown plugin", func(t *testing.T) {
		_, err := s.KeeperDeploy(KeeperDeployRequest{Plugin: "unknown", Name: "db1", Command: "docker run -d redis:7"})
		if err == nil {
			t.Fatal("expected error for unknown plugin, got nil")
		}
	})
}

func TestService_KeeperPostDeployRejectsUnknownVariables(t *testing.T) {
	s := newDeployTestService()

	logs := s.KeeperPostDeploy(KeeperDeployRequest{
		Name:       "db1",
		PostScript: `etcdctl user add {{dbUesr}}`,
		Connection: PlatformVaultConnection{Host: "db1", Port: 22},
	})
	if len(logs) != 1 || !strings.Contains(logs[0], "{{dbUesr}}") {
		t.Fatalf("expected one log line naming the unknown variable, got %v", logs)
	}
}

func TestService_KeeperPostDeploySkipsEmptyScript(t *testing.T) {
	s := newDeployTestService()

	if logs := s.KeeperPostDeploy(KeeperDeployRequest{Name: "db1"}); logs != nil {
		t.Fatalf("expected no logs without a post script, got %v", logs)
	}
}
