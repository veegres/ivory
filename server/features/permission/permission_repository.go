package permission

import "ivory/clients/storage"

type Repository struct {
	bucket *storage.DbBucket[PermissionMap]
}

func NewRepository(bucket *storage.DbBucket[PermissionMap]) *Repository {
	return &Repository{
		bucket: bucket,
	}
}

func (r *Repository) CreateOrUpdate(username string, permissions PermissionMap) error {
	return r.bucket.Update(username, permissions)
}

func (r *Repository) Get(username string) (PermissionMap, error) {
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
