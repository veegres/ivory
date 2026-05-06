package vault

import (
	"ivory/src/clients/ssh"
	"ivory/src/features/encryption"
	"ivory/src/features/secret"

	"github.com/google/uuid"
)

type Service struct {
	vaultRepository *Repository
	sshClient       *ssh.Client
	secretService   *secret.Service
	encryption      *encryption.Service

	hiddenSecret string
}

func NewService(
	vaultRepository *Repository,
	sshClient *ssh.Client,
	secretService *secret.Service,
	encryption *encryption.Service,
) *Service {
	return &Service{
		vaultRepository: vaultRepository,
		sshClient:       sshClient,
		secretService:   secretService,
		encryption:      encryption,
		hiddenSecret:    "********",
	}
}

func (s *Service) Create(vault Vault) (*uuid.UUID, *Vault, error) {
	sc, mt, err := s.generateVaultByType(vault)
	if err != nil {
		return nil, nil, err
	}
	encryptedCred := Vault{Type: vault.Type, Username: vault.Username, Secret: sc, Metadata: mt}
	key, errCreate := s.vaultRepository.Create(encryptedCred)

	cred := encryptedCred
	cred.Secret = s.hiddenSecret
	return &key, &cred, errCreate
}

func (s *Service) Update(key uuid.UUID, vault Vault) (*uuid.UUID, *Vault, error) {
	sc, mt, err := s.generateVaultByType(vault)
	if err != nil {
		return nil, nil, err
	}

	encryptedCred := Vault{Username: vault.Username, Secret: sc, Type: vault.Type, Metadata: mt}
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

func (s *Service) generateVaultByType(vault Vault) (string, *string, error) {
	switch vault.Type {
	case SSH_KEY:
		pubKey, prvKey, err := s.sshClient.GenerateKey()
		if err != nil {
			return "", nil, err
		}
		encryptedSecret, errEnc := s.encryption.Encrypt(prvKey, s.secretService.Get())
		if errEnc != nil {
			return "", nil, errEnc
		}
		return encryptedSecret, &pubKey, nil
	default:
		encryptedSecret, errEnc := s.encryption.Encrypt(vault.Secret, s.secretService.Get())
		if errEnc != nil {
			return "", nil, errEnc
		}
		return encryptedSecret, nil, errEnc
	}
}
