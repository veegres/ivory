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
// validation path and planDeploy without ever touching the network: no
// vaultService is wired in, so tests must stay on the side of Deploy that
// returns before any vault or SSH work happens.
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
		Nodes: []DeployNode{{NodeConfig: NodeConfig{Host: "db1"}}},
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

	t.Run("unknown keeper plugin fails planning", func(t *testing.T) {
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

		ok := s.deployNode(Request{}, nil, &node.KeeperDeployPlanResponse{}, node.KeeperDeployPlanNode{}, logsSend)
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
		ok := s.deployNode(cluster, nil, &node.KeeperDeployPlanResponse{}, node.KeeperDeployPlanNode{Host: "db1"}, logsSend)
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

func TestService_planDeploy(t *testing.T) {
	s := newDeployTestService(t)

	sshPort := 2222
	nodes := []DeployNode{
		{NodeConfig: NodeConfig{Host: "db1", SshPort: &sshPort}},
		{NodeConfig: NodeConfig{Host: "db2"}},
	}
	options := Options{Plugins: Plugins{Keeper: node.KeeperPlugin(keeper.PATRONI_POSTGRES)}}

	plan, err := s.planDeploy("test-cluster", options, false, "", nil, nodes)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if plan.Image == "" {
		t.Fatalf("expected the plugin's default image to be used")
	}
	if len(plan.Nodes) != 2 {
		t.Fatalf("expected 2 planned nodes, got %d", len(plan.Nodes))
	}
	if plan.Nodes[0].Host != "db1" || plan.Nodes[0].SshPort != sshPort {
		t.Errorf("expected the first node's host/ssh port to carry through, got %+v", plan.Nodes[0])
	}
	if plan.Nodes[1].Host != "db2" || plan.Nodes[1].SshPort != 22 {
		t.Errorf("expected the second node to default to ssh port 22, got %+v", plan.Nodes[1])
	}
}
