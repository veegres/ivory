package persistence

import (
	"bytes"
	"encoding/gob"
	"github.com/google/uuid"
	. "ivory/model"
)

type CertRepository struct {
	common common
	bucket []byte
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

func (r CertRepository) Create(fileName string) (*CertModel, error) {
	key := uuid.New()
	certModel := CertModel{FileName: fileName, Path: File.Certs.Create(key)}

	err := r.common.update(r.bucket, key.String(), certModel)
	if err != nil {
		return nil, err
	}
	return &certModel, nil
}

func (r CertRepository) Delete(certUuid uuid.UUID) error {
	var err error
	err = File.Certs.Delete(certUuid)
	err = r.common.delete(r.bucket, certUuid.String())
	return err
}
