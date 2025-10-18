package config

import (
	"context"
	"encoding/json"
	"errors"
	"ivory/src/features/encryption"
	"ivory/src/features/secret"
	"ivory/src/storage/files"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Service struct {
	configFiles       *files.Storage
	encryptionService *encryption.Service
	secretService     *secret.Service

	appConfigFileName string
	appConfig         *AppConfig

	oauthConfig   *oauth2.Config
	oauthVerifier *oidc.IDTokenVerifier
}

func NewService(
	configFiles *files.Storage,
	encryptionService *encryption.Service,
	secretService *secret.Service,
) *Service {
	return &Service{
		configFiles:       configFiles,
		encryptionService: encryptionService,
		secretService:     secretService,

		appConfigFileName: "application",
	}
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
		return nil, errors.New("config is not specified, or it cannot be read from file")
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

	appConfig.Auth = authDecrypt
	s.appConfig = &appConfig
	return s.appConfig, nil
}

func (s *Service) GetOAuthConfig() (*oauth2.Config, *oidc.IDTokenVerifier, error) {
	if s.oauthConfig != nil {
		return s.oauthConfig, s.oauthVerifier, nil
	}

	appConfig, errConfig := s.GetAppConfig()
	if errConfig != nil {
		return nil, nil, errConfig
	}
	authConfig := appConfig.Auth
	if authConfig.Oidc == nil {
		return nil, nil, errors.New("oidc is not configured")
	}
	provider, errProvider := oidc.NewProvider(context.Background(), authConfig.Oidc.IssuerURL)
	if errProvider != nil {
		return nil, nil, errProvider
	}

	oauth2Config := oauth2.Config{
		ClientID:     authConfig.Oidc.ClientID,
		ClientSecret: authConfig.Oidc.ClientSecret,
		RedirectURL:  authConfig.Oidc.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	s.oauthVerifier = provider.Verifier(&oidc.Config{ClientID: authConfig.Oidc.ClientID})
	s.oauthConfig = &oauth2Config
	return s.oauthConfig, s.oauthVerifier, nil
}

func (s *Service) SetAppConfig(newAppConfig AppConfig) error {
	if newAppConfig.Company == "" {
		return errors.New("company name cannot be empty")
	}
	errValid := s.validateAuthConfig(newAppConfig.Auth)
	if errValid != nil {
		return errValid
	}

	// NOTE: we skip error checking here, because we need to know only if config is not set up yet
	currentAppConfig, _ := s.GetAppConfig()
	if currentAppConfig != nil {
		return errors.New("config is already set up")
	}

	encryptedAuthConfig, errAuthConfig := s.encryptAuthConfig(newAppConfig.Auth)
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

func (s *Service) encryptAuthConfig(authConfig AuthConfig) (AuthConfig, error) {
	switch authConfig.Type {
	case BASIC:
		encrypted, err := s.encrypt(authConfig.Basic.Password)
		if err != nil {
			return authConfig, err
		}
		authConfig.Basic.Password = encrypted
		return authConfig, nil
	case OIDC:
		encrypted, err := s.encrypt(authConfig.Oidc.ClientSecret)
		if err != nil {
			return authConfig, err
		}
		authConfig.Oidc.ClientSecret = encrypted
		return authConfig, nil
	case LDAP:
		encrypted, err := s.encrypt(authConfig.Ldap.BindPass)
		if err != nil {
			return authConfig, err
		}
		authConfig.Ldap.BindPass = encrypted
		return authConfig, nil
	default:
		return authConfig, nil
	}
}

func (s *Service) decryptAuthConfig(authConfig AuthConfig) (AuthConfig, error) {
	switch authConfig.Type {
	case BASIC:
		decrypted, err := s.decrypt(authConfig.Basic.Password)
		if err != nil {
			return authConfig, err
		}
		authConfig.Basic.Password = decrypted
		return authConfig, nil
	case OIDC:
		decrypted, err := s.decrypt(authConfig.Oidc.ClientSecret)
		if err != nil {
			return authConfig, err
		}
		authConfig.Oidc.ClientSecret = decrypted
		return authConfig, nil
	case LDAP:
		decrypted, err := s.decrypt(authConfig.Ldap.BindPass)
		if err != nil {
			return authConfig, err
		}
		authConfig.Ldap.BindPass = decrypted
		return authConfig, nil
	default:
		return authConfig, nil
	}
}

func (s *Service) ReEncryptAuthConfig(authConfig AuthConfig, oldSecret [16]byte, newSecret [16]byte) (AuthConfig, error) {
	switch authConfig.Type {
	case BASIC:
		encrypted, err := s.reencrypt(authConfig.Basic.Password, oldSecret, newSecret)
		if err != nil {
			return authConfig, err
		}
		authConfig.Basic.Password = encrypted
		return authConfig, nil
	case OIDC:
		encrypted, err := s.reencrypt(authConfig.Oidc.ClientSecret, oldSecret, newSecret)
		if err != nil {
			return authConfig, err
		}
		authConfig.Oidc.ClientSecret = encrypted
		return authConfig, nil
	case LDAP:
		encrypted, err := s.reencrypt(authConfig.Ldap.BindPass, oldSecret, newSecret)
		if err != nil {
			return authConfig, err
		}
		authConfig.Ldap.BindPass = encrypted
		return authConfig, nil
	default:
		return authConfig, nil
	}
}

func (s *Service) DeleteAll() error {
	s.appConfig = nil
	return s.configFiles.DeleteAll()
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

func (s *Service) validateAuthConfig(authConfig AuthConfig) error {
	switch authConfig.Type {
	case NONE:
		if authConfig.Oidc == nil && authConfig.Basic == nil && authConfig.Ldap == nil {
			return nil
		}
		return errors.New("all configuration fields must be empty (oidc, basic, ldap)")
	case BASIC:
		if authConfig.Oidc == nil && authConfig.Basic != nil && authConfig.Ldap == nil {
			return nil
		}
		if authConfig.Basic.Password == "" || authConfig.Basic.Username == "" {
			return errors.New("username or password are not specified")
		}
		return errors.New("only basic field should be configured")
	case OIDC:
		if authConfig.Oidc != nil && authConfig.Basic == nil && authConfig.Ldap == nil {
			return nil
		}
		if authConfig.Oidc.IssuerURL == "" {
			return errors.New("IssuerURL is not specified")
		}
		if authConfig.Oidc.ClientID == "" {
			return errors.New("ClientID is not specified")
		}
		if authConfig.Oidc.ClientSecret == "" {
			return errors.New("ClientSecret is not specified")
		}
		if authConfig.Oidc.RedirectURL == "" {
			return errors.New("RedirectURL is not specified")
		}
		return errors.New("only oidc field should be configured")
	case LDAP:
		if authConfig.Oidc == nil && authConfig.Basic == nil && authConfig.Ldap != nil {
			return nil
		}
		if authConfig.Ldap.Url == "" {
			return errors.New("url is not specified")
		}
		if authConfig.Ldap.BindDN == "" {
			return errors.New("BindDN is not specified")
		}
		if authConfig.Ldap.BindPass == "" {
			return errors.New("BindPass is not specified")
		}
		if authConfig.Ldap.BaseDN == "" {
			return errors.New("BaseDN is not specified")
		}
		return errors.New("only ldap field should be configured")
	default:
		return errors.New("the type is unknown")
	}
}
