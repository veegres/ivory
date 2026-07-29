package cluster

import (
	"errors"
	"ivory/clients/storage"
	"ivory/core/utils"
	"ivory/features/node"
	"ivory/features/query"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"ivory/tools"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boltdb/bolt"
)

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

func (f *fakeKeeperAdapter) List(keeper.Request) ([]keeper.Response, int, error) {
	return f.listResponse, f.listStatus, f.listErr
}

// createTestServiceWithNode wires a real node.Service whose keeper registry
// resolves plugin "fake" to adapter, so Overview/Detect/Fix can be exercised
// end to end without a real keeper endpoint.
func createTestServiceWithNode(t *testing.T, adapter keeper.Adapter) *Service {
	t.Helper()
	s := createTestService(t)

	keeperRegistry := utils.NewRegistry[keeper.Plugin, keeper.Adapter]()
	keeperRegistry.Register("fake", adapter)
	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	s.nodeService = node.NewService(nil, keeperRegistry, keeperMetadataRegistry, nil, nil, nil)

	db, errOpen := bolt.Open(filepath.Join(t.TempDir(), "query.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() { db.Close() })
	s.queryService = query.NewService(
		query.NewRepository(storage.NewDbBucket[query.Response](db, "Query"), nil),
		utils.NewRegistry[database.Plugin, database.Adapter](),
		nil, nil, "ivory",
	)
	s.toolRegistry = utils.NewRegistry[tools.Tool, tools.Adapter]()
	return s
}

func TestService_Overview_EndToEnd(t *testing.T) {
	port := 8008
	host1, host2 := "host1", "host2"
	adapter := &fakeKeeperAdapter{
		listResponse: []keeper.Response{
			{Role: keeper.Leader, State: "running", DiscoveredHost: &host1, DiscoveredKeeperPort: &port},
			{Role: "replica", State: "running", DiscoveredHost: &host2, DiscoveredKeeperPort: &port},
		},
	}
	s := createTestServiceWithNode(t, adapter)

	if _, err := s.clusterRepository.Create(Request{
		Name:  "c1",
		Nodes: []NodeConfig{{Host: host1, KeeperPort: &port}, {Host: host2, KeeperPort: &port}},
		Options: Options{
			Plugins: Plugins{Keeper: "fake"},
		},
	}); err != nil {
		t.Fatalf("failed to seed cluster: %v", err)
	}

	overview, err := s.Overview("c1", "", 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(overview.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %v", overview.Nodes)
	}
	for key, n := range overview.Nodes {
		for _, w := range n.Warnings {
			t.Errorf("unexpected warning for %s: %s", key, w)
		}
	}

	t.Run("unknown cluster fails", func(t *testing.T) {
		if _, err := s.Overview("unknown", "", 0); err == nil {
			t.Fatalf("expected an error for an unknown cluster")
		}
	})
}

func TestService_Detect(t *testing.T) {
	port := 8008
	host1 := "host1"
	adapter := &fakeKeeperAdapter{
		listResponse: []keeper.Response{
			{Role: keeper.Leader, State: "running", DiscoveredHost: &host1, DiscoveredKeeperPort: &port},
		},
	}
	s := createTestServiceWithNode(t, adapter)

	detected, err := s.Detect(CreateAutoRequest{
		Name: "detected-cluster",
		Host: host1,
		Port: port,
		Options: Options{
			Plugins: Plugins{Keeper: "fake"},
			Tags:    []string{"PROD"},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(detected.Nodes) != 1 || detected.Nodes[0].Host != host1 {
		t.Fatalf("expected 1 detected node with host %s, got %v", host1, detected.Nodes)
	}
	if len(detected.Tags) != 1 || detected.Tags[0] != "prod" {
		t.Fatalf("expected the tag to be lowercased to 'prod', got %v", detected.Tags)
	}

	got, errGet := s.Get("detected-cluster")
	if errGet != nil {
		t.Fatalf("expected the cluster to be persisted, got %v", errGet)
	}
	if len(got.Nodes) != 1 {
		t.Fatalf("expected the persisted cluster to have 1 node, got %v", got.Nodes)
	}
}

func TestService_Fix(t *testing.T) {
	port := 8008
	host1, host2 := "host1", "host2"
	adapter := &fakeKeeperAdapter{
		listResponse: []keeper.Response{
			{Role: keeper.Leader, State: "running", DiscoveredHost: &host1, DiscoveredKeeperPort: &port},
			{Role: "replica", State: "running", DiscoveredHost: &host2, DiscoveredKeeperPort: &port},
		},
	}
	s := createTestServiceWithNode(t, adapter)

	if _, err := s.clusterRepository.Create(Request{
		Name:  "c1",
		Nodes: []NodeConfig{{Host: host1, KeeperPort: &port}},
		Options: Options{
			Plugins: Plugins{Keeper: "fake"},
		},
	}); err != nil {
		t.Fatalf("failed to seed cluster: %v", err)
	}

	fixed, err := s.Fix("c1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(fixed.Nodes) != 2 {
		t.Fatalf("expected the fixed cluster to report both members from the leader, got %v", fixed.Nodes)
	}

	t.Run("unknown cluster fails", func(t *testing.T) {
		if _, err := s.Fix("unknown"); err == nil {
			t.Fatalf("expected an error for an unknown cluster")
		}
	})

	t.Run("no leader found fails", func(t *testing.T) {
		noLeaderAdapter := &fakeKeeperAdapter{
			listResponse: []keeper.Response{
				{Role: "replica", State: "running", DiscoveredHost: &host1, DiscoveredKeeperPort: &port},
			},
		}
		s2 := createTestServiceWithNode(t, noLeaderAdapter)
		if _, err := s2.clusterRepository.Create(Request{
			Name:    "c2",
			Nodes:   []NodeConfig{{Host: host1, KeeperPort: &port}},
			Options: Options{Plugins: Plugins{Keeper: "fake"}},
		}); err != nil {
			t.Fatalf("failed to seed cluster: %v", err)
		}
		if _, err := s2.Fix("c2"); err == nil {
			t.Fatalf("expected an error when no leader is found")
		}
	})
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

func TestService_mergeKeeperNode_ResolvesUnknownPortByHost(t *testing.T) {
	s := &Service{}
	port1, port2, port3 := 5001, 5002, 5003
	configs := []NodeConfig{
		{Host: "patroni1", KeeperPort: &port1, DbPort: &port1},
		{Host: "patroni2", KeeperPort: &port2, DbPort: &port2},
		{Host: "patroni3", KeeperPort: &port3, DbPort: &port3},
	}

	t.Run("standby response without a discovered port merges into its configured node instead of a phantom one", func(t *testing.T) {
		host1 := "patroni1"
		nodeMap := s.getConfiguredNodeMap(configs, nil, nil)
		// NOTE: mirrors postgres.mapSyncStandby - Sync reported from the
		// primary's pg_stat_replication, with the host known (application_name)
		// but the port deliberately left unknown.
		s.mergeKeeperNode(nodeMap, node.KeeperOneResponse{Role: keeper.Replica, Sync: true, DiscoveredHost: &host1})

		if len(nodeMap) != 3 {
			t.Fatalf("expected exactly 3 nodes (no phantom entry), got %d: %v", len(nodeMap), nodeMap)
		}
		merged, ok := nodeMap["patroni1:5001"]
		if !ok {
			t.Fatalf("expected patroni1 to be resolved to its configured node, got %v", nodeMap)
		}
		if !merged.Keeper.Sync {
			t.Fatalf("expected the resolved response's Sync=true to be applied")
		}
		if len(merged.Warnings) != 0 {
			t.Fatalf("expected no warnings, got %v", merged.Warnings)
		}
	})

	t.Run("ambiguous host (declared more than once) falls back to a phantom entry instead of guessing", func(t *testing.T) {
		duplicateConfigs := append(configs, NodeConfig{Host: "patroni1", KeeperPort: &port2, DbPort: &port2})
		nodeMap := s.getConfiguredNodeMap(duplicateConfigs, nil, nil)
		host1 := "patroni1"
		s.mergeKeeperNode(nodeMap, node.KeeperOneResponse{Role: keeper.Replica, DiscoveredHost: &host1})

		if _, ok := nodeMap["patroni1"]; !ok {
			t.Fatalf("expected a fallback host-only entry when the host is ambiguous, got %v", nodeMap)
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

func TestService_addOverviewWarnings_DbPortMismatch(t *testing.T) {
	s := &Service{}
	configuredPort := 5002
	discoveredPort := 5001
	nodes := map[string]Node{
		"patroni1:5001": {
			Config: NodeConfig{DbPort: &discoveredPort},
			Keeper: node.KeeperOneResponse{Role: keeper.Leader, DiscoveredDbPort: &discoveredPort},
		},
		"patroni2:5001": {
			Config: NodeConfig{DbPort: &configuredPort},
			Keeper: node.KeeperOneResponse{Role: keeper.Replica, DiscoveredDbPort: &discoveredPort},
		},
	}

	s.addOverviewWarnings(nodes)

	if len(nodes["patroni1:5001"].Warnings) != 0 {
		t.Fatalf("Expected patroni1 without warnings, got %v", nodes["patroni1:5001"].Warnings)
	}
	if len(nodes["patroni2:5001"].Warnings) != 1 || nodes["patroni2:5001"].Warnings[0] != "database port in keeper response and cluster configuration mismatch" {
		t.Fatalf("Expected patroni2 to have a db port mismatch warning, got %v", nodes["patroni2:5001"].Warnings)
	}
}

func TestService_addOverviewWarnings_NoLeaderFound(t *testing.T) {
	s := &Service{}
	nodes := map[string]Node{
		"db1:8008": {Keeper: node.KeeperOneResponse{Role: keeper.Replica}},
		"db2:8008": {Keeper: node.KeeperOneResponse{Role: keeper.Replica}},
	}

	s.addOverviewWarnings(nodes)

	for key, n := range nodes {
		if len(n.Warnings) != 1 || n.Warnings[0] != "no leader node was found in Keeper response" {
			t.Fatalf("Expected %s to have a no-leader warning, got %v", key, n.Warnings)
		}
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
