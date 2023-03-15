package service

import (
	"context"
	"errors"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"ivory/src/constant"
	. "ivory/src/model"
	"ivory/src/persistence"
	"strconv"
)

type QueryService struct {
	queryRepository   *persistence.QueryRepository
	clusterRepository *persistence.ClusterRepository
	passwordService   *PasswordService
	secretService     *SecretService
}

func NewQueryService(
	queryRepository *persistence.QueryRepository,
	clusterRepository *persistence.ClusterRepository,
	passwordService *PasswordService,
	secretService *SecretService,
) *QueryService {
	queryService := &QueryService{
		queryRepository:   queryRepository,
		clusterRepository: clusterRepository,
		passwordService:   passwordService,
		secretService:     secretService,
	}
	err := queryService.createDefaultQueries()
	if err != nil {
		panic("Cannot create default queries: " + err.Error())
	}
	return queryService
}

func (s *QueryService) RunQuery(queryUuid uuid.UUID, clusterName string, db Database) (*QueryRunResponse, error) {
	query, errQuery := s.queryRepository.Get(queryUuid)
	if errQuery != nil {
		return nil, errQuery
	}
	if query.Custom == "" {
		return nil, errors.New("query is empty")
	}
	return s.getFields(clusterName, db, query.Custom)
}

func (s *QueryService) CommonChartQuery(clusterName string, db Database) (*[4]QueryChart, error) {
	// TODO think about parallel this queries
	dbCount, err := s.getOne(clusterName, db, "SELECT count(*) FROM pg_database;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get db count"), err)
	}
	connectionCount, err := s.getOne(clusterName, db, "SELECT count(*) FROM pg_stat_activity;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get connection count"), err)
	}
	dbSize, err := s.getOne(clusterName, db, "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_database_size(datname) AS size FROM pg_database) AS sizes;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get total db size"), err)
	}
	uptime, err := s.getOne(clusterName, db, "SELECT date_trunc('seconds', now() - pg_postmaster_start_time())::text;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get uptime"), err)
	}

	return &[4]QueryChart{
		{Name: "Databases", Value: dbCount},
		{Name: "Connections", Value: connectionCount},
		{Name: "Database Size", Value: dbSize},
		{Name: "Database Uptime", Value: uptime},
	}, nil
}

func (s *QueryService) DatabaseChartQuery(clusterName string, db Database) (*[4]QueryChart, error) {
	// TODO think about parallel this queries
	schemaCount, err := s.getOne(clusterName, db, "SELECT count(*) FROM pg_namespace;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get schema count"), err)
	}
	totalSize, err := s.getOne(clusterName, db, "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_total_relation_size(relid) AS size FROM pg_stat_all_tables) AS sizes;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get total db size"), err)
	}
	indexSize, err := s.getOne(clusterName, db, "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_indexes_size(relid) AS size FROM pg_stat_all_tables) AS sizes;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get indexes size"), err)
	}
	tableSize, err := s.getOne(clusterName, db, "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_table_size(relid) AS size FROM pg_stat_all_tables) AS sizes;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get tables size"), err)
	}

	return &[4]QueryChart{
		{Name: "Schemas", Value: schemaCount},
		{Name: "Tables Size", Value: tableSize},
		{Name: "Indexes Size", Value: indexSize},
		{Name: "Total Size", Value: totalSize},
	}, nil
}

func (s *QueryService) CancelQuery(pid int, clusterName string, db Database) error {
	_, _, err := s.sendRequest(clusterName, db, "SELECT pg_cancel_backend("+strconv.Itoa(pid)+")")
	return err
}

func (s *QueryService) TerminateQuery(pid int, clusterName string, db Database) error {
	_, _, err := s.sendRequest(clusterName, db, "SELECT pg_terminate_backend("+strconv.Itoa(pid)+")")
	return err
}

func (s *QueryService) GetList(queryType *QueryType) ([]Query, error) {
	if queryType == nil {
		return s.queryRepository.List()
	} else {
		return s.queryRepository.ListByType(*queryType)
	}
}

func (s *QueryService) Create(creation QueryCreation, query QueryRequest) (*uuid.UUID, *Query, error) {
	if query.Name == nil || query.Type == nil || query.Description == nil {
		return nil, nil, errors.New("all fields have to be filled")
	}

	return s.queryRepository.Create(Query{
		Name:        *query.Name,
		Type:        *query.Type,
		Creation:    creation,
		Description: *query.Description,
		Default:     query.Query,
		Custom:      query.Query,
	})
}

func (s *QueryService) Update(key uuid.UUID, query QueryRequest) (*uuid.UUID, *Query, error) {
	currentQuery, err := s.queryRepository.Get(key)
	if err != nil {
		return nil, nil, err
	}
	if currentQuery.Creation == System {
		if query.Name != nil {
			return nil, nil, errors.New("name change is not allowed for creation type system")
		}
		if query.Type != nil {
			return nil, nil, errors.New("type change is not allowed for creation type system")
		}
		if query.Description != nil {
			return nil, nil, errors.New("description change is not allowed for creation type system")
		}

		return s.queryRepository.Update(key, Query{
			Id:          key,
			Name:        currentQuery.Name,
			Type:        currentQuery.Type,
			Creation:    currentQuery.Creation,
			Description: currentQuery.Description,
			Default:     currentQuery.Default,
			Custom:      query.Query,
		})
	} else {
		n := currentQuery.Name
		t := currentQuery.Type
		d := currentQuery.Description

		if query.Name != nil {
			n = *query.Name
		}
		if query.Type != nil {
			t = *query.Type
		}
		if query.Description != nil {
			d = *query.Description
		}

		return s.queryRepository.Update(key, Query{
			Id:          key,
			Name:        n,
			Type:        t,
			Creation:    currentQuery.Creation,
			Description: d,
			Default:     currentQuery.Default,
			Custom:      query.Query,
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

func (s *QueryService) getOne(clusterName string, db Database, query string) (any, error) {
	rows, _, errReq := s.sendRequest(clusterName, db, query)
	if errReq != nil {
		return nil, errReq
	}
	defer rows.Close()

	var value any
	if rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		value = values[0]
	}

	return value, nil
}

func (s *QueryService) getFields(clusterName string, db Database, query string) (*QueryRunResponse, error) {
	rows, typeMap, errReq := s.sendRequest(clusterName, db, query)
	if errReq != nil {
		return nil, errReq
	}
	defer rows.Close()

	fields := make([]QueryField, 0)
	for _, field := range rows.FieldDescriptions() {
		dataType, ok := typeMap.TypeForOID(field.DataTypeOID)
		if !ok {
			fields = append(fields, QueryField{Name: field.Name, DataType: "unknown", DataTypeOID: field.DataTypeOID})
		} else {
			fields = append(fields, QueryField{Name: field.Name, DataType: dataType.Name, DataTypeOID: field.DataTypeOID})
		}
	}

	rowList := make([][]any, 0)
	for rows.Next() {
		val, err := rows.Values()
		if err != nil {
			return nil, err
		}
		rowList = append(rowList, val)
	}

	res := &QueryRunResponse{
		Fields: fields,
		Rows:   rowList,
	}
	return res, nil
}

func (s *QueryService) sendRequest(clusterName string, db Database, query string) (pgx.Rows, *pgtype.Map, error) {
	conn, errConn := s.getConnection(clusterName, db)
	if errConn != nil {
		return nil, nil, errConn
	}
	defer s.closeConnection(conn, context.Background())

	res, errRes := conn.Query(context.Background(), query)
	if errRes != nil {
		return nil, nil, errRes
	}
	return res, conn.TypeMap(), nil
}

func (s *QueryService) getConnection(clusterName string, db Database) (*pgx.Conn, error) {
	// TODO think about cache connections
	cluster, errCluster := s.clusterRepository.Get(clusterName)
	if errCluster != nil {
		return nil, errCluster
	}
	if cluster.Credentials.PostgresId == nil {
		return nil, errors.New("there is no password for this cluster")
	}

	cred, errCred := s.passwordService.GetDecrypted(*cluster.Credentials.PostgresId)
	if errCred != nil {
		return nil, errCred
	}

	dbName := "postgres"
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, errors.New("database host or port are not specified")
	}
	if db.Database != nil && *db.Database != "" {
		dbName = *db.Database
	}
	connString := "postgres://" + cred.Username + ":" + cred.Password + "@" + db.Host + ":" + strconv.Itoa(db.Port) + "/" + dbName
	return pgx.Connect(context.Background(), connString)
}

func (s *QueryService) closeConnection(conn *pgx.Conn, ctx context.Context) {
	err := conn.Close(ctx)
	if err != nil {
		glog.Warning(err)
	}
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
	_, _, err = s.Create(System, s.allLocks())
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
	_, _, err = s.Create(System, s.replication())
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
	_, _, err = s.Create(System, s.deadTuples())
	if err != nil {
		return err
	}

	return nil
}

func (s *QueryService) deadTuples() QueryRequest {
	n, t, d := "Numbers of dead tuples", BLOAT, "Shows number of dead tuples"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultNumberOfDeadTuples}
}

func (s *QueryService) runningQueries() QueryRequest {
	n, t, d := "Active Running Queries", ACTIVITY, "Shows running queries with duration information and his owner"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultActiveRunningQueries}
}

func (s *QueryService) allQueries() QueryRequest {
	n, t, d := "All Queries By State", ACTIVITY, "Shows all queries by state and database"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultAllQueriesByState}
}

func (s *QueryService) allLocks() QueryRequest {
	n, t, d := "All locks", ACTIVITY, "Shows all locks with lock duration, type, it's ids owner, etc"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultAllLocks}
}

func (s *QueryService) config() QueryRequest {
	n, t, d := "Config", OTHER, "Shows postgres config elements with it's values and information about restart"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultPostgresConfig}
}

func (s *QueryService) configDescription() QueryRequest {
	n, t, d := "Config Description", OTHER, "Shows description of postgres config elements"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultPostgresConfigDescription}
}

func (s *QueryService) allUsers() QueryRequest {
	n, t, d := "Users", OTHER, "Shows all locks with lock duration, type, it's ids owner, etc"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultPostgresUsers}
}

func (s *QueryService) replication() QueryRequest {
	n, t, d := "Replication", REPLICATION, "Shows pure replication table"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultPureReplication}
}

func (s *QueryService) databaseSize() QueryRequest {
	n, t, d := "Database Size", STATISTIC, "Shows all database sizes"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultDatabaseSize}
}

func (s *QueryService) tableSize() QueryRequest {
	n, t, d := "Table Size", STATISTIC, "Shows all table sizes, index size and total (index + table)"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultTableSize}
}

func (s *QueryService) indexesCache() QueryRequest {
	n, t, d := "Indexes in cache", STATISTIC, "Shows ratio indexes in cache"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultIndexInCache}
}

func (s *QueryService) indexesUnused() QueryRequest {
	n, t, d := "Unused indexes", STATISTIC, "Shows unused indexes and their size"
	return QueryRequest{Name: &n, Type: &t, Description: &d, Query: constant.DefaultIndexUnused}
}
