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

	hasLeader := s.nodeService.KeeperHasLeader(cluster.Plugins.Keeper)
	resultNodeMap := s.buildOverviewNodes(cluster.Nodes, keeperNodeMap, connectionErrors, requestError, hasLeader)
	supportedFeatures := s.getSupportedFeatures(cluster.Plugins.Keeper, cluster.Plugins.Database)
	return &Overview{resultNodeMap, supportedFeatures}, nil
}

func (s *Service) Detect(cluster CreateAutoRequest) (*Response, error) {
	keeperNodeMap, _, errOver := s.getKeeperListByOne(cluster.Host, cluster.Port, cluster.Options)
	if errOver != nil {
		return nil, errOver
	}
	nodes := mapKeeperResponseMap(keeperNodeMap)
	if err := s.validateNodeNames(nodes); err != nil {
		return nil, err
	}

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

	return s.Create(model)
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
	if err := s.validateNodeNames(nodes); err != nil {
		return nil, err
	}

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
		if s.hasLeaderEntry(response.Response) {
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
		if s.hasLeaderEntry(response.Response) {
			return response.Response, nil
		}
	}
	return nil, errors.Join(requestErrs, ErrNoLeaderFound)
}

func (s *Service) hasLeaderEntry(responses []node.KeeperOneResponse) bool {
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

func (s *Service) buildOverviewNodes(configs []NodeConfig, keeperNodes map[string]node.KeeperOneResponse, connectionErrors map[string]error, requestError error, hasLeader bool) map[string]Node {
	resultNodeMap := s.getConfiguredNodeMap(configs, connectionErrors, requestError)
	// NOTE: two passes, because one of these describes a node and the other only
	// describes an attribute of one. Merging them together would leave the
	// outcome to map iteration order - whichever reached a node first would
	// win, and a sync flag applied before the node's own response arrived would
	// then be overwritten by it.
	for _, kn := range keeperNodes {
		if s.hasKeeper(kn) {
			s.mergeKeeperNode(resultNodeMap, kn)
		}
	}
	for _, kn := range keeperNodes {
		if !s.hasKeeper(kn) {
			s.mergeKeeperSync(resultNodeMap, kn)
		}
	}
	s.addOverviewWarnings(resultNodeMap, hasLeader)
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
	// NOTE: a response that states no host names its node instead - an etcd
	// member the keeper knows about but has no client url for yet. It cannot
	// create a node, since there is nothing to reach it at, but it can still
	// report the state of one Ivory already has configured; it used to be
	// dropped before it got here.
	if kn.DiscoveredHost == nil {
		if nodeKey, ok := s.resolveKeyByName(nodeMap, kn.DiscoveredName); ok && !s.hasKeeper(nodeMap[nodeKey].Keeper) {
			cn := nodeMap[nodeKey]
			cn.Keeper = kn
			nodeMap[nodeKey] = cn
		}
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
		config := mapKeeperResponse(kn)
		warnings := []string{"node was found in Keeper response, but not in the cluster configuration"}
		nodeMap[nodeKey] = Node{config, kn, warnings}
	}
}

// mergeKeeperSync applies a response that claims no state of its own: it is an
// attribute of a node rather than a node, and Sync is the only such attribute -
// a replica's synchronous membership is visible from the primary alone, so a
// node's own response can never carry it (see postgres.mapSyncStandby).
//
// It is applied on top of whatever the node itself reported, never in place of
// it: the alternative reported a replica with no lag, no ports and a liveness
// the primary is in no position to vouch for, and only for as long as map
// iteration order happened to favour it.
func (s *Service) mergeKeeperSync(nodeMap map[string]Node, kn node.KeeperOneResponse) {
	nodeKey, ok := s.resolveKeyByName(nodeMap, kn.DiscoveredName)
	if !ok {
		return
	}
	cn := nodeMap[nodeKey]
	cn.Keeper.Sync = kn.Sync
	nodeMap[nodeKey] = cn
}

// resolveKeyByName finds the configured node with this name. Unlike a host it
// needs no ambiguity check: a name is unique within a cluster, enforced on every
// write by validateNodeNames, which is what lets it identify one node of a
// single-host cluster where the shared host cannot.
func (s *Service) resolveKeyByName(nodeMap map[string]Node, name *string) (string, bool) {
	if name == nil || *name == "" {
		return "", false
	}
	for nodeKey, n := range nodeMap {
		if n.Config.Name == *name {
			return nodeKey, true
		}
	}
	return "", false
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

// addOverviewWarnings annotates each node. hasLeader says whether the engine
// elects a single primary at all: where it does not, no node reports one and
// warning that none was found would fire on every healthy cluster.
func (s *Service) addOverviewWarnings(nodeMap map[string]Node, hasLeader bool) {
	leaderKeys := make([]string, 0)
	for nodeKey, cn := range nodeMap {
		// NOTE: the port check runs only against a keeper that answered - the
		// placeholder below reports no port, and a node nothing observed
		// cannot disagree with the configuration
		if s.hasKeeper(cn.Keeper) {
			if !s.isPortEqual(cn.Config.DbPort, cn.Keeper.DiscoveredDbPort) {
				cn.Warnings = append(cn.Warnings, "database port in keeper response and cluster configuration mismatch")
			}
		} else {
			cn.Keeper = node.KeeperOneResponse{Role: node.KeeperRoleUnknown, State: node.KeeperStateUnreachable}
			cn.Warnings = append(cn.Warnings, "node was not found in Keeper response")
		}
		if cn.Keeper.Role == node.KeeperRoleLeader {
			leaderKeys = append(leaderKeys, nodeKey)
		}
		nodeMap[nodeKey] = cn
	}
	if !hasLeader {
		return
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
		nodeKey, ok := s.getResponseKey(n)
		if !ok {
			continue
		}
		if _, exists := nodeMap[nodeKey]; exists {
			continue
		}
		nodeMap[nodeKey] = n
	}
}

// getResponseKey identifies the member a keeper response is about, for dedup
// across the several nodes whose view it may arrive in. An endpoint is preferred
// where the keeper reported one; a response that states only a name is kept
// under that name rather than dropped, which is what used to happen to native
// postgres' sync-standby responses and to an etcd member with no client url yet.
func (s *Service) getResponseKey(n node.KeeperOneResponse) (string, bool) {
	if n.DiscoveredHost != nil {
		return s.getNodeKey(*n.DiscoveredHost, n.DiscoveredKeeperPort), true
	}
	if n.DiscoveredName != nil {
		return *n.DiscoveredName, true
	}
	return "", false
}
