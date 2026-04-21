package cluster

import (
	"errors"
	"fmt"
	"ivory/src/clients/sidecar"
	"ivory/src/features/cert"
	"ivory/src/features/node"
	"ivory/src/features/tag"
)

var ErrClusterNameEmpty = errors.New("cluster name cannot be empty")
var ErrClusterSidecarsEmpty = errors.New("cluster sidecars cannot be empty")

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
	if cluster.Sidecars == nil {
		return nil, ErrClusterSidecarsEmpty
	}
	tags, err := s.saveTags(cluster.Name, cluster.Tags)
	if err != nil {
		return nil, err
	}
	cluster.Tags = tags
	errCluster := s.clusterRepository.Update(cluster)
	return &cluster, errCluster
}

func (s *Service) Overview(name string, side *sidecar.Sidecar) (*ClusterOverview, error) {
	cluster, clusterError := s.Get(name)
	if clusterError != nil {
		return nil, clusterError
	}
	detectedBy := side
	var nodes []sidecar.Instance
	var err error
	if side == nil {
		nodes, detectedBy, err = s.getOverviewAuto(cluster.Sidecars, cluster.ClusterOptions)
	} else {
		nodes, err = s.getOverview(*side, cluster.ClusterOptions)
	}
	if err != nil {
		return nil, err
	}

	detectedByDomain := fmt.Sprintf("%s:%d", detectedBy.Host, detectedBy.Port)
	var mainNode *Node
	nodesMap := make(map[string]*Node)
	for _, el := range nodes {
		domain := fmt.Sprintf("%s:%d", el.Sidecar.Host, el.Sidecar.Port)
		fullNode := Node{el, false, true}
		nodesMap[domain] = &fullNode

		if detectedByDomain == domain {
			mainNode = &fullNode
		}
		if fullNode.Role == sidecar.Leader {
			mainNode = &fullNode
		}
	}
	for _, sc := range cluster.Sidecars {
		domain := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
		if value, ok := nodesMap[domain]; ok && value != nil {
			nodesMap[domain] = &Node{value.Instance, true, value.InSidecar}
		} else {
			nodesMap[domain] = nil
		}
	}
	return &ClusterOverview{nodesMap, detectedBy, mainNode}, nil
}

func (s *Service) CreateAuto(cluster ClusterAuto) (Cluster, error) {
	overview, errOver := s.getOverview(cluster.Node, cluster.ClusterOptions)
	if errOver != nil {
		return Cluster{}, errOver
	}

	sidecars := make([]sidecar.Sidecar, 0)
	for _, item := range overview {
		sidecars = append(sidecars, item.Sidecar)
	}

	tags, errSave := s.saveTags(cluster.Name, cluster.Tags)
	if errSave != nil {
		return Cluster{}, errSave
	}
	cluster.Tags = tags

	model := Cluster{
		Name:     cluster.Name,
		Sidecars: sidecars,
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
	overview, _, err := s.getOverviewAuto(cluster.Sidecars, cluster.ClusterOptions)
	if err != nil {
		return nil, err
	}

	sidecars := make([]sidecar.Sidecar, 0)
	for _, item := range overview {
		sidecars = append(sidecars, item.Sidecar)
	}

	model := Cluster{
		Name:     cluster.Name,
		Sidecars: sidecars,
		ClusterOptions: ClusterOptions{
			Tls:         cluster.Tls,
			Certs:       cluster.Certs,
			Credentials: cluster.Credentials,
			Tags:        cluster.Tags,
		},
	}
	return &model, s.clusterRepository.Update(model)
}

func (s *Service) getOverview(sidecar sidecar.Sidecar, cluster ClusterOptions) ([]sidecar.Instance, error) {
	var certs *cert.Certs
	// NOTE: we want to rewrite `nil` only if tls is enabled
	if cluster.Tls.Sidecar {
		certs = &cluster.Certs
	}
	request := node.Request{
		Sidecar:      sidecar,
		CredentialId: cluster.Credentials.PatroniId,
		Certs:        certs,
	}
	overview, _, errOver := s.nodeService.Overview(request)
	return overview, errOver
}

func (s *Service) getOverviewAuto(sidecars []sidecar.Sidecar, cluster ClusterOptions) ([]sidecar.Instance, *sidecar.Sidecar, error) {
	var certs *cert.Certs
	// NOTE: we want to rewrite `nil` only if tls is enabled
	if cluster.Tls.Sidecar {
		certs = &cluster.Certs
	}
	request := node.NodeAutoRequest{
		Sidecars:     sidecars,
		CredentialId: cluster.Credentials.PatroniId,
		Certs:        certs,
	}
	overview, _, detectedBy, errOver := s.nodeService.OverviewAuto(request)
	if errOver != nil {
		return nil, nil, errOver
	}
	return overview, detectedBy, nil
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
