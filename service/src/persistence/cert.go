package persistence

import (
	"errors"
	"github.com/google/uuid"
	"ivory/src/config"
	. "ivory/src/model"
)

type CertRepository struct {
	bucket *config.Bucket[CertModel]
	file   *config.FileGateway
}

func NewCertRepository(bucket *config.Bucket[CertModel], file *config.FileGateway) *CertRepository {
	return &CertRepository{
		bucket: bucket,
		file:   file,
	}
}

func (r *CertRepository) Get(uuid uuid.UUID) (CertModel, error) {
	return r.bucket.Get(uuid.String())
}

func (r *CertRepository) GetFile(uuid uuid.UUID) ([]byte, error) {
	info, err := r.Get(uuid)
	if err != nil {
		return nil, err
	}
	return r.file.Read(info.Path)
}

func (r *CertRepository) List() (map[string]CertModel, error) {
	return r.bucket.GetMap(nil)
}

func (r *CertRepository) ListByType(certType CertType) (map[string]CertModel, error) {
	return r.bucket.GetMap(func(cert CertModel) bool {
		return cert.Type == certType
	})
}

func (r *CertRepository) Create(cert CertModel, pathStr string) (*CertModel, error) {
	key := uuid.New()

	switch cert.FileUsageType {
	case FileUsageType(UPLOAD):
		cert.Path = r.file.Create(key)
	case FileUsageType(PATH):
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
	if cert.FileUsageType == FileUsageType(UPLOAD) {
		// NOTE: we shouldn't check error here, if there is no file we should try to remove info
		_ = r.file.Delete(certUuid)
	}
	return r.bucket.Delete(certUuid.String())
}

func (r *CertRepository) DeleteAll() error {
	errFile := r.file.DeleteAll()
	errBucket := r.bucket.DeleteAll()
	return errors.Join(errFile, errBucket)
}
