package password

import (
	"ivory/src/storage/bolt"

	"github.com/google/uuid"
)

type PasswordRepository struct {
	bucket *bolt.Bucket[Password]
}

func NewPasswordRepository(bucket *bolt.Bucket[Password]) *PasswordRepository {
	return &PasswordRepository{
		bucket: bucket,
	}
}

func (r *PasswordRepository) Create(credential Password) (uuid.UUID, error) {
	key := uuid.New()
	return key, r.bucket.Update(key.String(), credential)
}

func (r *PasswordRepository) Update(key uuid.UUID, credential Password) (uuid.UUID, error) {
	return key, r.bucket.Update(key.String(), credential)
}

func (r *PasswordRepository) Delete(key uuid.UUID) error {
	return r.bucket.Delete(key.String())
}

func (r *PasswordRepository) DeleteAll() error {
	return r.bucket.DeleteAll()
}

func (r *PasswordRepository) Map() (PasswordMap, error) {
	return r.bucket.GetMap(nil)
}

func (r *PasswordRepository) MapByType(credentialType PasswordType) (PasswordMap, error) {
	return r.bucket.GetMap(func(credential Password) bool {
		return credential.Type == credentialType
	})
}

func (r *PasswordRepository) Get(uuid uuid.UUID) (Password, error) {
	return r.bucket.Get(uuid.String())
}
