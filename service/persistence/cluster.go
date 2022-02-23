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
		nodes := make([]string, 0)
		buff := bytes.NewBuffer(el.value)
		_ = gob.NewDecoder(buff).Decode(&nodes)
		modelList[i] = ClusterModel{Name: el.key, Nodes: nodes}
	}
	return modelList, err
}

func (r ClusterRepository) Get(key string) (ClusterModel, error) {
	value, err := r.common.get(r.bucket, key)
	var nodes []string
	buff := bytes.NewBuffer(value)
	_ = gob.NewDecoder(buff).Decode(&nodes)
	return ClusterModel{Name: key, Nodes: nodes}, err
}

func (r ClusterRepository) Update(cluster ClusterModel) error {
	return r.common.update(r.bucket, cluster.Name, cluster.Nodes)
}

func (r ClusterRepository) Delete(key string) error {
	return r.common.delete(r.bucket, key)
}
