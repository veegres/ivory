package vault

import (
	"ivory/src/storage/db"

	"github.com/google/uuid"
)

type Repository struct {
	bucket *db.Bucket[Vault]
}

func NewRepository(bucket *db.Bucket[Vault]) *Repository {
	return &Repository{
		bucket: bucket,
	}
}

func (r *Repository) Create(vault Vault) (uuid.UUID, error) {
	key := uuid.New()
	return key, r.bucket.Update(key.String(), vault)
}

func (r *Repository) Update(key uuid.UUID, vault Vault) (uuid.UUID, error) {
	return key, r.bucket.Update(key.String(), vault)
}

func (r *Repository) Delete(key uuid.UUID) error {
	return r.bucket.Delete(key.String())
}

func (r *Repository) DeleteAll() error {
	return r.bucket.DeleteAll()
}

func (r *Repository) Map() (VaultMap, error) {
	return r.bucket.GetMap(nil)
}

func (r *Repository) MapByType(vaultType VaultType) (VaultMap, error) {
	return r.bucket.GetMap(func(vault Vault) bool {
		return vault.Type == vaultType
	})
}

func (r *Repository) Get(uuid uuid.UUID) (Vault, error) {
	return r.bucket.Get(uuid.String())
}
