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
	key                    [16]byte
	ref                    string
	mutex                  *sync.Mutex
	passwordRepository     *persistence.PasswordRepository
	secretRepository       *persistence.SecretRepository
	clusterRepository      *persistence.ClusterRepository
	certRepository         *persistence.CertRepository
	tagRepository          *persistence.TagRepository
	compactTableRepository *persistence.CompactTableRepository
	encryption             *Encryption
}

func NewSecretService(
	passwordRepository *persistence.PasswordRepository,
	secretRepository *persistence.SecretRepository,
	clusterRepository *persistence.ClusterRepository,
	certRepository *persistence.CertRepository,
	tagRepository *persistence.TagRepository,
	compactTableRepository *persistence.CompactTableRepository,
	encryption *Encryption,
) *SecretService {
	return &SecretService{
		key:                    [16]byte{},
		ref:                    "",
		mutex:                  &sync.Mutex{},
		passwordRepository:     passwordRepository,
		secretRepository:       secretRepository,
		compactTableRepository: compactTableRepository,
		clusterRepository:      clusterRepository,
		certRepository:         certRepository,
		tagRepository:          tagRepository,
		encryption:             encryption,
	}
}

func (s *SecretService) Get() [16]byte {
	return s.key
}

func (s *SecretService) Set(key string, decrypted string) error {
	if decrypted == "" {
		var err error
		decrypted, err = s.secretRepository.GetDecryptedRef()
		if err != nil {
			return err
		}
	}
	keySha1 := md5.Sum([]byte(key))
	ref, err := s.encryption.Encrypt(decrypted, keySha1)
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
	err = s.secretRepository.UpdateRefs(ref, decrypted)
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
	if s.IsRefEmpty() {
		err = errors.New("there is not secret")
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
	err := s.secretRepository.UpdateRefs("", "")

	err = s.passwordRepository.DeleteAll()
	err = s.certRepository.DeleteAll()
	err = s.clusterRepository.DeleteAll()
	err = s.compactTableRepository.DeleteAll()
	err = s.tagRepository.DeleteAll()

	s.mutex.Unlock()
	return err
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
