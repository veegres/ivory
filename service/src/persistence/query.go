package persistence

import (
	"github.com/google/uuid"
	"ivory/src/config"
	. "ivory/src/model"
)

type QueryRepository struct {
	bucket *config.Bucket[Query]
}

func NewQueryRepository(bucket *config.Bucket[Query]) *QueryRepository {
	return &QueryRepository{
		bucket: bucket,
	}
}

func (r *QueryRepository) Get(uuid uuid.UUID) (Query, error) {
	return r.bucket.Get(uuid.String())
}

func (r *QueryRepository) Map() (map[string]Query, error) {
	return r.bucket.GetMap(nil)
}

func (r *QueryRepository) MapByType(queryType QueryType) (map[string]Query, error) {
	return r.bucket.GetMap(func(cert Query) bool {
		return cert.Type == queryType
	})
}

func (r *QueryRepository) Create(query Query) (*uuid.UUID, *Query, error) {
	key := uuid.New()
	err := r.bucket.Update(key.String(), query)
	return &key, &query, err
}

func (r *QueryRepository) Update(key uuid.UUID, query Query) (*uuid.UUID, *Query, error) {
	err := r.bucket.Update(key.String(), query)
	return &key, &query, err
}

func (r *QueryRepository) Delete(key uuid.UUID) error {
	return r.bucket.Delete(key.String())
}

func (r *QueryRepository) DeleteAll() error {
	return r.bucket.DeleteAll()
}
