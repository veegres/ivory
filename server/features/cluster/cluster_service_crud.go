package cluster

import (
	"ivory/features/node"
	"ivory/features/query"
)

func (s *Service) List() ([]Response, error) {
	return s.clusterRepository.List()
}

// SearchRequest narrows Search down to matching clusters. Tags is resolved
// to cluster names before hitting the repository; Keeper and Database are
// passed through as-is. A nil/empty field is skipped (no restriction).
type SearchRequest struct {
	Tags     []string
	Keeper   *node.KeeperPlugin
	Database *query.DbPlugin
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

func (s *Service) Update(cluster Request) (*Response, error) {
	if cluster.Name == "" {
		return nil, ErrClusterNameEmpty
	}
	if cluster.Nodes == nil || len(cluster.Nodes) == 0 {
		return nil, ErrClusterKeepersEmpty
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
