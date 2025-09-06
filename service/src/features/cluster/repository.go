package cluster

import (
	"ivory/src/storage/db"
)

type Repository struct {
	bucket *db.Bucket[Cluster]
}

func NewRepository(bucket *db.Bucket[Cluster]) *Repository {
	return &Repository{
		bucket: bucket,
	}
}

func (r *Repository) List() ([]Cluster, error) {
	return r.bucket.GetList(nil, nil)
}

func (r *Repository) ListByName(clusters []string) ([]Cluster, error) {
	clusterMap := make(map[string]bool)
	for _, c := range clusters {
		clusterMap[c] = true
	}
	return r.bucket.GetList(func(cert Cluster) bool {
		return clusterMap[cert.Name]
	}, nil)
}

func (r *Repository) Get(key string) (Cluster, error) {
	return r.bucket.Get(key)
}

func (r *Repository) Update(cluster Cluster) error {
	return r.bucket.Update(cluster.Name, cluster)
}

func (r *Repository) Create(cluster Cluster) (Cluster, error) {
	return r.bucket.Create(cluster.Name, cluster)
}

func (r *Repository) Delete(key string) error {
	return r.bucket.Delete(key)
}

func (r *Repository) DeleteAll() error {
	return r.bucket.DeleteAll()
}
