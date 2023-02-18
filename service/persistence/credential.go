package persistence

import (
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
	err = Update(r.secretBucket, r.encryptedRefKey, encrypted)
	err = Update(r.secretBucket, r.decryptedRefKey, decrypted)
	return err
}

func (r CredentialRepository) GetEncryptedRef() (string, error) {
	return Get[string](r.secretBucket, r.encryptedRefKey)
}

func (r CredentialRepository) GetDecryptedRef() (string, error) {
	return Get[string](r.secretBucket, r.decryptedRefKey)
}

func (r CredentialRepository) Create(credential Credential) (uuid.UUID, Credential, error) {
	key := uuid.New()
	err := Update(r.credentialBucket, key.String(), credential)
	cred := Credential{Username: credential.Username, Password: "configured", Type: credential.Type}
	return key, cred, err
}

func (r CredentialRepository) Update(key uuid.UUID, credential Credential) (uuid.UUID, Credential, error) {
	err := Update(r.credentialBucket, key.String(), credential)
	cred := Credential{Username: credential.Username, Password: "configured", Type: credential.Type}
	return key, cred, err
}

func (r CredentialRepository) Delete(key uuid.UUID) error {
	return Delete(r.credentialBucket, key.String())
}

func (r CredentialRepository) DeleteAll() error {
	return DeleteAll(r.credentialBucket)
}

func (r CredentialRepository) List() (map[string]Credential, error) {
	credentialMap, err := GetMap[Credential](r.credentialBucket, nil)
	credentialHiddenMap := r.hidePasswords(credentialMap)
	return credentialHiddenMap, err
}

func (r CredentialRepository) ListByType(credentialType CredentialType) (map[string]Credential, error) {
	credentialMap, err := GetMap[Credential](r.credentialBucket, func(credential Credential) bool {
		return credential.Type == credentialType
	})
	credentialHiddenMap := r.hidePasswords(credentialMap)
	return credentialHiddenMap, err
}

func (r CredentialRepository) Get(uuid uuid.UUID) (Credential, error) {
	return Get[Credential](r.credentialBucket, uuid.String())
}

func (r CredentialRepository) hidePasswords(credentialMap map[string]Credential) map[string]Credential {
	for key, credential := range credentialMap {
		credential.Password = "configured"
		credentialMap[key] = credential
	}
	return credentialMap
}
