package password

import (
	"ivory/src/storage/db"

	"github.com/google/uuid"
)

type Repository struct {
	bucket *db.Bucket[Password]
}

func NewRepository(bucket *db.Bucket[Password]) *Repository {
	return &Repository{
		bucket: bucket,
	}
}

func (r *Repository) Create(credential Password) (uuid.UUID, error) {
	key := uuid.New()
	return key, r.bucket.Update(key.String(), credential)
}

func (r *Repository) Update(key uuid.UUID, credential Password) (uuid.UUID, error) {
	return key, r.bucket.Update(key.String(), credential)
}

func (r *Repository) Delete(key uuid.UUID) error {
	return r.bucket.Delete(key.String())
}

func (r *Repository) DeleteAll() error {
	return r.bucket.DeleteAll()
}

func (r *Repository) Map() (PasswordMap, error) {
	return r.bucket.GetMap(nil)
}

func (r *Repository) MapByType(credentialType PasswordType) (PasswordMap, error) {
	return r.bucket.GetMap(func(credential Password) bool {
		return credential.Type == credentialType
	})
}

func (r *Repository) Get(uuid uuid.UUID) (Password, error) {
	return r.bucket.Get(uuid.String())
}
