package cluster

import (
	"ivory/src/clients/sidecar"
	"ivory/src/features/cert"
	"ivory/src/features/instance"
	"ivory/src/features/tag"
)

type Service struct {
	clusterRepository *Repository
	instanceService   *instance.Service
	tagService        *tag.Service
}

func NewService(
	clusterRepository *Repository,
	instanceService *instance.Service,
	tagService *tag.Service,
) *Service {
	return &Service{
		clusterRepository: clusterRepository,
		instanceService:   instanceService,
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
	tags, err := s.saveTags(cluster.Name, cluster.Tags)
	if err != nil {
		return nil, err
	}
	cluster.Tags = tags
	errCluster := s.clusterRepository.Update(cluster)
	return &cluster, errCluster
}

func (s *Service) CreateAuto(cluster ClusterAuto) (Cluster, error) {
	var requestTls *cert.Certs
	if cluster.Tls.Sidecar {
		// NOTE: we want to rewrite `nil` only if tls is enabled
		requestTls = &cluster.Certs
	}
	request := instance.InstanceRequest{
		Sidecar:      cluster.Instance,
		CredentialId: cluster.Credentials.PatroniId,
		Certs:        requestTls,
	}
	overview, _, errOver := s.instanceService.Overview(request)
	if errOver != nil {
		return Cluster{}, errOver
	}

	instances := make([]sidecar.Sidecar, 0)
	for _, item := range overview {
		instances = append(instances, item.Sidecar)
	}

	tags, errSave := s.saveTags(cluster.Name, cluster.Tags)
	if errSave != nil {
		return Cluster{}, errSave
	}
	cluster.Tags = tags

	model := Cluster{
		Name:      cluster.Name,
		Instances: instances,
		ClusterOptions: ClusterOptions{
			Tls:         cluster.Tls,
			Certs:       cluster.Certs,
			Credentials: cluster.Credentials,
			Tags:        tags,
		},
	}

	return s.clusterRepository.Create(model)
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
