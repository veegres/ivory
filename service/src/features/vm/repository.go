package vm

import (
	"ivory/src/storage/db"

	"github.com/google/uuid"
)

type Repository struct {
	bucket *db.Bucket[VM]
}

func NewRepository(bucket *db.Bucket[VM]) *Repository {
	return &Repository{
		bucket: bucket,
	}
}

func (r *Repository) List() ([]VM, error) {
	return r.bucket.GetList(nil, func(list []VM, i int, j int) bool {
		return list[i].Name < list[j].Name
	})
}

func (r *Repository) Get(id uuid.UUID) (VM, error) {
	return r.bucket.Get(id.String())
}

func (r *Repository) Create(model VM) (uuid.UUID, error) {
	id := uuid.New()
	model.Id = id
	_, err := r.bucket.Create(id.String(), model)
	return id, err
}

func (r *Repository) Update(id uuid.UUID, model VM) (uuid.UUID, error) {
	model.Id = id
	return id, r.bucket.Update(id.String(), model)
}

func (r *Repository) Delete(id uuid.UUID) error {
	return r.bucket.Delete(id.String())
}

func (r *Repository) DeleteAll() error {
	return r.bucket.DeleteAll()
}

func (r *Repository) Map() (map[string]VM, error) {
	return r.bucket.GetMap(nil)
}
