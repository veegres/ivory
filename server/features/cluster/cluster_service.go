package cluster

import (
	"errors"
	"fmt"
	"ivory/core/config"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/features/node"
	"ivory/features/query"
	"ivory/features/tag"
	"ivory/tools"
	"maps"
)

var ErrClusterNameEmpty = errors.New("cluster name cannot be empty")
var ErrClusterKeepersEmpty = errors.New("cluster keepers cannot be empty")
var ErrClusterNameNotProvided = errors.New("cluster name not provided")
var ErrClusterNodesNotProvided = errors.New("cluster nodes not provided")
var ErrClusterNodeNameNotProvided = errors.New("cluster node name not provided")
var ErrClusterNodeNameNotUnique = errors.New("cluster node name is not unique")
var ErrSshCredentialsRequired = errors.New("ssh credentials are required")
var ErrDatabaseCredentialsRequired = errors.New("database credentials are required")
var ErrSshCredentialsAmbiguous = errors.New("provide either an ssh vault or a username and password, not both")
var ErrDatabaseCredentialsAmbiguous = errors.New("provide either a database vault or a username and password, not both")
var ErrClusterNameTaken = errors.New("cluster name is already taken")
var ErrSshKeyVaultMissingMetadata = errors.New("ssh key from vault is missing metadata (public key)")
var ErrNoKeeperConnections = errors.New("no configured keeper connections can be requested")
var ErrNoLeaderFound = errors.New("no configured node reported being the cluster leader")

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

// getSupportedFeatures asks the platform for its capabilities alongside the
// keeper's and the database's. A cluster names no platform of its own yet, so
// it reports the one its connections resolve to.
func (s *Service) getSupportedFeatures(k node.KeeperPlugin, db query.DbPlugin) map[config.Feature]bool {
	features := make(map[config.Feature]bool)
	maps.Copy(features, s.nodeService.PlatformSupportedFeatures(node.DefaultPlatform))
	maps.Copy(features, s.nodeService.SupportedFeatures(k))
	maps.Copy(features, s.queryService.SupportedFeatures(db))
	maps.Copy(features, s.getToolSupportedFeatures(db))
	return features
}

func (s *Service) getToolSupportedFeatures(db query.DbPlugin) map[config.Feature]bool {
	features := make(map[config.Feature]bool)
	for _, tool := range s.toolRegistry.All() {
		maps.Copy(features, tool.SupportedFeatures(db))
	}
	return features
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

func (s *Service) getNodeKey(h string, kp *int) string {
	if kp == nil {
		return h
	}
	return fmt.Sprintf("%s:%d", h, *kp)
}
