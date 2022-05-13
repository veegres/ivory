package service

import (
	"crypto/md5"
	"errors"
	"ivory/persistence"
	"sync"
)

var Secret *secret

type secret struct {
	key   [16]byte
	ref   string
	mutex *sync.Mutex
}

func (s *secret) Create(ref string) {
	Secret = &secret{
		key:   [16]byte{},
		ref:   ref,
		mutex: &sync.Mutex{},
	}
}

func (s *secret) Get() [16]byte {
	return s.key
}

func (s *secret) Set(key string, decrypted string) error {
	if decrypted == "" {
		decrypted = persistence.Database.Credential.GetDecryptedRef()
	}
	keySha1 := md5.Sum([]byte(key))
	ref, err := Encrypt(decrypted, keySha1)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	if s.IsRefEmpty() {
		s.key = keySha1
		s.ref = ref
	} else {
		if s.ref == ref {
			s.key = keySha1
		} else {
			s.mutex.Unlock()
			return errors.New("the secret is not correct")
		}
	}
	err = persistence.Database.Credential.UpdateRefs(ref, decrypted)
	s.mutex.Unlock()
	if err != nil {
		s.Clean()
	}
	return err
}

func (s *secret) Update(previousKey string, newKey string) error {
	previousKeySha1 := md5.Sum([]byte(previousKey))
	newKeySha1 := md5.Sum([]byte(newKey))
	s.mutex.Lock()
	if s.IsRefEmpty() {
		return errors.New("there is not secret")
	}
	if previousKeySha1 != s.key {
		return errors.New("the secret is not correct")
	}

	s.key = newKeySha1
	s.mutex.Unlock()
	return nil
}

func (s *secret) Clean() {
	s.mutex.Lock()
	s.ref = ""
	s.key = [16]byte{}
	_ = persistence.Database.Credential.UpdateRefs("", "")
	s.mutex.Unlock()
}

func (s *secret) IsRefEmpty() bool {
	return s.ref == ""
}

func (s *secret) IsSecretEmpty() bool {
	for _, b := range s.key {
		if b != 0 {
			return false
		}
	}
	return true
}
