package password

import (
	"ivory/src/features/encryption"
	"ivory/src/features/secret"

	"github.com/google/uuid"
)

type Service struct {
	passwordRepository *Repository
	secretService      *secret.Service
	encryption         *encryption.Service

	defaultPasswordWord string
}

func NewService(
	passwordRepository *Repository,
	secretService *secret.Service,
	encryption *encryption.Service,
) *Service {
	return &Service{
		passwordRepository:  passwordRepository,
		secretService:       secretService,
		encryption:          encryption,
		defaultPasswordWord: "********",
	}
}

func (s *Service) Create(credential Password) (*uuid.UUID, *Password, error) {
	encryptedPassword, errEnc := s.encryption.Encrypt(credential.Password, s.secretService.Get())
	if errEnc != nil {
		return nil, nil, errEnc
	}

	encryptedCred := Password{Username: credential.Username, Password: encryptedPassword, Type: credential.Type}
	key, errCreate := s.passwordRepository.Create(encryptedCred)

	cred := encryptedCred
	cred.Password = s.defaultPasswordWord
	return &key, &cred, errCreate
}

func (s *Service) Update(key uuid.UUID, credential Password) (*uuid.UUID, *Password, error) {
	encryptedPassword, errEnc := s.encryption.Encrypt(credential.Password, s.secretService.Get())
	if errEnc != nil {
		return nil, nil, errEnc
	}

	encryptedCred := Password{Username: credential.Username, Password: encryptedPassword, Type: credential.Type}
	key, errUpdate := s.passwordRepository.Update(key, encryptedCred)

	cred := encryptedCred
	cred.Password = s.defaultPasswordWord
	return &key, &cred, errUpdate
}

func (s *Service) Delete(key uuid.UUID) error {
	return s.passwordRepository.Delete(key)
}

func (s *Service) DeleteAll() error {
	return s.passwordRepository.DeleteAll()
}

func (s *Service) Map(credentialType *PasswordType) (PasswordMap, error) {
	var err error
	var credentialMap PasswordMap
	if credentialType != nil {
		credentialMap, err = s.passwordRepository.MapByType(*credentialType)
	} else {
		credentialMap, err = s.passwordRepository.Map()
	}
	credentialHiddenMap := s.hidePasswords(credentialMap)
	return credentialHiddenMap, err
}

func (s *Service) Get(uuid uuid.UUID) (*Password, error) {
	cred, err := s.passwordRepository.Get(uuid)
	if err != nil {
		return nil, err
	}
	cred.Password = s.defaultPasswordWord
	return &cred, nil
}

func (s *Service) GetDecrypted(uuid uuid.UUID) (*Password, error) {
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

func (s *Service) Reencrypt(oldSecret [16]byte, newSecret [16]byte) error {
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

func (s *Service) hidePasswords(credentialMap PasswordMap) PasswordMap {
	for key, credential := range credentialMap {
		credential.Password = s.defaultPasswordWord
		credentialMap[key] = credential
	}
	return credentialMap
}
