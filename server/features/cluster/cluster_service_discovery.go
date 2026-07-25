package cluster

import (
	"errors"
	"fmt"
	"ivory/core/service/cert"
	"ivory/features/node"
)

func (s *Service) Overview(name string, host string, port int) (*Overview, error) {
	cluster, clusterError := s.Get(name)
	if clusterError != nil {
		return nil, clusterError
	}
	var keeperNodeMap map[string]node.KeeperOneResponse
	var connectionErrors map[string]error
	var requestError error

	// NOTE: if host is not set, search manually only for this host
	if host == "" {
		keeperNodeMap, connectionErrors, requestError = s.getKeeperListByManyAll(cluster.Nodes, cluster.Options)
	} else {
		keeperNodeMap, connectionErrors, requestError = s.getKeeperListByOne(host, port, cluster.Options)
	}

	resultNodeMap := s.buildOverviewNodes(cluster.Nodes, keeperNodeMap, connectionErrors, requestError)
	supportedFeatures := s.getSupportedFeatures(cluster.Plugins.Keeper, cluster.Plugins.Database)
	return &Overview{resultNodeMap, supportedFeatures}, nil
}

func (s *Service) Detect(cluster CreateAutoRequest) (*Response, error) {
	keeperNodeMap, _, errOver := s.getKeeperListByOne(cluster.Host, cluster.Port, cluster.Options)
	if errOver != nil {
		return nil, errOver
	}
	nodes := mapKeeperResponseMap(keeperNodeMap)

	tags, errSave := s.saveTags(cluster.Name, cluster.Tags)
	if errSave != nil {
		return nil, errSave
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

	r, err := s.clusterRepository.Create(model)

	return &r, err
}

func (s *Service) Fix(name string) (*Response, error) {
	cluster, clusterError := s.Get(name)
	if clusterError != nil {
		return nil, clusterError
	}
	keeperNodes, err := s.getKeeperListByLeader(cluster.Nodes, cluster.Options)
	if err != nil {
		return nil, err
	}
	nodes := mapKeeperResponseList(keeperNodes)

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

func (s *Service) getKeeperListByOne(host string, port int, cluster Options) (map[string]node.KeeperOneResponse, map[string]error, error) {
	var certs *cert.Certs
	// NOTE: we want to rewrite `nil` only if tls is enabled
	if cluster.Tls.Keeper {
		certs = &cluster.Certs
	}
	con := node.KeeperConnection{Host: host, Port: port}
	request := node.KeeperOneRequest{
		KeeperConnection: con,
		KeeperOptions: node.KeeperOptions{
			Plugin:  cluster.Plugins.Keeper,
			VaultId: cluster.Vaults.KeeperId,
			Certs:   certs,
		},
	}
	nodes, _, errOver := s.nodeService.KeeperNodeList(request)
	var connectionErrors map[string]error
	if errOver != nil {
		connectionErrors = make(map[string]error)
		connectionErrors[s.getNodeKey(host, &port)] = errOver
	}
	nodeMap := make(map[string]node.KeeperOneResponse)
	s.addKeeperResponsesToMap(nodeMap, nodes)
	return nodeMap, connectionErrors, errOver
}

func (s *Service) getKeeperListByManyAll(configs []NodeConfig, cluster Options) (map[string]node.KeeperOneResponse, map[string]error, error) {
	responses, connectionErrors, err := s.getKeeperListByManyResponse(configs, cluster)
	if err != nil {
		return nil, connectionErrors, err
	}

	var requestErrs error
	for _, response := range responses {
		if response.Error != "" {
			errResponse := errors.New(response.Error)
			connectionKey := s.getNodeKey(response.Connection.Host, &response.Connection.Port)
			connectionErrors[connectionKey] = errResponse
			requestErrs = errors.Join(requestErrs, errResponse)
		}
	}

	keeperNodeMap := make(map[string]node.KeeperOneResponse)
	for _, response := range responses {
		if hasLeaderEntry(response.Response) {
			s.addKeeperResponsesToMap(keeperNodeMap, response.Response)
		}
	}
	for _, response := range responses {
		s.addKeeperResponsesToMap(keeperNodeMap, response.Response)
	}
	return keeperNodeMap, connectionErrors, requestErrs
}

func (s *Service) getKeeperListByLeader(configs []NodeConfig, cluster Options) ([]node.KeeperOneResponse, error) {
	responses, _, err := s.getKeeperListByManyResponse(configs, cluster)
	if err != nil {
		return nil, err
	}

	var requestErrs error
	for _, response := range responses {
		if response.Error != "" {
			requestErrs = errors.Join(requestErrs, errors.New(response.Error))
			continue
		}
		if hasLeaderEntry(response.Response) {
			return response.Response, nil
		}
	}
	return nil, errors.Join(requestErrs, ErrNoLeaderFound)
}

func hasLeaderEntry(responses []node.KeeperOneResponse) bool {
	for _, r := range responses {
		if r.Role == node.KeeperRoleLeader {
			return true
		}
	}
	return false
}

func (s *Service) getKeeperListByManyResponse(configs []NodeConfig, cluster Options) ([]node.KeeperMultiResponse, map[string]error, error) {
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
		return nil, connectionErrors, ErrNoKeeperConnections
	}

	request := node.KeeperMultiRequest{
		Connections: connections,
		KeeperOptions: node.KeeperOptions{
			Plugin:  cluster.Plugins.Keeper,
			VaultId: cluster.Vaults.KeeperId,
			Certs:   certs,
		},
	}

	responses, err := s.nodeService.KeeperNodeListMulti(request)
	if err != nil {
		return nil, connectionErrors, err
	}
	return responses, connectionErrors, nil
}

func (s *Service) buildOverviewNodes(configs []NodeConfig, keeperNodes map[string]node.KeeperOneResponse, connectionErrors map[string]error, requestError error) map[string]Node {
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

func (s *Service) mergeKeeperNode(nodeMap map[string]Node, kn node.KeeperOneResponse) {
	if kn.DiscoveredHost == nil {
		return
	}
	// NOTE: a response that couldn't determine its own port (e.g. native
	// postgres reporting a standby's sync status from the primary's
	// pg_stat_replication - see postgres.mapSyncStandby) leaves
	// DiscoveredKeeperPort nil rather than guessing it. Resolve it to the
	// configured node with a matching host instead, but only when that
	// host is unambiguous (exactly one configured node), otherwise fall
	// through to the generic host-only key below rather than risk
	// attributing it to the wrong node.
	if kn.DiscoveredKeeperPort == nil {
		if resolved, ok := s.resolveConfigByHost(nodeMap, *kn.DiscoveredHost); ok {
			kn.DiscoveredKeeperPort = resolved.KeeperPort
			kn.DiscoveredDbPort = resolved.DbPort
		}
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

// resolveConfigByHost finds the configured node matching host, only when
// unambiguous (exactly one configured node has that host).
func (s *Service) resolveConfigByHost(nodeMap map[string]Node, host string) (NodeConfig, bool) {
	var match NodeConfig
	count := 0
	for _, n := range nodeMap {
		if n.Config.Host == host {
			match = n.Config
			count++
		}
	}
	if count == 1 {
		return match, true
	}
	return NodeConfig{}, false
}

func (s *Service) addOverviewWarnings(nodeMap map[string]Node) {
	leaderKeys := make([]string, 0)
	for nodeKey, cn := range nodeMap {
		if !s.hasKeeper(cn.Keeper) {
			cn.Keeper = node.KeeperOneResponse{Role: node.KeeperRoleUnknown, State: node.KeeperStateUnreachable}
			cn.Warnings = append(cn.Warnings, "node was not found in Keeper response")
		}
		if !s.isPortEqual(cn.Config.DbPort, cn.Keeper.DiscoveredDbPort) {
			cn.Warnings = append(cn.Warnings, "database port in keeper response and cluster configuration mismatch")
		}
		if cn.Keeper.Role == node.KeeperRoleLeader {
			leaderKeys = append(leaderKeys, nodeKey)
		}
		nodeMap[nodeKey] = cn
	}
	switch {
	case len(leaderKeys) == 0:
		for nodeKey, n := range nodeMap {
			n.Warnings = append(n.Warnings, "no leader node was found in Keeper response")
			nodeMap[nodeKey] = n
		}
	case len(leaderKeys) > 1:
		for _, nodeKey := range leaderKeys {
			n := nodeMap[nodeKey]
			n.Warnings = append(n.Warnings, "multiple leader nodes were found in Keeper response")
			nodeMap[nodeKey] = n
		}
	}
}

func (s *Service) addKeeperResponsesToMap(nodeMap map[string]node.KeeperOneResponse, nodes []node.KeeperOneResponse) {
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
