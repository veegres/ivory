package vm

import (
	"errors"
	"ivory/src/features/encryption"
	"ivory/src/features/secret"
	"strings"

	"github.com/google/uuid"
)

var ErrVmNameEmpty = errors.New("vm name cannot be empty")
var ErrVmHostEmpty = errors.New("vm host cannot be empty")
var ErrVmUsernameEmpty = errors.New("vm username cannot be empty")
var ErrVmPortInvalid = errors.New("vm ssh port should be greater than 0")
var ErrVmSshKeyEmpty = errors.New("vm ssh key cannot be empty")

type Service struct {
	repository   *Repository
	secret       *secret.Service
	encryption   *encryption.Service
	hiddenSshKey string
}

func NewService(
	repository *Repository,
	secret *secret.Service,
	encryption *encryption.Service,
) *Service {
	return &Service{
		repository:   repository,
		secret:       secret,
		encryption:   encryption,
		hiddenSshKey: "********",
	}
}

func (s *Service) List() ([]VM, error) {
	list, err := s.repository.List()
	if err != nil {
		return nil, err
	}
	return s.hideSshKeys(list), nil
}

func (s *Service) Get(id uuid.UUID) (*VM, error) {
	model, err := s.repository.Get(id)
	if err != nil {
		return nil, err
	}
	hidden := s.hideSshKey(model)
	return &hidden, nil
}

func (s *Service) Create(model VM) (*uuid.UUID, *VM, error) {
	if err := validate(model); err != nil {
		return nil, nil, err
	}

	model.Id = uuid.Nil
	encrypted, err := s.encryptKey(model.SshKey)
	if err != nil {
		return nil, nil, err
	}
	model.SshKey = encrypted

	id, err := s.repository.Create(model)
	if err != nil {
		return nil, nil, err
	}

	model.Id = id
	hidden := s.hideSshKey(model)
	return &id, &hidden, nil
}

func (s *Service) Update(id uuid.UUID, model VM) (*uuid.UUID, *VM, error) {
	stored, err := s.repository.Get(id)
	if err != nil {
		return nil, nil, err
	}

	if model.SshKey == "" || model.SshKey == s.hiddenSshKey {
		model.SshKey = stored.SshKey
	} else {
		encrypted, err := s.encryptKey(model.SshKey)
		if err != nil {
			return nil, nil, err
		}
		model.SshKey = encrypted
	}

	model.Id = id
	if err := validate(model); err != nil {
		return nil, nil, err
	}

	_, err = s.repository.Update(id, model)
	if err != nil {
		return nil, nil, err
	}

	hidden := s.hideSshKey(model)
	return &id, &hidden, nil
}

func (s *Service) Delete(id uuid.UUID) error {
	return s.repository.Delete(id)
}

func (s *Service) DeleteAll() error {
	return s.repository.DeleteAll()
}

func (s *Service) Reencrypt(oldSecret [16]byte, newSecret [16]byte) error {
	models, err := s.repository.Map()
	if err != nil {
		return err
	}

	for id, model := range models {
		decrypted, err := s.encryption.Decrypt(model.SshKey, oldSecret)
		if err != nil {
			return err
		}
		encrypted, err := s.encryption.Encrypt(decrypted, newSecret)
		if err != nil {
			return err
		}
		model.SshKey = encrypted
		parsedId, err := uuid.Parse(id)
		if err != nil {
			return err
		}
		_, err = s.repository.Update(parsedId, model)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetDecrypted(id uuid.UUID) (*VM, error) {
	model, err := s.repository.Get(id)
	if err != nil {
		return nil, err
	}
	sshKeyDecrypted, err := s.encryption.Decrypt(model.SshKey, s.secret.Get())
	if err != nil {
		return nil, err
	}
	model.SshKey = sshKeyDecrypted
	return &model, nil
}

func (s *Service) encryptKey(key string) (string, error) {
	return s.encryption.Encrypt(key, s.secret.Get())
}

func (s *Service) hideSshKeys(list []VM) []VM {
	result := make([]VM, 0, len(list))
	for _, model := range list {
		result = append(result, s.hideSshKey(model))
	}
	return result
}

func (s *Service) hideSshKey(model VM) VM {
	model.SshKey = s.hiddenSshKey
	return model
}

func validate(model VM) error {
	switch {
	case strings.TrimSpace(model.Name) == "":
		return ErrVmNameEmpty
	case strings.TrimSpace(model.Host) == "":
		return ErrVmHostEmpty
	case strings.TrimSpace(model.Username) == "":
		return ErrVmUsernameEmpty
	case model.SshPort <= 0:
		return ErrVmPortInvalid
	case strings.TrimSpace(model.SshKey) == "":
		return ErrVmSshKeyEmpty
	default:
		return nil
	}
}
