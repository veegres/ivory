package database

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
)

var DB *bolt.DB
var ClusterBucketName = []byte("Cluster")

func Start() {
	DB = setupDatabase("bolt", "cluster.db")
}

func ClusterRepository() ClusterRepositoryStruct {
	return ClusterRepositoryStruct{
		Delete: deleteCluster,
		Update: updateCluster,
		Get:    getCluster,
		List:   getClusterList,
	}
}

func setupDatabase(dirName string, dbName string) *bolt.DB {
	_ = os.Mkdir(dirName, os.ModePerm)
	db, err := bolt.Open(dirName+"/"+dbName, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	_ = db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucket(ClusterBucketName)
		return nil
	})

	return db
}
