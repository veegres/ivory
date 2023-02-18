package persistence

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	. "ivory/model"
	"path"
	"strings"
)

type CertRepository struct {
	common common
	bucket []byte
}

func (r CertRepository) Get(uuid uuid.UUID) (CertModel, error) {
	return Get[CertModel](r.bucket, uuid.String())
}

func (r CertRepository) List() (map[string]CertModel, error) {
	return GetMap[CertModel](r.bucket, nil)
}

func (r CertRepository) ListByType(certType CertType) (map[string]CertModel, error) {
	return GetMap[CertModel](r.bucket, func(cert CertModel) bool {
		return cert.Type == certType
	})
}

func (r CertRepository) Create(pathStr string, certType CertType, fileUsageType FileUsageType) (*CertModel, error) {
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
		cert.Path = File.Certs.Create(key)
	case FileUsageType(PATH):
		cert.FileName = fileName
		cert.Path = pathStr
	default:
		return nil, errors.New("unknown certificate type")
	}

	err := Update(r.bucket, key.String(), cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r CertRepository) Delete(certUuid uuid.UUID) error {
	var err error
	cert, err := r.Get(certUuid)
	if cert.FileUsageType == FileUsageType(UPLOAD) {
		err = File.Certs.Delete(certUuid)
	}
	err = Delete(r.bucket, certUuid.String())
	return err
}

func (r CertRepository) DeleteAll() error {
	err := File.Certs.DeleteAll()
	err = DeleteAll(r.bucket)
	return err
}
