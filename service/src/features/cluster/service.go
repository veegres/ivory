package cluster

import (
	"errors"
	"fmt"
	"ivory/src/features"
	"ivory/src/features/cert"
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

func (s *Service) List() ([]Response, error) {
	return s.clusterRepository.List()
}

func (s *Service) ListByTag(tags []string) ([]Response, error) {
	listMap := make(map[string]bool)
	for _, t := range tags {
		// NOTE: we shouldn't check the error here, we want to return an empty array if there is no such tag
		clusters, _ := s.tagService.Get(t)
		for _, c := range clusters {
			if !listMap[c] {
				listMap[c] = true
			}
		}
	}

	listName := make([]string, 0)
	for k := range listMap {
		listName = append(listName, k)
	}

	return s.ListByName(listName)
}

func (s *Service) ListByName(clusters []string) ([]Response, error) {
	return s.clusterRepository.ListByName(clusters)
}

func (s *Service) Get(cluster string) (Response, error) {
	return s.clusterRepository.Get(cluster)
}

func (s *Service) Update(cluster Request) (*Response, error) {
	if cluster.Name == "" {
		return nil, ErrClusterNameEmpty
	}
	if cluster.Nodes == nil {
		return nil, ErrClusterKeepersEmpty
	}
	tags, err := s.saveTags(cluster.Name, cluster.Tags)
	if err != nil {
		return nil, err
	}
	cluster.Tags = tags
	errCluster := s.clusterRepository.Update(cluster)
	return (*Response)(&cluster), errCluster
}

func (s *Service) Delete(cluster string) error {
	_, errTag := s.tagService.UpdateCluster(cluster, nil)
	if errTag != nil {
		return errTag
	}
	return s.clusterRepository.Delete(cluster)
}

func (s *Service) DeleteAll() error {
	return s.clusterRepository.DeleteAll()
}

func (s *Service) Overview(name string, host string, port int) (*Overview, error) {
	cluster, clusterError := s.Get(name)
	if clusterError != nil {
		return nil, clusterError
	}
	var keeperNodes []node.KeeperResponse
	var detectedCon *node.KeeperConnection
	var err error

	// if host is not set, search manually only for this host
	if host == "" {
		keeperNodes, detectedCon, err = s.getKeeperListAuto(cluster.Nodes, cluster.Options)
	} else {
		keeperNodes, detectedCon, err = s.getKeeperList(host, port, cluster.Options)
	}
	if keeperNodes == nil || detectedCon == nil {
		return nil, err
	}

	matcherNodeMap := make(map[string]Node)
	for _, n := range cluster.Nodes {
		nodeKey := s.getNodeKey(n.Host, n.KeeperPort)
		warnings := make([]string, 0)
		if v, ok := matcherNodeMap[nodeKey]; ok {
			v.Warnings = append(warnings, "node is declared more than once in the cluster configuration")
			matcherNodeMap[nodeKey] = v
		} else {
			matcherNodeMap[nodeKey] = Node{Config: n, Warnings: warnings}
		}
	}

	resultNodeMap := make(map[string]Node)
	for _, kn := range keeperNodes {
		nodeKey := s.getNodeKey(*kn.DiscoveredHost, kn.DiscoveredKeeperPort)
		if cn, ok := matcherNodeMap[nodeKey]; ok {
			cn.Keeper = kn
			resultNodeMap[nodeKey] = cn
			if !s.isPortEqual(cn.Config.DbPort, kn.DiscoveredDbPort) {
				cn.Warnings = append(cn.Warnings, "database port in keeper response and cluster configuration mismatch")
			}
			delete(matcherNodeMap, nodeKey)
		} else {
			config := NodeConfig{Host: *kn.DiscoveredHost, KeeperPort: kn.DiscoveredKeeperPort, DbPort: kn.DiscoveredDbPort}
			warnings := []string{"node was found in Keeper response, but not in the cluster configuration"}
			resultNodeMap[nodeKey] = Node{config, kn, warnings}
		}
	}

	for nodeKey, cn := range matcherNodeMap {
		cn.Keeper = node.KeeperResponse{Role: keeper.Unknown, State: "-"}
		cn.Warnings = append(cn.Warnings, "node was not found in Keeper response")
		resultNodeMap[nodeKey] = cn
	}

	supportedFeatures := s.getSupportedFeatures(cluster.Plugins.Keeper, cluster.Plugins.Database)
	detectionNodeKey := s.getNodeKey(detectedCon.Host, &detectedCon.Port)
	return &Overview{resultNodeMap, detectionNodeKey, supportedFeatures}, nil
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

func (s *Service) CreateAuto(cluster AutoRequest) (Response, error) {
	keeperNodes, _, errOver := s.getKeeperList(cluster.Host, cluster.Port, cluster.Options)
	if errOver != nil {
		return Response{}, errOver
	}
	nodes := s.mapKeeperResponseList(keeperNodes)

	tags, errSave := s.saveTags(cluster.Name, cluster.Tags)
	if errSave != nil {
		return Response{}, errSave
	}
	cluster.Tags = tags

	model := Request{
		Name:  cluster.Name,
		Nodes: nodes,
		Options: Options{
			Plugins: cluster.Plugins,
			Tls:     cluster.Tls,
			Certs:   cluster.Certs,
			Vaults:  cluster.Vaults,
			Tags:    tags,
		},
	}

	return s.clusterRepository.Create(model)
}

func (s *Service) FixAuto(name string) (*Response, error) {
	cluster, clusterError := s.Get(name)
	if clusterError != nil {
		return nil, clusterError
	}
	keeperNodes, _, err := s.getKeeperListAuto(cluster.Nodes, cluster.Options)
	if err != nil {
		return nil, err
	}
	nodes := s.mapKeeperResponseList(keeperNodes)

	model := Request{
		Name:  cluster.Name,
		Nodes: nodes,
		Options: Options{
			Plugins: cluster.Plugins,
			Tls:     cluster.Tls,
			Certs:   cluster.Certs,
			Vaults:  cluster.Vaults,
			Tags:    cluster.Tags,
		},
	}
	return (*Response)(&model), s.clusterRepository.Update(model)
}

func (s *Service) mapKeeperResponseList(keeperNodes []node.KeeperResponse) []NodeConfig {
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

func (s *Service) getKeeperList(host string, port int, cluster Options) ([]node.KeeperResponse, *node.KeeperConnection, error) {
	var certs *cert.Certs
	// NOTE: we want to rewrite `nil` only if tls is enabled
	if cluster.Tls.Keeper {
		certs = &cluster.Certs
	}
	con := node.KeeperConnection{Host: host, Port: port}
	request := node.KeeperRequest{
		Plugin:           cluster.Plugins.Keeper,
		KeeperConnection: con,
		VaultId:          cluster.Vaults.KeeperId,
		Certs:            certs,
	}
	nodes, _, errOver := s.nodeService.List(request)
	return nodes, &con, errOver
}

func (s *Service) getKeeperListAuto(configs []NodeConfig, cluster Options) ([]node.KeeperResponse, *node.KeeperConnection, error) {
	connections := make([]node.KeeperConnection, 0)
	var portErrs error
	for i, config := range configs {
		if config.KeeperPort == nil {
			portErrs = errors.Join(portErrs, fmt.Errorf("#%d keeper port for the host `%s` isn't set", i, config.Host))
		} else {
			connections = append(connections, node.KeeperConnection{
				Host: config.Host,
				Port: *config.KeeperPort,
			})
		}
	}
	var certs *cert.Certs
	// NOTE: we want to rewrite `nil` only if tls is enabled
	if cluster.Tls.Keeper {
		certs = &cluster.Certs
	}
	request := node.KeeperAutoRequest{
		Connections: connections,
		Plugin:      cluster.Plugins.Keeper,
		VaultId:     cluster.Vaults.KeeperId,
		Certs:       certs,
	}
	nodes, _, detected, err := s.nodeService.ListAuto(request)
	return nodes, detected, errors.Join(portErrs, err)
}

func (s *Service) getSupportedFeatures(k keeper.Plugin, db database.Plugin) []features.Feature {
	fk := s.nodeService.SupportedFeatures(k)
	fdb := s.queryService.SupportedFeatures(db)
	ft := s.toolsService.SupportedFeatures(db)
	return slices.Concat(fk, fdb, ft)
}

func (s *Service) saveTags(name string, tags []string) ([]string, error) {
	// NOTE: remove duplicates
	tagMap := make(map[string]bool)
	for _, t := range tags {
		tagMap[t] = true
	}
	tagList := make([]string, 0)
	for key := range tagMap {
		tagList = append(tagList, key)
	}

	// NOTE: create tags in db with a cluster name
	return s.tagService.UpdateCluster(name, tagList)
}
