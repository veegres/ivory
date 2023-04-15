package service

import (
	"errors"
	"github.com/google/uuid"
	"ivory/src/constant"
	. "ivory/src/model"
	"ivory/src/persistence"
)

type QueryService struct {
	queryRepository *persistence.QueryRepository
	secretService   *SecretService
	postgresGateway *PostgresClient
}

func NewQueryService(
	queryRepository *persistence.QueryRepository,
	secretService *SecretService,
	postgresGateway *PostgresClient,
) *QueryService {
	queryService := &QueryService{
		postgresGateway: postgresGateway,
		queryRepository: queryRepository,
		secretService:   secretService,
	}
	err := queryService.createDefaultQueries()
	if err != nil {
		panic("Cannot create default queries: " + err.Error())
	}
	return queryService
}

func (s *QueryService) RunQuery(queryUuid uuid.UUID, credentialId uuid.UUID, db Database) (*QueryRunResponse, error) {
	query, errQuery := s.queryRepository.Get(queryUuid)
	if errQuery != nil {
		return nil, errQuery
	}
	if query.Custom == "" {
		return nil, errors.New("query is empty")
	}
	return s.postgresGateway.GetFields(credentialId, db, query.Custom)
}

func (s *QueryService) DatabasesQuery(credentialId uuid.UUID, db Database, name string) ([]string, error) {
	return s.postgresGateway.GetMany(credentialId, db, constant.GetAllDatabases, "%"+name+"%")
}

func (s *QueryService) SchemasQuery(credentialId uuid.UUID, db Database, name string) ([]string, error) {
	if db.Database == nil || *db.Database == "" {
		return []string{}, nil
	}
	return s.postgresGateway.GetMany(credentialId, db, constant.GetAllSchemas, "%"+name+"%")
}

func (s *QueryService) TablesQuery(credentialId uuid.UUID, db Database, schema string, name string) ([]string, error) {
	if db.Database == nil || *db.Database == "" || schema == "" {
		return []string{}, nil
	}
	return s.postgresGateway.GetMany(credentialId, db, constant.GetAllTables, schema, "%"+name+"%")
}

func (s *QueryService) CommonChartQuery(credentialId uuid.UUID, db Database) (*[4]QueryChart, error) {
	// TODO think about parallel this queries
	dbCount, err := s.postgresGateway.GetOne(credentialId, db, "SELECT count(*) FROM pg_database;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get db count"), err)
	}
	connectionCount, err := s.postgresGateway.GetOne(credentialId, db, "SELECT count(*) FROM pg_stat_activity;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get connection count"), err)
	}
	dbSize, err := s.postgresGateway.GetOne(credentialId, db, "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_database_size(datname) AS size FROM pg_database) AS sizes;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get total db size"), err)
	}
	uptime, err := s.postgresGateway.GetOne(credentialId, db, "SELECT date_trunc('seconds', now() - pg_postmaster_start_time())::text;")
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

func (s *QueryService) DatabaseChartQuery(credentialId uuid.UUID, db Database) (*[4]QueryChart, error) {
	// TODO think about parallel this queries
	schemaCount, err := s.postgresGateway.GetOne(credentialId, db, "SELECT count(*) FROM pg_namespace;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get schema count"), err)
	}
	totalSize, err := s.postgresGateway.GetOne(credentialId, db, "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_total_relation_size(relid) AS size FROM pg_stat_all_tables) AS sizes;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get total db size"), err)
	}
	indexSize, err := s.postgresGateway.GetOne(credentialId, db, "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_indexes_size(relid) AS size FROM pg_stat_all_tables) AS sizes;")
	if err != nil {
		return nil, errors.Join(errors.New("cannot get indexes size"), err)
	}
	tableSize, err := s.postgresGateway.GetOne(credentialId, db, "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_table_size(relid) AS size FROM pg_stat_all_tables) AS sizes;")
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
