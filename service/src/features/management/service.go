package management

import (
	"errors"
	"ivory/src/features/auth"
	"ivory/src/features/backup"
	"ivory/src/features/cert"
	"ivory/src/features/cluster"
	"ivory/src/features/config"
	"ivory/src/features/node"
	"ivory/src/features/permission"
	"ivory/src/features/query"
	"ivory/src/features/secret"
	"ivory/src/features/tag"
	"ivory/src/features/tools"
	"ivory/src/features/vault"
	"ivory/src/storage/env"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type Service struct {
	env               *env.AppEnv
	authService       *auth.Service
	vaultService      *vault.Service
	clusterService    *cluster.Service
	certService       *cert.Service
	tagService        *tag.Service
	toolsService      *tools.Service
	queryService      *query.Service
	nodeService       *node.Service
	secretService     *secret.Service
	configService     *config.Service
	permissionService *permission.Service
	backupService     *backup.Service
}

func NewService(
	env *env.AppEnv,
	authService *auth.Service,
	vaultService *vault.Service,
	clusterService *cluster.Service,
	certService *cert.Service,
	tagService *tag.Service,
	toolsService *tools.Service,
	queryService *query.Service,
	nodeService *node.Service,
	secretService *secret.Service,
	configService *config.Service,
	permissionService *permission.Service,
	backupService *backup.Service,
) *Service {
	return &Service{
		env:               env,
		authService:       authService,
		vaultService:      vaultService,
		toolsService:      toolsService,
		clusterService:    clusterService,
		certService:       certService,
		tagService:        tagService,
		queryService:      queryService,
		nodeService:       nodeService,
		secretService:     secretService,
		configService:     configService,
		permissionService: permissionService,
		backupService:     backupService,
	}
}

func (s *Service) Free() error {
	errComTable := s.toolsService.DeleteAll()
	errQuery := s.queryService.DeleteAllLogs()
	return errors.Join(errComTable, errQuery)
}

func (s *Service) Erase() error {
	errSecret := s.secretService.Clean()
	errCred := s.vaultService.DeleteAll()
	errCert := s.certService.DeleteAll()
	errCluster := s.clusterService.DeleteAll()
	errTools := s.toolsService.DeleteAll()
	errTag := s.tagService.DeleteAll()
	errQuery := s.queryService.DeleteAll()
	errConfig := s.configService.DeleteAll()
	errPerm := s.permissionService.DeleteAll()
	return errors.Join(errSecret, errCred, errCert, errCluster, errTools, errTag, errQuery, errConfig, errPerm)
}

func (s *Service) ChangeSecret(previousKey string, newKey string) error {
	prevSha, newSha, err := s.secretService.Update(previousKey, newKey)
	if err != nil {
		return err
	}
	errEnc := s.vaultService.Reencrypt(prevSha, newSha)
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

func (s *Service) BackupFileName() string {
	return s.backupService.GetFileName()
}

func (s *Service) BackupExport() ([]byte, error) {
	return s.backupService.Export()
}

func (s *Service) BackupImport(file *multipart.FileHeader) error {
	return s.backupService.Import(file)
}
