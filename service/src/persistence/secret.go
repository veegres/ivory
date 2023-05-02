package persistence

import (
	"errors"
	"ivory/src/config"
)

type SecretRepository struct {
	bucket          *config.Bucket[string]
	encryptedRefKey string
	decryptedRefKey string
}

func NewSecretRepository(bucket *config.Bucket[string]) *SecretRepository {
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
	return &SecretRepository{
		bucket:          bucket,
		encryptedRefKey: encryptedRef,
		decryptedRefKey: decryptedRef,
	}
}

func (r *SecretRepository) UpdateRefs(encrypted string, decrypted string) error {
	errEnc := r.bucket.Update(r.encryptedRefKey, encrypted)
	errDec := r.bucket.Update(r.decryptedRefKey, decrypted)
	return errors.Join(errEnc, errDec)
}

func (r *SecretRepository) GetEncryptedRef() (string, error) {
	return r.bucket.Get(r.encryptedRefKey)
}

func (r *SecretRepository) GetDecryptedRef() (string, error) {
	return r.bucket.Get(r.decryptedRefKey)
}
