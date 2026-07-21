package cluster

import (
	"errors"
	"fmt"
	"ivory/core/utils"
	"ivory/features/node"
	"ivory/plugins/keeper"
	"net/http"
	"strings"
	"testing"
)

func TestService_getNodeKey(t *testing.T) {
	s := &Service{}
	host := "localhost"
	port := 8008

	t.Run("should return host and port when port is not nil", func(t *testing.T) {
		key := s.getNodeKey(host, &port)
		expected := fmt.Sprintf("%s:%d", host, port)
		if key != expected {
			t.Errorf("Expected key '%s', got '%s'", expected, key)
		}
	})

	t.Run("should return only host when port is nil", func(t *testing.T) {
		key := s.getNodeKey(host, nil)
		if key != host {
			t.Errorf("Expected key '%s', got '%s'", host, key)
		}
	})
}

func TestService_Overview_Mapping(t *testing.T) {
	s := &Service{}

	t.Run("should correctly map keeper responses to cluster config", func(t *testing.T) {
		port8008 := 8008
		port5432 := 5432

		clusterNodes := []NodeConfig{
			{Host: "db1", KeeperPort: &port8008, DbPort: &port5432},
		}

		host := "db1"
		keeperNodes := map[string]node.KeeperOneResponse{
			"db1:8008": {
				Role:                 keeper.Leader,
				State:                "running",
				DiscoveredHost:       &host,
				DiscoveredKeeperPort: &port8008,
				DiscoveredDbPort:     &port5432,
			},
		}

		resultNodeMap := s.buildOverviewNodes(clusterNodes, keeperNodes, nil, nil)

		if len(resultNodeMap) != 1 {
			t.Fatalf("Expected 1 mapped node, got %d", len(resultNodeMap))
		}

		mappedNode := resultNodeMap["db1:8008"]
		if mappedNode.Config.Host != "db1" {
			t.Errorf("Expected config host 'db1', got '%s'", mappedNode.Config.Host)
		}
		if mappedNode.Keeper.Role != keeper.Leader {
			t.Errorf("Expected keeper role 'leader', got '%s'", mappedNode.Keeper.Role)
		}
	})
}

// fakeKeeperAdapter lets tests control exactly what List() reports, in
// particular returning a usable Response alongside a non-nil error the way
// a real adapter does when it can still describe a node's state despite a
// connection problem (e.g. postgres starting up).
type fakeKeeperAdapter struct {
	keeper.Adapter
	listResponse []keeper.Response
	listStatus   int
	listErr      error
}

func (f *fakeKeeperAdapter) List(request keeper.Request) ([]keeper.Response, int, error) {
	return f.listResponse, f.listStatus, f.listErr
}

func TestService_getKeeperListByManyAll_KeepsResponseAlongsideError(t *testing.T) {
	host := "db1"
	port := 8008
	errMessage := "the database system is starting up"

	keeperRegistry := utils.NewRegistry[keeper.Plugin, keeper.Adapter]()
	keeperRegistry.Register("fake", &fakeKeeperAdapter{
		listResponse: []keeper.Response{{
			State:                keeper.StateStarting,
			Role:                 keeper.Unknown,
			DiscoveredHost:       &host,
			DiscoveredKeeperPort: &port,
		}},
		listStatus: http.StatusServiceUnavailable,
		listErr:    errors.New(errMessage),
	})
	s := &Service{nodeService: node.NewService(nil, keeperRegistry, nil, nil, nil, nil)}

	configs := []NodeConfig{{Host: host, KeeperPort: &port}}
	keeperNodes, connectionErrors, err := s.getKeeperListByManyAll(configs, Options{Plugins: Plugins{Keeper: "fake"}})
	nodes := s.buildOverviewNodes(configs, keeperNodes, connectionErrors, err)

	nodeKey := "db1:8008"
	if err == nil || !strings.Contains(err.Error(), errMessage) {
		t.Fatalf("expected returned error to contain %q, got %v", errMessage, err)
	}
	if nodes[nodeKey].Keeper.State != keeper.StateStarting {
		t.Fatalf("expected the degraded response to still be merged in, got state %q", nodes[nodeKey].Keeper.State)
	}
	warnings := nodes[nodeKey].Warnings
	found := false
	for _, w := range warnings {
		if strings.Contains(w, errMessage) {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected the handled error to still be surfaced as a warning, got %v", warnings)
	}
}

// multiHostFakeKeeperAdapter lets tests simulate an adapter (like native
// postgres) whose List() response differs depending on which node is being
// queried - in particular, a primary's response can report extra peers with
// data a peer can't determine about itself, such as Sync status.
type multiHostFakeKeeperAdapter struct {
	keeper.Adapter
	responses map[string][]keeper.Response
}

func (f *multiHostFakeKeeperAdapter) List(request keeper.Request) ([]keeper.Response, int, error) {
	return f.responses[request.Host], http.StatusOK, nil
}

func TestService_getKeeperListByManyAll_PrefersLeaderReportedSyncData(t *testing.T) {
	port := 5432
	db1, db2, db3 := "db1", "db2", "db3"

	responses := map[string][]keeper.Response{
		db1: {
			{Role: keeper.Leader, State: keeper.StateRunning, DiscoveredHost: &db1, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port},
			{Role: keeper.Replica, Sync: true, State: keeper.StateRunning, DiscoveredHost: &db2, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port},
			{Role: keeper.Replica, Sync: false, State: keeper.StateRunning, DiscoveredHost: &db3, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port},
		},
		// NOTE: each replica's own connection can't determine its sync
		// status by querying itself, so it always self-reports Sync: false.
		db2: {
			{Role: keeper.Replica, Sync: false, State: keeper.StateRunning, DiscoveredHost: &db2, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port},
		},
		db3: {
			{Role: keeper.Replica, Sync: false, State: keeper.StateRunning, DiscoveredHost: &db3, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port},
		},
	}

	keeperRegistry := utils.NewRegistry[keeper.Plugin, keeper.Adapter]()
	keeperRegistry.Register("fake", &multiHostFakeKeeperAdapter{responses: responses})
	s := &Service{nodeService: node.NewService(nil, keeperRegistry, nil, nil, nil, nil)}

	// NOTE: replicas are listed before the leader on purpose - responses
	// preserve config order (node_service_keeper.go's KeeperNodeListMulti
	// writes each goroutine's result by index), so this would make db2's
	// own plain self-report win by "first response wins" if the leader's
	// richer response weren't deliberately preferred regardless of order.
	configs := []NodeConfig{
		{Host: db2, KeeperPort: &port, DbPort: &port},
		{Host: db3, KeeperPort: &port, DbPort: &port},
		{Host: db1, KeeperPort: &port, DbPort: &port},
	}

	keeperNodes, connectionErrors, err := s.getKeeperListByManyAll(configs, Options{Plugins: Plugins{Keeper: "fake"}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(connectionErrors) != 0 {
		t.Fatalf("expected no connection errors, got %v", connectionErrors)
	}
	if len(keeperNodes) != 3 {
		t.Fatalf("expected exactly 3 nodes, got %d: %v", len(keeperNodes), keeperNodes)
	}

	db2Node, ok := keeperNodes["db2:5432"]
	if !ok {
		t.Fatalf("expected db2 to be resolved to its configured node, got %v", keeperNodes)
	}
	if !db2Node.Sync {
		t.Fatalf("expected the leader-reported sync=true to win over db2's own self-report")
	}

	db3Node, ok := keeperNodes["db3:5432"]
	if !ok {
		t.Fatalf("expected db3 to be resolved to its configured node, got %v", keeperNodes)
	}
	if db3Node.Sync {
		t.Fatalf("expected db3 to stay async")
	}
}

func TestService_getKeeperListByLeader(t *testing.T) {
	port := 5432
	db1, db2, db3 := "db1", "db2", "db3"

	t.Run("rebuilds the node list from the leader's response, not whichever node answers first", func(t *testing.T) {
		responses := map[string][]keeper.Response{
			// NOTE: db2 is listed first in configs below, so a "first
			// success wins" strategy would incorrectly collapse the
			// cluster down to db2's own 1-entry self-report.
			db2: {
				{Role: keeper.Replica, State: keeper.StateRunning, DiscoveredHost: &db2, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port},
			},
			db1: {
				{Role: keeper.Leader, State: keeper.StateRunning, DiscoveredHost: &db1, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port},
				{Role: keeper.Replica, Sync: true, State: keeper.StateRunning, DiscoveredHost: &db2, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port},
				{Role: keeper.Replica, State: keeper.StateRunning, DiscoveredHost: &db3, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port},
			},
			db3: {
				{Role: keeper.Replica, State: keeper.StateRunning, DiscoveredHost: &db3, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port},
			},
		}

		keeperRegistry := utils.NewRegistry[keeper.Plugin, keeper.Adapter]()
		keeperRegistry.Register("fake", &multiHostFakeKeeperAdapter{responses: responses})
		s := &Service{nodeService: node.NewService(nil, keeperRegistry, nil, nil, nil, nil)}

		configs := []NodeConfig{
			{Host: db2, KeeperPort: &port, DbPort: &port},
			{Host: db1, KeeperPort: &port, DbPort: &port},
			{Host: db3, KeeperPort: &port, DbPort: &port},
		}

		keeperNodes, err := s.getKeeperListByLeader(configs, Options{Plugins: Plugins{Keeper: "fake"}})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(keeperNodes) != 3 {
			t.Fatalf("expected the leader's full 3-node view, got %d: %v", len(keeperNodes), keeperNodes)
		}
		for _, n := range keeperNodes {
			if n.DiscoveredHost != nil && *n.DiscoveredHost == db2 && (n.DiscoveredKeeperPort == nil || *n.DiscoveredKeeperPort != port) {
				t.Fatalf("expected db2's port to already be known from the leader's own response, got %v", n.DiscoveredKeeperPort)
			}
		}
	})

	t.Run("errors when no configured node reports being the leader", func(t *testing.T) {
		responses := map[string][]keeper.Response{
			db1: {{Role: keeper.Replica, State: keeper.StateRunning, DiscoveredHost: &db1, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port}},
			db2: {{Role: keeper.Replica, State: keeper.StateRunning, DiscoveredHost: &db2, DiscoveredKeeperPort: &port, DiscoveredDbPort: &port}},
		}

		keeperRegistry := utils.NewRegistry[keeper.Plugin, keeper.Adapter]()
		keeperRegistry.Register("fake", &multiHostFakeKeeperAdapter{responses: responses})
		s := &Service{nodeService: node.NewService(nil, keeperRegistry, nil, nil, nil, nil)}

		configs := []NodeConfig{
			{Host: db1, KeeperPort: &port, DbPort: &port},
			{Host: db2, KeeperPort: &port, DbPort: &port},
		}

		keeperNodes, err := s.getKeeperListByLeader(configs, Options{Plugins: Plugins{Keeper: "fake"}})
		if !errors.Is(err, ErrNoLeaderFound) {
			t.Fatalf("expected ErrNoLeaderFound, got %v", err)
		}
		if keeperNodes != nil {
			t.Fatalf("expected no node list on error, got %v", keeperNodes)
		}
	})
}

func TestService_buildOverviewNodes(t *testing.T) {
	s := &Service{}
	host1 := "db1"
	host2 := "db2"
	host3 := "db3"
	port := 8008
	nodes := s.buildOverviewNodes(nil, map[string]node.KeeperOneResponse{
		"db1:8008": {Role: keeper.Leader, DiscoveredHost: &host1, DiscoveredKeeperPort: &port},
		"db2:8008": {Role: keeper.Leader, DiscoveredHost: &host2, DiscoveredKeeperPort: &port},
		"db3:8008": {Role: keeper.Replica, DiscoveredHost: &host3, DiscoveredKeeperPort: &port},
	}, nil, nil)

	t.Run("should add multi leader warnings", func(t *testing.T) {
		if len(nodes["db1:8008"].Warnings) != 2 {
			t.Fatalf("Expected db1 leader warning, got %v", nodes["db1:8008"].Warnings)
		}
		if len(nodes["db2:8008"].Warnings) != 2 {
			t.Fatalf("Expected db2 leader warning, got %v", nodes["db2:8008"].Warnings)
		}
		if len(nodes["db3:8008"].Warnings) != 1 {
			t.Fatalf("Expected db3 without leader warning, got %v", nodes["db3:8008"].Warnings)
		}
	})

	t.Run("should add unconfigured node warnings", func(t *testing.T) {
		if !strings.Contains(nodes["db1:8008"].Warnings[0], "not in the cluster configuration") {
			t.Fatalf("Expected unconfigured node warning, got %v", nodes["db1:8008"].Warnings)
		}
	})
}

func TestService_addOverviewWarnings(t *testing.T) {
	s := &Service{}
	nodes := map[string]Node{
		"db1:8008": {Keeper: node.KeeperOneResponse{Role: keeper.Leader}},
		"db2:8008": {Keeper: node.KeeperOneResponse{Role: keeper.Leader}},
		"db3:8008": {Keeper: node.KeeperOneResponse{Role: keeper.Replica}},
	}

	s.addOverviewWarnings(nodes)

	if len(nodes["db1:8008"].Warnings) != 1 {
		t.Fatalf("Expected db1 leader warning, got %v", nodes["db1:8008"].Warnings)
	}
	if len(nodes["db2:8008"].Warnings) != 1 {
		t.Fatalf("Expected db2 leader warning, got %v", nodes["db2:8008"].Warnings)
	}
	if len(nodes["db3:8008"].Warnings) != 0 {
		t.Fatalf("Expected db3 without leader warning, got %v", nodes["db3:8008"].Warnings)
	}
}

func TestService_getKeeperListAutoMerge(t *testing.T) {
	s := &Service{}

	t.Run("should return concrete errors when no configured nodes can be requested", func(t *testing.T) {
		configs := []NodeConfig{{Host: "db1"}}
		keeperNodes, connectionErrors, err := s.getKeeperListByManyAll(configs, Options{})
		nodes := s.buildOverviewNodes(configs, keeperNodes, connectionErrors, err)

		if len(nodes) != 1 {
			t.Fatalf("Expected one configured node, got %v", nodes)
		}
		if err == nil || !strings.Contains(err.Error(), "no configured keeper connections can be requested") {
			t.Fatalf("Expected no keeper connections error, got %v", err)
		}
		if len(nodes["db1"].Warnings) == 0 || !strings.Contains(nodes["db1"].Warnings[0], "missing a keeper port") {
			t.Fatalf("Expected concrete db1 warning, got %v", nodes["db1"].Warnings)
		}
	})
}
