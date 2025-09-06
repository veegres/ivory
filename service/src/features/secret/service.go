package secret

import (
	"crypto/md5"
	"errors"
	"ivory/src/features/encryption"
	"sync"
)

type SecretService struct {
	key              [16]byte
	ref              string
	mutex            *sync.Mutex
	secretRepository *SecretRepository
	encryption       *encryption.EncryptionService
}

func NewSecretService(
	secretRepository *SecretRepository,
	encryption *encryption.EncryptionService,
) *SecretService {
	encryptedRef, err := secretRepository.GetEncryptedRef()
	if err != nil {
		panic(err)
	}

	return &SecretService{
		key:              [16]byte{},
		ref:              encryptedRef,
		mutex:            &sync.Mutex{},
		secretRepository: secretRepository,
		encryption:       encryption,
	}
}

func (s *SecretService) Get() [16]byte {
	return s.key
}

func (s *SecretService) GetByte() []byte {
	return s.key[:]
}

func (s *SecretService) Set(decryptedKey string, decryptedRef string) error {
	if decryptedKey == "" {
		return errors.New("secret word cannot be empty")
	}
	if decryptedRef == "" {
		var err error
		decryptedRef, err = s.secretRepository.GetDecryptedRef()
		if err != nil {
			return err
		}
		if decryptedRef == "" {
			return errors.New("reference word cannot be empty")
		}
	}
	encryptedKey := md5.Sum([]byte(decryptedKey))
	encryptedRef, err := s.encryption.Encrypt(decryptedRef, encryptedKey)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	if s.IsRefEmpty() {
		s.key = encryptedKey
		s.ref = encryptedRef
	} else {
		if s.ref == encryptedRef {
			s.key = encryptedKey
		} else {
			s.mutex.Unlock()
			return errors.New("the secret is not correct")
		}
	}
	err = s.secretRepository.UpdateRefs(encryptedRef, decryptedRef)
	s.mutex.Unlock()
	if err != nil {
		err = s.Clean()
	}
	return err
}

func (s *SecretService) Update(prevSecret string, newSecret string) ([16]byte, [16]byte, error) {
	if prevSecret == newSecret {
		return [16]byte{}, [16]byte{}, errors.New("the secrets are the same")
	}
	if s.IsSecretEmpty() {
		return [16]byte{}, [16]byte{}, errors.New("there is no secret")
	}

	previousKeyMd5 := md5.Sum([]byte(prevSecret))
	newKeyMd5 := md5.Sum([]byte(newSecret))
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if previousKeyMd5 != s.key {
		return [16]byte{}, [16]byte{}, errors.New("the secret is not correct")
	}

	decryptedRef, errDec := s.secretRepository.GetDecryptedRef()
	if errDec != nil {
		return [16]byte{}, [16]byte{}, errDec
	}
	encryptedRef, errEnc := s.encryption.Encrypt(decryptedRef, newKeyMd5)
	if errEnc != nil {
		return [16]byte{}, [16]byte{}, errEnc
	}
	errUpdate := s.secretRepository.UpdateRefs(encryptedRef, decryptedRef)
	if errUpdate != nil {
		return [16]byte{}, [16]byte{}, errUpdate
	}

	s.key = newKeyMd5
	return previousKeyMd5, newKeyMd5, nil
}

func (s *SecretService) Clean() error {
	s.mutex.Lock()

	s.ref = ""
	s.key = [16]byte{}
	errUpdateRef := s.secretRepository.UpdateRefs("", "")

	s.mutex.Unlock()
	return errUpdateRef
}

func (s *SecretService) IsRefEmpty() bool {
	return s.ref == ""
}

func (s *SecretService) IsSecretEmpty() bool {
	for _, b := range s.key {
		if b != 0 {
			return false
		}
	}
	return true
}

func (s *SecretService) Status() SecretStatus {
	return SecretStatus{
		Key: !s.IsSecretEmpty(),
		Ref: !s.IsRefEmpty(),
	}
}
