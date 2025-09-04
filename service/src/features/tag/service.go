package tag

import (
	"strings"
)

type TagService struct {
	tagRepository *TagRepository
}

func NewTagService(tagRepository *TagRepository) *TagService {
	return &TagService{tagRepository: tagRepository}
}

func (s *TagService) Get(tag string) ([]string, error) {
	return s.tagRepository.Get(tag)
}

func (s *TagService) GetMap() (map[string][]string, error) {
	return s.tagRepository.GetMap()
}

func (s *TagService) List() ([]string, error) {
	return s.tagRepository.List()
}

func (s *TagService) UpdateCluster(cluster string, tags []string) ([]string, error) {
	var tagsLower []string
	for _, tag := range tags {
		tagsLower = append(tagsLower, strings.ToLower(tag))
	}

	tagMap, err := s.tagRepository.GetMap()
	if err != nil {
		return nil, err
	}

	// NOTE: remove cluster from all tags
	for k, v := range tagMap {
		tmp := make([]string, 0)
		for _, c := range v {
			if c != cluster {
				tmp = append(tmp, c)
			}
		}

		if len(tmp) == 0 {
			tagMap[k] = nil
		} else {
			tagMap[k] = tmp
		}
	}

	// NOTE: add cluster to tags
	for _, v := range tagsLower {
		tagMap[v] = append(tagMap[v], cluster)
	}

	// NOTE: update tags
	for k, v := range tagMap {
		if v != nil {
			err := s.tagRepository.Update(k, v)
			if err != nil {
				return nil, err
			}
		} else {
			err := s.tagRepository.Delete(k)
			if err != nil {
				return nil, err
			}
		}
	}

	return tagsLower, nil
}

func (s *TagService) Delete(tag string) error {
	tagLower := strings.ToLower(tag)
	return s.tagRepository.Delete(tagLower)
}

func (s *TagService) DeleteAll() error {
	return s.tagRepository.DeleteAll()
}
