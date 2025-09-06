package cluster

import (
	"ivory/src/storage/db"
)

type ClusterRepository struct {
	bucket *db.Bucket[Cluster]
}

func NewClusterRepository(bucket *db.Bucket[Cluster]) *ClusterRepository {
	return &ClusterRepository{
		bucket: bucket,
	}
}

func (r *ClusterRepository) List() ([]Cluster, error) {
	return r.bucket.GetList(nil, nil)
}

func (r *ClusterRepository) ListByName(clusters []string) ([]Cluster, error) {
	clusterMap := make(map[string]bool)
	for _, c := range clusters {
		clusterMap[c] = true
	}
	return r.bucket.GetList(func(cert Cluster) bool {
		return clusterMap[cert.Name]
	}, nil)
}

func (r *ClusterRepository) Get(key string) (Cluster, error) {
	return r.bucket.Get(key)
}

func (r *ClusterRepository) Update(cluster Cluster) error {
	return r.bucket.Update(cluster.Name, cluster)
}

func (r *ClusterRepository) Create(cluster Cluster) (Cluster, error) {
	return r.bucket.Create(cluster.Name, cluster)
}

func (r *ClusterRepository) Delete(key string) error {
	return r.bucket.Delete(key)
}

func (r *ClusterRepository) DeleteAll() error {
	return r.bucket.DeleteAll()
}
