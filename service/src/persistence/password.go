package persistence

import (
	"github.com/google/uuid"
	"ivory/src/config"
	. "ivory/src/model"
)

type PasswordRepository struct {
	bucket *config.Bucket[Credential]
}

func NewPasswordRepository(bucket *config.Bucket[Credential]) *PasswordRepository {
	return &PasswordRepository{
		bucket: bucket,
	}
}

func (r *PasswordRepository) Create(credential Credential) (uuid.UUID, Credential, error) {
	key := uuid.New()
	err := r.bucket.Update(key.String(), credential)
	cred := Credential{Username: credential.Username, Password: "configured", Type: credential.Type}
	return key, cred, err
}

func (r *PasswordRepository) Update(key uuid.UUID, credential Credential) (uuid.UUID, Credential, error) {
	err := r.bucket.Update(key.String(), credential)
	cred := Credential{Username: credential.Username, Password: "configured", Type: credential.Type}
	return key, cred, err
}

func (r *PasswordRepository) Delete(key uuid.UUID) error {
	return r.bucket.Delete(key.String())
}

func (r *PasswordRepository) DeleteAll() error {
	return r.bucket.DeleteAll()
}

func (r *PasswordRepository) List() (map[string]Credential, error) {
	credentialMap, err := r.bucket.GetMap(nil)
	credentialHiddenMap := r.hidePasswords(credentialMap)
	return credentialHiddenMap, err
}

func (r *PasswordRepository) ListByType(credentialType CredentialType) (map[string]Credential, error) {
	credentialMap, err := r.bucket.GetMap(func(credential Credential) bool {
		return credential.Type == credentialType
	})
	credentialHiddenMap := r.hidePasswords(credentialMap)
	return credentialHiddenMap, err
}

func (r *PasswordRepository) Get(uuid uuid.UUID) (Credential, error) {
	return r.bucket.Get(uuid.String())
}

func (r *PasswordRepository) hidePasswords(credentialMap map[string]Credential) map[string]Credential {
	for key, credential := range credentialMap {
		credential.Password = "configured"
		credentialMap[key] = credential
	}
	return credentialMap
}
