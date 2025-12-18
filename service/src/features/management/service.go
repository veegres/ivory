package management

import (
	"errors"
	"ivory/src/features/auth"
	"ivory/src/features/bloat"
	"ivory/src/features/cert"
	"ivory/src/features/cluster"
	"ivory/src/features/config"
	"ivory/src/features/password"
	"ivory/src/features/permission"
	"ivory/src/features/query"
	"ivory/src/features/secret"
	"ivory/src/features/tag"
	"ivory/src/storage/env"

	"github.com/gin-gonic/gin"
)

type Service struct {
	env               *env.AppEnv
	authService       *auth.Service
	passwordService   *password.Service
	clusterService    *cluster.Service
	certService       *cert.Service
	tagService        *tag.Service
	bloatService      *bloat.Service
	queryService      *query.Service
	queryLogService   *query.LogService
	secretService     *secret.Service
	configService     *config.Service
	permissionService *permission.Service
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
	queryLogService *query.LogService,
	secretService *secret.Service,
	configService *config.Service,
	permissionService *permission.Service,
) *Service {
	return &Service{
		env:               env,
		authService:       authService,
		passwordService:   passwordService,
		bloatService:      bloatService,
		clusterService:    clusterService,
		certService:       certService,
		tagService:        tagService,
		queryService:      queryService,
		queryLogService:   queryLogService,
		secretService:     secretService,
		configService:     configService,
		permissionService: permissionService,
	}
}

func (s *Service) Free() error {
	errComTable := s.bloatService.DeleteAll()
	errQuery := s.queryLogService.DeleteAll()
	return errors.Join(errComTable, errQuery)
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
	errPerm := s.permissionService.DeleteAll()
	return errors.Join(errSecret, errPass, errCert, errCluster, errComTable, errTag, errQuery, errConfig, errPerm)
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
	errConfig := s.configService.Reencrypt()
	if errConfig != nil {
		return errConfig
	}
	return nil
}

func (s *Service) GetAppInfo(context *gin.Context) *AppInfo {
	appConfig, errConfig := s.configService.GetAppConfig()
	configConfigured := s.configService.GetIsConfigured()
	authSupported := s.authService.GetSupportedTypes()
	if errConfig != nil {
		return &AppInfo{
			Config: ConfigInfo{
				Configured: configConfigured,
				Company:    "Ivory",
				Error:      errConfig.Error(),
			},
			Secret:  s.secretService.Status(),
			Version: s.env.Version,
			Auth: AuthInfo{
				Supported:  authSupported,
				Authorised: false,
				User:       nil,
				Error:      "",
			},
		}
	}

	authorised, user, authError := s.getAuthInfo(context)
	return &AppInfo{
		Config: ConfigInfo{
			Configured: configConfigured,
			Company:    appConfig.Company,
			Error:      "",
		},
		Secret:  s.secretService.Status(),
		Version: s.env.Version,
		Auth: AuthInfo{
			Supported:  authSupported,
			Authorised: authorised,
			User:       user,
			Error:      authError,
		},
	}
}

func (s *Service) getAuthInfo(context *gin.Context) (bool, *UserInfo, string) {
	authorised, username, authType, errParse := s.authService.ParseAuthToken(context)
	if username == "" {
		if errParse != nil {
			return authorised, nil, errParse.Error()
		}
		return authorised, nil, ""
	}
	prefix := ""
	if authType != nil {
		prefix = authType.String()
	}
	permissions, errPerm := s.permissionService.GetUserPermissions(prefix, username, errors.Is(errParse, auth.ErrAuthDisabled))
	user := &UserInfo{Username: username, Permissions: permissions}
	if errPerm != nil {
		return authorised, user, errPerm.Error()
	}
	return authorised, user, ""
}
