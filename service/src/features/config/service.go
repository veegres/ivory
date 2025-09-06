package config

import (
	"encoding/json"
	"errors"
	"ivory/src/features/encryption"
	"ivory/src/features/secret"
	"ivory/src/storage/files"
)

type Service struct {
	configFiles       *files.Storage
	encryptionService *encryption.Service
	secretService     *secret.Service

	appConfigFileName string
	appConfig         *AppConfig
}

func NewService(
	configFiles *files.Storage,
	encryptionService *encryption.Service,
	secretService *secret.Service,
) *Service {
	service := &Service{
		configFiles:       configFiles,
		encryptionService: encryptionService,
		secretService:     secretService,

		appConfigFileName: "application",
	}

	appConfig, _ := service.GetAppConfig()
	service.appConfig = appConfig

	return service
}

func (s *Service) UpdateAppConfig(config AppConfig) error {
	file, errOpen := s.configFiles.OpenByName(s.appConfigFileName)
	if errOpen != nil {
		return errOpen
	}

	errTruncate := file.Truncate(0)
	if errTruncate != nil {
		return errTruncate
	}

	jsonAuth, errMarshall := json.MarshalIndent(config, "", "  ")
	if errMarshall != nil {
		return errMarshall
	}

	_, errWrite := file.Write(jsonAuth)
	if errWrite != nil {
		return errWrite
	}
	return nil
}

func (s *Service) GetAppConfig() (*AppConfig, error) {
	if s.appConfig != nil {
		return s.appConfig, nil
	}

	read, err := s.configFiles.ReadByName(s.appConfigFileName)
	if err != nil {
		return nil, errors.New("config is not specified or it cannot be read from file")
	}

	var appConfig AppConfig
	errUnmarshal := json.Unmarshal(read, &appConfig)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	s.appConfig = &appConfig
	return s.appConfig, nil
}

func (s *Service) SetAppConfig(newAppConfig AppConfig) error {
	if newAppConfig.Company == "" {
		return errors.New("company name cannot be empty")
	}

	// NOTE: we skip error checking here, because we need to know only if config is not set up yet
	currentAppConfig, _ := s.GetAppConfig()
	if currentAppConfig != nil {
		return errors.New("config is already set up")
	}

	encryptedAuthConfig, errAuthConfig := s.EncryptAuthConfig(newAppConfig.Auth)
	if errAuthConfig != nil {
		return errAuthConfig
	}

	updatedAppConfig := AppConfig{
		Company:      newAppConfig.Company,
		Availability: newAppConfig.Availability,
		Auth:         encryptedAuthConfig,
	}

	_, errCreate := s.configFiles.CreateByName(s.appConfigFileName)
	if errCreate != nil {
		return errCreate
	}

	return s.UpdateAppConfig(updatedAppConfig)
}

func (s *Service) EncryptAuthConfig(authConfig Auth) (Auth, error) {
	switch authConfig.Type {
	case NONE:
		return authConfig, nil
	case BASIC:
		user := authConfig.Body["username"]
		pass := authConfig.Body["password"]
		if pass == "" || user == "" {
			return authConfig, errors.New("username or password are not specified")
		}

		encryptedPassword, errEnc := s.encryptionService.Encrypt(pass, s.secretService.Get())
		if errEnc != nil {
			return authConfig, errEnc
		}
		authConfig.Body["password"] = encryptedPassword

		return authConfig, nil
	default:
		return authConfig, errors.New("type is unknown")
	}
}

func (s *Service) ReEncryptAuthConfig(authConfig Auth, oldSecret [16]byte, newSecret [16]byte) (Auth, error) {
	switch authConfig.Type {
	case NONE:
		return authConfig, nil
	case BASIC:
		user := authConfig.Body["username"]
		pass := authConfig.Body["password"]
		if pass == "" || user == "" {
			return authConfig, errors.New("username or password are not specified")
		}

		decryptedPassword, errDec := s.encryptionService.Decrypt(pass, oldSecret)
		if errDec != nil {
			return authConfig, errDec
		}
		encryptedPassword, errEnc := s.encryptionService.Encrypt(decryptedPassword, newSecret)
		if errEnc != nil {
			return authConfig, errEnc
		}
		authConfig.Body["password"] = encryptedPassword

		return authConfig, nil
	default:
		return authConfig, errors.New("the type is unknown")
	}
}

func (s *Service) DeleteAll() error {
	s.appConfig = nil
	return s.configFiles.DeleteAll()
}
