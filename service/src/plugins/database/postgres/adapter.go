package postgres

import (
	"context"
	"fmt"
	"ivory/src/features"
	"ivory/src/plugins/database"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slog"
)

// NOTE: validate that is matches interface in compile-time
var _ database.Adapter = (*Adapter)(nil)

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) SupportedFeatures() []features.Feature {
	return []features.Feature{
		features.ViewQueryDbInfo,
		features.ViewQueryDbChart,
		features.ManageQueryDbTemplate,
		features.ManageQueryDbConsole,
		features.ManageQueryDbCancel,
		features.ManageQueryDbTerminate,
	}
}

func (a *Adapter) GetMany(ctx database.Context, query string, queryParams []any) ([]string, error) {
	values := make([]string, 0)
	errReq := a.sendRequest(ctx, query, queryParams, func(rows pgx.Rows, _ *pgtype.Map, _ string) error {
		for rows.Next() {
			var value string
			err := rows.Scan(&value)
			if err != nil {
				return err
			}
			values = append(values, value)
		}
		return nil
	})
	if errReq != nil {
		return nil, errReq
	}
	return values, nil
}

func (a *Adapter) GetOne(ctx database.Context, query string) (any, error) {
	var value any
	errReq := a.sendRequest(ctx, query, nil, func(rows pgx.Rows, _ *pgtype.Map, _ string) error {
		if rows.Next() {
			values, err := rows.Values()
			if err != nil {
				return err
			}
			value = values[0]
		}
		return nil
	})
	if errReq != nil {
		return nil, errReq
	}
	return value, nil
}

func (a *Adapter) GetFields(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error) {
	startTime := time.Now().UnixMilli()

	fields := make([]database.QueryField, 0)
	rowList := make([][]any, 0)
	url := "-"

	// NOTE: we need this object ot avoid fatal errors when the option variable is `nil`
	tmpOptions := database.QueryOptions{}
	if options != nil {
		tmpOptions.Limit = options.Limit
		tmpOptions.Trim = options.Trim
		tmpOptions.Params = options.Params
	}

	normQuery, normLimit, errNormQuery := a.normalizeQuery(query, tmpOptions.Trim, tmpOptions.Limit)
	if errNormQuery != nil {
		return nil, errNormQuery
	}

	errReq := a.sendRequest(ctx, normQuery, tmpOptions.Params, func(rows pgx.Rows, typeMap *pgtype.Map, connUrl string) error {
		url = connUrl
		for _, field := range rows.FieldDescriptions() {
			dataType, ok := typeMap.TypeForOID(field.DataTypeOID)
			if !ok {
				fields = append(fields, database.QueryField{Name: field.Name, DataType: "unknown", DataTypeOID: field.DataTypeOID})
			} else {
				fields = append(fields, database.QueryField{Name: field.Name, DataType: dataType.Name, DataTypeOID: field.DataTypeOID})
			}
		}
		for rows.Next() {
			val, err := rows.Values()
			if err != nil {
				return err
			}
			rowList = append(rowList, val)
		}
		return nil
	})
	if errReq != nil {
		return nil, errReq
	}

	if options != nil {
		options.Limit = normLimit
	}
	endTime := time.Now().UnixMilli()
	res := &database.QueryFields{
		Fields:    fields,
		Rows:      rowList,
		Url:       url,
		StartTime: startTime,
		EndTime:   endTime,
		Options:   options,
	}
	return res, nil
}

func (a *Adapter) Cancel(ctx database.Context, pid int) error {
	return a.sendRequest(ctx, "SELECT pg_cancel_backend("+strconv.Itoa(pid)+")", nil, nil)
}

func (a *Adapter) Terminate(ctx database.Context, pid int) error {
	return a.sendRequest(ctx, "SELECT pg_terminate_backend("+strconv.Itoa(pid)+")", nil, nil)
}

func (a *Adapter) ListDatabases(ctx database.Context, name string) ([]string, error) {
	return a.GetMany(ctx, GetAllDatabases, []any{"%" + name + "%"})
}

func (a *Adapter) ListSchemas(ctx database.Context, name string) ([]string, error) {
	return a.GetMany(ctx, GetAllSchemas, []any{"%" + name + "%"})
}

func (a *Adapter) ListTables(ctx database.Context, schema string, name string) ([]string, error) {
	return a.GetMany(ctx, GetAllTables, []any{schema, "%" + name + "%"})
}

func (a *Adapter) ActiveQueries(ctx database.Context, options *database.QueryOptions) (*database.QueryFields, error) {
	return a.GetFields(ctx, GetAllActiveQueriesByApplicationName, options)
}

func (a *Adapter) SystemRequests() []database.SystemRequest {
	return []database.SystemRequest{
		{
			Name: "Active running queries", Type: database.ACTIVITY,
			Description: "Shows running queries. It can be useful if you want to check your queries that is long.",
			Query:       DefaultActiveRunningQueries,
		},
		{
			Name: "All running queries", Type: database.ACTIVITY,
			Description: "Shows all queries. Just can help clarify what is going on postgres side.",
			Query:       DefaultAllRunningQueries,
		},
		{
			Name: "Active vacuums in progress", Type: database.ACTIVITY,
			Description: "Shows list of active vacuums and their progress",
			Query:       DefaultActiveVacuums,
		},
		{
			Name: "Number of queries by state and database", Type: database.ACTIVITY,
			Description: "Shows all queries by state and database",
			Query:       DefaultAllQueriesByState,
		},
		{
			Name: "All locks", Type: database.ACTIVITY,
			Description: "Shows all locks with lock duration, type, it's ids owner, etc",
			Query:       DefaultAllLocks,
		},
		{
			Name: "Number of locks by lock type", Type: database.ACTIVITY,
			Description: "Shows all locks by lock type",
			Query:       DefaultAllLocksByLock,
		},
		{
			Name: "Config", Type: database.OTHER,
			Description: "Shows postgres config elements with it's values and information about restart",
			Query:       DefaultPostgresConfig,
		},
		{
			Name: "Config description", Type: database.OTHER,
			Description: "Shows description of postgres config elements",
			Query:       DefaultPostgresConfigDescription,
		},
		{
			Name: "Users", Type: database.OTHER,
			Description: "Shows all users",
			Query:       DefaultPostgresUsers,
		},
		{
			Name: "Simple replication", Type: database.REPLICATION,
			Description: "Shows simple replication table only with lsn info",
			Varieties:   []database.SystemRequestVariety{database.MasterOnly},
			Query:       DefaultSimpleReplication,
		},
		{
			Name: "Pretty replication", Type: database.REPLICATION,
			Description: "Shows pretty replication table with data in mb",
			Varieties:   []database.SystemRequestVariety{database.MasterOnly},
			Query:       DefaultPrettyReplication,
		},
		{
			Name: "Pure replication", Type: database.REPLICATION,
			Description: "Shows pure replication table",
			Varieties:   []database.SystemRequestVariety{database.MasterOnly},
			Query:       DefaultPureReplication,
		},
		{
			Name: "Database size", Type: database.STATISTIC,
			Description: "Shows all database sizes",
			Query:       DefaultDatabaseSize,
		},
		{
			Name: "Table size", Type: database.STATISTIC,
			Description: "Shows all table sizes, index size and total (index + table)",
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultTableSize,
		},
		{
			Name: "Indexes in cache", Type: database.STATISTIC,
			Description: "Shows ratio indexes in cache",
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultIndexInCache,
		},
		{
			Name: "Unused indexes", Type: database.STATISTIC,
			Description: "Shows unused indexes and their size",
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultIndexUnused,
		},
		{
			Name: "Ratio of dead and live tuples", Type: database.BLOAT,
			Description: "Shows 100 tables with biggest number of dead tuples and ratio of dead tuples divided by total numbers of tuples",
			Query:       DefaultRatioOfDeadTuples,
		},
		{
			Name: "Dead tuples and live tuples with last vacuum and analyze Time", Type: database.BLOAT,
			Description: "Shows 100 tables with biggest number of dead tuples and their last vacuum and analyze time",
			Query:       DefaultPureNumberOfDeadTuples,
		},
		{
			Name: "Table bloat approximate", Type: database.BLOAT,
			Description: "This query will read tables using pgstattuple extension and return 20 bloated approximate results and doesn't read whole table (but reads toast tables). WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)",
			Params:      []string{"schema", "table"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive, database.ReplicaRecommended},
			Query:       DefaultTableBloatApproximate,
		},
		{
			Name: "Table bloat", Type: database.BLOAT,
			Description: "This query will read tables using pgstattuple extension and return top 20 bloated tables. WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)",
			Params:      []string{"schema", "table"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive, database.ReplicaRecommended},
			Query:       DefaultTableBloat,
		},
		{
			Name: "Index bloat", Type: database.BLOAT,
			Description: "This query will read indexes with pgstattuple extension and return top 100 bloated indexes. WARNING: without index mask query will read all available indexes which could cause I/O spikes. Please enter mask for index name (check all indexes if nothing is specified)",
			Params:      []string{"schema", "table", "index"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive, database.ReplicaRecommended},
			Query:       DefaultIndexBloat,
		},
		{
			Name: "Check specific table bloat", Type: database.BLOAT,
			Description: "Shows one table bloat, you need to edit query and provide table name to see information about it",
			Params:      []string{"schema.table"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultCheckTableBloat,
		},
		{
			Name: "Check specific index bloat", Type: database.BLOAT,
			Description: "Shows one index bloat, you need to edit query and provide index name to see information about it",
			Params:      []string{"schema.index"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultCheckIndexBloat,
		},
		{
			Name: "Invalid indexes", Type: database.STATISTIC,
			Description: "Shows invalid indexes. It can happen when concurrent index creation failed. It means that postgres doesn't use this index. You need to reindex it concurrently.",
			Query:       DefaultIndexInvalid,
		},
	}
}

func (a *Adapter) SystemCharts() map[database.SystemChartType]string {
	return map[database.SystemChartType]string{
		database.Databases:      "SELECT count(*) FROM pg_database;",
		database.Connections:    "SELECT count(*) FROM pg_stat_activity;",
		database.DatabaseSize:   "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_database_size(datname) AS size FROM pg_database) AS sizes;",
		database.DatabaseUptime: "SELECT date_trunc('seconds', now() - pg_postmaster_start_time())::text;",
		database.Schemas:        "SELECT count(*) FROM pg_namespace;",
		database.TablesSize:     "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_table_size(relid) AS size FROM pg_stat_all_tables) AS sizes;",
		database.IndexesSize:    "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_indexes_size(relid) AS size FROM pg_stat_all_tables) AS sizes;",
		database.TotalSize:      "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_total_relation_size(relid) AS size FROM pg_stat_all_tables) AS sizes;",
	}
}

type fn func(pgx.Rows, *pgtype.Map, string) error

func (a *Adapter) sendRequest(ctx database.Context, query string, queryParams []any, parse fn) error {
	conn, connUrl, errConn := a.getConnection(ctx)
	if errConn != nil {
		return errConn
	}
	defer a.closeConnection(conn, context.Background())

	txCtx := context.Background()
	tx, errTx := conn.Begin(txCtx)
	if errTx != nil {
		return errTx
	}
	defer a.closeTransaction(tx, txCtx)

	if ctx.Connection.Config.Schema != nil {
		safeSchema := pgx.Identifier{*ctx.Connection.Config.Schema}.Sanitize()
		_, errSchema := tx.Exec(txCtx, fmt.Sprintf("SET LOCAL search_path TO %s", safeSchema))
		if errSchema != nil {
			return fmt.Errorf("failed to set schema: %w", errSchema)
		}
	}

	var rows pgx.Rows
	var err error
	if queryParams == nil {
		rows, err = tx.Query(txCtx, query)
	} else {
		rows, err = tx.Query(txCtx, query, queryParams...)
	}
	if err != nil {
		return err
	}

	defer rows.Close()
	if parse != nil {
		errParse := parse(rows, conn.TypeMap(), connUrl)
		if errParse != nil {
			return errParse
		}
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

func (a *Adapter) normalizeQuery(query string, trim *bool, limit *string) (string, *string, error) {
	if trim == nil || *trim == false {
		if limit != nil {
			return "", limit, database.ErrCannotLimitWithoutTrim
		}
		return query, limit, nil
	}
	trimmedQuery := a.trimQuery(query)
	if limit == nil {
		return trimmedQuery, limit, nil
	}
	parsed := a.parseQuery(trimmedQuery)
	newQuery, newLimit := a.addLimitToQuery(trimmedQuery, parsed, *limit)
	return newQuery, newLimit, nil
}

func (a *Adapter) addLimitToQuery(query string, queryAnalysis database.QueryAnalysis, limit string) (string, *string) {
	if queryAnalysis.LIMIT == 0 && queryAnalysis.SELECT > 0 && queryAnalysis.FROM > 0 && queryAnalysis.EXPLAIN == 0 &&
		queryAnalysis.DELETE == 0 && queryAnalysis.UPDATE == 0 && queryAnalysis.INSERT == 0 {
		replace := " LIMIT " + limit + ";"
		if queryAnalysis.Semicolon {
			// NOTE: removing all spaces and semicolon at the end
			regex := regexp.MustCompile("\\s*;\\s*$")
			return regex.ReplaceAllString(query, replace), &limit
		}
		return query + replace, &limit
	}
	return query, nil
}

func (a *Adapter) trimQuery(query string) string {
	// Remove comments
	commentRegex := regexp.MustCompile("--.*")
	query = commentRegex.ReplaceAllString(query, " ")

	// Normalize whitespace (including tabs and newlines)
	return strings.Join(strings.Fields(query), " ")
}

func (a *Adapter) parseQuery(query string) database.QueryAnalysis {
	lowerQuery := strings.ToLower(query)
	words := strings.Fields(lowerQuery)
	parsed := database.QueryAnalysis{LIMIT: 0, UPDATE: 0, SELECT: 0, INSERT: 0, DELETE: 0, Semicolon: false}
	for i, word := range words {
		// NOTE: we need this check to avoid params rename confusion
		if i-1 > 0 && words[i-1] == "as" {
			continue
		}

		switch word {
		case "limit":
			parsed.LIMIT += 1
		case "update":
			parsed.UPDATE += 1
		case "insert":
			parsed.INSERT += 1
		case "delete":
			parsed.DELETE += 1
		case "select":
			parsed.SELECT += 1
		case "from":
			parsed.FROM += 1
		case "explain":
			parsed.EXPLAIN += 1
		}
	}
	lastWord := words[len(words)-1]
	if lastWord[len(lastWord)-1:] == ";" {
		parsed.Semicolon = true
	}
	return parsed
}

func (a *Adapter) getConnection(ctx database.Context) (*pgx.Conn, string, error) {
	connection := ctx.Connection
	db := connection.Config
	if db.Port == 0 || db.Host == "" || db.Host == "-" {
		return nil, "unknown", database.ErrDatabaseHostOrPortNotSpecified
	}

	dbName := "postgres"
	if db.Name != nil && *db.Name != "" {
		dbName = *db.Name
	}

	credentials := connection.Credentials
	if credentials == nil {
		return nil, "unknown", database.ErrPasswordNotSet
	}

	connProtocol := "postgres://"
	connHost := db.Host + ":" + strconv.Itoa(db.Port) + "/" + dbName
	connUrl := connProtocol + connHost

	tlsConfig := connection.TlsConfig
	if tlsConfig != nil {
		// NOTE: verify-ca was chosen because it potentially can protect from machine-in-the-middle attack if
		// it has the right CA policy. More info can be found here https://www.postgresql.org/docs/16/libpq-ssl.html#LIBPQ-SSL-PROTECTION
		connUrl += "?sslmode=verify-ca"
	}

	conConfig, errConfig := pgx.ParseConfig(connUrl)
	if errConfig != nil {
		return nil, connUrl, errConfig
	}
	conConfig.User = credentials.Username
	conConfig.Password = credentials.Password
	conConfig.RuntimeParams = map[string]string{
		"application_name": ctx.Application,
	}
	if tlsConfig != nil {
		// NOTE: we rewrite only RootCAs and Certificates, because pgx.ParseConfig creates proper
		//  tlsConfig for different `sslmode`. For example `verify-ca` should mark `InsecureSkipVerify=true`
		//  and it always sets `ServerName` it required for `verify-full` mode.
		conConfig.TLSConfig.RootCAs = tlsConfig.RootCAs
		conConfig.TLSConfig.Certificates = tlsConfig.Certificates
	}

	conn, err := pgx.ConnectConfig(context.Background(), conConfig)
	return conn, connUrl, err
}

func (a *Adapter) closeConnection(conn *pgx.Conn, ctx context.Context) {
	err := conn.Close(ctx)
	if err != nil {
		slog.Warn("postgres close connection", err)
	}
}

func (a *Adapter) closeTransaction(tx pgx.Tx, txCtx context.Context) {
	err := tx.Rollback(txCtx)
	if err != nil {
		slog.Warn("postgres rollback", err)
	}
}
