package query

import (
	"errors"
	"ivory/src/storage/db"
	"ivory/src/storage/files"
)

var ErrCannotParseFileCorrupted = errors.New("cannot parse file, it is corrupted")

type Repository struct {
	bucket            *db.Bucket[Response]
	queryLogFiles     *files.Storage
	maxBufferCapacity int
	maxLogElements    int
}

func NewRepository(
	bucket *db.Bucket[Response],
	queryLogFiles *files.Storage,
) *Repository {
	return &Repository{
		bucket:            bucket,
		queryLogFiles:     queryLogFiles,
		maxBufferCapacity: 1024 * 1024,
		maxLogElements:    10,
	}
}
