package config

import (
	"encoding/json"
	"errors"
	"ivory/src/clients/auth/basic"
	"ivory/src/clients/auth/ldap"
	"ivory/src/clients/auth/oidc"
	"ivory/src/features/auth"
	"ivory/src/features/encryption"
	"ivory/src/features/secret"
	"ivory/src/storage/files"
)

type Service struct {
	configFiles       *files.Storage
	encryptionService *encryption.Service
	secretService     *secret.Service
	authService       *auth.Service
	basicProvider     *basic.Provider
	ldapProvider      *ldap.Provider
	oidcProvider      *oidc.Provider

	appConfig           *AppConfig
	appConfigFileName   string
	appConfigConfigured bool
}

func NewService(
	configFiles *files.Storage,
	encryptionService *encryption.Service,
	secretService *secret.Service,
	authService *auth.Service,
	basicProvider *basic.Provider,
	ldapProvider *ldap.Provider,
	oidcProvider *oidc.Provider,
) *Service {
	return &Service{
		configFiles:       configFiles,
		encryptionService: encryptionService,
		secretService:     secretService,
		authService:       authService,
		basicProvider:     basicProvider,
		ldapProvider:      ldapProvider,
		oidcProvider:      oidcProvider,

		appConfigFileName:   "application",
		appConfigConfigured: false,
	}
}

func (s *Service) initializeAppConfig() (*AppConfig, error) {
	s.appConfigConfigured = s.configFiles.ExistByName(s.appConfigFileName)
	read, err := s.configFiles.ReadByName(s.appConfigFileName)
	if err != nil {
		return nil, errors.New("config is not specified")
	}
	var appConfig AppConfig
	errUnmarshal := json.Unmarshal(read, &appConfig)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	authDecrypt, errDecrypt := s.decryptAuthConfig(appConfig.Auth)
	if errDecrypt != nil {
		return nil, errDecrypt
	}
	errSetAuthConfig := s.setAuthConfig(authDecrypt)
	if errSetAuthConfig != nil {
		return nil, errSetAuthConfig
	}
	appConfig.Auth = authDecrypt
	return &appConfig, nil
}

func (s *Service) saveAppConfig(config AppConfig) error {
	if !s.configFiles.ExistByName(s.appConfigFileName) {
		_, errCreate := s.configFiles.CreateByName(s.appConfigFileName)
		if errCreate != nil {
			return errCreate
		}
	}

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

func (s *Service) GetIsConfigured() bool {
	return s.appConfigConfigured
}

func (s *Service) GetAppConfig() (*AppConfig, error) {
	if s.appConfig != nil {
		return s.appConfig, nil
	}
	config, err := s.initializeAppConfig()
	if err != nil {
		return nil, err
	}
	s.appConfig = config
	return s.appConfig, nil
}

func (s *Service) SetAppConfig(newAppConfig NewAppConfig) error {
	if s.GetIsConfigured() && !s.secretService.Verify(newAppConfig.Secret) {
		return errors.New("config is already set; to change it, you need to provide the correct secret")
	}
	appConfig := newAppConfig.AppConfig
	if appConfig.Company == "" {
		return errors.New("company name cannot be empty")
	}

	errValid := s.setAuthConfig(appConfig.Auth)
	if errValid != nil {
		return errValid
	}
	encryptedAuthConfig, errAuthConfig := s.encryptAuthConfig(appConfig.Auth)
	if errAuthConfig != nil {
		return errAuthConfig
	}

	appConfig.Auth = encryptedAuthConfig
	return s.saveAppConfig(appConfig)
}

func (s *Service) Reencrypt(oldSecret [16]byte, newSecret [16]byte) error {
	newAuthConfig, err := s.reencryptAuthConfig(s.appConfig.Auth, oldSecret, newSecret)
	if err != nil {
		return err
	}
	// NOTE: we do not to touch cashed `s.appConfig` because it is decrypted
	newAppConfig := *s.appConfig
	newAppConfig.Auth = newAuthConfig
	err = s.saveAppConfig(newAppConfig)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteAll() error {
	s.appConfig = nil
	s.authService.SetAuthorisation(false)
	s.basicProvider.DeleteConfig()
	s.oidcProvider.DeleteConfig()
	s.ldapProvider.DeleteConfig()
	return s.configFiles.DeleteAll()
}

func (s *Service) reencryptAuthConfig(authConfig AuthConfig, oldSecret [16]byte, newSecret [16]byte) (AuthConfig, error) {
	if authConfig.Basic != nil {
		encrypted, err := s.reencrypt(authConfig.Basic.Password, oldSecret, newSecret)
		if err != nil {
			return authConfig, err
		}
		tmp := *authConfig.Basic
		tmp.Password = encrypted
		authConfig.Basic = &tmp
	}
	if authConfig.Oidc != nil {
		encrypted, err := s.reencrypt(authConfig.Oidc.ClientSecret, oldSecret, newSecret)
		if err != nil {
			return authConfig, err
		}
		tmp := *authConfig.Oidc
		tmp.ClientSecret = encrypted
		authConfig.Oidc = &tmp
	}
	if authConfig.Ldap != nil {
		encrypted, err := s.reencrypt(authConfig.Ldap.BindPass, oldSecret, newSecret)
		if err != nil {
			return authConfig, err
		}
		tmp := *authConfig.Ldap
		tmp.BindPass = encrypted
		authConfig.Ldap = &tmp
	}
	return authConfig, nil
}

func (s *Service) encryptAuthConfig(authConfig AuthConfig) (AuthConfig, error) {
	if authConfig.Basic != nil {
		encrypted, err := s.encrypt(authConfig.Basic.Password)
		if err != nil {
			return authConfig, err
		}
		tmp := *authConfig.Basic
		tmp.Password = encrypted
		authConfig.Basic = &tmp
	}
	if authConfig.Oidc != nil {
		encrypted, err := s.encrypt(authConfig.Oidc.ClientSecret)
		if err != nil {
			return authConfig, err
		}
		tmp := *authConfig.Oidc
		tmp.ClientSecret = encrypted
		authConfig.Oidc = &tmp
	}
	if authConfig.Ldap != nil {
		encrypted, err := s.encrypt(authConfig.Ldap.BindPass)
		if err != nil {
			return authConfig, err
		}
		tmp := *authConfig.Ldap
		tmp.BindPass = encrypted
		authConfig.Ldap = &tmp
	}
	return authConfig, nil
}

func (s *Service) decryptAuthConfig(authConfig AuthConfig) (AuthConfig, error) {
	if authConfig.Basic != nil {
		decrypted, err := s.decrypt(authConfig.Basic.Password)
		if err != nil {
			return authConfig, err
		}
		tmp := *authConfig.Basic
		tmp.Password = decrypted
		authConfig.Basic = &tmp
	}
	if authConfig.Oidc != nil {
		decrypted, err := s.decrypt(authConfig.Oidc.ClientSecret)
		if err != nil {
			return authConfig, err
		}
		tmp := *authConfig.Oidc
		tmp.ClientSecret = decrypted
		authConfig.Oidc = &tmp
	}
	if authConfig.Ldap != nil {
		decrypted, err := s.decrypt(authConfig.Ldap.BindPass)
		if err != nil {
			return authConfig, err
		}
		tmp := *authConfig.Ldap
		tmp.BindPass = decrypted
		authConfig.Ldap = &tmp
	}
	return authConfig, nil
}

func (s *Service) encrypt(str string) (string, error) {
	encrypted, errEnc := s.encryptionService.Encrypt(str, s.secretService.Get())
	if errEnc != nil {
		return "", errEnc
	}
	return encrypted, nil
}

func (s *Service) decrypt(str string) (string, error) {
	decrypted, errEnc := s.encryptionService.Decrypt(str, s.secretService.Get())
	if errEnc != nil {
		return "", errEnc
	}
	return decrypted, nil
}

func (s *Service) reencrypt(str string, oldSecret [16]byte, newSecret [16]byte) (string, error) {
	decrypted, errDec := s.encryptionService.Decrypt(str, oldSecret)
	if errDec != nil {
		return "", errDec
	}
	encrypted, errEnc := s.encryptionService.Encrypt(decrypted, newSecret)
	if errEnc != nil {
		return "", errEnc
	}
	return encrypted, nil
}

func (s *Service) setAuthConfig(authConfig AuthConfig) error {
	if authConfig.Type == auth.NONE {
		return nil
	}
	s.authService.SetAuthorisation(true)
	var err error
	if authConfig.Basic != nil {
		errTmp := s.basicProvider.SetConfig(*authConfig.Basic)
		if errTmp != nil {
			err = errors.Join(err, errTmp)
		}

	}
	if authConfig.Oidc != nil {
		errTmp := s.oidcProvider.SetConfig(*authConfig.Oidc)
		if errTmp != nil {
			err = errors.Join(err, errTmp)
		}
	}
	if authConfig.Ldap != nil {
		errTmp := s.ldapProvider.SetConfig(*authConfig.Ldap)
		if errTmp != nil {
			err = errors.Join(err, errTmp)
		}
	}
	return err
}
