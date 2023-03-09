package config

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"sort"
)

func NewBoltDB(name string) *bolt.DB {
	path := "data/bolt"
	err := os.MkdirAll(path, os.ModePerm)
	db, err := bolt.Open(path+"/"+name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type Bucket[T any] struct {
	db   *bolt.DB
	name []byte
}

func NewBoltBucket[T any](db *bolt.DB, name string) *Bucket[T] {
	byteName := []byte(name)
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(byteName)
		return err
	})
	if err != nil {
		panic(err)
	}
	return &Bucket[T]{
		db:   db,
		name: byteName,
	}
}

func (b *Bucket[T]) GetList(f func(el T) bool, s func(list []T, i, j int) bool) ([]T, error) {
	result := make([]T, 0)
	err := b.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(b.name).Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			var el T
			buff := bytes.NewBuffer(value)
			err := gob.NewDecoder(buff).Decode(&el)
			if err != nil {
				return err
			}
			if f == nil || f(el) {
				result = append(result, el)
			}
		}
		return nil
	})
	if s != nil {
		sort.Slice(result, func(i, j int) bool {
			return s(result, i, j)
		})
	}
	return result, err
}

func (b *Bucket[T]) GetKeyList() ([]string, error) {
	result := make([]string, 0)
	err := b.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(b.name).Cursor()
		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			result = append(result, string(key))
		}
		return nil
	})
	return result, err
}

func (b *Bucket[T]) GetMap(filter func(el T) bool) (map[string]T, error) {
	result := make(map[string]T)
	err := b.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(b.name).Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			var el T
			buff := bytes.NewBuffer(value)
			err := gob.NewDecoder(buff).Decode(&el)
			if err != nil {
				return err
			}
			if filter == nil || filter(el) {
				result[string(key)] = el
			}
		}
		return nil
	})
	return result, err
}

func (b *Bucket[T]) Get(key string) (T, error) {
	var value T
	err := b.db.View(func(tx *bolt.Tx) error {
		el := tx.Bucket(b.name).Get([]byte(key))
		if el == nil {
			return errors.New("element doesn't exist")
		}
		buff := bytes.NewBuffer(el)
		err := gob.NewDecoder(buff).Decode(&value)
		if err != nil {
			return err
		}
		return nil
	})
	return value, err
}

func (b *Bucket[T]) Update(key string, value T) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		var buff bytes.Buffer
		err := gob.NewEncoder(&buff).Encode(value)
		if err != nil {
			return err
		}
		return tx.Bucket(b.name).Put([]byte(key), buff.Bytes())
	})
}

func (b *Bucket[T]) Delete(key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.name).Delete([]byte(key))
	})
}

func (b *Bucket[T]) DeleteAll() error {
	return b.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(b.name)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket(b.name)
		return err
	})
}
