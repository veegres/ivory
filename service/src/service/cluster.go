package service

import (
	"errors"
	. "ivory/src/model"
	"ivory/src/persistence"
)

type ClusterService struct {
	clusterRepository *persistence.ClusterRepository
	passwordService   *PasswordService
	instanceService   InstanceService
}

func NewClusterService(
	clusterRepository *persistence.ClusterRepository,
	passwordService *PasswordService,
	instanceService InstanceService,
) *ClusterService {
	return &ClusterService{
		clusterRepository: clusterRepository,
		passwordService:   passwordService,
		instanceService:   instanceService,
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

func (s *ClusterService) GetPatroniCredential(name string) (*Credential, error) {
	cluster, errCluster := s.Get(name)
	if errCluster != nil {
		return nil, errCluster
	}
	if cluster.Credentials.PatroniId == nil {
		return nil, errors.New("there is no password for this cluster")
	}
	cred, errCred := s.passwordService.GetDecrypted(*cluster.Credentials.PatroniId)
	if errCred != nil {
		return nil, errCred
	}
	return cred, nil
}

func (s *ClusterService) GetPostgresCredential(name string) (*Credential, error) {
	cluster, errCluster := s.Get(name)
	if errCluster != nil {
		return nil, errCluster
	}
	if cluster.Credentials.PostgresId == nil {
		return nil, errors.New("there is no password for this cluster")
	}
	cred, errCred := s.passwordService.GetDecrypted(*cluster.Credentials.PostgresId)
	if errCred != nil {
		return nil, errCred
	}
	return cred, nil
}

func (s *ClusterService) Update(cluster ClusterModel) error {
	return s.clusterRepository.Update(cluster)
}

func (s *ClusterService) AddAuto(cluster ClusterAutoModel) ClusterModel {
	return ClusterModel{}
}

func (s *ClusterService) Delete(key string) error {
	return s.clusterRepository.Delete(key)
}

func (s *ClusterService) DeleteAll() error {
	return s.clusterRepository.DeleteAll()
}
