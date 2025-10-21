package secret

import (
	"ivory/src/storage/db"
)

type Repository struct {
	bucket          *db.Bucket[string]
	encryptedRefKey string
}

func NewRepository(bucket *db.Bucket[string]) *Repository {
	encryptedRefKey := "encryptedRef"
	if _, err := bucket.Get(encryptedRefKey); err != nil {
		err := bucket.Update(encryptedRefKey, "")
		if err != nil {
			panic(err)
		}
	}
	return &Repository{
		bucket:          bucket,
		encryptedRefKey: encryptedRefKey,
	}
}

func (r *Repository) Update(ref string) error {
	return r.bucket.Update(r.encryptedRefKey, ref)
}

func (r *Repository) Get() (string, error) {
	return r.bucket.Get(r.encryptedRefKey)
}
