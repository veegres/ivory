package cluster

import (
	"errors"
	"fmt"
	"ivory/src/clients/keeper"
	"ivory/src/features/cert"
	"ivory/src/features/node"
	"ivory/src/features/tag"
)

var ErrClusterNameEmpty = errors.New("cluster name cannot be empty")
var ErrClusterKeepersEmpty = errors.New("cluster keepers cannot be empty")

type Service struct {
	clusterRepository *Repository
	nodeService       *node.Service
	tagService        *tag.Service
}

func NewService(
	clusterRepository *Repository,
	nodeService *node.Service,
	tagService *tag.Service,
) *Service {
	return &Service{
		clusterRepository: clusterRepository,
		nodeService:       nodeService,
		tagService:        tagService,
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

func (s *Service) Overview(name string, k *node.Connection) (*ClusterOverview, error) {
	cluster, clusterError := s.Get(name)
	if clusterError != nil {
		return nil, clusterError
	}
	var nodes []node.Node
	var detectedBy *keeper.Keeper
	var err error
	if k == nil {
		nodes, _, detectedBy, err = s.getOverviewAuto(cluster.Nodes, cluster.ClusterOptions)
	} else {
		nodes, err = s.getOverview(*k, cluster.ClusterOptions)
		detectedBy = &keeper.Keeper{Host: k.Host, Port: k.KeeperPort}
	}
	if err != nil {
		return nil, err
	}

	detectedByDomain := fmt.Sprintf("%s:%d", detectedBy.Host, detectedBy.Port)
	var mainNode *Node
	nodesMap := make(map[string]*Node)
	for _, el := range nodes {
		domain := fmt.Sprintf("%s:%d", el.Connection.Host, el.Connection.KeeperPort)
		fullNode := Node{el, false, true}
		nodesMap[domain] = &fullNode

		if detectedByDomain == domain {
			mainNode = &fullNode
		}
		if fullNode.Keeper.Role == keeper.Leader {
			mainNode = &fullNode
		}
	}
	for _, n := range cluster.Nodes {
		domain := fmt.Sprintf("%s:%d", n.Host, n.KeeperPort)
		if value, ok := nodesMap[domain]; ok && value != nil {
			nodesMap[domain] = &Node{value.Node, true, value.InKeeper}
		} else {
			nodesMap[domain] = nil
		}
	}
	return &ClusterOverview{nodesMap, detectedBy, mainNode}, nil
}

func (s *Service) CreateAuto(cluster ClusterAuto) (Cluster, error) {
	nodes, errOver := s.getOverview(cluster.Node, cluster.ClusterOptions)
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
			Tls:         cluster.Tls,
			Certs:       cluster.Certs,
			Credentials: cluster.Credentials,
			Tags:        tags,
		},
	}

	return s.clusterRepository.Create(model)
}

func (s *Service) FixAuto(name string) (*Cluster, error) {
	cluster, clusterError := s.Get(name)
	if clusterError != nil {
		return nil, clusterError
	}
	nodes, _, _, err := s.getOverviewAuto(cluster.Nodes, cluster.ClusterOptions)
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
			Tls:         cluster.Tls,
			Certs:       cluster.Certs,
			Credentials: cluster.Credentials,
			Tags:        cluster.Tags,
		},
	}
	return &model, s.clusterRepository.Update(model)
}

func (s *Service) getOverview(connection node.Connection, cluster ClusterOptions) ([]node.Node, error) {
	var certs *cert.Certs
	// NOTE: we want to rewrite `nil` only if tls is enabled
	if cluster.Tls.Keeper {
		certs = &cluster.Certs
	}
	request := node.Request{
		Connection:   connection,
		CredentialId: cluster.Credentials.PatroniId,
		Certs:        certs,
	}
	nodes, _, errOver := s.nodeService.Overview(request)
	return nodes, errOver
}

func (s *Service) getOverviewAuto(connections []node.Connection, cluster ClusterOptions) ([]node.Node, int, *keeper.Keeper, error) {
	var certs *cert.Certs
	// NOTE: we want to rewrite `nil` only if tls is enabled
	if cluster.Tls.Keeper {
		certs = &cluster.Certs
	}
	request := node.AutoRequest{
		Connections:  connections,
		CredentialId: cluster.Credentials.PatroniId,
		Certs:        certs,
	}
	return s.nodeService.OverviewAuto(request)
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
