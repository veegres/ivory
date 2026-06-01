package query

import (
	"errors"
	"ivory/core/store"
)

var ErrCannotParseFileCorrupted = errors.New("cannot parse file, it is corrupted")

type Repository struct {
	bucket            *store.DbBucket[Response]
	queryLogFiles     *store.FileStorage
	maxBufferCapacity int
	maxLogElements    int
}

func NewRepository(
	bucket *store.DbBucket[Response],
	queryLogFiles *store.FileStorage,
) *Repository {
	return &Repository{
		bucket:            bucket,
		queryLogFiles:     queryLogFiles,
		maxBufferCapacity: 1024 * 1024,
		maxLogElements:    10,
	}
}
