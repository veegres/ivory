package persistence

import (
	"bytes"
	"encoding/gob"
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
	value, err := r.common.get(r.bucket, uuid.String())
	var credential CertModel
	buff := bytes.NewBuffer(value)
	err = gob.NewDecoder(buff).Decode(&credential)
	return credential, err
}

func (r CertRepository) List() (map[string]CertModel, error) {
	bytesList, err := r.common.getList(r.bucket)
	certMap := make(map[string]CertModel)
	for _, el := range bytesList {
		var cert CertModel
		buff := bytes.NewBuffer(el.value)
		err = gob.NewDecoder(buff).Decode(&cert)
		certMap[el.key] = cert
	}
	return certMap, err
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

	err := r.common.update(r.bucket, key.String(), cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r CertRepository) Delete(certUuid uuid.UUID) error {
	var err error
	err = File.Certs.Delete(certUuid)
	err = r.common.delete(r.bucket, certUuid.String())
	return err
}
