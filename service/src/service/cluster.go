package service

import (
	. "ivory/src/model"
	"ivory/src/persistence"
)

type ClusterService struct {
	clusterRepository *persistence.ClusterRepository
	tagService        *TagService
	instanceService   InstanceGateway
}

func NewClusterService(
	clusterRepository *persistence.ClusterRepository,
	tagService *TagService,
	instanceService InstanceGateway,
) *ClusterService {
	return &ClusterService{
		clusterRepository: clusterRepository,
		tagService:        tagService,
		instanceService:   instanceService,
	}
}

func (s *ClusterService) List() ([]Cluster, error) {
	return s.clusterRepository.List()
}

func (s *ClusterService) ListByTag(tags []string) ([]Cluster, error) {
	listMap := make(map[string]bool)
	for _, tag := range tags {
		// NOTE: we shouldn't check error here, we want to return empty array if there is no such tag
		clusters, _ := s.tagService.Get(tag)

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

func (s *ClusterService) ListByName(clusters []string) ([]Cluster, error) {
	return s.clusterRepository.ListByName(clusters)
}

func (s *ClusterService) Get(cluster string) (Cluster, error) {
	return s.clusterRepository.Get(cluster)
}

func (s *ClusterService) Update(cluster Cluster) (*Cluster, error) {
	tags, err := s.saveTags(cluster.Name, cluster.Tags)
	if err != nil {
		return nil, err
	}
	cluster.Tags = tags
	errCluster := s.clusterRepository.Update(cluster)
	return &cluster, errCluster
}

func (s *ClusterService) CreateAuto(cluster ClusterAuto) (Cluster, error) {
	request := InstanceRequest{
		Sidecar:      cluster.Instance,
		CredentialId: cluster.Credentials.PatroniId,
		Certs:        cluster.Certs,
	}
	overview, _, errOver := s.instanceService.Overview(request)
	if errOver != nil {
		return Cluster{}, errOver
	}

	instances := make([]Sidecar, 0)
	for _, instance := range overview {
		instances = append(instances, instance.Sidecar)
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
			Certs:       cluster.Certs,
			Credentials: cluster.Credentials,
			Tags:        tags,
		},
	}

	return s.clusterRepository.Create(model)
}

func (s *ClusterService) Delete(cluster string) error {
	_, errTag := s.tagService.UpdateCluster(cluster, nil)
	if errTag != nil {
		return errTag
	}
	return s.clusterRepository.Delete(cluster)
}

func (s *ClusterService) DeleteAll() error {
	return s.clusterRepository.DeleteAll()
}

func (s *ClusterService) saveTags(name string, tags []string) ([]string, error) {
	// NOTE: remove duplicates
	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[tag] = true
	}
	tagList := make([]string, 0)
	for key := range tagMap {
		tagList = append(tagList, key)
	}

	// NOTE: create tags in db with cluster name
	return s.tagService.UpdateCluster(name, tagList)
}
