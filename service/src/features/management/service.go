package management

import (
	"errors"
	"ivory/src/features/auth"
	"ivory/src/features/bloat"
	"ivory/src/features/cert"
	"ivory/src/features/cluster"
	"ivory/src/features/config"
	"ivory/src/features/encryption"
	"ivory/src/features/password"
	"ivory/src/features/query"
	"ivory/src/features/secret"
	"ivory/src/features/tag"
	"ivory/src/storage/env"
)

type Service struct {
	env             *env.AppInfo
	authService     *auth.Service
	passwordService *password.Service
	clusterService  *cluster.Service
	certService     *cert.Service
	tagService      *tag.Service
	bloatService    *bloat.Service
	queryService    *query.Service
	secretService   *secret.Service
	encryption      *encryption.Service
	config          *config.Service
}

func NewService(
	env *env.AppInfo,
	authService *auth.Service,
	passwordService *password.Service,
	clusterService *cluster.Service,
	certService *cert.Service,
	tagService *tag.Service,
	bloatService *bloat.Service,
	queryService *query.Service,
	secretService *secret.Service,
	config *config.Service,
) *Service {
	return &Service{
		env:             env,
		authService:     authService,
		passwordService: passwordService,
		bloatService:    bloatService,
		clusterService:  clusterService,
		certService:     certService,
		tagService:      tagService,
		queryService:    queryService,
		secretService:   secretService,
		config:          config,
	}
}

func (s *Service) Erase() error {
	errSecret := s.secretService.Clean()
	errPass := s.passwordService.DeleteAll()
	errCert := s.certService.DeleteAll()
	errCluster := s.clusterService.DeleteAll()
	errComTable := s.bloatService.DeleteAll()
	errTag := s.tagService.DeleteAll()
	errQuery := s.queryService.DeleteAll()
	errConfig := s.config.DeleteAll()

	return errors.Join(errSecret, errPass, errCert, errCluster, errComTable, errTag, errQuery, errConfig)
}

func (s *Service) ChangeSecret(previousKey string, newKey string) error {
	prevSha, newSha, err := s.secretService.Update(previousKey, newKey)
	if err != nil {
		return err
	}

	errEnc := s.passwordService.ReEncryptPasswords(prevSha, newSha)
	if errEnc != nil {
		return errEnc
	}

	appConfig, errAppConfig := s.config.GetAppConfig()
	if errAppConfig != nil {
		return errAppConfig
	}
	reEncryptAuthConfig, errAuthConfig := s.config.ReEncryptAuthConfig(appConfig.Auth, prevSha, newSha)
	if errAuthConfig != nil {
		return errAuthConfig
	}
	updatedAppConfig := config.AppConfig{
		Company:      appConfig.Company,
		Availability: appConfig.Availability,
		Auth:         reEncryptAuthConfig,
	}
	errUpdateAppConfig := s.config.UpdateAppConfig(updatedAppConfig)
	if errUpdateAppConfig != nil {
		return errUpdateAppConfig
	}

	return nil
}

func (s *Service) GetAppInfo(authHeader string) *AppInfo {
	appConfig, err := s.config.GetAppConfig()
	if err != nil {
		return &AppInfo{
			Company:      "",
			Configured:   false,
			Secret:       s.secretService.Status(),
			Version:      s.env.Version,
			Availability: config.Availability{ManualQuery: false},
			Auth: auth.Info{
				Type:       config.NONE,
				Authorised: false,
				Error:      "",
			},
		}
	}
	authorised, authError := s.authService.ValidateAuthHeader(authHeader, appConfig.Auth)

	return &AppInfo{
		Company:      appConfig.Company,
		Configured:   true,
		Secret:       s.secretService.Status(),
		Version:      s.env.Version,
		Availability: appConfig.Availability,
		Auth: auth.Info{
			Type:       appConfig.Auth.Type,
			Authorised: authorised,
			Error:      authError,
		},
	}
}
