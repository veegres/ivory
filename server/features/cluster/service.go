package cluster

import (
	"errors"
	"fmt"
	env "ivory/core/config"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/features/node"
	"ivory/features/query"
	"ivory/features/tag"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"ivory/tools"
	"slices"
)

var ErrClusterNameEmpty = errors.New("cluster name cannot be empty")
var ErrClusterKeepersEmpty = errors.New("cluster keepers cannot be empty")

type Service struct {
	clusterRepository *Repository
	nodeService       *node.Service
	tagService        *tag.Service
	queryService      *query.Service
	vaultService      *vault.Service
	toolRegistry      *utils.Registry[tools.Tool, tools.Adapter]
}

func NewService(
	clusterRepository *Repository,
	nodeService *node.Service,
	tagService *tag.Service,
	queryService *query.Service,
	vaultService *vault.Service,
	toolRegistry *utils.Registry[tools.Tool, tools.Adapter],
) *Service {
	return &Service{
		clusterRepository: clusterRepository,
		nodeService:       nodeService,
		tagService:        tagService,
		queryService:      queryService,
		vaultService:      vaultService,
		toolRegistry:      toolRegistry,
	}
}

func (s *Service) getSupportedFeatures(k keeper.Plugin, db database.Plugin) []env.Feature {
	fk := s.nodeService.SupportedFeatures(k)
	fdb := s.queryService.SupportedFeatures(db)
	ft := s.getToolSupportedFeatures(db)
	return slices.Concat(fk, fdb, ft)
}

func (s *Service) getToolSupportedFeatures(db database.Plugin) []env.Feature {
	allFeatures := make([]env.Feature, 0)
	for _, tool := range s.toolRegistry.All() {
		allFeatures = append(allFeatures, tool.SupportedFeatures(db)...)
	}
	return allFeatures
}

func (s *Service) hasKeeper(k node.KeeperOneResponse) bool {
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

func (s *Service) mapKeeperResponseList(keeperNodes []node.KeeperOneResponse) []NodeConfig {
	nodes := make([]NodeConfig, 0)
	for _, item := range keeperNodes {
		nodes = append(nodes, s.mapKeeperResponse(item))
	}
	return nodes
}

func (s *Service) mapKeeperResponseMap(keeperNodes map[string]node.KeeperOneResponse) []NodeConfig {
	nodes := make([]NodeConfig, 0)
	for _, item := range keeperNodes {
		nodes = append(nodes, s.mapKeeperResponse(item))
	}
	return nodes
}

func (s *Service) mapKeeperResponse(keeper node.KeeperOneResponse) NodeConfig {
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
