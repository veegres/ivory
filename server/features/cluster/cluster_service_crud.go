package cluster

import (
	"errors"
	"fmt"
	"ivory/clients/storage"
)

func (s *Service) List() ([]Response, error) {
	return s.clusterRepository.List()
}

func (s *Service) Search(request SearchRequest) ([]Response, error) {
	criteria := SearchCriteria{Keeper: request.Keeper, Database: request.Database}
	if len(request.Tags) > 0 {
		criteria.Names = s.resolveTagNames(request.Tags)
	}
	return s.clusterRepository.Search(criteria)
}

func (s *Service) resolveTagNames(tags []string) []string {
	listMap := make(map[string]bool)
	for _, t := range tags {
		// NOTE: we shouldn't check the error here, we want to return an empty array if there is no such tag
		clusters, _ := s.tagService.Get(t)
		for _, c := range clusters {
			listMap[c] = true
		}
	}

	listName := make([]string, 0, len(listMap))
	for k := range listMap {
		listName = append(listName, k)
	}
	return listName
}

func (s *Service) Get(cluster string) (Response, error) {
	return s.clusterRepository.Get(cluster)
}

func (s *Service) Create(cluster Request) (*Response, error) {
	if err := s.validateWritableCluster(cluster); err != nil {
		return nil, err
	}
	if _, err := s.Get(cluster.Name); err == nil {
		return nil, ErrClusterNameTaken
	} else if !errors.Is(err, storage.ErrNotFound) {
		return nil, err
	}
	tags, err := s.saveTags(cluster.Name, cluster.Tags)
	if err != nil {
		return nil, err
	}
	cluster.Tags = tags
	created, errCluster := s.clusterRepository.Create(cluster)
	if errors.Is(errCluster, storage.ErrAlreadyExists) {
		return nil, ErrClusterNameTaken
	}
	return &created, errCluster
}

func (s *Service) Update(cluster Request) (*Response, error) {
	if err := s.validateWritableCluster(cluster); err != nil {
		return nil, err
	}
	tags, err := s.saveTags(cluster.Name, cluster.Tags)
	if err != nil {
		return nil, err
	}
	cluster.Tags = tags
	errCluster := s.clusterRepository.Update(cluster)
	return (*Response)(&cluster), errCluster
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

func (s *Service) validateWritableCluster(cluster Request) error {
	if cluster.Name == "" {
		return ErrClusterNameEmpty
	}
	if cluster.Nodes == nil || len(cluster.Nodes) == 0 {
		return ErrClusterKeepersEmpty
	}
	if err := s.validateNodeNames(cluster.Nodes); err != nil {
		return err
	}
	return nil
}

// validateNodeNames enforces the naming rules the cluster depends on: a
// node's name identifies its deployment on the platform, so it must exist
// and be unique within the cluster.
func (s *Service) validateNodeNames(nodes []NodeConfig) error {
	seen := make(map[string]bool, len(nodes))
	for _, n := range nodes {
		if n.Name == "" {
			return ErrClusterNodeNameNotProvided
		}
		if seen[n.Name] {
			return fmt.Errorf("%w: %s", ErrClusterNodeNameNotUnique, n.Name)
		}
		seen[n.Name] = true
	}
	return nil
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
