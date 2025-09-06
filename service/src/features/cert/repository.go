package cert

import (
	"errors"
	"ivory/src/storage/db"
	"ivory/src/storage/files"

	"github.com/google/uuid"
)

type CertRepository struct {
	bucket *db.Bucket[Cert]
	file   *files.Storage
}

func NewCertRepository(bucket *db.Bucket[Cert], file *files.Storage) *CertRepository {
	return &CertRepository{
		bucket: bucket,
		file:   file,
	}
}

func (r *CertRepository) Get(uuid uuid.UUID) (Cert, error) {
	return r.bucket.Get(uuid.String())
}

func (r *CertRepository) List() (CertMap, error) {
	return r.bucket.GetMap(nil)
}

func (r *CertRepository) ListByType(certType CertType) (CertMap, error) {
	return r.bucket.GetMap(func(cert Cert) bool {
		return cert.Type == certType
	})
}

func (r *CertRepository) Create(cert Cert, pathStr string) (*Cert, error) {
	key := uuid.New()

	switch cert.FileUsageType {
	case UPLOAD:
		path, errCreate := r.file.CreateByName(key.String())
		if errCreate != nil {
			return nil, errCreate
		}
		cert.Path = path
	case PATH:
		if !r.file.ExistByPath(pathStr) {
			return nil, errors.New("there is no such file")
		}
		cert.Path = pathStr
	default:
		return nil, errors.New("unknown certificate type")
	}

	err := r.bucket.Update(key.String(), cert)
	return &cert, err
}

func (r *CertRepository) Delete(certUuid uuid.UUID) error {
	cert, err := r.Get(certUuid)
	if err != nil {
		return err
	}
	if cert.FileUsageType == UPLOAD {
		// NOTE: we shouldn't check the error here, if there is no file, we should try to remove info
		_ = r.file.DeleteByName(certUuid.String())
	}
	return r.bucket.Delete(certUuid.String())
}

func (r *CertRepository) DeleteAll() error {
	errFile := r.file.DeleteAll()
	errBucket := r.bucket.DeleteAll()
	return errors.Join(errFile, errBucket)
}
