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
	var keeperNodeMap map[string]node.KeeperResponse
	var connectionErrors map[string]error
	var requestError error

	// NOTE: if host is not set, search manually only for this host
	if host == "" {
		keeperNodeMap, connectionErrors, requestError = s.getKeeperListAutoAll(cluster.Nodes, cluster.Options)
	} else {
		keeperNodeMap, connectionErrors, requestError = s.getKeeperList(host, port, cluster.Options)
	}

	resultNodeMap := s.buildOverviewNodes(cluster.Nodes, keeperNodeMap, connectionErrors, requestError)
	supportedFeatures := s.getSupportedFeatures(cluster.Plugins.Keeper, cluster.Plugins.Database)
	return &Overview{resultNodeMap, supportedFeatures}, nil
}

func (s *Service) buildOverviewNodes(configs []NodeConfig, keeperNodes map[string]node.KeeperResponse, connectionErrors map[string]error, requestError error) map[string]Node {
	resultNodeMap := s.getConfiguredNodeMap(configs, connectionErrors, requestError)
	for _, kn := range keeperNodes {
		s.mergeKeeperNode(resultNodeMap, kn)
	}
	s.addOverviewWarnings(resultNodeMap)
	return resultNodeMap
}

func (s *Service) getConfiguredNodeMap(configs []NodeConfig, connectionErrors map[string]error, requestError error) map[string]Node {
	nodeMap := make(map[string]Node)
	hasSpecificErrors := len(connectionErrors) > 0
	for _, n := range configs {
		nodeKey := s.getNodeKey(n.Host, n.KeeperPort)
		if v, ok := nodeMap[nodeKey]; ok {
			v.Warnings = append(v.Warnings, "node is declared more than once in the cluster configuration")
			nodeMap[nodeKey] = v
		} else {
			nn := Node{Config: n, Warnings: make([]string, 0)}
			if errCon, ok := connectionErrors[nodeKey]; ok {
				nn.Warnings = append(nn.Warnings, fmt.Sprintf("failed to get Keeper response: %s", errCon.Error()))
			} else if requestError != nil && !hasSpecificErrors {
				nn.Warnings = append(nn.Warnings, fmt.Sprintf("failed to get Keeper response: %s", requestError.Error()))
			}
			nodeMap[nodeKey] = nn
		}
	}
	return nodeMap
}

func (s *Service) mergeKeeperNode(nodeMap map[string]Node, kn node.KeeperResponse) {
	if kn.DiscoveredHost == nil {
		return
	}
	nodeKey := s.getNodeKey(*kn.DiscoveredHost, kn.DiscoveredKeeperPort)
	if s.hasKeeper(nodeMap[nodeKey].Keeper) {
		return
	}
	if cn, ok := nodeMap[nodeKey]; ok {
		cn.Keeper = kn
		nodeMap[nodeKey] = cn
	} else {
		config := NodeConfig{Host: *kn.DiscoveredHost, KeeperPort: kn.DiscoveredKeeperPort, DbPort: kn.DiscoveredDbPort}
		warnings := []string{"node was found in Keeper response, but not in the cluster configuration"}
		nodeMap[nodeKey] = Node{config, kn, warnings}
	}
}

func (s *Service) addOverviewWarnings(nodeMap map[string]Node) {
	leaderKeys := make([]string, 0)
	for nodeKey, cn := range nodeMap {
		if !s.hasKeeper(cn.Keeper) {
			cn.Keeper = node.KeeperResponse{Role: keeper.Unknown, State: "-"}
			cn.Warnings = append(cn.Warnings, "node was not found in Keeper response")
			nodeMap[nodeKey] = cn
		}
		if !s.isPortEqual(cn.Config.DbPort, cn.Keeper.DiscoveredDbPort) {
			cn.Warnings = append(cn.Warnings, "database port in keeper response and cluster configuration mismatch")
		}
		if cn.Keeper.Role == keeper.Leader {
			leaderKeys = append(leaderKeys, nodeKey)
		}
	}
	if len(leaderKeys) < 2 {
		return
	}
	for _, nodeKey := range leaderKeys {
		n := nodeMap[nodeKey]
		n.Warnings = append(n.Warnings, "multiple leader nodes were found in Keeper response")
		nodeMap[nodeKey] = n
	}
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

func (s *Service) CreateAuto(cluster AutoRequest) (Response, error) {
	keeperNodeMap, _, errOver := s.getKeeperList(cluster.Host, cluster.Port, cluster.Options)
	if errOver != nil {
		return Response{}, errOver
	}
	nodes := s.mapKeeperResponseMap(keeperNodeMap)

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
	keeperNodes, err := s.getKeeperListAutoFirstSuccess(cluster.Nodes, cluster.Options)
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

func (s *Service) addKeeperResponsesToMap(nodeMap map[string]node.KeeperResponse, nodes []node.KeeperResponse) {
	for _, n := range nodes {
		if n.DiscoveredHost == nil {
			continue
		}
		nodeKey := s.getNodeKey(*n.DiscoveredHost, n.DiscoveredKeeperPort)
		if _, ok := nodeMap[nodeKey]; ok {
			continue
		}
		nodeMap[nodeKey] = n
	}
}

func (s *Service) getKeeperList(host string, port int, cluster Options) (map[string]node.KeeperResponse, map[string]error, error) {
	var certs *cert.Certs
	// NOTE: we want to rewrite `nil` only if tls is enabled
	if cluster.Tls.Keeper {
		certs = &cluster.Certs
	}
	con := node.KeeperConnection{Host: host, Port: port}
	request := node.KeeperRequest{
		KeeperConnection: con,
		KeeperOptions: node.KeeperOptions{
			Plugin:  cluster.Plugins.Keeper,
			VaultId: cluster.Vaults.KeeperId,
			Certs:   certs,
		},
	}
	nodes, _, errOver := s.nodeService.List(request)
	var connectionErrors map[string]error
	if errOver != nil {
		connectionErrors = make(map[string]error)
		connectionErrors[s.getNodeKey(host, &port)] = errOver
	}
	nodeMap := make(map[string]node.KeeperResponse)
	s.addKeeperResponsesToMap(nodeMap, nodes)
	return nodeMap, connectionErrors, errOver
}

func (s *Service) getKeeperListAutoAll(configs []NodeConfig, cluster Options) (map[string]node.KeeperResponse, map[string]error, error) {
	responses, connectionErrors, err := s.getKeeperListAutoResponse(configs, cluster)
	if err != nil {
		return nil, connectionErrors, err
	}

	keeperNodeMap := make(map[string]node.KeeperResponse)
	var requestErrs error
	for _, response := range responses {
		connection := response.Connection
		connectionKey := s.getNodeKey(connection.Host, &connection.Port)
		if response.Error != nil {
			connectionErrors[connectionKey] = response.Error
			requestErrs = errors.Join(requestErrs, response.Error)
			continue
		}
		s.addKeeperResponsesToMap(keeperNodeMap, response.Response)
	}
	return keeperNodeMap, connectionErrors, requestErrs
}

func (s *Service) getKeeperListAutoFirstSuccess(configs []NodeConfig, cluster Options) ([]node.KeeperResponse, error) {
	responses, _, err := s.getKeeperListAutoResponse(configs, cluster)
	if err != nil {
		return nil, err
	}

	var requestErrs error
	for _, response := range responses {
		if response.Error != nil {
			requestErrs = errors.Join(requestErrs, response.Error)
			continue
		}
		return response.Response, nil
	}
	return nil, requestErrs
}

func (s *Service) getKeeperListAutoResponse(configs []NodeConfig, cluster Options) ([]node.KeeperParallelResponse, map[string]error, error) {
	connections := make([]node.KeeperConnection, 0)
	connectionErrors := make(map[string]error)
	for _, config := range configs {
		if config.KeeperPort == nil {
			key := s.getNodeKey(config.Host, config.KeeperPort)
			err := fmt.Errorf("host %q is missing a keeper port", config.Host)
			connectionErrors[key] = err
		} else {
			connections = append(connections, node.KeeperConnection{
				Host: config.Host,
				Port: *config.KeeperPort,
			})
		}
	}
	// NOTE: we want to rewrite `nil` only if tls is enabled
	var certs *cert.Certs
	if cluster.Tls.Keeper {
		certs = &cluster.Certs
	}
	if len(connections) == 0 {
		return nil, connectionErrors, errors.New("no configured keeper connections can be requested")
	}

	request := node.KeeperParallelRequest{
		Connections: connections,
		KeeperOptions: node.KeeperOptions{
			Plugin:  cluster.Plugins.Keeper,
			VaultId: cluster.Vaults.KeeperId,
			Certs:   certs,
		},
	}

	responses, err := s.nodeService.ListParallel(request)
	if err != nil {
		return nil, connectionErrors, err
	}
	return responses, connectionErrors, nil
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
