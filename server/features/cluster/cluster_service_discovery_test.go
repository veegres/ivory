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
