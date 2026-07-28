package cluster

import (
	"ivory/clients/storage"
	"ivory/core/config"
	"ivory/core/utils"
	"ivory/features/node"
	"ivory/features/query"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"ivory/tools"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

// fakeKeeperMetadata is a minimal keeper.Metadata test double: it only needs
// to answer SupportedFeatures for these tests, so DeploymentSpec is unused.
type fakeKeeperMetadata struct {
	features map[env.Feature]bool
}

func (f fakeKeeperMetadata) SupportedFeatures() map[env.Feature]bool { return f.features }
func (f fakeKeeperMetadata) DeploymentSpec() keeper.DeploymentSpec   { return keeper.DeploymentSpec{} }

// fakeDbAdapter is a minimal database.Adapter test double. Embedding the nil
// interface satisfies the (large) query-execution surface without stubbing
// every method by hand; only the two methods query.Service's constructor and
// SupportedFeatures actually invoke are overridden.
type fakeDbAdapter struct {
	database.Adapter
	features map[env.Feature]bool
}

func (f fakeDbAdapter) SupportedFeatures() map[env.Feature]bool           { return f.features }
func (f fakeDbAdapter) SystemRequests() []database.SystemRequest          { return nil }
func (f fakeDbAdapter) SystemCharts() map[database.SystemChartType]string { return nil }

// fakeToolAdapter is a minimal tools.Adapter test double.
type fakeToolAdapter struct {
	features map[env.Feature]bool
}

func (f fakeToolAdapter) SupportedFeatures(env.Plugin) map[env.Feature]bool { return f.features }
func (f fakeToolAdapter) DeleteAll() error                                  { return nil }

// createFeatureTestService wires a cluster Service whose node/query/tool
// dependencies each expose exactly one fake plugin, so getSupportedFeatures
// and getToolSupportedFeatures can be exercised without pulling in any real
// keeper/database/tool plugin implementation.
func createFeatureTestService(t *testing.T) (*Service, node.KeeperPlugin, query.DbPlugin) {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "cluster-service-features-test-*")
	if errDir != nil {
		t.Fatalf("failed to create temp dir: %v", errDir)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	db, errOpen := bolt.Open(filepath.Join(tmpDir, "test.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() {
		db.Close()
	})

	keeperPlugin := node.KeeperPlugin("fake-keeper")
	dbPlugin := query.DbPlugin("fake-db")

	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	keeperMetadataRegistry.Register(keeper.Plugin(keeperPlugin), fakeKeeperMetadata{
		features: map[env.Feature]bool{env.ViewNodeKeeperOverview: true},
	})
	nodeService := node.NewService(nil, nil, keeperMetadataRegistry, nil, nil, nil)

	databaseRegistry := utils.NewRegistry[database.Plugin, database.Adapter]()
	databaseRegistry.Register(database.Plugin(dbPlugin), fakeDbAdapter{
		features: map[env.Feature]bool{env.ViewQueryDbInfo: true},
	})
	queryRepository := query.NewRepository(
		storage.NewDbBucket[query.Response](db, "Query"),
		storage.NewFileStorage("cluster-service-test-query-logs", ".log"),
	)
	queryService := query.NewService(queryRepository, databaseRegistry, nil, nil, "ivory")

	toolRegistry := utils.NewRegistry[tools.Tool, tools.Adapter]()
	toolRegistry.Register(tools.PgCompactTable, fakeToolAdapter{
		features: map[env.Feature]bool{env.ViewToolPgCompactTableList: true},
	})

	return &Service{
		nodeService:  nodeService,
		queryService: queryService,
		toolRegistry: toolRegistry,
	}, keeperPlugin, dbPlugin
}

func TestService_getSupportedFeatures(t *testing.T) {
	s, keeperPlugin, dbPlugin := createFeatureTestService(t)

	t.Run("merges keeper, database and tool features", func(t *testing.T) {
		features := s.getSupportedFeatures(keeperPlugin, dbPlugin)
		if !features[env.ViewNodeKeeperOverview] {
			t.Errorf("expected keeper feature to be present and true, got %v", features)
		}
		if !features[env.ViewQueryDbInfo] {
			t.Errorf("expected database feature to be present and true, got %v", features)
		}
		if !features[env.ViewToolPgCompactTableList] {
			t.Errorf("expected tool feature to be present and true, got %v", features)
		}
	})

	t.Run("unknown keeper plugin contributes no keeper features but keeps the rest", func(t *testing.T) {
		features := s.getSupportedFeatures(node.KeeperPlugin("unknown"), dbPlugin)
		if _, ok := features[env.ViewNodeKeeperOverview]; ok {
			t.Errorf("expected no keeper feature for an unknown plugin, got %v", features)
		}
		if !features[env.ViewQueryDbInfo] {
			t.Errorf("expected database feature to still be present, got %v", features)
		}
	})
}

func TestService_getToolSupportedFeatures(t *testing.T) {
	s, _, dbPlugin := createFeatureTestService(t)

	features := s.getToolSupportedFeatures(dbPlugin)
	if !features[env.ViewToolPgCompactTableList] {
		t.Errorf("expected tool feature to be present and true, got %v", features)
	}
}

func TestService_hasKeeper(t *testing.T) {
	s := &Service{}
	host := "db1"

	tests := []struct {
		name     string
		response node.KeeperOneResponse
		expected bool
	}{
		{"empty response has no keeper", node.KeeperOneResponse{}, false},
		{"role set means a keeper responded", node.KeeperOneResponse{Role: node.KeeperRoleLeader}, true},
		{"state set means a keeper responded", node.KeeperOneResponse{State: keeper.StateRunning}, true},
		{"discovered host set means a keeper responded", node.KeeperOneResponse{DiscoveredHost: &host}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.hasKeeper(tt.response); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestService_isPortEqual(t *testing.T) {
	s := &Service{}
	p8008 := 8008
	p5432 := 5432
	otherP8008 := 8008

	tests := []struct {
		name     string
		p1       *int
		p2       *int
		expected bool
	}{
		{"both nil are equal", nil, nil, true},
		{"nil and non-nil are not equal", nil, &p8008, false},
		{"non-nil and nil are not equal", &p8008, nil, false},
		{"equal values with different pointers are equal", &p8008, &otherP8008, true},
		{"different values are not equal", &p8008, &p5432, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.isPortEqual(tt.p1, tt.p2); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestService_getNodeKey(t *testing.T) {
	s := &Service{}
	host := "localhost"
	port := 8008

	t.Run("should return host and port when port is not nil", func(t *testing.T) {
		key := s.getNodeKey(host, &port)
		expected := "localhost:8008"
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
