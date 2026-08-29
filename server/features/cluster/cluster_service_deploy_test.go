package cluster

import (
	"errors"
	"ivory/clients/console/ssh"
	"ivory/clients/storage"
	"ivory/core/utils"
	"ivory/features/node"
	"ivory/features/tag"
	"ivory/plugins/keeper"
	"ivory/plugins/keeper/etcd"
	"ivory/plugins/keeper/patroni"
	"ivory/plugins/keeper/postgres"
	"ivory/plugins/platform"
	"ivory/plugins/platform/docker"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

// newDeployTestService builds a cluster Service with a real repository and
// real keeper/platform plugin registries, enough to exercise Deploy's
// validation path without ever touching the network: no vaultService is wired
// in, so tests must stay on the side of Deploy that returns before any vault
// or SSH work happens.
func newDeployTestService(t *testing.T) *Service {
	t.Helper()

	tmpDir := t.TempDir()
	db, errOpen := bolt.Open(filepath.Join(tmpDir, "test.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() {
		db.Close()
	})

	clusterRepository := NewRepository(storage.NewDbBucket[Response](db, "Cluster"))
	tagRepository := tag.NewRepository(storage.NewDbBucket[[]string](db, "Tag"))
	tagService := tag.NewService(tagRepository)

	platformRegistry := utils.NewRegistry[platform.PluginType, platform.Plugin]()
	platformRegistry.Register(platform.Docker, docker.NewAdapter(ssh.NewClient()))
	keeperRegistry := utils.NewRegistry[keeper.PluginType, keeper.Plugin]()
	keeperRegistry.Register(keeper.PATRONI_POSTGRES, patroni.NewAdapter(nil))
	keeperRegistry.Register(keeper.NATIVE_POSTGRES, postgres.NewAdapter())
	keeperRegistry.Register(keeper.NATIVE_ETCD, etcd.NewAdapter())
	nodeService := node.NewService(platformRegistry, keeperRegistry, nil, nil, nil)

	return &Service{
		clusterRepository: clusterRepository,
		tagService:        tagService,
		nodeService:       nodeService,
	}
}

func validDeployRequest() DeployRequest {
	sshPort, keeperPort, dbPort := 22, 8008, 5432
	return DeployRequest{
		Nodes: []DeployNode{{
			NodeConfig: NodeConfig{Name: "db-1", Host: "db1", SshPort: &sshPort, KeeperPort: &keeperPort, DbPort: &dbPort},
			Command:    `docker run -d --name {{name}} -e ETCD3_HOSTS="etcd1:2379" spilo`,
		}},
		CommonConfig: CommonConfig{
			Cluster: "test-cluster",
			SshUser: "root",
			SshPass: "secret",
			DbPass:  "secret",
		},
		ClusterOptions: Options{Plugins: Plugins{Keeper: node.KeeperPlugin(keeper.PATRONI_POSTGRES)}},
	}
}

// etcdDeployRequest is a request for a plugin that consumes keeper credentials
// as well as database ones, which patroni's shipped deployment does not.
func etcdDeployRequest() DeployRequest {
	r := validDeployRequest()
	r.ClusterOptions.Plugins.Keeper = node.KeeperPlugin(keeper.NATIVE_ETCD)
	r.CommonConfig.KeeperPass = "secret"
	return r
}

func TestService_Deploy_ValidationErrors(t *testing.T) {
	t.Run("cluster name not provided", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		r.CommonConfig.Cluster = ""

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrClusterNameNotProvided) {
			t.Fatalf("expected ErrClusterNameNotProvided, got %v", err)
		}
	})

	t.Run("no nodes provided", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		r.Nodes = nil

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrClusterNodesNotProvided) {
			t.Fatalf("expected ErrClusterNodesNotProvided, got %v", err)
		}
	})

	t.Run("node name not provided", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		r.Nodes[0].Name = ""

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrClusterNodeNameNotProvided) {
			t.Fatalf("expected ErrClusterNodeNameNotProvided, got %v", err)
		}
	})

	t.Run("node names must be unique within the cluster", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		duplicate := r.Nodes[0]
		duplicate.Host = "db2"
		r.Nodes = append(r.Nodes, duplicate)

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrClusterNodeNameNotUnique) {
			t.Fatalf("expected ErrClusterNodeNameNotUnique, got %v", err)
		}
	})

	t.Run("a node without a keeper port is rejected", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		r.Nodes[0].KeeperPort = nil

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrClusterNodePortsNotProvided) {
			t.Fatalf("expected ErrClusterNodePortsNotProvided, got %v", err)
		}
	})

	t.Run("unknown keeper plugin is rejected", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		r.ClusterOptions.Plugins.Keeper = "unknown-plugin"

		_, err := s.Deploy(r)
		if err == nil {
			t.Fatalf("expected an error for an unknown keeper plugin, got none")
		}
	})

	t.Run("ssh credentials required when no vault and no user/pass", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		r.CommonConfig.SshUser = ""
		r.CommonConfig.SshPass = ""

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrSshCredentialsRequired) {
			t.Fatalf("expected ErrSshCredentialsRequired, got %v", err)
		}
	})

	t.Run("database credentials required for a plugin that needs them", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		r.CommonConfig.DbPass = ""

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrDatabaseCredentialsRequired) {
			t.Fatalf("expected ErrDatabaseCredentialsRequired, got %v", err)
		}
	})

	t.Run("keeper credentials required for a plugin whose keeper endpoint needs them", func(t *testing.T) {
		s := newDeployTestService(t)
		r := etcdDeployRequest()
		r.CommonConfig.KeeperPass = ""

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrKeeperCredentialsRequired) {
			t.Fatalf("expected ErrKeeperCredentialsRequired, got %v", err)
		}
	})

	t.Run("keeper vault and inline credentials cannot both be given", func(t *testing.T) {
		s := newDeployTestService(t)
		r := etcdDeployRequest()
		id := uuid.New()
		r.ClusterOptions.Vaults.KeeperId = &id

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrKeeperCredentialsAmbiguous) {
			t.Fatalf("expected ErrKeeperCredentialsAmbiguous, got %v", err)
		}
	})

	t.Run("locked keeper username cannot be overridden", func(t *testing.T) {
		s := newDeployTestService(t)
		r := etcdDeployRequest()
		r.CommonConfig.KeeperUser = "someone-else"

		_, err := s.Deploy(r)
		if err == nil {
			t.Fatalf("expected an error for an unauthorized username override, got none")
		}
	})

	t.Run("locked database username cannot be overridden", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		r.CommonConfig.DbUser = "someone-else"

		_, err := s.Deploy(r)
		if err == nil {
			t.Fatalf("expected an error for an unauthorized username override, got none")
		}
	})

	// NOTE: a vault and inline credentials are two answers to one question, so
	// the request is refused rather than one of them silently winning
	t.Run("ssh vault and inline credentials cannot both be given", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		id := uuid.New()
		r.ClusterOptions.Vaults.SshKeyId = &id

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrSshCredentialsAmbiguous) {
			t.Fatalf("expected ErrSshCredentialsAmbiguous, got %v", err)
		}
	})

	t.Run("database vault and inline credentials cannot both be given", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		id := uuid.New()
		r.ClusterOptions.Vaults.DatabaseId = &id

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrDatabaseCredentialsAmbiguous) {
			t.Fatalf("expected ErrDatabaseCredentialsAmbiguous, got %v", err)
		}
	})

	t.Run("cluster name already taken", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		if _, err := s.clusterRepository.Create(Request{Name: r.CommonConfig.Cluster}); err != nil {
			t.Fatalf("failed to seed existing cluster: %v", err)
		}

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrClusterNameTaken) {
			t.Fatalf("expected ErrClusterNameTaken, got %v", err)
		}
	})

}

func TestService_deployNode(t *testing.T) {
	s := &Service{}

	t.Run("missing host fails without calling node service", func(t *testing.T) {
		var logs []string
		logsSend := func(ctx string, msg string) { logs = append(logs, ctx+" | "+msg) }

		ok := s.deployNode(Request{}, DeployNode{}, logsSend)
		if ok {
			t.Fatalf("expected deployNode to fail for a missing host")
		}
		if len(logs) != 1 {
			t.Fatalf("expected exactly one log line, got %v", logs)
		}
	})

	t.Run("missing ssh key vault id fails without calling node service", func(t *testing.T) {
		var logs []string
		logsSend := func(ctx string, msg string) { logs = append(logs, ctx+" | "+msg) }

		cluster := Request{Options: Options{Vaults: Vaults{SshKeyId: nil}}}
		ok := s.deployNode(cluster, DeployNode{NodeConfig: NodeConfig{Name: "db-1", Host: "db1"}}, logsSend)
		if ok {
			t.Fatalf("expected deployNode to fail without an ssh key vault id")
		}
		if len(logs) != 1 {
			t.Fatalf("expected exactly one log line, got %v", logs)
		}
	})
}

func TestGetVaultId(t *testing.T) {
	t.Run("nil vault id returns uuid.Nil", func(t *testing.T) {
		if got := getVaultId(nil); got != uuid.Nil {
			t.Errorf("expected uuid.Nil, got %v", got)
		}
	})

	t.Run("set vault id is returned as a value", func(t *testing.T) {
		id := uuid.New()
		if got := getVaultId(&id); got != id {
			t.Errorf("expected %v, got %v", id, got)
		}
	})
}

func TestService_getNodeVaults(t *testing.T) {
	s := &Service{}
	sshKeyId, databaseId := uuid.New(), uuid.New()

	t.Run("each vault is carried through as chosen", func(t *testing.T) {
		keeperId := uuid.New()
		cluster := Request{Options: Options{Vaults: Vaults{KeeperId: &keeperId, DatabaseId: &databaseId, SshKeyId: &sshKeyId}}}
		got := s.getNodeVaults(cluster)
		if got.KeeperId != keeperId || got.DatabaseId != databaseId || got.SshKeyId != sshKeyId {
			t.Errorf("expected every vault to carry through, got %+v", got)
		}
	})

	t.Run("a missing keeper vault is never taken from the database one", func(t *testing.T) {
		cluster := Request{Options: Options{Vaults: Vaults{DatabaseId: &databaseId, SshKeyId: &sshKeyId}}}
		if got := s.getNodeVaults(cluster); got.KeeperId != uuid.Nil {
			t.Errorf("expected no keeper vault, got %v", got.KeeperId)
		}
	})
}

func TestService_resolveLockedUsername(t *testing.T) {
	s := &Service{}

	t.Run("a free username is whatever the user typed", func(t *testing.T) {
		got, err := s.resolveLockedUsername("someone", "")
		if err != nil || got != "someone" {
			t.Errorf("expected the typed username, got %q (%v)", got, err)
		}
	})

	t.Run("an empty username is prefilled with the locked one", func(t *testing.T) {
		got, err := s.resolveLockedUsername("", "root")
		if err != nil || got != "root" {
			t.Errorf("expected the locked username, got %q (%v)", got, err)
		}
	})

	t.Run("overriding a locked username is rejected", func(t *testing.T) {
		if _, err := s.resolveLockedUsername("someone-else", "root"); err == nil {
			t.Error("expected an error for an unauthorized username override, got none")
		}
	})
}

// TestService_mapDeployRequest covers what replaced the deploy plan: the
// cluster no longer computes anything for a node, it just hands node its own
// command and connection details.
func TestService_mapDeployRequest(t *testing.T) {
	s := &Service{}
	sshKeyId := uuid.New()
	cluster := Request{
		Name:    "test-cluster",
		Options: Options{Plugins: Plugins{Keeper: node.KeeperPlugin(keeper.PATRONI_POSTGRES)}, Vaults: Vaults{SshKeyId: &sshKeyId}},
	}

	t.Run("carries the node's own command, ports and defaults", func(t *testing.T) {
		sshPort, keeperPort, dbPort := 2222, 8008, 5432
		got := s.mapDeployRequest(cluster, DeployNode{
			NodeConfig: NodeConfig{Name: "db-1", Host: "db1", SshPort: &sshPort, KeeperPort: &keeperPort, DbPort: &dbPort},
			Command:    "docker run -d spilo",
			PostScript: "echo done",
		})

		if got.Name != "db-1" || got.Cluster != "test-cluster" {
			t.Errorf("expected the node and cluster names to carry through, got %+v", got)
		}
		if got.Connection.Host != "db1" || got.Connection.Port != sshPort {
			t.Errorf("expected the ssh connection to come from the node, got %+v", got.Connection)
		}
		if got.KeeperPort != keeperPort || got.DbPort != dbPort {
			t.Errorf("expected the node's own ports, got %d/%d", got.KeeperPort, got.DbPort)
		}
		if got.Command != "docker run -d spilo" || got.PostScript != "echo done" {
			t.Errorf("expected the node's own command and post script, got %+v", got)
		}
	})

	// NOTE: Deploy rejects an unset port before ever mapping the node, so no
	// port is assumed here - not even the conventional ssh 22
	t.Run("an unset port is carried through as zero rather than assumed", func(t *testing.T) {
		got := s.mapDeployRequest(cluster, DeployNode{NodeConfig: NodeConfig{Name: "db-2", Host: "db2"}})
		if got.Connection.Port != 0 {
			t.Errorf("expected the ssh port to stay zero, got %d", got.Connection.Port)
		}
		if got.KeeperPort != 0 || got.DbPort != 0 {
			t.Errorf("expected unset ports to stay zero, got %d/%d", got.KeeperPort, got.DbPort)
		}
	})
}

func TestService_validateNodePorts(t *testing.T) {
	s := &Service{}
	port := 5432

	t.Run("every port provided is accepted", func(t *testing.T) {
		nodes := []NodeConfig{{Name: "db-1", Host: "db1", SshPort: &port, KeeperPort: &port, DbPort: &port}}
		if err := s.validateNodePorts(nodes); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("a missing port is rejected", func(t *testing.T) {
		zero := 0
		tests := map[string]NodeConfig{
			"keeper port unset":   {Name: "db-1", Host: "db1", SshPort: &port, DbPort: &port},
			"database port unset": {Name: "db-1", Host: "db1", SshPort: &port, KeeperPort: &port},
			"ssh port unset":      {Name: "db-1", Host: "db1", KeeperPort: &port, DbPort: &port},
			"keeper port is zero": {Name: "db-1", Host: "db1", SshPort: &port, KeeperPort: &zero, DbPort: &port},
			"every port unset":    {Name: "db-1", Host: "db1"},
		}
		for name, config := range tests {
			t.Run(name, func(t *testing.T) {
				if err := s.validateNodePorts([]NodeConfig{config}); !errors.Is(err, ErrClusterNodePortsNotProvided) {
					t.Fatalf("expected ErrClusterNodePortsNotProvided, got %v", err)
				}
			})
		}
	})
}
