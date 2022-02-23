package persistence

import (
	"bytes"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

var Database database

type database struct {
	Cluster      *ClusterRepository
	CompactTable *CompactTableRepository
}

type common struct {
	bolt *bolt.DB
}

type Element struct {
	key   string
	value []byte
}

func (d database) Build(dbName string) {
	dirName := "data/bolt"
	_ = os.MkdirAll(dirName, os.ModePerm)
	db, err := bolt.Open(dirName+"/"+dbName, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	clusterName := []byte("Cluster")
	_ = db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucket(clusterName)
		return nil
	})

	compactTableName := []byte("CompactTable")
	_ = db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucket(compactTableName)
		return nil
	})

	common := common{bolt: db}
	d.Cluster = &ClusterRepository{
		common: common,
		bucket: clusterName,
	}
	d.CompactTable = &CompactTableRepository{
		common: common,
		bucket: compactTableName,
	}
	Database = d
}

func (d common) getList(bucket []byte) ([]Element, error) {
	list := make([]Element, 0)
	err := d.bolt.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(bucket).Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			list = append(list, Element{string(key), value})
		}
		return nil
	})
	return list, err
}

func (d common) get(bucket []byte, key string) ([]byte, error) {
	var value []byte
	err := d.bolt.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(bucket).Get([]byte(key))
		return nil
	})
	return value, err
}

func (d common) delete(bucket []byte, key string) error {
	return d.bolt.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucket).Delete([]byte(key))
	})
}

func (d common) update(bucket []byte, key string, value interface{}) error {
	return d.bolt.Update(func(tx *bolt.Tx) error {
		var buff bytes.Buffer
		_ = gob.NewEncoder(&buff).Encode(value)
		return tx.Bucket(bucket).Put([]byte(key), buff.Bytes())
	})
}
