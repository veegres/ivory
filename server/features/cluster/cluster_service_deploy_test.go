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
	"ivory/plugins/platform/linux"
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

	platformRegistry := utils.NewRegistry[platform.Plugin, platform.Adapter]()
	platformRegistry.Register(platform.Linux, linux.NewAdapter(ssh.NewClient()))
	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	keeperMetadataRegistry.Register(keeper.PATRONI_POSTGRES, patroni.NewAdapter(nil))
	keeperMetadataRegistry.Register(keeper.NATIVE_POSTGRES, postgres.NewAdapter())
	keeperMetadataRegistry.Register(keeper.NATIVE_ETCD, etcd.NewAdapter())
	nodeService := node.NewService(platformRegistry, nil, keeperMetadataRegistry, nil, nil, nil)

	return &Service{
		clusterRepository: clusterRepository,
		tagService:        tagService,
		nodeService:       nodeService,
	}
}

func validDeployRequest() DeployRequest {
	return DeployRequest{
		Nodes: []DeployNode{{
			NodeConfig: NodeConfig{Name: "db-1", Host: "db1"},
			Command:    `docker run -d --name {{name}} -e ETCD3_HOSTS="etcd-1:2379" spilo`,
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
		r.Nodes = append(r.Nodes, DeployNode{
			NodeConfig: NodeConfig{Name: "db-1", Host: "db2"},
			Command:    "docker run -d spilo",
		})

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrClusterNodeNameNotUnique) {
			t.Fatalf("expected ErrClusterNodeNameNotUnique, got %v", err)
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

func TestService_getDatabaseId(t *testing.T) {
	s := &Service{}

	t.Run("nil database vault id returns uuid.Nil", func(t *testing.T) {
		if got := s.getDatabaseId(Request{}); got != uuid.Nil {
			t.Errorf("expected uuid.Nil, got %v", got)
		}
	})

	t.Run("set database vault id is returned as a value", func(t *testing.T) {
		id := uuid.New()
		cluster := Request{Options: Options{Vaults: Vaults{DatabaseId: &id}}}
		if got := s.getDatabaseId(cluster); got != id {
			t.Errorf("expected %v, got %v", id, got)
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

	t.Run("an unset ssh port defaults to 22", func(t *testing.T) {
		got := s.mapDeployRequest(cluster, DeployNode{NodeConfig: NodeConfig{Name: "db-2", Host: "db2"}})
		if got.Connection.Port != 22 {
			t.Errorf("expected ssh port 22, got %d", got.Connection.Port)
		}
		if got.KeeperPort != 0 || got.DbPort != 0 {
			t.Errorf("expected unset ports to stay zero, got %d/%d", got.KeeperPort, got.DbPort)
		}
	})
}
