package persistence

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

var BoltDB *database

type database struct {
	db           *bolt.DB
	Cluster      *ClusterRepository
	CompactTable *CompactTableRepository
	Credential   *CredentialRepository
	Cert         *CertRepository
	Tag          *TagRepository
}

func (d *database) Build(dbName string) {
	path := "data/bolt"
	_ = os.MkdirAll(path, os.ModePerm)
	db, err := bolt.Open(path+"/"+dbName, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	clusterName := []byte("Cluster")
	_ = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(clusterName)
		return err
	})

	compactTableName := []byte("CompactTable")
	_ = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(compactTableName)
		return err
	})

	certName := []byte("Cert")
	_ = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(certName)
		return err
	})

	tagName := []byte("Tag")
	_ = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(tagName)
		return err
	})

	credentialSecretName := []byte("Secret")
	credentialCredentialName := []byte("Credential")
	credentialEncryptedRef := "refEncrypted"
	credentialDecryptedRef := "refDecrypted"
	_ = db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucketIfNotExists(credentialCredentialName)
		secret, _ := tx.CreateBucketIfNotExists(credentialSecretName)
		encByte := []byte(credentialEncryptedRef)
		decByte := []byte(credentialDecryptedRef)
		if secret.Get(encByte) == nil {
			_ = secret.Put(encByte, []byte(""))
		}
		if secret.Get(decByte) == nil {
			_ = secret.Put(decByte, []byte(""))
		}
		return nil
	})

	common := common{bolt: db}
	BoltDB = &database{
		db: db,
		Cluster: &ClusterRepository{
			common: common,
			bucket: clusterName,
		},
		CompactTable: &CompactTableRepository{
			common: common,
			bucket: compactTableName,
		},
		Cert: &CertRepository{
			common: common,
			bucket: certName,
		},
		Tag: &TagRepository{
			common: common,
			bucket: tagName,
		},
		Credential: &CredentialRepository{
			common:           common,
			secretBucket:     credentialSecretName,
			credentialBucket: credentialCredentialName,
			encryptedRefKey:  credentialEncryptedRef,
			decryptedRefKey:  credentialDecryptedRef,
		},
	}
}

type common struct {
	bolt *bolt.DB
}

type element struct {
	key   string
	value []byte
}

func (d common) getList(bucket []byte) ([]element, error) {
	list := make([]element, 0)
	err := d.bolt.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(bucket).Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			list = append(list, element{string(key), value})
		}
		return nil
	})
	return list, err
}

func (d common) get(bucket []byte, key string) ([]byte, error) {
	var value []byte
	err := d.bolt.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(bucket).Get([]byte(key))
		if value == nil {
			return errors.New("no such element")
		} else {
			return nil
		}
	})
	return value, err
}

func (d common) delete(bucket []byte, key string) error {
	return d.bolt.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucket).Delete([]byte(key))
	})
}

func (d common) deleteAll(bucket []byte) error {
	return d.bolt.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucket)
		_, err = tx.CreateBucket(bucket)
		return err
	})
}

func (d common) update(bucket []byte, key string, value interface{}) error {
	return d.bolt.Update(func(tx *bolt.Tx) error {
		var buff bytes.Buffer
		_ = gob.NewEncoder(&buff).Encode(value)
		return tx.Bucket(bucket).Put([]byte(key), buff.Bytes())
	})
}

// TODO change all implementations to these generic methods

func GetList[T any](bucket []byte) ([]T, error) {
	result := make([]T, 0)
	err := BoltDB.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(bucket).Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			var el T
			buff := bytes.NewBuffer(value)
			err := gob.NewDecoder(buff).Decode(&el)
			if err != nil {
				return err
			}
			result = append(result, el)
		}
		return nil
	})
	return result, err
}

func GetKeyList(bucket []byte) ([]string, error) {
	result := make([]string, 0)
	err := BoltDB.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(bucket).Cursor()
		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			result = append(result, string(key))
		}
		return nil
	})
	return result, err
}

func GetMap[T any](bucket []byte) (map[string]T, error) {
	result := make(map[string]T)
	err := BoltDB.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(bucket).Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			var el T
			buff := bytes.NewBuffer(value)
			err := gob.NewDecoder(buff).Decode(&el)
			if err != nil {
				return err
			}
			result[string(key)] = el
		}
		return nil
	})
	return result, err
}

func Get[T any](bucket []byte, key string) (T, error) {
	var value T
	err := BoltDB.db.View(func(tx *bolt.Tx) error {
		el := tx.Bucket(bucket).Get([]byte(key))
		if el == nil {
			return nil
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

func Delete(bucket []byte, key string) error {
	return BoltDB.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucket).Delete([]byte(key))
	})
}

func Update[T any](bucket []byte, key string, value T) error {
	return BoltDB.db.Update(func(tx *bolt.Tx) error {
		var buff bytes.Buffer
		err := gob.NewEncoder(&buff).Encode(value)
		if err != nil {
			return err
		}
		return tx.Bucket(bucket).Put([]byte(key), buff.Bytes())
	})
}

func DeleteAll(bucket []byte) error {
	return BoltDB.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket(bucket)
		return err
	})
}
