package persistence

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"ivory/src/config"
	. "ivory/src/model"
	"path"
	"strings"
)

type CertRepository struct {
	bucket *config.Bucket[CertModel]
	file   *config.FilePersistence
}

func NewCertRepository(bucket *config.Bucket[CertModel], file *config.FilePersistence) *CertRepository {
	return &CertRepository{
		bucket: bucket,
		file:   file,
	}
}

func (r *CertRepository) Get(uuid uuid.UUID) (CertModel, error) {
	return r.bucket.Get(uuid.String())
}

func (r *CertRepository) List() (map[string]CertModel, error) {
	return r.bucket.GetMap(nil)
}

func (r *CertRepository) ListByType(certType CertType) (map[string]CertModel, error) {
	return r.bucket.GetMap(func(cert CertModel) bool {
		return cert.Type == certType
	})
}

func (r *CertRepository) Create(pathStr string, certType CertType, fileUsageType FileUsageType) (*CertModel, error) {
	formats := []string{".crt", ".key", ".chain"}
	key := uuid.New()
	cert := CertModel{FileUsageType: fileUsageType, Type: certType}

	fileName := path.Base(pathStr)
	if fileName == "" {
		return nil, errors.New("file name cannot be empty")
	}

	fileFormat := path.Ext(pathStr)
	idx := slices.IndexFunc(formats, func(s string) bool { return s == fileFormat })
	if idx == -1 {
		return nil, errors.New("file format is not correct, allowed formats: " + strings.Join(formats, ", "))
	}

	switch fileUsageType {
	case FileUsageType(UPLOAD):
		cert.FileName = fileName
		cert.Path = r.file.Create(key)
	case FileUsageType(PATH):
		cert.FileName = fileName
		cert.Path = pathStr
	default:
		return nil, errors.New("unknown certificate type")
	}

	err := r.bucket.Update(key.String(), cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *CertRepository) Delete(certUuid uuid.UUID) error {
	var err error
	cert, err := r.Get(certUuid)
	if cert.FileUsageType == FileUsageType(UPLOAD) {
		err = r.file.Delete(certUuid)
	}
	err = r.bucket.Delete(certUuid.String())
	return err
}

func (r *CertRepository) DeleteAll() error {
	err := r.file.DeleteAll()
	err = r.bucket.DeleteAll()
	return err
}
