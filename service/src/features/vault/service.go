package vault

import (
	"ivory/src/features/encryption"
	"ivory/src/features/secret"

	"github.com/google/uuid"
)

type Service struct {
	vaultRepository *Repository
	secretService   *secret.Service
	encryption      *encryption.Service

	hiddenSecret string
}

func NewService(
	vaultRepository *Repository,
	secretService *secret.Service,
	encryption *encryption.Service,
) *Service {
	return &Service{
		vaultRepository: vaultRepository,
		secretService:   secretService,
		encryption:      encryption,
		hiddenSecret:    "********",
	}
}

func (s *Service) Create(vault Vault) (*uuid.UUID, *Vault, error) {
	encryptedSecret, errEnc := s.encryption.Encrypt(vault.Secret, s.secretService.Get())
	if errEnc != nil {
		return nil, nil, errEnc
	}

	encryptedCred := Vault{Username: vault.Username, Secret: encryptedSecret, Type: vault.Type}
	key, errCreate := s.vaultRepository.Create(encryptedCred)

	cred := encryptedCred
	cred.Secret = s.hiddenSecret
	return &key, &cred, errCreate
}

func (s *Service) Update(key uuid.UUID, vault Vault) (*uuid.UUID, *Vault, error) {
	encryptedSecret, errEnc := s.encryption.Encrypt(vault.Secret, s.secretService.Get())
	if errEnc != nil {
		return nil, nil, errEnc
	}

	encryptedCred := Vault{Username: vault.Username, Secret: encryptedSecret, Type: vault.Type}
	key, errUpdate := s.vaultRepository.Update(key, encryptedCred)

	cred := encryptedCred
	cred.Secret = s.hiddenSecret
	return &key, &cred, errUpdate
}

func (s *Service) Delete(key uuid.UUID) error {
	return s.vaultRepository.Delete(key)
}

func (s *Service) DeleteAll() error {
	return s.vaultRepository.DeleteAll()
}

func (s *Service) Map(vaultType *VaultType) (VaultMap, error) {
	var err error
	var vaultMap VaultMap
	if vaultType != nil {
		vaultMap, err = s.vaultRepository.MapByType(*vaultType)
	} else {
		vaultMap, err = s.vaultRepository.Map()
	}
	vaultHiddenMap := s.hideSecrets(vaultMap)
	return vaultHiddenMap, err
}

func (s *Service) Get(uuid uuid.UUID) (*Vault, error) {
	cred, err := s.vaultRepository.Get(uuid)
	if err != nil {
		return nil, err
	}
	cred.Secret = s.hiddenSecret
	return &cred, nil
}

func (s *Service) GetDecrypted(uuid uuid.UUID) (*Vault, error) {
	cred, errCred := s.vaultRepository.Get(uuid)
	if errCred != nil {
		return nil, errCred
	}
	decryptedSecret, errDec := s.encryption.Decrypt(cred.Secret, s.secretService.Get())
	if errDec != nil {
		return nil, errDec
	}
	cred.Secret = decryptedSecret
	return &cred, nil
}

func (s *Service) Reencrypt(oldSecret [16]byte, newSecret [16]byte) error {
	vaultMap, errGet := s.vaultRepository.Map()
	if errGet != nil {
		return errGet
	}
	for id, vault := range vaultMap {
		id, errParse := uuid.Parse(id)
		if errParse != nil {
			return errParse
		}
		oldEncryptedSecret := vault.Secret
		decryptedSecret, errDec := s.encryption.Decrypt(oldEncryptedSecret, oldSecret)
		if errDec != nil {
			return errDec
		}
		newEncryptedSecret, errEnc := s.encryption.Encrypt(decryptedSecret, newSecret)
		if errEnc != nil {
			return errEnc
		}
		vault.Secret = newEncryptedSecret
		_, errUpdate := s.vaultRepository.Update(id, vault)
		if errUpdate != nil {
			return errUpdate
		}
	}
	return nil
}

func (s *Service) hideSecrets(vaultMap VaultMap) VaultMap {
	for key, vault := range vaultMap {
		vault.Secret = s.hiddenSecret
		vaultMap[key] = vault
	}
	return vaultMap
}
