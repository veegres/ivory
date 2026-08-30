package cluster

import (
	"errors"
	"fmt"
	"ivory/clients/console/ssh"
	"ivory/clients/storage"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/features/node"
	"ivory/features/tag"
	"ivory/plugins/keeper"
	"ivory/plugins/keeper/etcd"
	"ivory/plugins/keeper/patroni"
	"ivory/plugins/keeper/postgres"
	"ivory/plugins/platform"
	"ivory/plugins/platform/docker"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

// newDeployTestService builds a cluster Service with a real repository and
// real keeper/platform plugin registries, enough to exercise Deploy's
// validation path without ever touching the network: no vaultService is wired
// in, so tests must stay on the side of Deploy that returns before any vault
// or SSH work happens. newDeployTestServiceWithVault adds a real vault for the
// tests that need to run past that point.
func newDeployTestService(t *testing.T) *Service {
	t.Helper()
	s, _ := newDeployTestServiceWithVault(t, false)
	return s
}

func newDeployTestServiceWithVault(t *testing.T, withVault bool) (*Service, *vault.Service) {
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
	platformRegistry.Register(platform.Docker, docker.NewPlugin(ssh.NewClient()))
	keeperRegistry := utils.NewRegistry[keeper.PluginType, keeper.Plugin]()
	keeperRegistry.Register(keeper.PATRONI_POSTGRES, patroni.NewPlugin(nil))
	keeperRegistry.Register(keeper.NATIVE_POSTGRES, postgres.NewPlugin())
	keeperRegistry.Register(keeper.NATIVE_ETCD, etcd.NewPlugin())
	var vaultService *vault.Service
	if withVault {
		oldWd, _ := os.Getwd()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}
		t.Cleanup(func() { os.Chdir(oldWd) })

		secretService := secret.NewService(
			secret.NewRepository(storage.NewDbBucket[string](db, "Secret")),
			encryption.NewService(),
		)
		if err := secretService.SetDefault(); err != nil {
			t.Fatalf("failed to set default secret: %v", err)
		}
		vaultService = vault.NewService(
			vault.NewRepository(storage.NewDbBucket[vault.Vault](db, "Vault")),
			ssh.NewClient(),
			secretService,
			encryption.NewService(),
		)
	}

	nodeService := node.NewService(platformRegistry, keeperRegistry, vaultService, nil, nil)

	return &Service{
		clusterRepository: clusterRepository,
		tagService:        tagService,
		nodeService:       nodeService,
		vaultService:      vaultService,
	}, vaultService
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
			DbUser:  "postgres",
			DbPass:  "secret",
		},
		ClusterOptions: Options{Plugins: Plugins{Keeper: node.KeeperPlugin(keeper.PATRONI_POSTGRES)}},
	}
}

// etcdDeployRequest answers the keeper credential pair as well as the database
// one, which etcd's shipped deployment asks for and patroni's does not.
func etcdDeployRequest() DeployRequest {
	r := validDeployRequest()
	r.ClusterOptions.Plugins.Keeper = node.KeeperPlugin(keeper.NATIVE_ETCD)
	r.CommonConfig.KeeperUser = "root"
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

	t.Run("a database username without a password is half an answer", func(t *testing.T) {
		s := newDeployTestService(t)
		r := validDeployRequest()
		r.CommonConfig.DbPass = ""

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrDatabaseCredentialsIncomplete) {
			t.Fatalf("expected ErrDatabaseCredentialsIncomplete, got %v", err)
		}
	})

	t.Run("a keeper password without a username is half an answer", func(t *testing.T) {
		s := newDeployTestService(t)
		r := etcdDeployRequest()
		r.CommonConfig.KeeperUser = ""

		_, err := s.Deploy(r)
		if !errors.Is(err, ErrKeeperCredentialsIncomplete) {
			t.Fatalf("expected ErrKeeperCredentialsIncomplete, got %v", err)
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

	// NOTE: no username is locked any more - a template names a default the
	// deploy screen opens on, and what the user types instead is their answer.
	// The pair is still refused alongside a vault, which is two answers to one
	// question rather than an unexpected one.
	t.Run("usernames the engine did not suggest are stored as typed", func(t *testing.T) {
		s, vaultService := newDeployTestServiceWithVault(t, true)
		r := etcdDeployRequest()
		r.Nodes[0].Host = ""
		r.CommonConfig.KeeperUser = "someone-else"
		r.CommonConfig.DbUser = "someone-else-too"

		if _, err := s.Deploy(r); err != nil {
			t.Fatalf("expected the deploy to run past validation, got %v", err)
		}
		stored, err := s.Get(r.CommonConfig.Cluster)
		if err != nil {
			t.Fatalf("expected the cluster to be registered, got %v", err)
		}
		for _, want := range []struct {
			label    string
			vaultId  *uuid.UUID
			username string
		}{
			{label: "keeper", vaultId: stored.Vaults.KeeperId, username: "someone-else"},
			{label: "database", vaultId: stored.Vaults.DatabaseId, username: "someone-else-too"},
		} {
			if want.vaultId == nil {
				t.Fatalf("expected a %s vault entry for the pair that was typed", want.label)
			}
			credentials, err := vaultService.Get(*want.vaultId)
			if err != nil {
				t.Fatalf("failed to read the %s vault: %v", want.label, err)
			}
			if credentials.Username != want.username {
				t.Errorf("expected the typed %s username %q, got %q", want.label, want.username, credentials.Username)
			}
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

// TestService_isHalfCredential covers what replaced resolveLockedUsername: no
// username is required or locked any more, but half a pair would be written to
// a vault entry that authenticates nothing.
func TestService_isHalfCredential(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		expected bool
	}{
		{name: "a whole pair is an answer", username: "root", password: "secret", expected: false},
		{name: "neither half is an answer too", username: "", password: "", expected: false},
		{name: "a username alone is half", username: "root", password: "", expected: true},
		{name: "a password alone is half", username: "", password: "secret", expected: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHalfCredential(tt.username, tt.password); got != tt.expected {
				t.Errorf("isHalfCredential(%q, %q) = %v, want %v", tt.username, tt.password, got, tt.expected)
			}
		})
	}
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
			NodeConfig:  NodeConfig{Name: "db-1", Host: "db1", SshPort: &sshPort, KeeperPort: &keeperPort, DbPort: &dbPort},
			Command:     "docker run -d spilo",
			PostScripts: []string{"echo done"},
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
		if got.Command != "docker run -d spilo" || len(got.PostScripts) != 1 || got.PostScripts[0] != "echo done" {
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

// TestService_Deploy_RegistersClusterBeforeDeploying covers the order a deploy
// runs in and why: the cluster and its vaults are written first, so a
// deployment that never came up is still reachable to be opened and read. The
// nodes here carry no host, so every one of them fails inside deployNode
// before any connection is attempted - and the cluster must survive that
// intact rather than be rolled back out from under the containers.
func TestService_Deploy_RegistersClusterBeforeDeploying(t *testing.T) {
	s, vaultService := newDeployTestServiceWithVault(t, true)

	r := etcdDeployRequest()
	r.Nodes[0].Host = ""

	logs, err := s.Deploy(r)
	if err != nil {
		t.Fatalf("expected the failure to be reported through the logs, got %v", err)
	}
	if len(logs) == 0 {
		t.Fatal("expected the operator to get the logs of the attempt")
	}

	stored, err := s.Get(r.CommonConfig.Cluster)
	if err != nil {
		t.Fatalf("the cluster must be registered even when no node deployed: %v", err)
	}
	if len(stored.Nodes) != len(r.Nodes) {
		t.Errorf("expected every configured node to be registered, got %d of %d", len(stored.Nodes), len(r.Nodes))
	}
	if stored.Vaults.SshKeyId == nil {
		t.Error("the cluster must keep the ssh vault it was deployed with")
	}

	// NOTE: the vaults are the cluster's configuration now - removing them
	// would leave the failed deployments unreachable
	vaults, err := vaultService.Map(nil)
	if err != nil {
		t.Fatalf("failed to list vaults: %v", err)
	}
	if len(vaults) == 0 {
		t.Error("expected the vaults created for this deploy to be kept")
	}
}

// TestService_sshTargets covers who the generated key has to be installed on:
// each machine once. Several nodes of a single-host cluster share a VM, so
// visiting them per node would copy the same key three times over.
func TestService_sshTargets(t *testing.T) {
	s := &Service{}
	port := func(p int) *int { return &p }

	tests := []struct {
		name     string
		nodes    []DeployNode
		expected []string
	}{
		{
			name: "one target per host",
			nodes: []DeployNode{
				{NodeConfig: NodeConfig{Host: "10.0.0.1", SshPort: port(22)}},
				{NodeConfig: NodeConfig{Host: "10.0.0.2", SshPort: port(22)}},
			},
			expected: []string{"10.0.0.1:22", "10.0.0.2:22"},
		},
		{
			name: "three nodes on one VM are one target",
			nodes: []DeployNode{
				{NodeConfig: NodeConfig{Host: "10.0.0.1", SshPort: port(22)}},
				{NodeConfig: NodeConfig{Host: "10.0.0.1", SshPort: port(22)}},
				{NodeConfig: NodeConfig{Host: "10.0.0.1", SshPort: port(22)}},
			},
			expected: []string{"10.0.0.1:22"},
		},
		{
			name: "one host on two ssh ports is two targets",
			nodes: []DeployNode{
				{NodeConfig: NodeConfig{Host: "10.0.0.1", SshPort: port(22)}},
				{NodeConfig: NodeConfig{Host: "10.0.0.1", SshPort: port(2222)}},
			},
			expected: []string{"10.0.0.1:22", "10.0.0.1:2222"},
		},
		{
			name: "a node with no host is left to deployNode to report",
			nodes: []DeployNode{
				{NodeConfig: NodeConfig{Host: "", SshPort: port(22)}},
				{NodeConfig: NodeConfig{Host: "10.0.0.1", SshPort: port(22)}},
			},
			expected: []string{"10.0.0.1:22"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			targets := s.sshTargets(test.nodes, "root", "secret")
			if len(targets) != len(test.expected) {
				t.Fatalf("expected %d target(s), got %d", len(test.expected), len(targets))
			}
			for i, want := range test.expected {
				got := fmt.Sprintf("%s:%d", targets[i].Host, targets[i].Port)
				if got != want {
					t.Errorf("target %d: expected %q, got %q", i, want, got)
				}
				if targets[i].Username != "root" || targets[i].Password != "secret" {
					t.Errorf("target %d must carry the typed ssh credentials, the only thing they are used for", i)
				}
			}
		})
	}
}
