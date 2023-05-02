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
	encryptedPassword, errEnc := s.encryption.Encrypt(credential.Password, s.secretService.Get())
	if errEnc != nil {
		return nil, nil, errEnc
	}

	encryptedCred := Credential{Username: credential.Username, Password: encryptedPassword, Type: credential.Type}
	key, errCreate := s.passwordRepository.Create(encryptedCred)

	cred := encryptedCred
	cred.Password = s.defaultPasswordWord
	return &key, &cred, errCreate
}

func (s *PasswordService) Update(key uuid.UUID, credential Credential) (*uuid.UUID, *Credential, error) {
	encryptedPassword, errEnc := s.encryption.Encrypt(credential.Password, s.secretService.Get())
	if errEnc != nil {
		return nil, nil, errEnc
	}

	encryptedCred := Credential{Username: credential.Username, Password: encryptedPassword, Type: credential.Type}
	key, errUpdate := s.passwordRepository.Update(key, encryptedCred)

	cred := encryptedCred
	cred.Password = s.defaultPasswordWord
	return &key, &cred, errUpdate
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
	cred, errCred := s.passwordRepository.Get(uuid)
	if errCred != nil {
		return nil, errCred
	}
	decryptedPassword, errDec := s.encryption.Decrypt(cred.Password, s.secretService.Get())
	if errDec != nil {
		return nil, errDec
	}
	cred.Password = decryptedPassword
	return &cred, nil
}

func (s *PasswordService) ReEncryptPasswords(oldSecret [16]byte, newSecret [16]byte) error {
	credentialMap, errGet := s.passwordRepository.Map()
	if errGet != nil {
		return errGet
	}
	for id, credential := range credentialMap {
		id, errParse := uuid.Parse(id)
		if errParse != nil {
			return errParse
		}
		oldEncryptedPass := credential.Password
		decryptedPass, errDec := s.encryption.Decrypt(oldEncryptedPass, oldSecret)
		if errDec != nil {
			return errDec
		}
		newEncryptedPass, errEnc := s.encryption.Encrypt(decryptedPass, newSecret)
		if errEnc != nil {
			return errEnc
		}
		credential.Password = newEncryptedPass
		_, errUpdate := s.passwordRepository.Update(id, credential)
		if errUpdate != nil {
			return errUpdate
		}
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
