package persistence

import (
	"ivory/src/config"
)

type TagRepository struct {
	bucket *config.Bucket[[]string]
}

func NewTagRepository(bucket *config.Bucket[[]string]) *TagRepository {
	return &TagRepository{
		bucket: bucket,
	}
}

func (t *TagRepository) List() ([]string, error) {
	return t.bucket.GetKeyList()
}

func (t *TagRepository) Get(tag string) ([]string, error) {
	return t.bucket.Get(tag)
}

func (t *TagRepository) GetMap() (map[string][]string, error) {
	return t.bucket.GetMap(nil)
}

func (t *TagRepository) UpdateCluster(cluster string, tags []string) error {
	tagMap, err := t.GetMap()
	if err != nil {
		return err
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
	for _, v := range tags {
		tagMap[v] = append(tagMap[v], cluster)
	}

	// NOTE: update tags
	for k, v := range tagMap {
		if v != nil {
			err := t.bucket.Update(k, v)
			if err != nil {
				return err
			}
		} else {
			err := t.bucket.Delete(k)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TagRepository) Delete(tag string) error {
	return t.bucket.Delete(tag)
}

func (t *TagRepository) DeleteAll() error {
	return t.bucket.DeleteAll()
}
