package query

import (
	"errors"
	"ivory/src/clients/database"
	"ivory/src/clients/database/postgres"
	"ivory/src/features"
	"ivory/src/features/cert"
	"ivory/src/features/secret"
	"ivory/src/features/vault"
)

var ErrQueryEmpty = errors.New("query is empty")
var ErrVaultProblems = errors.New("vault problems, check if it exists")
var ErrAllFieldsRequired = errors.New("all fields have to be filled")
var ErrNameChangeNotAllowed = errors.New("name change is not allowed for system queries")
var ErrTypeChangeNotAllowed = errors.New("type change is not allowed for system queries")
var ErrDescriptionChangeNotAllowed = errors.New("description change is not allowed for system queries")
var ErrDeletionOfSystemQueriesRestricted = errors.New("deletion of system queries is restricted")

type Service struct {
	repository       *Repository
	databaseRegistry *database.Registry
	vaultService     *vault.Service
	certService      *cert.Service
	secretService    *secret.Service

	chartMap map[ChartType]Request
}

func NewService(
	repository *Repository,
	databaseRegistry *database.Registry,
	vaultService *vault.Service,
	certService *cert.Service,
	secretService *secret.Service,
) *Service {
	queryService := &Service{
		repository:       repository,
		databaseRegistry: databaseRegistry,
		vaultService:     vaultService,
		certService:      certService,
		secretService:    secretService,
	}
	queryService.initializeChartsMap()
	err := queryService.initializeDefaultQueries()
	if err != nil {
		panic("Cannot create default queries: " + err.Error())
	}
	return queryService
}

func (s *Service) SupportedFeatures(t database.Type) []features.Feature {
	c, e := s.databaseRegistry.Get(t)
	if e != nil {
		return []features.Feature{}
	}
	return c.SupportedFeatures()
}

func (s *Service) initializeChartsMap() {
	s.chartMap = map[ChartType]Request{
		Databases:      {Name: string(Databases), Query: "SELECT count(*) FROM pg_database;"},
		Connections:    {Name: Connections, Query: "SELECT count(*) FROM pg_stat_activity;"},
		DatabaseSize:   {Name: DatabaseSize, Query: "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_database_size(datname) AS size FROM pg_database) AS sizes;"},
		DatabaseUptime: {Name: DatabaseUptime, Query: "SELECT date_trunc('seconds', now() - pg_postmaster_start_time())::text;"},
		Schemas:        {Name: Schemas, Query: "SELECT count(*) FROM pg_namespace;"},
		TablesSize:     {Name: TablesSize, Query: "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_table_size(relid) AS size FROM pg_stat_all_tables) AS sizes;"},
		IndexesSize:    {Name: IndexesSize, Query: "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_indexes_size(relid) AS size FROM pg_stat_all_tables) AS sizes;"},
		TotalSize:      {Name: TotalSize, Query: "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_total_relation_size(relid) AS size FROM pg_stat_all_tables) AS sizes;"},
	}
}

func (s *Service) initializeDefaultQueries() error {
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

func (s *Service) simpleDeadTuples() Request {
	n, t, d := "Ratio of dead and live tuples", BLOAT, "Shows 100 tables with biggest number of dead tuples and ratio of dead tuples divided by total numbers of tuples"
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultRatioOfDeadTuples}
}

func (s *Service) pureDeadTuples() Request {
	n, t, d := "Dead tuples and live tuples with last vacuum and analyze Time", BLOAT, "Shows 100 tables with biggest number of dead tuples and their last vacuum and analyze time"
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPureNumberOfDeadTuples}
}

func (s *Service) tableBloat() Request {
	n, t, d := "Table bloat", BLOAT, "This query will read tables using pgstattuple extension and return top 20 bloated tables. WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)"
	p, v := []string{"schema", "table"}, []VarietyType{DatabaseSensitive, ReplicaRecommended}
	return Request{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultTableBloat}
}

func (s *Service) tableBloatApproximate() Request {
	n, t, d := "Table bloat approximate", BLOAT, "This query will read tables using pgstattuple extension and return 20 bloated approximate results and doesn't read whole table (but reads toast tables). WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)"
	p, v := []string{"schema", "table"}, []VarietyType{DatabaseSensitive, ReplicaRecommended}
	return Request{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultTableBloatApproximate}
}

func (s *Service) indexBloat() Request {
	n, t, d := "Index bloat", BLOAT, "This query will read indexes with pgstattuple extension and return top 100 bloated indexes. WARNING: without index mask query will read all available indexes which could cause I/O spikes. Please enter mask for index name (check all indexes if nothing is specified)"
	p, v := []string{"schema", "table", "index"}, []VarietyType{DatabaseSensitive, ReplicaRecommended}
	return Request{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultIndexBloat}
}

func (s *Service) checkTableBloat() Request {
	n, t, d := "Check specific table bloat", BLOAT, "Shows one table bloat, you need to edit query and provide table name to see information about it"
	p, v := []string{"schema.table"}, []VarietyType{DatabaseSensitive}
	return Request{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultCheckTableBloat}
}

func (s *Service) checkIndexBloat() Request {
	n, t, d := "Check specific index bloat", BLOAT, "Shows one index bloat, you need to edit query and provide index name to see information about it"
	p, v := []string{"schema.index"}, []VarietyType{DatabaseSensitive}
	return Request{Name: n, Type: &t, Description: &d, Params: p, Varieties: v, Query: postgres.DefaultCheckIndexBloat}
}

func (s *Service) runningQueries() Request {
	n, t, d := "Active running queries", ACTIVITY, "Shows running queries. It can be useful if you want to check your queries that is long."
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultActiveRunningQueries}
}

func (s *Service) allQueries() Request {
	n, t, d := "All running queries", ACTIVITY, "Shows all queries. Just can help clarify what is going on postgres side."
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllRunningQueries}
}

func (s *Service) runningVacuums() Request {
	n, t, d := "Active vacuums in progress", ACTIVITY, "Shows list of active vacuums and their progress"
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultActiveVacuums}
}

func (s *Service) allQueriesByState() Request {
	n, t, d := "Number of queries by state and database", ACTIVITY, "Shows all queries by state and database"
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllQueriesByState}
}

func (s *Service) allLocks() Request {
	n, t, d := "All locks", ACTIVITY, "Shows all locks with lock duration, type, it's ids owner, etc"
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllLocks}
}

func (s *Service) allLocksByType() Request {
	n, t, d := "Number of locks by lock type", ACTIVITY, "Shows all locks by lock type"
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultAllLocksByLock}
}

func (s *Service) config() Request {
	n, t, d := "Config", OTHER, "Shows postgres config elements with it's values and information about restart"
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPostgresConfig}
}

func (s *Service) configDescription() Request {
	n, t, d := "Config description", OTHER, "Shows description of postgres config elements"
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPostgresConfigDescription}
}

func (s *Service) allUsers() Request {
	n, t, d := "Users", OTHER, "Shows all users"
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultPostgresUsers}
}

func (s *Service) pureReplication() Request {
	n, t, d := "Pure replication", REPLICATION, "Shows pure replication table"
	v := []VarietyType{MasterOnly}
	return Request{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultPureReplication}
}

func (s *Service) simpleReplication() Request {
	n, t, d := "Simple replication", REPLICATION, "Shows simple replication table only with lsn info"
	v := []VarietyType{MasterOnly}
	return Request{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultSimpleReplication}
}

func (s *Service) prettyReplication() Request {
	n, t, d := "Pretty replication", REPLICATION, "Shows pretty replication table with data in mb"
	v := []VarietyType{MasterOnly}
	return Request{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultPrettyReplication}
}

func (s *Service) databaseSize() Request {
	n, t, d := "Database size", STATISTIC, "Shows all database sizes"
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultDatabaseSize}
}

func (s *Service) tableSize() Request {
	n, t, d := "Table size", STATISTIC, "Shows all table sizes, index size and total (index + table)"
	v := []VarietyType{DatabaseSensitive}
	return Request{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultTableSize}
}

func (s *Service) indexesCache() Request {
	n, t, d := "Indexes in cache", STATISTIC, "Shows ratio indexes in cache"
	v := []VarietyType{DatabaseSensitive}
	return Request{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultIndexInCache}
}

func (s *Service) indexesUnused() Request {
	n, t, d := "Unused indexes", STATISTIC, "Shows unused indexes and their size"
	v := []VarietyType{DatabaseSensitive}
	return Request{Name: n, Type: &t, Description: &d, Varieties: v, Query: postgres.DefaultIndexUnused}
}

func (s *Service) indexesInvalid() Request {
	n, t, d := "Invalid indexes", STATISTIC, "Shows invalid indexes. It can happen when concurrent index creation failed. It means that postgres doesn't use this index. You need to reindex it concurrently."
	return Request{Name: n, Type: &t, Description: &d, Query: postgres.DefaultIndexInvalid}
}
