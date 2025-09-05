package tag

import (
	"ivory/src/storage/bolt"
)

type TagRepository struct {
	bucket *bolt.Bucket[[]string]
}

func NewTagRepository(bucket *bolt.Bucket[[]string]) *TagRepository {
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

func (t *TagRepository) Update(key string, value []string) error {
	return t.bucket.Update(key, value)
}

func (t *TagRepository) Delete(tag string) error {
	return t.bucket.Delete(tag)
}

func (t *TagRepository) DeleteAll() error {
	return t.bucket.DeleteAll()
}
