package database

import (
	"bytes"
	"encoding/gob"
	"github.com/boltdb/bolt"
	. "ivory/model"
)

type ClusterRepositoryStruct struct {
	Delete func(key string)
	Update func(cluster Cluster)
	Get    func(key string) *Cluster
	List   func() []Cluster
}

func deleteCluster(key string) {
	_ = DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(ClusterBucketName).Delete([]byte(key))
	})
}

func updateCluster(cluster Cluster) {
	_ = DB.Update(func(tx *bolt.Tx) error {
		var buff bytes.Buffer
		_ = gob.NewEncoder(&buff).Encode(cluster.Nodes)
		return tx.Bucket(ClusterBucketName).Put([]byte(cluster.Name), buff.Bytes())
	})
}

func getCluster(key string) *Cluster {
	var cluster *Cluster
	_ = DB.View(func(tx *bolt.Tx) error {
		value := tx.Bucket(ClusterBucketName).Get([]byte(key))
		if value == nil {
			return nil
		}
		buff := bytes.NewBuffer(value)
		var nodes []string
		_ = gob.NewDecoder(buff).Decode(&nodes)
		tmp := Cluster{Name: key, Nodes: nodes}
		cluster = &tmp
		return nil
	})
	return cluster
}

func getClusterList() []Cluster {
	list := make([]Cluster, 0)
	err := DB.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(ClusterBucketName).Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			buff := bytes.NewBuffer(value)
			var nodes []string
			_ = gob.NewDecoder(buff).Decode(&nodes)
			list = append(list, Cluster{Name: string(key), Nodes: nodes})
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return list
}
