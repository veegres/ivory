package management

import (
	"errors"
	"ivory/src/features/auth"
	"ivory/src/features/bloat"
	"ivory/src/features/cert"
	"ivory/src/features/cluster"
	"ivory/src/features/config"
	"ivory/src/features/password"
	"ivory/src/features/query"
	"ivory/src/features/secret"
	"ivory/src/features/tag"
	"ivory/src/storage/env"

	"github.com/gin-gonic/gin"
)

type Service struct {
	env             *env.AppEnv
	authService     *auth.Service
	passwordService *password.Service
	clusterService  *cluster.Service
	certService     *cert.Service
	tagService      *tag.Service
	bloatService    *bloat.Service
	queryService    *query.Service
	secretService   *secret.Service
	configService   *config.Service
}

func NewService(
	env *env.AppEnv,
	authService *auth.Service,
	passwordService *password.Service,
	clusterService *cluster.Service,
	certService *cert.Service,
	tagService *tag.Service,
	bloatService *bloat.Service,
	queryService *query.Service,
	secretService *secret.Service,
	configService *config.Service,
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
		configService:   configService,
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
	errConfig := s.configService.DeleteAll()

	return errors.Join(errSecret, errPass, errCert, errCluster, errComTable, errTag, errQuery, errConfig)
}

func (s *Service) ChangeSecret(previousKey string, newKey string) error {
	prevSha, newSha, err := s.secretService.Update(previousKey, newKey)
	if err != nil {
		return err
	}
	errEnc := s.passwordService.Reencrypt(prevSha, newSha)
	if errEnc != nil {
		return errEnc
	}
	errConfig := s.configService.Reencrypt(prevSha, newSha)
	if errConfig != nil {
		return errConfig
	}
	return nil
}

func (s *Service) GetAppInfo(context *gin.Context) *AppInfo {
	appConfig, err := s.configService.GetAppConfig()
	if err != nil {
		return &AppInfo{
			Config: config.Info{
				Configured:   false,
				Company:      "Ivory",
				Availability: config.Availability{ManualQuery: false},
				Error:        err.Error(),
			},
			Secret:  s.secretService.Status(),
			Version: s.env.Version,
			Auth: auth.Info{
				Type:       auth.NONE,
				Authorised: false,
				Error:      "",
			},
		}
	}
	authorised, authError := s.authService.ValidateAuthToken(context)
	errorMessage := ""
	if authError != nil {
		errorMessage = authError.Error()
	}

	return &AppInfo{
		Config: config.Info{
			Configured:   true,
			Company:      appConfig.Company,
			Availability: appConfig.Availability,
			Error:        "",
		},
		Secret:  s.secretService.Status(),
		Version: s.env.Version,
		Auth: auth.Info{
			Type:       appConfig.Auth.Type,
			Authorised: authorised,
			Error:      errorMessage,
		},
	}
}
