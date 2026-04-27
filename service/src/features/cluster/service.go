package cluster

import (
	"errors"
	"fmt"
	"ivory/src/clients/database"
	"ivory/src/clients/keeper"
	"ivory/src/features"
	"ivory/src/features/cert"
	"ivory/src/features/node"
	"ivory/src/features/query"
	"ivory/src/features/tag"
	"ivory/src/features/tools"
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

func (s *Service) List() ([]Cluster, error) {
	return s.clusterRepository.List()
}

func (s *Service) ListByTag(tags []string) ([]Cluster, error) {
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

func (s *Service) ListByName(clusters []string) ([]Cluster, error) {
	return s.clusterRepository.ListByName(clusters)
}

func (s *Service) Get(cluster string) (Cluster, error) {
	return s.clusterRepository.Get(cluster)
}

func (s *Service) Update(cluster Cluster) (*Cluster, error) {
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
	return &cluster, errCluster
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

func (s *Service) Overview(name string, host string, port int) (*ClusterOverview, error) {
	cluster, clusterError := s.Get(name)
	if clusterError != nil {
		return nil, clusterError
	}
	var nodes []node.KeeperResponse
	var detectedDomain string
	var err error

	// if host is set, search manually only for this host
	if host == "" {
		var detected *node.Connection
		nodes, detected, err = s.getOverviewAuto(cluster.Nodes, cluster.ClusterOptions)
		if detected != nil {
			detectedDomain = s.getKeyDomain(detected.Host, detected.KeeperPort)
		}
	} else {
		nodes, err = s.getOverview(host, port, cluster.ClusterOptions)
		detectedDomain = s.getKeyDomain(host, port)
	}
	if err != nil {
		return nil, err
	}

	clusterNodeMap := make(map[string]Node)
	for _, n := range cluster.Nodes {
		domain := s.getKeyDomain(n.Host, n.KeeperPort)
		if v, ok := clusterNodeMap[domain]; ok {
			v.Warnings = append(v.Warnings, "node is declared more than once in the cluster configuration")
			clusterNodeMap[domain] = v
		} else {
			keeperResponse := node.KeeperResponse{Connection: n}
			clusterNodeMap[domain] = Node{KeeperResponse: keeperResponse}
		}
	}

	resultNodeMap := make(map[string]Node)
	for _, el := range nodes {
		domain := s.getKeyDomain(el.Connection.Host, el.Connection.KeeperPort)
		if v, ok := clusterNodeMap[domain]; ok {
			v.Keeper = el.Keeper
			resultNodeMap[domain] = v
			delete(clusterNodeMap, domain)
		} else {
			warnings := []string{"node was found in Keeper response, but not in the cluster configuration"}
			resultNodeMap[domain] = Node{el, warnings}
		}
	}

	for domain, n := range clusterNodeMap {
		n.Keeper = keeper.Response{Role: keeper.Unknown, State: "-"}
		n.Warnings = append(n.Warnings, "node was not found in Keeper response")
		resultNodeMap[domain] = n
	}

	supportedFeatures := s.getSupportedFeatures(cluster.KeeperType, cluster.DbType)
	return &ClusterOverview{resultNodeMap, detectedDomain, supportedFeatures}, nil
}

func (s *Service) CreateAuto(cluster ClusterAuto) (Cluster, error) {
	nodes, errOver := s.getOverview(cluster.Host, cluster.Port, cluster.ClusterOptions)
	if errOver != nil {
		return Cluster{}, errOver
	}

	nodesConnections := make([]node.Connection, 0)
	for _, item := range nodes {
		nodesConnections = append(nodesConnections, item.Connection)
	}

	tags, errSave := s.saveTags(cluster.Name, cluster.Tags)
	if errSave != nil {
		return Cluster{}, errSave
	}
	cluster.Tags = tags

	model := Cluster{
		Name:  cluster.Name,
		Nodes: nodesConnections,
		ClusterOptions: ClusterOptions{
			DbType:     cluster.DbType,
			KeeperType: cluster.KeeperType,
			Tls:        cluster.Tls,
			Certs:      cluster.Certs,
			Vaults:     cluster.Vaults,
			Tags:       tags,
		},
	}

	return s.clusterRepository.Create(model)
}

func (s *Service) FixAuto(name string) (*Cluster, error) {
	cluster, clusterError := s.Get(name)
	if clusterError != nil {
		return nil, clusterError
	}
	nodes, _, err := s.getOverviewAuto(cluster.Nodes, cluster.ClusterOptions)
	if err != nil {
		return nil, err
	}

	nodesConnections := make([]node.Connection, 0)
	for _, item := range nodes {
		nodesConnections = append(nodesConnections, item.Connection)
	}

	model := Cluster{
		Name:  cluster.Name,
		Nodes: nodesConnections,
		ClusterOptions: ClusterOptions{
			DbType:     cluster.DbType,
			KeeperType: cluster.KeeperType,
			Tls:        cluster.Tls,
			Certs:      cluster.Certs,
			Vaults:     cluster.Vaults,
			Tags:       cluster.Tags,
		},
	}
	return &model, s.clusterRepository.Update(model)
}

func (s *Service) getKeyDomain(h string, p int) string {
	return fmt.Sprintf("%s:%d", h, p)
}

func (s *Service) getOverview(host string, port int, cluster ClusterOptions) ([]node.KeeperResponse, error) {
	var certs *cert.Certs
	// NOTE: we want to rewrite `nil` only if tls is enabled
	if cluster.Tls.Keeper {
		certs = &cluster.Certs
	}
	request := node.KeeperRequest{
		Host:    host,
		Port:    port,
		VaultId: cluster.Vaults.KeeperId,
		Certs:   certs,
	}
	nodes, _, errOver := s.nodeService.Overview(request)
	return nodes, errOver
}

func (s *Service) getOverviewAuto(connections []node.Connection, cluster ClusterOptions) ([]node.KeeperResponse, *node.Connection, error) {
	var certs *cert.Certs
	// NOTE: we want to rewrite `nil` only if tls is enabled
	if cluster.Tls.Keeper {
		certs = &cluster.Certs
	}
	request := node.KeeperAutoRequest{
		Connections: connections,
		VaultId:     cluster.Vaults.KeeperId,
		Certs:       certs,
	}
	nodes, _, detected, err := s.nodeService.OverviewAuto(request)
	return nodes, detected, err
}

func (s *Service) getSupportedFeatures(k keeper.Type, db database.Type) []features.Feature {
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
