package management

import (
	"errors"
	coreConfig "ivory/core/config"
	"ivory/core/service/cert"
	"ivory/core/service/job"
	"ivory/core/service/secret"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/features/auth"
	"ivory/features/backup"
	"ivory/features/cluster"
	"ivory/features/config"
	"ivory/features/deployment"
	"ivory/features/node"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/features/tag"
	"ivory/features/user"
	"ivory/tools"
	"log/slog"
	"mime/multipart"
)

type Service struct {
	env               *coreConfig.Environment
	authService       *auth.Service
	vaultService      *vault.Service
	clusterService    *cluster.Service
	certService       *cert.Service
	tagService        *tag.Service
	queryService      *query.Service
	deploymentService *deployment.Service
	nodeService       *node.Service
	secretService     *secret.Service
	configService     *config.Service
	permissionService *permission.Service
	backupService     *backup.Service
	userService       *user.Service
	toolRegistry      *utils.Registry[tools.Tool, tools.Adapter]
}

func NewService(
	env *coreConfig.Environment,
	authService *auth.Service,
	vaultService *vault.Service,
	clusterService *cluster.Service,
	certService *cert.Service,
	tagService *tag.Service,
	jobService *job.Service,
	queryService *query.Service,
	deploymentService *deployment.Service,
	nodeService *node.Service,
	secretService *secret.Service,
	configService *config.Service,
	permissionService *permission.Service,
	backupService *backup.Service,
	userService *user.Service,
	toolRegistry *utils.Registry[tools.Tool, tools.Adapter],
) *Service {
	return &Service{
		env:               env,
		authService:       authService,
		vaultService:      vaultService,
		clusterService:    clusterService,
		certService:       certService,
		tagService:        tagService,
		queryService:      queryService,
		deploymentService: deploymentService,
		nodeService:       nodeService,
		secretService:     secretService,
		configService:     configService,
		permissionService: permissionService,
		backupService:     backupService,
		userService:       userService,
		toolRegistry:      toolRegistry,
	}
}

func (s *Service) Free() error {
	errTools := s.deleteAllTools()
	errQuery := s.queryService.DeleteAllLogs()
	return errors.Join(errTools, errQuery)
}

func (s *Service) Erase() error {
	errSecret := s.secretService.Clean()
	errCred := s.vaultService.DeleteAll()
	errCert := s.certService.DeleteAll()
	errCluster := s.clusterService.DeleteAll()
	errTag := s.tagService.DeleteAll()
	errQuery := s.queryService.DeleteAll()
	errDeployment := s.deploymentService.DeleteAll()
	errConfig := s.configService.DeleteAll()
	errPerm := s.permissionService.DeleteAll()
	errUser := s.userService.DeleteAll()
	errUserLink := s.userService.LinkDeleteAll()
	errTools := s.deleteAllTools()

	return errors.Join(errSecret, errCred, errCert, errCluster, errTag, errQuery, errDeployment, errConfig, errPerm, errUser, errUserLink, errTools)
}

func (s *Service) deleteAllTools() error {
	errs := make([]error, 0)
	for _, tool := range s.toolRegistry.All() {
		errs = append(errs, tool.DeleteAll())
	}
	return errors.Join(errs...)
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
	return s.userService.Reencrypt(prevSha, newSha)
}

func (s *Service) GetAppInfo(authorised bool, authEnabled bool, username string, authType string, authError string) *AppInfo {
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

	authorisedResult, user, resultError := s.getAuthInfo(authorised, authEnabled, username, authType, authError)
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
			Authorised: authorisedResult,
			User:       user,
			Error:      resultError,
		},
	}
}

// isSuperuser reports the rights the UI should offer. Without authentication
// there is no identity to check and everything is permitted, which is what the
// permissions themselves already say.
func (s *Service) isSuperuser(username string, authEnabled bool) bool {
	if !authEnabled {
		return true
	}
	isSuperuser, err := s.userService.IsSuperuser(username)
	if err != nil {
		slog.Error("failed to check whether the user is a superuser", "error", err)
		return false
	}
	return isSuperuser
}

func (s *Service) getAuthInfo(authorised bool, authEnabled bool, username string, authType string, authError string) (bool, *UserInfo, string) {
	if authError != "" {
		return authorised, nil, authError
	}
	permissions, errPerm := s.permissionService.GetUserPermissions(permission.Prefix(authType), username, !authEnabled)
	user := &UserInfo{Username: username, Superuser: s.isSuperuser(username, authEnabled), Permissions: permissions}
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
