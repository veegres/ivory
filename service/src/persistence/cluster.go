package persistence

import (
	"ivory/src/config"
	. "ivory/src/model"
)

type ClusterRepository struct {
	bucket *config.Bucket[ClusterModel]
}

func NewClusterRepository(bucket *config.Bucket[ClusterModel]) *ClusterRepository {
	return &ClusterRepository{
		bucket: bucket,
	}
}

func (r *ClusterRepository) List() ([]ClusterModel, error) {
	return r.bucket.GetList(nil, nil)
}

func (r *ClusterRepository) ListByName(clusters []string) ([]ClusterModel, error) {
	clusterMap := make(map[string]bool)
	for _, c := range clusters {
		clusterMap[c] = true
	}
	return r.bucket.GetList(func(cert ClusterModel) bool {
		return clusterMap[cert.Name]
	}, nil)
}

func (r *ClusterRepository) Get(key string) (ClusterModel, error) {
	return r.bucket.Get(key)
}

func (r *ClusterRepository) Update(cluster ClusterModel) error {
	return r.bucket.Update(cluster.Name, cluster)
}

func (r *ClusterRepository) Create(cluster ClusterModel) (ClusterModel, error) {
	return r.bucket.Create(cluster.Name, cluster)
}

func (r *ClusterRepository) Delete(key string) error {
	return r.bucket.Delete(key)
}

func (r *ClusterRepository) DeleteAll() error {
	return r.bucket.DeleteAll()
}
