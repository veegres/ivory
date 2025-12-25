package management

import (
	"encoding/json"
	"errors"
	"fmt"
	"ivory/src/clients/database"
	"ivory/src/clients/sidecar"
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
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

var ErrInvalidQueryType = errors.New("invalid query type")
var ErrInvalidQueryVariety = errors.New("invalid query variety")
var ErrInvalidPermissionStatus = errors.New("invalid permission status")
var ErrInvalidPermission = errors.New("invalid permission")

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

func (s *Service) Export() (*Backup, error) {
	clusters, errCluster := s.clusterService.List()
	if errCluster != nil {
		return nil, errCluster
	}
	queries, errQuery := s.queryService.GetList(nil)
	if errQuery != nil {
		return nil, errQuery
	}
	permissions, errPermission := s.permissionService.GetAllUserPermissions()
	if errPermission != nil {
		return nil, errPermission
	}

	// Map clusters to backup format
	backupClusters := make([]backupCluster, 0)
	for _, c := range clusters {
		backupCluster := clusterToBackup(c)
		if backupCluster != nil {
			backupClusters = append(backupClusters, *backupCluster)
		}
	}

	// Map queries to backup format
	backupQueries := make([]backupQuery, 0)
	for _, q := range queries {
		backupQuery, err := queryToBackup(q)
		if err != nil {
			return nil, err
		}
		if backupQuery != nil {
			backupQueries = append(backupQueries, *backupQuery)
		}
	}

	// Map permissions to backup format
	backupPermissions := make([]backupPermissions, 0)
	for _, p := range permissions {
		backupPerm, err := userPermissionsToBackup(p)
		if err != nil {
			return nil, err
		}
		if backupPerm != nil {
			backupPermissions = append(backupPermissions, *backupPerm)
		}
	}

	return &Backup{
		Clusters:    backupClusters,
		Queries:     backupQueries,
		Permissions: backupPermissions,
	}, nil
}

// Mapper: Cluster to backupCluster
func clusterToBackup(c cluster.Cluster) *backupCluster {
	sidecars := make([]backupSidecar, len(c.Sidecars))
	for i, sc := range c.Sidecars {
		sidecars[i] = backupSidecar{
			Host: sc.Host,
			Port: sc.Port,
		}
	}
	return &backupCluster{
		Name:     c.Name,
		Tags:     c.Tags,
		Sidecars: sidecars,
	}
}

// Mapper: Query to backupQuery
func queryToBackup(q query.Query) (*backupQuery, error) {
	if q.Creation == query.System {
		return nil, nil
	}
	varieties := make([]backupQueryVariety, len(q.Varieties))
	for i, v := range q.Varieties {
		variety, err := queryVarietyToBackup(v)
		if err != nil {
			return nil, err
		}
		varieties[i] = variety
	}

	queryType, err := queryTypeToBackup(q.Type)
	if err != nil {
		return nil, err
	}

	return &backupQuery{
		Name:        q.Name,
		Type:        queryType,
		Varieties:   varieties,
		Params:      q.Params,
		Description: q.Description,
		Default:     q.Default,
		Custom:      q.Custom,
	}, nil
}

// Mapper: QueryType to backupQueryType
func queryTypeToBackup(qt database.QueryType) (backupQueryType, error) {
	switch qt {
	case database.BLOAT:
		return BLOAT, nil
	case database.ACTIVITY:
		return ACTIVITY, nil
	case database.REPLICATION:
		return REPLICATION, nil
	case database.STATISTIC:
		return STATISTIC, nil
	case database.OTHER:
		return OTHER, nil
	default:
		return 0, ErrInvalidQueryType
	}
}

// Mapper: QueryVariety to backupQueryVariety
func queryVarietyToBackup(qv database.QueryVariety) (backupQueryVariety, error) {
	switch qv {
	case database.DatabaseSensitive:
		return DatabaseSensitive, nil
	case database.MasterOnly:
		return MasterOnly, nil
	case database.ReplicaRecommended:
		return ReplicaRecommended, nil
	default:
		return 0, ErrInvalidQueryVariety
	}
}

// Mapper: UserPermissions to backupPermissions
func userPermissionsToBackup(up permission.UserPermissions) (*backupPermissions, error) {
	perms := make(map[string]backupPermissionType)
	for k, v := range up.Permissions {
		status, err := permissionStatusToBackup(v)
		if err != nil {
			return nil, err
		}
		perms[string(k)] = status
	}
	return &backupPermissions{
		Username:    up.Username,
		Permissions: perms,
	}, nil
}

// Mapper: PermissionStatus to backupPermissionType
func permissionStatusToBackup(ps permission.PermissionStatus) (backupPermissionType, error) {
	switch ps {
	case permission.NOT_PERMITTED:
		return NOT_PERMITTED, nil
	case permission.PENDING:
		return PENDING, nil
	case permission.GRANTED:
		return GRANTED, nil
	default:
		return 0, ErrInvalidPermissionStatus
	}
}

func (s *Service) Import(file *multipart.FileHeader) error {
	f, errOpen := file.Open()
	if errOpen != nil {
		return errOpen
	}
	defer f.Close()

	var bkp Backup
	decoder := json.NewDecoder(f)
	if errDecode := decoder.Decode(&bkp); errDecode != nil {
		return errDecode
	}

	var err error
	// Save clusters
	for i, bc := range bkp.Clusters {
		clusterModel := backupToCluster(bc)
		_, errMut := s.clusterService.Update(clusterModel)
		if errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "cluster", i, errMut))
		}
	}
	// Save queries
	for i, bq := range bkp.Queries {
		queryModel, errMap := backupToQuery(bq)
		if errMap != nil {
			continue
		}
		_, _, errMut := s.queryService.Create(query.Manual, queryModel)
		if errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "query", i, errMut))
		}
	}
	// Save permissions
	for i, bp := range bkp.Permissions {
		permModel := backupToUserPermissions(bp)
		errMut := s.permissionService.UpdateUserPermissions(permModel.Username, permModel.Permissions)
		if errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "permission", i, errMut))
		}
	}

	return err
}

// Mapper: backupCluster to Cluster
func backupToCluster(bc backupCluster) cluster.Cluster {
	sidecars := make([]sidecar.Sidecar, len(bc.Sidecars))
	for i, sc := range bc.Sidecars {
		sidecars[i] = sidecar.Sidecar{
			Host: sc.Host,
			Port: sc.Port,
		}
	}
	return cluster.Cluster{
		Name: bc.Name,
		ClusterOptions: cluster.ClusterOptions{
			Tags: bc.Tags,
		},
		Sidecars: sidecars,
	}
}

// Mapper: backupQuery to Query
func backupToQuery(bq backupQuery) (database.Query, error) {
	varieties := make([]database.QueryVariety, 0, len(bq.Varieties))
	for _, v := range bq.Varieties {
		variety, err := syncQueryVariety(v)
		if err == nil {
			varieties = append(varieties, variety)
		}
	}

	queryType, err := syncQueryType(bq.Type)
	if err != nil {
		return database.Query{}, err
	}

	return database.Query{
		Name:        bq.Name,
		Type:        &queryType,
		Description: bq.Description,
		Query:       bq.Default,
		Varieties:   varieties,
		Params:      bq.Params,
	}, nil
}

// Sync: backupQueryType to QueryType
func syncQueryType(bqt backupQueryType) (database.QueryType, error) {
	switch bqt {
	case BLOAT:
		return database.BLOAT, nil
	case ACTIVITY:
		return database.ACTIVITY, nil
	case REPLICATION:
		return database.REPLICATION, nil
	case STATISTIC:
		return database.STATISTIC, nil
	case OTHER:
		return database.OTHER, nil
	default:
		return 0, ErrInvalidQueryType
	}
}

// Sync: backupQueryVariety to QueryVariety
func syncQueryVariety(bqv backupQueryVariety) (database.QueryVariety, error) {
	switch bqv {
	case DatabaseSensitive:
		return database.DatabaseSensitive, nil
	case MasterOnly:
		return database.MasterOnly, nil
	case ReplicaRecommended:
		return database.ReplicaRecommended, nil
	default:
		return 0, ErrInvalidQueryVariety
	}
}

// Mapper: backupPermissions to UserPermissions
func backupToUserPermissions(bp backupPermissions) permission.UserPermissions {
	perms := make(permission.PermissionMap)
	for k, v := range bp.Permissions {
		perm, err := syncPermission(k)
		if err != nil {
			continue
		}
		status, err := syncPermissionStatus(v)
		if err != nil {
			continue
		}
		perms[perm] = status
	}
	return permission.UserPermissions{
		Username:    bp.Username,
		Permissions: perms,
	}
}

// Sync: string to Permission
func syncPermission(p string) (permission.Permission, error) {
	for _, validPerm := range permission.Permissions {
		if string(validPerm) == p {
			return validPerm, nil
		}
	}
	return "", ErrInvalidPermission
}

// Sync: backupPermissionType to PermissionStatus
func syncPermissionStatus(bpt backupPermissionType) (permission.PermissionStatus, error) {
	switch bpt {
	case NOT_PERMITTED:
		return permission.NOT_PERMITTED, nil
	case PENDING:
		return permission.PENDING, nil
	case GRANTED:
		return permission.GRANTED, nil
	default:
		return 0, ErrInvalidPermissionStatus
	}
}
