package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"sort"

	"github.com/boltdb/bolt"
)

func NewStorage(name string) *bolt.DB {
	path := "data/bolt"
	errMk := os.MkdirAll(path, os.ModePerm)
	if errMk != nil {
		panic(errMk)
	}
	db, errOpen := bolt.Open(path+"/"+name, 0600, nil)
	if errOpen != nil {
		panic(errOpen)
	}
	return db
}

type Bucket[T any] struct {
	storage *bolt.DB
	name    []byte
}

func NewBucket[T any](storage *bolt.DB, name string) *Bucket[T] {
	byteName := []byte(name)
	err := storage.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(byteName)
		return err
	})
	if err != nil {
		panic(err)
	}
	return &Bucket[T]{
		storage: storage,
		name:    byteName,
	}
}

func (b *Bucket[T]) GetList(f func(el T) bool, s func(list []T, i, j int) bool) ([]T, error) {
	result := make([]T, 0)
	err := b.storage.View(func(tx *bolt.Tx) error {
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
	err := b.storage.View(func(tx *bolt.Tx) error {
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
	err := b.storage.View(func(tx *bolt.Tx) error {
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
	err := b.storage.View(func(tx *bolt.Tx) error {
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

func (b *Bucket[T]) Create(key string, value T) (T, error) {
	if key == "" {
		return value, errors.New("element identifier cannot be empty")
	}
	err := b.storage.Update(func(tx *bolt.Tx) error {
		el := tx.Bucket(b.name).Get([]byte(key))
		if el != nil {
			return errors.New("such an element already exists")
		}
		var buff bytes.Buffer
		err := gob.NewEncoder(&buff).Encode(value)
		if err != nil {
			return err
		}
		return tx.Bucket(b.name).Put([]byte(key), buff.Bytes())
	})
	return value, err
}

func (b *Bucket[T]) Update(key string, value T) error {
	if key == "" {
		return errors.New("element identifier cannot be empty")
	}
	return b.storage.Update(func(tx *bolt.Tx) error {
		var buff bytes.Buffer
		err := gob.NewEncoder(&buff).Encode(value)
		if err != nil {
			return err
		}
		return tx.Bucket(b.name).Put([]byte(key), buff.Bytes())
	})
}

func (b *Bucket[T]) Delete(key string) error {
	return b.storage.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.name).Delete([]byte(key))
	})
}

func (b *Bucket[T]) DeleteAll() error {
	return b.storage.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(b.name)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket(b.name)
		return err
	})
}
