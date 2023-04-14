package service

import (
	"github.com/google/uuid"
	. "ivory/src/model"
	"ivory/src/persistence"
)

type PasswordService struct {
	passwordRepository *persistence.PasswordRepository
	secretService      *SecretService
	encryption         *EncryptionService

	defaultPasswordWord string
}

func NewPasswordService(
	passwordRepository *persistence.PasswordRepository,
	secretService *SecretService,
	encryption *EncryptionService,
) *PasswordService {
	return &PasswordService{
		passwordRepository:  passwordRepository,
		secretService:       secretService,
		encryption:          encryption,
		defaultPasswordWord: "********",
	}
}

func (s *PasswordService) Create(credential Credential) (*uuid.UUID, *Credential, error) {
	encryptedPassword, err := s.encryption.Encrypt(credential.Password, s.secretService.Get())
	if err != nil {
		return nil, nil, err
	}
	encryptedCred := Credential{Username: credential.Username, Password: encryptedPassword, Type: credential.Type}
	key, err := s.passwordRepository.Create(encryptedCred)

	cred := encryptedCred
	cred.Password = s.defaultPasswordWord
	return &key, &cred, err
}

func (s *PasswordService) Update(key uuid.UUID, credential Credential) (*uuid.UUID, *Credential, error) {
	encryptedPassword, err := s.encryption.Encrypt(credential.Password, s.secretService.Get())
	if err != nil {
		return nil, nil, err
	}
	encryptedCred := Credential{Username: credential.Username, Password: encryptedPassword, Type: credential.Type}
	key, err = s.passwordRepository.Update(key, encryptedCred)

	cred := encryptedCred
	cred.Password = s.defaultPasswordWord
	return &key, &cred, err
}

func (s *PasswordService) Delete(key uuid.UUID) error {
	return s.passwordRepository.Delete(key)
}

func (s *PasswordService) DeleteAll() error {
	return s.passwordRepository.DeleteAll()
}

func (s *PasswordService) Map(credentialType *CredentialType) (map[string]Credential, error) {
	var err error
	var credentialMap map[string]Credential
	if credentialType != nil {
		credentialMap, err = s.passwordRepository.MapByType(*credentialType)
	} else {
		credentialMap, err = s.passwordRepository.Map()
	}
	credentialHiddenMap := s.hidePasswords(credentialMap)
	return credentialHiddenMap, err
}

func (s *PasswordService) Get(uuid uuid.UUID) (*Credential, error) {
	cred, err := s.passwordRepository.Get(uuid)
	cred.Password = s.defaultPasswordWord
	return &cred, err
}

func (s *PasswordService) GetDecrypted(uuid uuid.UUID) (*Credential, error) {
	cred, err := s.passwordRepository.Get(uuid)
	if err != nil {
		return nil, err
	}
	decryptedPassword, err := s.encryption.Decrypt(cred.Password, s.secretService.Get())
	cred.Password = decryptedPassword
	return &cred, err
}

func (s *PasswordService) ReEncryptPasswords(oldSecret [16]byte, newSecret [16]byte) error {
	credentialMap, err := s.passwordRepository.Map()
	if err != nil {
		return err
	}
	for id, credential := range credentialMap {
		id, _ := uuid.Parse(id)
		oldEncryptedPass := credential.Password
		decryptedPass, _ := s.encryption.Decrypt(oldEncryptedPass, oldSecret)
		newEncryptedPass, _ := s.encryption.Encrypt(decryptedPass, newSecret)
		credential.Password = newEncryptedPass
		_, err = s.passwordRepository.Update(id, credential)
	}
	return nil
}

func (s *PasswordService) hidePasswords(credentialMap map[string]Credential) map[string]Credential {
	for key, credential := range credentialMap {
		credential.Password = s.defaultPasswordWord
		credentialMap[key] = credential
	}
	return credentialMap
}
