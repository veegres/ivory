package query

import (
	"errors"
	"ivory/src/clients/database"
	"ivory/src/clients/database/postgres"

	"github.com/google/uuid"
)

func (s *Service) GetList(queryType *database.QueryType) ([]Query, error) {
	if queryType == nil {
		return s.repository.List()
	}

	return s.repository.ListByType(*queryType)
}

func (s *Service) Create(creation QueryCreation, query database.Query) (*uuid.UUID, *Query, error) {
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

func (s *Service) Update(key uuid.UUID, query database.Query) (*uuid.UUID, *Query, error) {
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
	}

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

func (s *Service) Delete(key uuid.UUID) error {
	currentQuery, err := s.repository.Get(key)
	if err != nil {
		return err
	}
	if currentQuery.Creation == System {
		return ErrDeletionOfSystemQueriesRestricted
	}

	var errLog error
	if s.HasLog(key) {
		errLog = s.DeleteLog(key)
	}
	errBucket := s.repository.Delete(key)
	return errors.Join(errLog, errBucket)
}

func (s *Service) DeleteAll() error {
	errLog := s.DeleteAllLogs()
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

func (s *Service) simpleDeadTuples() database.Query {
	n, t, d := "Ratio of dead and live tuples", database.BLOAT, "Shows 100 tables with biggest number of dead tuples and ratio of dead tuples divided by total numbers of tuples"
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultRatioOfDeadTuples}
}

func (s *Service) pureDeadTuples() database.Query {
	n, t, d := "Dead tuples and live tuples with last vacuum and analyze Time", database.BLOAT, "Shows 100 tables with biggest number of dead tuples and their last vacuum and analyze time"
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPureNumberOfDeadTuples}
}

func (s *Service) tableBloat() database.Query {
	n, t, d := "Table bloat", database.BLOAT, "This query will read tables using pgstattuple extension and return top 20 bloated tables. WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)"
	p, v := []string{"schema", "table"}, []database.QueryVariety{database.DatabaseSensitive, database.ReplicaRecommended}
	return database.Query{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultTableBloat}
}

func (s *Service) tableBloatApproximate() database.Query {
	n, t, d := "Table bloat approximate", database.BLOAT, "This query will read tables using pgstattuple extension and return 20 bloated approximate results and doesn't read whole table (but reads toast tables). WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)"
	p, v := []string{"schema", "table"}, []database.QueryVariety{database.DatabaseSensitive, database.ReplicaRecommended}
	return database.Query{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultTableBloatApproximate}
}

func (s *Service) indexBloat() database.Query {
	n, t, d := "Index bloat", database.BLOAT, "This query will read indexes with pgstattuple extension and return top 100 bloated indexes. WARNING: without index mask query will read all available indexes which could cause I/O spikes. Please enter mask for index name (check all indexes if nothing is specified)"
	p, v := []string{"schema", "table", "index"}, []database.QueryVariety{database.DatabaseSensitive, database.ReplicaRecommended}
	return database.Query{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultIndexBloat}
}

func (s *Service) checkTableBloat() database.Query {
	n, t, d := "Check specific table bloat", database.BLOAT, "Shows one table bloat, you need to edit query and provide table name to see information about it"
	p, v := []string{"schema.table"}, []database.QueryVariety{database.DatabaseSensitive}
	return database.Query{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultCheckTableBloat}
}

func (s *Service) checkIndexBloat() database.Query {
	n, t, d := "Check specific index bloat", database.BLOAT, "Shows one index bloat, you need to edit query and provide index name to see information about it"
	p, v := []string{"schema.index"}, []database.QueryVariety{database.DatabaseSensitive}
	return database.Query{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultCheckIndexBloat}
}

func (s *Service) runningQueries() database.Query {
	n, t, d := "Active running queries", database.ACTIVITY, "Shows running queries. It can be useful if you want to check your queries that is long."
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultActiveRunningQueries}
}

func (s *Service) allQueries() database.Query {
	n, t, d := "All running queries", database.ACTIVITY, "Shows all queries. Just can help clarify what is going on postgres side."
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllRunningQueries}
}

func (s *Service) runningVacuums() database.Query {
	n, t, d := "Active vacuums in progress", database.ACTIVITY, "Shows list of active vacuums and their progress"
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultActiveVacuums}
}

func (s *Service) allQueriesByState() database.Query {
	n, t, d := "Number of queries by state and database", database.ACTIVITY, "Shows all queries by state and database"
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllQueriesByState}
}

func (s *Service) allLocks() database.Query {
	n, t, d := "All locks", database.ACTIVITY, "Shows all locks with lock duration, type, it's ids owner, etc"
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllLocks}
}

func (s *Service) allLocksByType() database.Query {
	n, t, d := "Number of locks by lock type", database.ACTIVITY, "Shows all locks by lock type"
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllLocksByLock}
}

func (s *Service) config() database.Query {
	n, t, d := "Config", database.OTHER, "Shows postgres config elements with it's values and information about restart"
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPostgresConfig}
}

func (s *Service) configDescription() database.Query {
	n, t, d := "Config description", database.OTHER, "Shows description of postgres config elements"
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPostgresConfigDescription}
}

func (s *Service) allUsers() database.Query {
	n, t, d := "Users", database.OTHER, "Shows all users"
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPostgresUsers}
}

func (s *Service) pureReplication() database.Query {
	n, t, d := "Pure replication", database.REPLICATION, "Shows pure replication table"
	v := []database.QueryVariety{database.MasterOnly}
	return database.Query{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultPureReplication}
}

func (s *Service) simpleReplication() database.Query {
	n, t, d := "Simple replication", database.REPLICATION, "Shows simple replication table only with lsn info"
	v := []database.QueryVariety{database.MasterOnly}
	return database.Query{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultSimpleReplication}
}

func (s *Service) prettyReplication() database.Query {
	n, t, d := "Pretty replication", database.REPLICATION, "Shows pretty replication table with data in mb"
	v := []database.QueryVariety{database.MasterOnly}
	return database.Query{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultPrettyReplication}
}

func (s *Service) databaseSize() database.Query {
	n, t, d := "Database size", database.STATISTIC, "Shows all database sizes"
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultDatabaseSize}
}

func (s *Service) tableSize() database.Query {
	n, t, d := "Table size", database.STATISTIC, "Shows all table sizes, index size and total (index + table)"
	v := []database.QueryVariety{database.DatabaseSensitive}
	return database.Query{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultTableSize}
}

func (s *Service) indexesCache() database.Query {
	n, t, d := "Indexes in cache", database.STATISTIC, "Shows ratio indexes in cache"
	v := []database.QueryVariety{database.DatabaseSensitive}
	return database.Query{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultIndexInCache}
}

func (s *Service) indexesUnused() database.Query {
	n, t, d := "Unused indexes", database.STATISTIC, "Shows unused indexes and their size"
	v := []database.QueryVariety{database.DatabaseSensitive}
	return database.Query{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultIndexUnused}
}

func (s *Service) indexesInvalid() database.Query {
	n, t, d := "Invalid indexes", database.STATISTIC, "Shows invalid indexes. It can happen when concurrent index creation failed. It means that postgres doesn't use this index. You need to reindex it concurrently."
	return database.Query{Name: n, Type: &t, Description: &d, Query: postgres.DefaultIndexInvalid}
}
