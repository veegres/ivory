package service

import (
	"crypto/md5"
	"errors"
	"github.com/google/uuid"
	. "ivory/src/model"
	"ivory/src/persistence"
	"sync"
)

type SecretService struct {
	key                [16]byte
	ref                string
	mutex              *sync.Mutex
	secretRepository   *persistence.SecretRepository
	passwordRepository *persistence.PasswordRepository
	encryption         *Encryption
}

func NewSecretService(
	secretRepository *persistence.SecretRepository,
	passwordRepository *persistence.PasswordRepository,
	encryption *Encryption,
) *SecretService {
	encryptedRef, err := secretRepository.GetEncryptedRef()
	if err != nil {
		panic(err)
	}

	return &SecretService{
		key:                [16]byte{},
		ref:                encryptedRef,
		mutex:              &sync.Mutex{},
		secretRepository:   secretRepository,
		passwordRepository: passwordRepository,
		encryption:         encryption,
	}
}

func (s *SecretService) Get() [16]byte {
	return s.key
}

func (s *SecretService) Set(decryptedKey string, decryptedRef string) error {
	if decryptedRef == "" {
		var err error
		decryptedRef, err = s.secretRepository.GetDecryptedRef()
		if err != nil {
			return err
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

func (s *SecretService) Update(previousKey string, newKey string) error {
	previousKeySha1 := md5.Sum([]byte(previousKey))
	newKeySha1 := md5.Sum([]byte(newKey))
	s.mutex.Lock()

	var err error
	if s.IsSecretEmpty() {
		err = errors.New("there is no secret")
		return err
	}
	if previousKeySha1 != s.key {
		err = errors.New("the secret is not correct")
		return err
	}

	s.key = newKeySha1
	credentialMap, err := s.passwordRepository.List()
	if err != nil {
		return err
	}
	for id, credential := range credentialMap {
		id, _ := uuid.Parse(id)
		oldEncodedPass := credential.Password
		decodedPass, _ := s.encryption.Decrypt(oldEncodedPass, previousKeySha1)
		newEncodedPass, _ := s.encryption.Encrypt(decodedPass, newKeySha1)
		credential.Password = newEncodedPass
		_, _, err = s.passwordRepository.Update(id, credential)
	}

	s.mutex.Unlock()
	return err
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
