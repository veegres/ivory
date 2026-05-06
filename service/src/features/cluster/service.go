package cluster

import (
	"errors"
	"fmt"
	"ivory/src/features"
	"ivory/src/features/node"
	"ivory/src/features/query"
	"ivory/src/features/tag"
	"ivory/src/features/tools"
	"ivory/src/plugins/database"
	"ivory/src/plugins/keeper"
	"slices"
)

var ErrClusterNameEmpty = errors.New("cluster name cannot be empty")
var ErrClusterKeepersEmpty = errors.New("cluster keepers cannot be empty")

type Service struct {
	clusterRepository *Repository
	nodeService       *node.Service
	tagService        *tag.Service
	queryService      *query.Service
	toolsService      *tools.Service
}

func NewService(
	clusterRepository *Repository,
	nodeService *node.Service,
	tagService *tag.Service,
	queryService *query.Service,
	toolsService *tools.Service,
) *Service {
	return &Service{
		clusterRepository: clusterRepository,
		nodeService:       nodeService,
		tagService:        tagService,
		queryService:      queryService,
		toolsService:      toolsService,
	}
}

func (s *Service) getSupportedFeatures(k keeper.Plugin, db database.Plugin) []features.Feature {
	fk := s.nodeService.SupportedFeatures(k)
	fdb := s.queryService.SupportedFeatures(db)
	ft := s.toolsService.SupportedFeatures(db)
	return slices.Concat(fk, fdb, ft)
}

func (s *Service) hasKeeper(k node.KeeperResponse) bool {
	return k.Role != "" || k.State != "" || k.DiscoveredHost != nil
}

func (s *Service) isPortEqual(p1 *int, p2 *int) bool {
	if p1 == nil && p2 == nil {
		return true
	}
	if p1 == nil || p2 == nil {
		return false
	}
	return *p1 == *p2
}

func (s *Service) mapKeeperResponseList(keeperNodes []node.KeeperResponse) []NodeConfig {
	nodes := make([]NodeConfig, 0)
	for _, item := range keeperNodes {
		nodes = append(nodes, s.mapKeeperResponse(item))
	}
	return nodes
}

func (s *Service) mapKeeperResponseMap(keeperNodes map[string]node.KeeperResponse) []NodeConfig {
	nodes := make([]NodeConfig, 0)
	for _, item := range keeperNodes {
		nodes = append(nodes, s.mapKeeperResponse(item))
	}
	return nodes
}

func (s *Service) mapKeeperResponse(keeper node.KeeperResponse) NodeConfig {
	return NodeConfig{
		Host:       *keeper.DiscoveredHost,
		KeeperPort: keeper.DiscoveredKeeperPort,
		DbPort:     keeper.DiscoveredDbPort,
	}
}

func (s *Service) getNodeKey(h string, kp *int) string {
	if kp == nil {
		return h
	}
	return fmt.Sprintf("%s:%d", h, *kp)
}
