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

func (r *QueryRepository) List() ([]Query, error) {
	return r.bucket.GetList(nil, r.sortByType)
}

func (r *QueryRepository) ListByType(queryType QueryType) ([]Query, error) {
	return r.bucket.GetList(func(cert Query) bool {
		return cert.Type == queryType
	}, r.sortByType)
}

func (r *QueryRepository) Create(query Query) (*uuid.UUID, *Query, error) {
	key := uuid.New()
	query.Id = key
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

func (r *QueryRepository) sortByType(list []Query, i, j int) bool {
	return list[i].Creation != list[j].Creation && list[i].Creation > list[j].Creation
}
