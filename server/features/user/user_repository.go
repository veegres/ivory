package user

import "ivory/clients/storage"

type Repository struct {
	bucket     *storage.DbBucket[User]
	linkBucket *storage.DbBucket[Link]
}

func NewRepository(bucket *storage.DbBucket[User], linkBucket *storage.DbBucket[Link]) *Repository {
	return &Repository{bucket: bucket, linkBucket: linkBucket}
}
