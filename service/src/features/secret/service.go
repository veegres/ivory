package secret

import (
	"crypto/md5"
	"errors"
	"ivory/src/features/encryption"
	"sync"
)

type Service struct {
	key          [16]byte
	encryptedRef string
	decryptedRef string
	mutex        *sync.Mutex
	repository   *Repository
	encryption   *encryption.Service
}

func NewService(
	repository *Repository,
	encryption *encryption.Service,
) *Service {
	encryptedRef, err := repository.Get()
	if err != nil {
		panic(err)
	}

	return &Service{
		key:          [16]byte{},
		encryptedRef: encryptedRef,
		decryptedRef: "ivory",
		mutex:        &sync.Mutex{},
		repository:   repository,
		encryption:   encryption,
	}
}

func (s *Service) Get() [16]byte {
	return s.key
}

func (s *Service) GetByte() []byte {
	return s.key[:]
}

func (s *Service) Set(decryptedKey string) error {
	if decryptedKey == "" {
		return errors.New("secret word cannot be empty")
	}
	encryptedKey := md5.Sum([]byte(decryptedKey))
	s.mutex.Lock()
	if s.IsRefEmpty() {
		encryptedRef, err := s.encryption.Encrypt(s.decryptedRef, encryptedKey)
		if err != nil {
			s.mutex.Unlock()
			return err
		}
		err = s.repository.Update(encryptedRef)
		if err != nil {
			s.mutex.Unlock()
			return err
		}
		s.key = encryptedKey
		s.encryptedRef = encryptedRef
	} else {
		decryptedRef, errDec := s.encryption.Decrypt(s.encryptedRef, encryptedKey)
		if errDec != nil {
			s.mutex.Unlock()
			return errDec
		}
		if s.decryptedRef == decryptedRef {
			s.key = encryptedKey
		} else {
			s.mutex.Unlock()
			return errors.New("the secret is not correct")
		}
	}
	s.mutex.Unlock()
	return nil
}

func (s *Service) Verify(secret string) bool {
	return s.key == md5.Sum([]byte(secret))
}

func (s *Service) Update(prevSecret string, newSecret string) ([16]byte, [16]byte, error) {
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

	decryptedRef, errDec := s.repository.Get()
	if errDec != nil {
		return [16]byte{}, [16]byte{}, errDec
	}
	encryptedRef, errEnc := s.encryption.Encrypt(decryptedRef, newKeyMd5)
	if errEnc != nil {
		return [16]byte{}, [16]byte{}, errEnc
	}
	errUpdate := s.repository.Update(encryptedRef)
	if errUpdate != nil {
		return [16]byte{}, [16]byte{}, errUpdate
	}

	s.key = newKeyMd5
	return previousKeyMd5, newKeyMd5, nil
}

func (s *Service) Clean() error {
	s.mutex.Lock()

	s.encryptedRef = ""
	s.key = [16]byte{}
	errUpdateRef := s.repository.Update("")

	s.mutex.Unlock()
	return errUpdateRef
}

func (s *Service) IsRefEmpty() bool {
	return s.encryptedRef == ""
}

func (s *Service) IsSecretEmpty() bool {
	for _, b := range s.key {
		if b != 0 {
			return false
		}
	}
	return true
}

func (s *Service) Status() SecretStatus {
	return SecretStatus{
		Key: !s.IsSecretEmpty(),
		Ref: !s.IsRefEmpty(),
	}
}
