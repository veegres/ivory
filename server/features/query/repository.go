package query

import (
	"errors"
	"ivory/clients/storage"
)

var ErrCannotParseFileCorrupted = errors.New("cannot parse file, it is corrupted")

type Repository struct {
	bucket            *storage.DbBucket[Response]
	queryLogFiles     *storage.FileStorage
	maxBufferCapacity int
	maxLogElements    int
}

func NewRepository(
	bucket *storage.DbBucket[Response],
	queryLogFiles *storage.FileStorage,
) *Repository {
	return &Repository{
		bucket:            bucket,
		queryLogFiles:     queryLogFiles,
		maxBufferCapacity: 1024 * 1024,
		maxLogElements:    10,
	}
}
