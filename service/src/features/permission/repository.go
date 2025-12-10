package permission

import (
	"ivory/src/storage/db"
)

type Repository struct {
	bucket *db.Bucket[map[string]PermissionStatus]
}

func NewRepository(bucket *db.Bucket[map[string]PermissionStatus]) *Repository {
	return &Repository{
		bucket: bucket,
	}
}

func (r *Repository) CreateOrUpdate(username string, permissions map[string]PermissionStatus) error {
	return r.bucket.Update(username, permissions)
}

func (r *Repository) Get(username string) (map[string]PermissionStatus, error) {
	return r.bucket.Get(username)
}

func (r *Repository) GetAll() (UserPermissionsMap, error) {
	return r.bucket.GetMap(nil)
}

func (r *Repository) Delete(username string) error {
	return r.bucket.Delete(username)
}

func (r *Repository) DeleteAll() error {
	return r.bucket.DeleteAll()
}
