package cluster

import (
	"ivory/src/storage/db"
)

type Repository struct {
	bucket *db.Bucket[Response]
}

func NewRepository(bucket *db.Bucket[Response]) *Repository {
	return &Repository{
		bucket: bucket,
	}
}

func (r *Repository) List() ([]Response, error) {
	return r.bucket.GetList(nil, nil)
}

func (r *Repository) ListByName(clusters []string) ([]Response, error) {
	clusterMap := make(map[string]bool)
	for _, c := range clusters {
		clusterMap[c] = true
	}
	return r.bucket.GetList(func(cert Response) bool {
		return clusterMap[cert.Name]
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
