package persistence

import (
	"bytes"
	"encoding/gob"
	"github.com/google/uuid"
	. "ivory/model"
)

type CredentialRepository struct {
	common           common
	secretBucket     []byte
	credentialBucket []byte
	encryptedRefKey  string
	decryptedRefKey  string
}

func (r CredentialRepository) UpdateRefs(encrypted string, decrypted string) error {
	var err error
	err = r.common.update(r.secretBucket, r.encryptedRefKey, encrypted)
	err = r.common.update(r.secretBucket, r.decryptedRefKey, decrypted)
	return err
}

func (r CredentialRepository) GetEncryptedRef() string {
	get, err := r.common.get(r.secretBucket, r.encryptedRefKey)
	if err != nil {
		return ""
	}
	var str string
	buff := bytes.NewBuffer(get)
	_ = gob.NewDecoder(buff).Decode(&str)
	return str
}

func (r CredentialRepository) GetDecryptedRef() string {
	get, err := r.common.get(r.secretBucket, r.decryptedRefKey)
	if err != nil {
		return ""
	}
	var str string
	buff := bytes.NewBuffer(get)
	_ = gob.NewDecoder(buff).Decode(&str)
	return str
}

func (r CredentialRepository) UpdateCredential(credential Credential) (uuid.UUID, error) {
	key, err := uuid.NewUUID()
	err = r.common.update(r.credentialBucket, key.String(), credential)
	return key, err
}

func (r CredentialRepository) DeleteCredential(key uuid.UUID) error {
	return r.common.delete(r.credentialBucket, key.String())
}

func (r CredentialRepository) GetCredentialMap() map[string]Credential {
	bytesList, _ := r.common.getList(r.credentialBucket)
	credentialMap := make(map[string]Credential)
	for _, el := range bytesList {
		var credential Credential
		buff := bytes.NewBuffer(el.value)
		_ = gob.NewDecoder(buff).Decode(&credential)
		credentialMap[el.key] = credential
	}
	return credentialMap
}

func (r CredentialRepository) GetCredential(uuid uuid.UUID) (Credential, error) {
	value, err := r.common.get(r.credentialBucket, uuid.String())
	var credential Credential
	buff := bytes.NewBuffer(value)
	_ = gob.NewDecoder(buff).Decode(&credential)
	return credential, err
}
