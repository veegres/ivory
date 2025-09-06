package secret

import (
	"errors"
	"ivory/src/storage/db"
)

type Repository struct {
	bucket          *db.Bucket[string]
	encryptedRefKey string
	decryptedRefKey string
}

func NewRepository(bucket *db.Bucket[string]) *Repository {
	encryptedRef := "encryptedRef"
	decryptedRef := "decryptedRef"
	if _, err := bucket.Get(encryptedRef); err != nil {
		err := bucket.Update(encryptedRef, "")
		if err != nil {
			panic(err)
		}
	}
	if _, err := bucket.Get(decryptedRef); err != nil {
		err := bucket.Update(decryptedRef, "")
		if err != nil {
			panic(err)
		}
	}
	return &Repository{
		bucket:          bucket,
		encryptedRefKey: encryptedRef,
		decryptedRefKey: decryptedRef,
	}
}

func (r *Repository) UpdateRefs(encrypted string, decrypted string) error {
	errEnc := r.bucket.Update(r.encryptedRefKey, encrypted)
	errDec := r.bucket.Update(r.decryptedRefKey, decrypted)
	return errors.Join(errEnc, errDec)
}

func (r *Repository) GetEncryptedRef() (string, error) {
	return r.bucket.Get(r.encryptedRefKey)
}

func (r *Repository) GetDecryptedRef() (string, error) {
	return r.bucket.Get(r.decryptedRefKey)
}
