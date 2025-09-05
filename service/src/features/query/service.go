package query

import (
	"errors"
	"ivory/src/features/secret"

	"github.com/google/uuid"
)

type QueryService struct {
	queryRepository *QueryRepository
	secretService   *secret.SecretService
	databaseClient  DatabaseClient

	chartMap map[QueryChartType]QueryRequest
}

func NewQueryService(
	queryRepository *QueryRepository,
	secretService *secret.SecretService,
	databaseClient DatabaseClient,
) *QueryService {
	queryService := &QueryService{
		databaseClient:  databaseClient,
		queryRepository: queryRepository,
		secretService:   secretService,

		chartMap: CreateChartsMap(),
	}
	err := queryService.createDefaultQueries()
	if err != nil {
		panic("Cannot create default queries: " + err.Error())
	}
	return queryService
}

func (s *QueryService) GetLog(queryUuid uuid.UUID) ([]QueryFields, error) {
	return s.queryRepository.GetLog(queryUuid)
}

func (s *QueryService) DeleteLog(queryUuid uuid.UUID) error {
	return s.queryRepository.DeleteLog(queryUuid)
}

func (s *QueryService) RunTemplateQuery(ctx QueryContext, uuid uuid.UUID, options *QueryOptions) (*QueryFields, error) {
	query, errQuery := s.queryRepository.Get(uuid)
	if errQuery != nil {
		return nil, errQuery
	}
	if query.Custom == "" {
		return nil, errors.New("query is empty")
	}
	response, errRun := s.RunQuery(ctx, query.Custom, options)
	if errRun == nil && len(response.Rows) > 0 {
		// NOTE: we don't want fail request if there is some problem with writing to the file
		_ = s.queryRepository.AddLog(uuid, response)
	}
	return response, errRun
}

func (s *QueryService) Cancel(ctx QueryContext, pid int) error {
	return s.databaseClient.Cancel(ctx, pid)
}

func (s *QueryService) Terminate(ctx QueryContext, pid int) error {
	return s.databaseClient.Terminate(ctx, pid)
}

func (s *QueryService) RunQuery(ctx QueryContext, query string, options *QueryOptions) (*QueryFields, error) {
	return s.databaseClient.GetFields(ctx, query, options)
}

func (s *QueryService) GetAllRunningQueriesByApplicationName(ctx QueryContext) (*QueryFields, error) {
	options := &QueryOptions{Params: []any{s.databaseClient.GetApplicationName(ctx.Session)}}
	return s.databaseClient.GetFields(ctx, GetAllRunningQueriesByApplicationName, options)
}

func (s *QueryService) DatabasesQuery(ctx QueryContext, name string) ([]string, error) {
	return s.databaseClient.GetMany(ctx, GetAllDatabases, []any{"%" + name + "%"})
}

func (s *QueryService) SchemasQuery(ctx QueryContext, name string) ([]string, error) {
	db := ctx.Connection.Db
	if db.Name == nil || *db.Name == "" {
		return []string{}, nil
	}
	return s.databaseClient.GetMany(ctx, GetAllSchemas, []any{"%" + name + "%"})
}

func (s *QueryService) TablesQuery(ctx QueryContext, schema string, name string) ([]string, error) {
	db := ctx.Connection.Db
	if db.Name == nil || *db.Name == "" || schema == "" {
		return []string{}, nil
	}
	return s.databaseClient.GetMany(ctx, GetAllTables, []any{schema, "%" + name + "%"})
}

func (s *QueryService) ChartQuery(ctx QueryContext, chartType QueryChartType) (*QueryChart, error) {
	request, ok := s.chartMap[chartType]
	if !ok {
		return nil, errors.New("chart " + string(chartType) + " is not supported")
	}
	response, err := s.databaseClient.GetOne(ctx, request.Query)
	if err != nil {
		return nil, errors.Join(errors.New("cannot get "+request.Name), err)
	}
	return &QueryChart{Name: request.Name, Value: response}, nil
}

func (s *QueryService) GetList(queryType *QueryType) ([]Query, error) {
	if queryType == nil {
		return s.queryRepository.List()
	} else {
		return s.queryRepository.ListByType(*queryType)
	}
}

func (s *QueryService) Create(creation QueryCreation, query QueryRequest) (*uuid.UUID, *Query, error) {
	if query.Name == "" || query.Type == nil || query.Query == "" {
		return nil, nil, errors.New("all fields have to be filled")
	}

	return s.queryRepository.Create(Query{
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

func (s *QueryService) Update(key uuid.UUID, query QueryRequest) (*uuid.UUID, *Query, error) {
	currentQuery, err := s.queryRepository.Get(key)
	if err != nil {
		return nil, nil, err
	}
	if currentQuery.Creation == System {
		if query.Name != currentQuery.Name {
			return nil, nil, errors.New("name change is not allowed for system queries")
		}
		if *query.Type != currentQuery.Type {
			return nil, nil, errors.New("type change is not allowed for system queries")
		}
		if *query.Description != *currentQuery.Description {
			return nil, nil, errors.New("description change is not allowed for system queries")
		}

		return s.queryRepository.Update(key, Query{
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

		return s.queryRepository.Update(key, Query{
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

func (s *QueryService) Delete(key uuid.UUID) error {
	currentQuery, err := s.queryRepository.Get(key)
	if err != nil {
		return err
	}
	if currentQuery.Creation == Manual {
		return s.queryRepository.Delete(key)
	} else {
		return errors.New("deletion of system queries is restricted")
	}
}

func (s *QueryService) DeleteAll() error {
	errDel := s.queryRepository.DeleteAll()
	errDefQueries := s.createDefaultQueries()
	return errors.Join(errDel, errDefQueries)
}

func (s *QueryService) createDefaultQueries() error {
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

func (s *QueryService) simpleDeadTuples() QueryRequest {
	n, t, d := "Ratio of dead and live tuples", BLOAT, "Shows 100 tables with biggest number of dead tuples and ratio of dead tuples divided by total numbers of tuples"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultRatioOfDeadTuples}
}

func (s *QueryService) pureDeadTuples() QueryRequest {
	n, t, d := "Dead tuples and live tuples with last vacuum and analyze Time", BLOAT, "Shows 100 tables with biggest number of dead tuples and their last vacuum and analyze time"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultPureNumberOfDeadTuples}
}

func (s *QueryService) tableBloat() QueryRequest {
	n, t, d := "Table bloat", BLOAT, "This query will read tables using pgstattuple extension and return top 20 bloated tables. WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)"
	p, v := []string{"schema", "table"}, []QueryVariety{DatabaseSensitive, ReplicaRecommended}
	return QueryRequest{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: DefaultTableBloat}
}
func (s *QueryService) tableBloatApproximate() QueryRequest {
	n, t, d := "Table bloat approximate", BLOAT, "This query will read tables using pgstattuple extension and return 20 bloated approximate results and doesn't read whole table (but reads toast tables). WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)"
	p, v := []string{"schema", "table"}, []QueryVariety{DatabaseSensitive, ReplicaRecommended}
	return QueryRequest{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: DefaultTableBloatApproximate}
}
func (s *QueryService) indexBloat() QueryRequest {
	n, t, d := "Index bloat", BLOAT, "This query will read indexes with pgstattuple extension and return top 100 bloated indexes. WARNING: without index mask query will read all available indexes which could cause I/O spikes. Please enter mask for index name (check all indexes if nothing is specified)"
	p, v := []string{"index", "schema", "table"}, []QueryVariety{DatabaseSensitive, ReplicaRecommended}
	return QueryRequest{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: DefaultIndexBloat}
}

func (s *QueryService) checkTableBloat() QueryRequest {
	n, t, d := "Check specific table bloat", BLOAT, "Shows one table bloat, you need to edit query and provide table name to see information about it"
	p, v := []string{"schema.table"}, []QueryVariety{DatabaseSensitive}
	return QueryRequest{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: DefaultCheckTableBloat}
}

func (s *QueryService) checkIndexBloat() QueryRequest {
	n, t, d := "Check specific index bloat", BLOAT, "Shows one index bloat, you need to edit query and provide index name to see information about it"
	p, v := []string{"schema.index"}, []QueryVariety{DatabaseSensitive}
	return QueryRequest{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: DefaultCheckIndexBloat}
}

func (s *QueryService) runningQueries() QueryRequest {
	n, t, d := "Active running queries", ACTIVITY, "Shows running queries. It can be useful if you want to check your queries that is long."
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultActiveRunningQueries}
}

func (s *QueryService) allQueries() QueryRequest {
	n, t, d := "All running queries", ACTIVITY, "Shows all queries. Just can help clarify what is going on postgres side."
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultAllRunningQueries}
}

func (s *QueryService) runningVacuums() QueryRequest {
	n, t, d := "Active vacuums in progress", ACTIVITY, "Shows list of active vacuums and their progress"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultActiveVacuums}
}
func (s *QueryService) allQueriesByState() QueryRequest {
	n, t, d := "Number of queries by state and database", ACTIVITY, "Shows all queries by state and database"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultAllQueriesByState}
}

func (s *QueryService) allLocks() QueryRequest {
	n, t, d := "All locks", ACTIVITY, "Shows all locks with lock duration, type, it's ids owner, etc"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultAllLocks}
}

func (s *QueryService) allLocksByType() QueryRequest {
	n, t, d := "Number of locks by lock type", ACTIVITY, "Shows all locks by lock type"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultAllLocksByLock}
}

func (s *QueryService) config() QueryRequest {
	n, t, d := "Config", OTHER, "Shows postgres config elements with it's values and information about restart"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultPostgresConfig}
}

func (s *QueryService) configDescription() QueryRequest {
	n, t, d := "Config description", OTHER, "Shows description of postgres config elements"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultPostgresConfigDescription}
}

func (s *QueryService) allUsers() QueryRequest {
	n, t, d := "Users", OTHER, "Shows all users"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultPostgresUsers}
}

func (s *QueryService) pureReplication() QueryRequest {
	n, t, d := "Pure replication", REPLICATION, "Shows pure replication table"
	v := []QueryVariety{MasterOnly}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: DefaultPureReplication}
}

func (s *QueryService) simpleReplication() QueryRequest {
	n, t, d := "Simple replication", REPLICATION, "Shows simple replication table only with lsn info"
	v := []QueryVariety{MasterOnly}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: DefaultSimpleReplication}
}

func (s *QueryService) prettyReplication() QueryRequest {
	n, t, d := "Pretty replication", REPLICATION, "Shows pretty replication table with data in mb"
	v := []QueryVariety{MasterOnly}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: DefaultPrettyReplication}
}

func (s *QueryService) databaseSize() QueryRequest {
	n, t, d := "Database size", STATISTIC, "Shows all database sizes"
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultDatabaseSize}
}

func (s *QueryService) tableSize() QueryRequest {
	n, t, d := "Table size", STATISTIC, "Shows all table sizes, index size and total (index + table)"
	v := []QueryVariety{DatabaseSensitive}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: DefaultTableSize}
}

func (s *QueryService) indexesCache() QueryRequest {
	n, t, d := "Indexes in cache", STATISTIC, "Shows ratio indexes in cache"
	v := []QueryVariety{DatabaseSensitive}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: DefaultIndexInCache}
}

func (s *QueryService) indexesUnused() QueryRequest {
	n, t, d := "Unused indexes", STATISTIC, "Shows unused indexes and their size"
	v := []QueryVariety{DatabaseSensitive}
	return QueryRequest{Name: n, Type: &t, Description: &d, Varieties: v, Query: DefaultIndexUnused}
}

func (s *QueryService) indexesInvalid() QueryRequest {
	n, t, d := "Invalid indexes", STATISTIC, "Shows invalid indexes. It can happen when concurrent index creation failed. It means that postgres doesn't use this index. You need to reindex it concurrently."
	return QueryRequest{Name: n, Type: &t, Description: &d, Query: DefaultIndexInvalid}
}
