package persistence

import (
	"bytes"
	"encoding/gob"
	. "ivory/model"
)

type ClusterRepository struct {
	common common
	bucket []byte
}

func (r ClusterRepository) List() ([]ClusterModel, error) {
	bytesList, err := r.common.getList(r.bucket)
	modelList := make([]ClusterModel, len(bytesList))
	for i, el := range bytesList {
		var cluster ClusterModel
		buff := bytes.NewBuffer(el.value)
		_ = gob.NewDecoder(buff).Decode(&cluster)
		modelList[i] = cluster
	}
	return modelList, err
}

func (r ClusterRepository) ListByName(clusters []string) ([]ClusterModel, error) {
	bytesList, err := r.common.getList(r.bucket)
	modelList := make([]ClusterModel, 0)
	for _, el := range bytesList {
		var cluster ClusterModel
		for _, c := range clusters {
			if c == el.key {
				buff := bytes.NewBuffer(el.value)
				_ = gob.NewDecoder(buff).Decode(&cluster)
				modelList = append(modelList, cluster)
			}
		}
	}
	return modelList, err
}

func (r ClusterRepository) Get(key string) (ClusterModel, error) {
	value, err := r.common.get(r.bucket, key)
	var cluster ClusterModel
	buff := bytes.NewBuffer(value)
	_ = gob.NewDecoder(buff).Decode(&cluster)
	return cluster, err
}

func (r ClusterRepository) Update(cluster ClusterModel) error {
	return r.common.update(r.bucket, cluster.Name, cluster)
}

func (r ClusterRepository) Delete(key string) error {
	return r.common.delete(r.bucket, key)
}

func (r ClusterRepository) DeleteAll() error {
	return r.common.deleteAll(r.bucket)
}
