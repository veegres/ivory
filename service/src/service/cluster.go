package service

import (
	. "ivory/src/model"
	"ivory/src/persistence"
)

type ClusterService struct {
	clusterRepository *persistence.ClusterRepository
	tagRepository     *persistence.TagRepository
	instanceService   InstanceService
}

func NewClusterService(
	clusterRepository *persistence.ClusterRepository,
	tagRepository *persistence.TagRepository,
	instanceService InstanceService,
) *ClusterService {
	return &ClusterService{
		clusterRepository: clusterRepository,
		tagRepository:     tagRepository,
		instanceService:   instanceService,
	}
}

func (s *ClusterService) List() ([]ClusterModel, error) {
	return s.clusterRepository.List()
}

func (s *ClusterService) ListByTag(tags []string) ([]ClusterModel, error) {
	listMap := make(map[string]bool)
	for _, tag := range tags {
		clusters, err := s.tagRepository.Get(tag)
		if err != nil {
			return nil, err
		}

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

func (s *ClusterService) ListByName(clusters []string) ([]ClusterModel, error) {
	return s.clusterRepository.ListByName(clusters)
}

func (s *ClusterService) Get(key string) (ClusterModel, error) {
	return s.clusterRepository.Get(key)
}

func (s *ClusterService) Update(cluster ClusterModel) (*ClusterModel, error) {
	// NOTE: remove duplicates
	tagMap := make(map[string]bool)
	for _, tag := range cluster.Tags {
		tagMap[tag] = true
	}
	tagList := make([]string, 0)
	for key := range tagMap {
		tagList = append(tagList, key)
	}

	// NOTE: create tags in db with cluster name
	err := s.tagRepository.UpdateCluster(cluster.Name, tagList)
	if err != nil {
		return nil, err
	}
	cluster.Tags = tagList

	errCluster := s.clusterRepository.Update(cluster)
	return &cluster, errCluster
}

func (s *ClusterService) CreateAuto(cluster ClusterAutoModel) (*ClusterModel, error) {
	request := InstanceRequest{
		Host:         cluster.Instance.Host,
		Port:         cluster.Instance.Port,
		CredentialId: cluster.Credentials.PatroniId,
		Certs:        cluster.Certs,
	}
	overview, _, err := s.instanceService.Overview(request)
	if err != nil {
		return nil, err
	}

	instances := make([]Sidecar, 0)
	for _, instance := range overview {
		instances = append(instances, instance.Sidecar)
	}

	model := ClusterModel{
		Name:        cluster.Name,
		Certs:       cluster.Certs,
		Credentials: cluster.Credentials,
		Tags:        cluster.Tags,
		Instances:   instances,
	}

	return s.Update(model)
}

func (s *ClusterService) Delete(key string) error {
	errTag := s.tagRepository.UpdateCluster(key, nil)
	if errTag != nil {
		return errTag
	}
	return s.clusterRepository.Delete(key)
}

func (s *ClusterService) DeleteAll() error {
	return s.clusterRepository.DeleteAll()
}
