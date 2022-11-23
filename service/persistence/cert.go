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

func (r CertRepository) List() ([]CertModel, error) {
	bytesList, err := r.common.getList(r.bucket)
	modelList := make([]CertModel, len(bytesList))
	for i, el := range bytesList {
		var cluster CertModel
		buff := bytes.NewBuffer(el.value)
		_ = gob.NewDecoder(buff).Decode(&cluster)
		modelList[i] = cluster
	}
	return modelList, err
}

func (r CertRepository) Create(fileName string) (*CertModel, error) {
	certUuid := uuid.New()
	certModel := CertModel{
		FileId:   certUuid,
		FileName: fileName,
		Path:     File.Certs.Create(certUuid),
	}

	err := r.common.update(r.bucket, certUuid.String(), certModel)
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
