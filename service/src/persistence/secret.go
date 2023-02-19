package persistence

import (
	"ivory/src/config"
)

type SecretRepository struct {
	bucket          *config.Bucket[string]
	encryptedRefKey string
	decryptedRefKey string
}

func NewSecretRepository(bucket *config.Bucket[string]) *SecretRepository {
	encryptedRefKey := "refEncrypted"
	decryptedRefKey := "refDecrypted"
	if _, err := bucket.Get(encryptedRefKey); err == nil {
		err := bucket.Update(encryptedRefKey, "")
		if err != nil {
			panic(err)
		}
	}
	if _, err := bucket.Get(decryptedRefKey); err == nil {
		err := bucket.Update(decryptedRefKey, "")
		if err != nil {
			panic(err)
		}
	}
	return &SecretRepository{
		bucket:          bucket,
		encryptedRefKey: encryptedRefKey,
		decryptedRefKey: decryptedRefKey,
	}
}

func (r *SecretRepository) UpdateRefs(encrypted string, decrypted string) error {
	var err error
	err = r.bucket.Update(r.encryptedRefKey, encrypted)
	err = r.bucket.Update(r.decryptedRefKey, decrypted)
	return err
}

func (r *SecretRepository) GetEncryptedRef() (string, error) {
	return r.bucket.Get(r.encryptedRefKey)
}

func (r *SecretRepository) GetDecryptedRef() (string, error) {
	return r.bucket.Get(r.decryptedRefKey)
}
