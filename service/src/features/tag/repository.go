package tag

import (
	"ivory/src/storage/db"
)

type Repository struct {
	bucket *db.Bucket[[]string]
}

func NewRepository(bucket *db.Bucket[[]string]) *Repository {
	return &Repository{
		bucket: bucket,
	}
}

func (t *Repository) List() ([]string, error) {
	return t.bucket.GetKeyList()
}

func (t *Repository) Get(tag string) ([]string, error) {
	return t.bucket.Get(tag)
}

func (t *Repository) GetMap() (map[string][]string, error) {
	return t.bucket.GetMap(nil)
}

func (t *Repository) Update(key string, value []string) error {
	return t.bucket.Update(key, value)
}

func (t *Repository) Delete(tag string) error {
	return t.bucket.Delete(tag)
}

func (t *Repository) DeleteAll() error {
	return t.bucket.DeleteAll()
}
