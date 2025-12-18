package query

import (
	"errors"
	"fmt"
	. "ivory/src/clients/database"
	"ivory/src/clients/database/postgres"
	"ivory/src/features/cert"
	"ivory/src/features/password"
	"ivory/src/features/secret"

	"github.com/google/uuid"
)

var ErrQueryEmpty = errors.New("query is empty")
var ErrPasswordProblems = errors.New("password problems, check if it exists")
var ErrAllFieldsRequired = errors.New("all fields have to be filled")
var ErrNameChangeNotAllowed = errors.New("name change is not allowed for system queries")
var ErrTypeChangeNotAllowed = errors.New("type change is not allowed for system queries")
var ErrDescriptionChangeNotAllowed = errors.New("description change is not allowed for system queries")
var ErrDeletionOfSystemQueriesRestricted = errors.New("deletion of system queries is restricted")

type ExecuteService struct {
	queryRepository *Repository
	databaseClient  Client
	logService      *LogService
	passwordService *password.Service
	certService     *cert.Service

	chartMap map[QueryChartType]QueryRequest
}

func NewExecuteService(
	queryRepository *Repository,
	databaseClient Client,
	logService *LogService,
	passwordService *password.Service,
	certService *cert.Service,
) *ExecuteService {
	return &ExecuteService{
		queryRepository: queryRepository,
		databaseClient:  databaseClient,
		logService:      logService,
		passwordService: passwordService,
		certService:     certService,

		chartMap: postgres.CreateChartsMap(),
	}
}

// QUERY EXECUTION

func (s *ExecuteService) ConsoleQuery(queryCtx QueryContext, query string, options *QueryOptions) (*QueryFields, error) {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	return s.databaseClient.GetFields(ctx, query, options)
}

func (s *ExecuteService) TemplateQuery(ctx QueryContext, uuid uuid.UUID, options *QueryOptions) (*QueryFields, error) {
	query, errQuery := s.queryRepository.Get(uuid)
	if errQuery != nil {
		return nil, errQuery
	}
	if query.Custom == "" {
		return nil, ErrQueryEmpty
	}
	response, errRun := s.ConsoleQuery(ctx, query.Custom, options)
	if errRun == nil && len(response.Rows) > 0 {
		// NOTE: we don't want fail request if there is some problem with writing to the file
		_ = s.logService.Add(uuid, response)
	}
	return response, errRun
}

func (s *ExecuteService) CancelQuery(queryCtx QueryContext, pid int) error {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return err
	}
	return s.databaseClient.Cancel(ctx, pid)
}

func (s *ExecuteService) TerminateQuery(queryCtx QueryContext, pid int) error {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return err
	}
	return s.databaseClient.Terminate(ctx, pid)
}

func (s *ExecuteService) RunningQueriesByApplicationName(queryCtx QueryContext) (*QueryFields, error) {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	options := &QueryOptions{Params: []any{s.databaseClient.GetApplicationName(ctx.Session)}}
	return s.databaseClient.GetFields(ctx, postgres.GetAllRunningQueriesByApplicationName, options)
}

func (s *ExecuteService) DatabasesQuery(queryCtx QueryContext, name string) ([]string, error) {
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	return s.databaseClient.GetMany(ctx, postgres.GetAllDatabases, []any{"%" + name + "%"})
}

func (s *ExecuteService) SchemasQuery(queryCtx QueryContext, name string) ([]string, error) {
	db := queryCtx.Connection.Db
	if db.Name == nil || *db.Name == "" {
		return []string{}, nil
	}
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	return s.databaseClient.GetMany(ctx, postgres.GetAllSchemas, []any{"%" + name + "%"})
}

func (s *ExecuteService) TablesQuery(queryCtx QueryContext, schema string, name string) ([]string, error) {
	db := queryCtx.Connection.Db
	if db.Name == nil || *db.Name == "" || schema == "" {
		return []string{}, nil
	}
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	return s.databaseClient.GetMany(ctx, postgres.GetAllTables, []any{schema, "%" + name + "%"})
}

func (s *ExecuteService) ChartQuery(queryCtx QueryContext, chartType QueryChartType) (*QueryChart, error) {
	request, ok := s.chartMap[chartType]
	if !ok {
		return nil, fmt.Errorf("chart %s is not supported", chartType)
	}
	ctx, err := s.mapContext(queryCtx)
	if err != nil {
		return nil, err
	}
	response, err := s.databaseClient.GetOne(ctx, request.Query)
	if err != nil {
		return nil, fmt.Errorf("cannot get %s: %w", request.Name, err)
	}
	return &QueryChart{Name: request.Name, Value: response}, nil
}

func (s *ExecuteService) mapContext(queryCtx QueryContext) (Context, error) {
	con := Connection{Database: queryCtx.Connection.Db}
	ctx := Context{Connection: &con, Session: queryCtx.Session}
	if queryCtx.Connection.CredentialId != nil {
		cred, errCred := s.passwordService.GetDecrypted(*queryCtx.Connection.CredentialId)
		if errCred != nil {
			return ctx, ErrPasswordProblems
		}
		con.Credentials = &Credentials{Username: cred.Username, Password: cred.Password}
	}
	if queryCtx.Connection.Certs != nil {
		errTls := s.certService.EnrichTLSConfig(&con.TlsConfig, queryCtx.Connection.Certs)
		if errTls != nil {
			return ctx, errTls
		}
	}
	return ctx, nil
}

// QUERY LOGS CRUD

type LogService struct {
	logRepository *LogRepository
}

func NewLogService(
	logRepository *LogRepository,
) *LogService {
	return &LogService{
		logRepository: logRepository,
	}
}

func (s *LogService) Get(queryUuid uuid.UUID) ([]QueryFields, error) {
	return s.logRepository.Get(queryUuid)
}

func (s *LogService) Add(uuid uuid.UUID, element any) error {
	return s.logRepository.Add(uuid, element)
}

func (s *LogService) Delete(queryUuid uuid.UUID) error {
	return s.logRepository.Delete(queryUuid)
}

func (s *LogService) Exist(queryUuid uuid.UUID) bool {
	return s.logRepository.Exist(queryUuid)
}

func (s *LogService) DeleteAll() error {
	return s.logRepository.DeleteAll()
}

// QUERY CRUD

type Service struct {
	repository    *Repository
	secretService *secret.Service
	logService    *LogService
}

func NewService(
	repository *Repository,
	logService *LogService,
	secretService *secret.Service,
) *Service {
	queryService := &Service{
		repository:    repository,
		logService:    logService,
		secretService: secretService,
	}
	err := queryService.createDefaultQueries()
	if err != nil {
		panic("Cannot create default queries: " + err.Error())
	}
	return queryService
}

func (s *Service) GetList(queryType *QueryType) ([]Query, error) {
	if queryType == nil {
		return s.repository.List()
	}

	return s.repository.ListByType(*queryType)
}

func (s *Service) Create(creation QueryCreation, query QueryRequest) (*uuid.UUID, *Query, error) {
	if query.Name == "" || query.Type == nil || query.Query == "" {
		return nil, nil, ErrAllFieldsRequired
	}

	return s.repository.Create(Query{
		Name:        query.Name,
		Type:        *query.Type,
		Creation:    creation,
		Description: query.Description,
		Default:     query.Query,
		Custom:      query.Query,
		Params:      query.Params,
		Varieties:   query.Varieties,
	})
}

func (s *Service) Update(key uuid.UUID, query QueryRequest) (*uuid.UUID, *Query, error) {
	currentQuery, err := s.repository.Get(key)
	if err != nil {
		return nil, nil, err
	}
	if currentQuery.Creation == System {
		if query.Name != currentQuery.Name {
			return nil, nil, ErrNameChangeNotAllowed
		}
		if *query.Type != currentQuery.Type {
			return nil, nil, ErrTypeChangeNotAllowed
		}
		if *query.Description != *currentQuery.Description {
			return nil, nil, ErrDescriptionChangeNotAllowed
		}

		return s.repository.Update(key, Query{
			Id:          key,
			Name:        currentQuery.Name,
			Type:        currentQuery.Type,
			Creation:    currentQuery.Creation,
			Description: currentQuery.Description,
			Default:     currentQuery.Default,
			Custom:      query.Query,
			Varieties:   query.Varieties,
			Params:      query.Params,
			CreatedAt:   currentQuery.CreatedAt,
		})
	} else {
		n := currentQuery.Name
		t := currentQuery.Type
		d := currentQuery.Description

		if query.Name != "" {
			n = query.Name
		}
		if query.Type != nil {
			t = *query.Type
		}
		if query.Description != nil {
			d = query.Description
		}

		return s.repository.Update(key, Query{
			Id:          key,
			Name:        n,
			Type:        t,
			Creation:    currentQuery.Creation,
			Description: d,
			Default:     currentQuery.Default,
			Custom:      query.Query,
			Varieties:   query.Varieties,
			Params:      query.Params,
			CreatedAt:   currentQuery.CreatedAt,
		})
	}
}

func (s *Service) Delete(key uuid.UUID) error {
	currentQuery, err := s.repository.Get(key)
	if err != nil {
		return err
	}
	if currentQuery.Creation == System {
		return ErrDeletionOfSystemQueriesRestricted
	}

	var errLog error
	if s.logService.Exist(key) {
		errLog = s.logService.Delete(key)
	}
	errBucket := s.repository.Delete(key)
	return errors.Join(errLog, errBucket)
}

func (s *Service) DeleteAll() error {
	errLog := s.logService.DeleteAll()
	errDel := s.repository.DeleteAll()
	errDefQueries := s.createDefaultQueries()
	return errors.Join(errLog, errDel, errDefQueries)
}

func (s *Service) createDefaultQueries() error {
	if !s.secretService.IsRefEmpty() {
		return nil
	}

	_, _, err := s.Create(System, s.runningQueries())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.allQueries())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.runningVacuums())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.allQueriesByState())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.allLocks())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.allLocksByType())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.config())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.configDescription())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.allUsers())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.simpleReplication())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.prettyReplication())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.pureReplication())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.databaseSize())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.tableSize())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.indexesCache())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.indexesUnused())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.simpleDeadTuples())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.pureDeadTuples())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.tableBloatApproximate())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.tableBloat())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.indexBloat())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.checkTableBloat())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.checkIndexBloat())
	if err != nil {
		return err
	}
	_, _, err = s.Create(System, s.indexesInvalid())
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) simpleDeadTuples() QueryRequest {
	n, t, d := "Ratio of dead and live tuples", BLOAT, "Shows 100 tables with biggest number of dead tuples and ratio of dead tuples divided by total numbers of tuples"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultRatioOfDeadTuples}
}

func (s *Service) pureDeadTuples() QueryRequest {
	n, t, d := "Dead tuples and live tuples with last vacuum and analyze Time", BLOAT, "Shows 100 tables with biggest number of dead tuples and their last vacuum and analyze time"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPureNumberOfDeadTuples}
}

func (s *Service) tableBloat() QueryRequest {
	n, t, d := "Table bloat", BLOAT, "This query will read tables using pgstattuple extension and return top 20 bloated tables. WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)"
	p, v := []string{"schema", "table"}, []QueryVariety{DatabaseSensitive, ReplicaRecommended}
	return QueryRequest{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultTableBloat}
}
func (s *Service) tableBloatApproximate() QueryRequest {
	n, t, d := "Table bloat approximate", BLOAT, "This query will read tables using pgstattuple extension and return 20 bloated approximate results and doesn't read whole table (but reads toast tables). WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)"
	p, v := []string{"schema", "table"}, []QueryVariety{DatabaseSensitive, ReplicaRecommended}
	return QueryRequest{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultTableBloatApproximate}
}
func (s *Service) indexBloat() QueryRequest {
	n, t, d := "Index bloat", BLOAT, "This query will read indexes with pgstattuple extension and return top 100 bloated indexes. WARNING: without index mask query will read all available indexes which could cause I/O spikes. Please enter mask for index name (check all indexes if nothing is specified)"
	p, v := []string{"schema", "table", "index"}, []QueryVariety{DatabaseSensitive, ReplicaRecommended}
	return QueryRequest{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultIndexBloat}
}

func (s *Service) checkTableBloat() QueryRequest {
	n, t, d := "Check specific table bloat", BLOAT, "Shows one table bloat, you need to edit query and provide table name to see information about it"
	p, v := []string{"schema.table"}, []QueryVariety{DatabaseSensitive}
	return QueryRequest{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultCheckTableBloat}
}

func (s *Service) checkIndexBloat() QueryRequest {
	n, t, d := "Check specific index bloat", BLOAT, "Shows one index bloat, you need to edit query and provide index name to see information about it"
	p, v := []string{"schema.index"}, []QueryVariety{DatabaseSensitive}
	return QueryRequest{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultCheckIndexBloat}
}

func (s *Service) runningQueries() QueryRequest {
	n, t, d := "Active running queries", ACTIVITY, "Shows running queries. It can be useful if you want to check your queries that is long."
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultActiveRunningQueries}
}

func (s *Service) allQueries() QueryRequest {
	n, t, d := "All running queries", ACTIVITY, "Shows all queries. Just can help clarify what is going on postgres side."
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllRunningQueries}
}

func (s *Service) runningVacuums() QueryRequest {
	n, t, d := "Active vacuums in progress", ACTIVITY, "Shows list of active vacuums and their progress"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultActiveVacuums}
}
func (s *Service) allQueriesByState() QueryRequest {
	n, t, d := "Number of queries by state and database", ACTIVITY, "Shows all queries by state and database"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllQueriesByState}
}

func (s *Service) allLocks() QueryRequest {
	n, t, d := "All locks", ACTIVITY, "Shows all locks with lock duration, type, it's ids owner, etc"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllLocks}
}

func (s *Service) allLocksByType() QueryRequest {
	n, t, d := "Number of locks by lock type", ACTIVITY, "Shows all locks by lock type"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllLocksByLock}
}

func (s *Service) config() QueryRequest {
	n, t, d := "Config", OTHER, "Shows postgres config elements with it's values and information about restart"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPostgresConfig}
}

func (s *Service) configDescription() QueryRequest {
	n, t, d := "Config description", OTHER, "Shows description of postgres config elements"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPostgresConfigDescription}
}

func (s *Service) allUsers() QueryRequest {
	n, t, d := "Users", OTHER, "Shows all users"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPostgresUsers}
}

func (s *Service) pureReplication() QueryRequest {
	n, t, d := "Pure replication", REPLICATION, "Shows pure replication table"
	v := []QueryVariety{MasterOnly}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultPureReplication}
}

func (s *Service) simpleReplication() QueryRequest {
	n, t, d := "Simple replication", REPLICATION, "Shows simple replication table only with lsn info"
	v := []QueryVariety{MasterOnly}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultSimpleReplication}
}

func (s *Service) prettyReplication() QueryRequest {
	n, t, d := "Pretty replication", REPLICATION, "Shows pretty replication table with data in mb"
	v := []QueryVariety{MasterOnly}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultPrettyReplication}
}

func (s *Service) databaseSize() QueryRequest {
	n, t, d := "Database size", STATISTIC, "Shows all database sizes"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultDatabaseSize}
}

func (s *Service) tableSize() QueryRequest {
	n, t, d := "Table size", STATISTIC, "Shows all table sizes, index size and total (index + table)"
	v := []QueryVariety{DatabaseSensitive}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultTableSize}
}

func (s *Service) indexesCache() QueryRequest {
	n, t, d := "Indexes in cache", STATISTIC, "Shows ratio indexes in cache"
	v := []QueryVariety{DatabaseSensitive}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultIndexInCache}
}

func (s *Service) indexesUnused() QueryRequest {
	n, t, d := "Unused indexes", STATISTIC, "Shows unused indexes and their size"
	v := []QueryVariety{DatabaseSensitive}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultIndexUnused}
}

func (s *Service) indexesInvalid() QueryRequest {
	n, t, d := "Invalid indexes", STATISTIC, "Shows invalid indexes. It can happen when concurrent index creation failed. It means that postgres doesn't use this index. You need to reindex it concurrently."
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: postgres.DefaultIndexInvalid}
}
