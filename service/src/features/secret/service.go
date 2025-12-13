package secret

import (
	"crypto/md5"
	"errors"
	"ivory/src/features/encryption"
	"sync"
)

var ErrWrongSecret = errors.New("the secret is not correct")

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
	decryptedRef := "ivory"
	encryptedRef, errRef := repository.Get()

	if errRef != nil {
		panic(errRef)
	}

	s := &Service{
		key:          [16]byte{},
		encryptedRef: encryptedRef,
		decryptedRef: decryptedRef,
		mutex:        &sync.Mutex{},
		repository:   repository,
		encryption:   encryption,
	}

	if !s.IsRefEmpty() {
		encryptedKey := md5.Sum([]byte(decryptedRef))
		_ = s.checkRef(encryptedKey)
	}

	return s
}

func (s *Service) Get() [16]byte {
	return s.key
}

func (s *Service) GetByte() []byte {
	return s.key[:]
}

func (s *Service) SetDefault() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.IsRefEmpty() {
		return errors.New("secret is already set")
	}
	encryptedKey := md5.Sum([]byte(s.decryptedRef))
	return s.setRef(encryptedKey)
}

func (s *Service) Set(secret string) error {
	if secret == "" {
		return errors.New("secret cannot be empty")
	}
	secretMd5 := md5.Sum([]byte(secret))
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.IsRefEmpty() {
		err := s.setRef(secretMd5)
		if err != nil {
			return err
		}
	} else {
		err := s.checkRef(secretMd5)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) setRef(key [16]byte) error {
	encryptedRef, err := s.encryption.Encrypt(s.decryptedRef, key)
	if err != nil {
		return err
	}
	err = s.repository.Update(encryptedRef)
	if err != nil {
		return err
	}
	s.key = key
	s.encryptedRef = encryptedRef
	return nil
}

func (s *Service) checkRef(key [16]byte) error {
	decryptedRef, errDec := s.encryption.Decrypt(s.encryptedRef, key)
	if errDec != nil {
		return errors.Join(ErrWrongSecret, errDec)
	}
	if s.decryptedRef != decryptedRef {
		return ErrWrongSecret
	}
	s.key = key
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

	if newSecret == "" {
		newSecret = s.decryptedRef
	}
	if prevSecret == "" {
		prevSecret = s.decryptedRef
	}
	previousKeyMd5 := md5.Sum([]byte(prevSecret))
	newKeyMd5 := md5.Sum([]byte(newSecret))
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if previousKeyMd5 != s.key {
		return [16]byte{}, [16]byte{}, ErrWrongSecret
	}

	encryptedRef, errEnc := s.encryption.Encrypt(s.decryptedRef, newKeyMd5)
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
