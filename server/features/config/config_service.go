package config

import (
	"encoding/json"
	"errors"
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/clients/storage"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/features/auth"
	"ivory/features/user"
)

var ErrConfigNotSpecified = errors.New("config is not specified")
var ErrConfigAlreadySet = errors.New("config is already set; to change it, you need to provide the correct secret")
var ErrCompanyNameEmpty = errors.New("company name cannot be empty")
var ErrSuperuserRequired = errors.New("a superuser is required to enable authentication")
var ErrSuperuserPasswordRequired = errors.New("a superuser registered for basic auth needs a password")

type Service struct {
	configFiles       *storage.FileStorage
	encryptionService *encryption.Service
	secretService     *secret.Service
	authService       *auth.Service
	userService       *user.Service
	basicProvider     *basic.Provider
	ldapProvider      *ldap.Provider
	oidcProvider      *oidc.Provider

	appConfig           *AppConfig
	appConfigFileName   string
	appConfigConfigured bool
}

func NewService(
	configFiles *storage.FileStorage,
	encryptionService *encryption.Service,
	secretService *secret.Service,
	authService *auth.Service,
	userService *user.Service,
	basicProvider *basic.Provider,
	ldapProvider *ldap.Provider,
	oidcProvider *oidc.Provider,
) *Service {
	return &Service{
		configFiles:       configFiles,
		encryptionService: encryptionService,
		secretService:     secretService,
		authService:       authService,
		userService:       userService,
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
		return nil, ErrConfigNotSpecified
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
		return ErrConfigAlreadySet
	}
	appConfig := newAppConfig.AppConfig
	if appConfig.Company == "" {
		return ErrCompanyNameEmpty
	}

	errValid := s.setAuthConfig(appConfig.Auth)
	if errValid != nil {
		s.clearCache()
		return errValid
	}
	// NOTE: the superuser is created before the config is saved. The other way
	// round, a failure here would leave authentication switched on with nobody
	// able to sign in, and no way left to fix it from the outside
	errSuperuser := s.createSetupUser(appConfig.Auth, newAppConfig.User)
	if errSuperuser != nil {
		s.clearCache()
		return errSuperuser
	}
	encryptedAuthConfig, errAuthConfig := s.encryptAuthConfig(appConfig.Auth)
	if errAuthConfig != nil {
		s.clearCache()
		return errAuthConfig
	}

	appConfig.Auth = encryptedAuthConfig
	errSave := s.saveAppConfig(appConfig)
	if errSave != nil {
		s.clearCache()
		return errSave
	}
	return nil
}

// createSetupUser makes sure enabling authentication leaves somebody who can
// administer Ivory. Setup is the only place a user is registered with a password
// somebody else typed, and the password is asked for only when that user is
// registered for basic - an Ivory signed into through LDAP or SSO alone has no
// use for one.
func (s *Service) createSetupUser(authConfig AuthConfig, request *UserSetupRequest) error {
	if !authConfig.enabled() {
		return nil
	}
	if request == nil {
		hasSuperuser, err := s.userService.HasSuperuser()
		if err != nil {
			return err
		}
		if !hasSuperuser {
			return ErrSuperuserRequired
		}
		return nil
	}
	password := ""
	for _, authType := range request.AuthTypes {
		if authType == user.AuthBasic {
			if request.Password == "" {
				return ErrSuperuserPasswordRequired
			}
			password = request.Password
		}
	}
	_, err := s.userService.CreateOutright(request.Username, password, request.AuthTypes, true)
	return err
}

func (s *Service) Reencrypt() error {
	newAuthConfig, err := s.encryptAuthConfig(s.appConfig.Auth)
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
	s.clearCache()
	return s.configFiles.DeleteAll()
}

func (s *Service) clearCache() {
	s.appConfig = nil
	s.basicProvider.DeleteConfig()
	s.oidcProvider.DeleteConfig()
	s.ldapProvider.DeleteConfig()
}

func (s *Service) encryptAuthConfig(authConfig AuthConfig) (AuthConfig, error) {
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

func (s *Service) setAuthConfig(authConfig AuthConfig) error {
	if authConfig.Basic == nil && authConfig.Ldap == nil && authConfig.Oidc == nil {
		return nil
	}
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
