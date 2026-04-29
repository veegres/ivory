package query

import (
	"time"

	"github.com/google/uuid"
)

func (r *Repository) Get(uuid uuid.UUID) (Response, error) {
	return r.bucket.Get(uuid.String())
}

func (r *Repository) List() ([]Response, error) {
	return r.bucket.GetList(nil, r.sortAscByCreatedAt)
}

func (r *Repository) ListByType(queryType Type) ([]Response, error) {
	return r.bucket.GetList(func(cert Response) bool {
		return cert.Type == queryType
	}, r.sortAscByCreatedAt)
}

func (r *Repository) Create(query Response) (*uuid.UUID, *Response, error) {
	key := uuid.New()
	query.Id = key
	query.CreatedAt = time.Now().UnixNano()
	err := r.bucket.Update(key.String(), query)
	return &key, &query, err
}

func (r *Repository) Update(key uuid.UUID, query Response) (*uuid.UUID, *Response, error) {
	err := r.bucket.Update(key.String(), query)
	return &key, &query, err
}

func (r *Repository) Delete(key uuid.UUID) error {
	return r.bucket.Delete(key.String())
}

func (r *Repository) DeleteAll() error {
	return r.bucket.DeleteAll()
}

func (r *Repository) sortAscByCreatedAt(list []Response, i, j int) bool {
	return list[i].CreatedAt < list[j].CreatedAt
}
