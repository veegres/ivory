package service

import (
	. "ivory/src/model"
	"ivory/src/persistence"
)

type ClusterService struct {
	clusterRepository *persistence.ClusterRepository
}

func NewClusterService(repository *persistence.ClusterRepository) *ClusterService {
	return &ClusterService{
		clusterRepository: repository,
	}
}

func (s *ClusterService) List() ([]ClusterModel, error) {
	return s.clusterRepository.List()
}

func (s *ClusterService) ListByName(clusters []string) ([]ClusterModel, error) {
	return s.clusterRepository.ListByName(clusters)
}

func (s *ClusterService) Get(key string) (ClusterModel, error) {
	return s.clusterRepository.Get(key)
}

func (s *ClusterService) Update(cluster ClusterModel) error {
	return s.clusterRepository.Update(cluster)
}

func (s *ClusterService) Delete(key string) error {
	return s.clusterRepository.Delete(key)
}

func (s *ClusterService) DeleteAll() error {
	return s.clusterRepository.DeleteAll()
}
