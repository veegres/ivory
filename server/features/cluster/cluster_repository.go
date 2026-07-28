package cluster

import (
	"ivory/clients/storage"
)

type Repository struct {
	bucket *storage.DbBucket[Response]
}

func NewRepository(bucket *storage.DbBucket[Response]) *Repository {
	return &Repository{
		bucket: bucket,
	}
}

func (r *Repository) List() ([]Response, error) {
	return r.bucket.GetList(nil, nil)
}

func (r *Repository) Search(criteria SearchCriteria) ([]Response, error) {
	var names map[string]bool
	if criteria.Names != nil {
		names = make(map[string]bool, len(criteria.Names))
		for _, name := range criteria.Names {
			names[name] = true
		}
	}
	return r.bucket.GetList(func(c Response) bool {
		if names != nil && !names[c.Name] {
			return false
		}
		if criteria.Keeper != nil && c.Plugins.Keeper != *criteria.Keeper {
			return false
		}
		if criteria.Database != nil && c.Plugins.Database != *criteria.Database {
			return false
		}
		return true
	}, nil)
}

func (r *Repository) Get(key string) (Response, error) {
	return r.bucket.Get(key)
}

func (r *Repository) Update(cluster Request) error {
	return r.bucket.Update(cluster.Name, Response(cluster))
}

func (r *Repository) Create(cluster Request) (Response, error) {
	return r.bucket.Create(cluster.Name, Response(cluster))
}

func (r *Repository) Delete(key string) error {
	return r.bucket.Delete(key)
}

func (r *Repository) DeleteAll() error {
	return r.bucket.DeleteAll()
}
