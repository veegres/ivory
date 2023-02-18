package persistence

import (
	. "ivory/model"
)

type ClusterRepository struct {
	common common
	bucket []byte
}

func (r ClusterRepository) List() ([]ClusterModel, error) {
	return GetList[ClusterModel](r.bucket, nil)
}

func (r ClusterRepository) ListByName(clusters []string) ([]ClusterModel, error) {
	clusterMap := make(map[string]bool)
	for _, c := range clusters {
		clusterMap[c] = true
	}
	return GetList[ClusterModel](r.bucket, func(cert ClusterModel) bool {
		return clusterMap[cert.Name]
	})
}

func (r ClusterRepository) Get(key string) (ClusterModel, error) {
	return Get[ClusterModel](r.bucket, key)
}

func (r ClusterRepository) Update(cluster ClusterModel) error {
	return Update(r.bucket, cluster.Name, cluster)
}

func (r ClusterRepository) Delete(key string) error {
	return Delete(r.bucket, key)
}

func (r ClusterRepository) DeleteAll() error {
	return DeleteAll(r.bucket)
}
