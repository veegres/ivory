package postgres

import (
	"ivory/core/config"
	"ivory/plugins/database"
)

func (a *Adapter) SupportedFeatures() map[config.Feature]bool {
	return map[config.Feature]bool{
		config.ViewQueryDbInfo:        true,
		config.ViewQueryDbChart:       true,
		config.ManageQueryDbTemplate:  true,
		config.ManageQueryDbConsole:   true,
		config.ManageQueryDbCancel:    true,
		config.ManageQueryDbTerminate: true,
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
			Name: "Table pg_compacttable approximate", Type: database.BLOAT,
			Description: "This query will read tables using pgstattuple extension and return 20 bloated approximate results and doesn't read whole table (but reads toast tables). WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)",
			Params:      []string{"schema", "table"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive, database.ReplicaRecommended},
			Query:       DefaultTableBloatApproximate,
		},
		{
			Name: "Table pg_compacttable", Type: database.BLOAT,
			Description: "This query will read tables using pgstattuple extension and return top 20 bloated tables. WARNING: without table mask/name, query will read all available tables which could cause I/O spikes. Please enter mask for table name (check all tables if nothing is specified)",
			Params:      []string{"schema", "table"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive, database.ReplicaRecommended},
			Query:       DefaultTableBloat,
		},
		{
			Name: "Index pg_compacttable", Type: database.BLOAT,
			Description: "This query will read indexes with pgstattuple extension and return top 100 bloated indexes. WARNING: without index mask query will read all available indexes which could cause I/O spikes. Please enter mask for index name (check all indexes if nothing is specified)",
			Params:      []string{"schema", "table", "index"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive, database.ReplicaRecommended},
			Query:       DefaultIndexBloat,
		},
		{
			Name: "Check specific table pg_compacttable", Type: database.BLOAT,
			Description: "Shows one table pg_compacttable, you need to edit query and provide table name to see information about it",
			Params:      []string{"schema.table"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultCheckTableBloat,
		},
		{
			Name: "Check specific index pg_compacttable", Type: database.BLOAT,
			Description: "Shows one index pg_compacttable, you need to edit query and provide index name to see information about it",
			Params:      []string{"schema.index"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultCheckIndexBloat,
		},
		{
			Name: "Invalid indexes", Type: database.STATISTIC,
			Description: "Shows invalid indexes. It can happen when concurrent index creation failed. It means that postgres doesn't use this index. You need to reindex it concurrently.",
			Query:       DefaultIndexInvalid,
		},
		{
			Name: "Replication slots", Type: database.REPLICATION,
			Description: "Shows every replication slot with the amount of WAL it is holding on to. An inactive slot never stops retaining WAL, so it is the usual reason a primary's disk fills up after a replica was removed or has been down for a while. wal_status and safe_wal_size need PostgreSQL 13 or later.",
			Varieties:   []database.SystemRequestVariety{database.MasterOnly},
			Query:       DefaultReplicationSlots,
		},
		{
			Name: "Blocking sessions", Type: database.ACTIVITY,
			Description: "Shows every waiting session next to the session actually holding the lock it wants. blocking_state 'idle in transaction' means the blocker is not doing any work at all and is only holding the lock open.",
			Query:       DefaultBlockingSessions,
		},
		{
			Name: "Long transactions", Type: database.ACTIVITY,
			Description: "Shows transactions open for more than a minute, including idle ones. A long open transaction holds back the xmin horizon, which stops vacuum from cleaning dead tuples anywhere in the database.",
			Query:       DefaultLongTransactions,
		},
		{
			Name: "Cache hit ratio", Type: database.STATISTIC,
			Description: "Shows how much of the table reads were served from shared buffers instead of disk, since the last statistics reset. A ratio that drops well below the usual figure for this database means the working set no longer fits.",
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultCacheHitRatio,
		},
		{
			Name: "Vacuum progress", Type: database.BLOAT,
			Description: "Shows how far every running vacuum has got, table by table. It is what tells a vacuum that is slowly working through a large table apart from one that is stuck.",
			Query:       DefaultVacuumProgress,
		},
		{
			Name: "Settings pending restart", Type: database.OTHER,
			Description: "Shows settings that were changed in a config file but need a restart to take effect, so the running server is not using the value the file states.",
			Query:       DefaultSettingsPendingRestart,
		},
		{
			Name: "Checkpoints and background writer", Type: database.STATISTIC,
			Description: "Shows background writer counters, and on PostgreSQL 16 and earlier the checkpoint counters too - checkpoints triggered by WAL volume rather than by the timeout mean max_wal_size is too small for the write rate, which shows up as periodic write stalls. From PostgreSQL 17 those counters live in pg_stat_checkpointer instead.",
			Query:       DefaultCheckpoints,
		},
		{
			Name: "WAL archiver", Type: database.REPLICATION,
			Description: "Shows whether WAL archiving is keeping up and what it last failed on. A failing archiver retains WAL on disk indefinitely and silently invalidates the backups taken from it.",
			Query:       DefaultArchiver,
		},
		{
			Name: "Sequential scans", Type: database.STATISTIC,
			Description: "Shows the tables read sequentially most often, with how many rows those scans read. A large table with many sequential scans and few index scans is usually a missing index.",
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultSequentialScans,
		},
		{
			Name: "Replica status", Type: database.REPLICATION,
			Description: "Shows what this standby has received against what it has replayed, and whether replay is paused. Received but not replayed is a standby that is connected and still falling behind, which the primary's own view cannot show.",
			Varieties:   []database.SystemRequestVariety{database.ReplicaRecommended},
			Query:       DefaultReplicaStatus,
		},
		{
			Name: "Queries by total time", Type: database.STATISTIC,
			Description: "Shows the queries that consumed the most execution time overall, which is usually a cheap query run very often rather than the slow one being hunted. Requires the pg_stat_statements extension to be installed and loaded.",
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultStatementsByTime,
		},
	}
}

const GetAllDatabases = `SELECT datname AS name FROM pg_database WHERE datistemplate = false AND datname LIKE $1 LIMIT 100;`
const GetAllSchemas = `SELECT nspname AS NAME FROM pg_namespace WHERE nspname LIKE $1 LIMIT 100;`
const GetAllTables = `SELECT relname AS name FROM pg_stat_all_tables WHERE schemaname = $1 AND relname LIKE $2 LIMIT 100;`
const GetAllActiveQueriesByApplicationName = `SELECT
    pid,
    (now() - pg_stat_activity.query_start)::text AS query_duration,
    query
FROM pg_stat_activity
WHERE now() - pg_stat_activity.query_start IS NOT NULL
  AND state <> 'idle'
  AND backend_type = 'client backend'
  AND application_name = $1
ORDER BY now() - pg_stat_activity.query_start DESC;`

const DefaultActiveRunningQueries = `SELECT
    pid,
    state,
    wait_event_type || '.' || wait_event AS wait,
    (now() - pg_stat_activity.backend_start)::text AS transaction_duration,
    (now() - pg_stat_activity.query_start)::text AS query_duration,
    query,
    usename AS username,
    application_name AS application,
    client_addr AS ip,
    pg_blocking_pids(pid) AS blocked_by_process_id
FROM pg_stat_activity
WHERE now() - pg_stat_activity.query_start IS NOT NULL
  AND state <> 'idle'
  AND backend_type = 'client backend'
ORDER BY now() - pg_stat_activity.query_start DESC;`

const DefaultAllRunningQueries = `SELECT
    pid,
    state,
    wait_event_type || '.' || wait_event AS wait,
    (now() - pg_stat_activity.backend_start)::text AS transaction_duration,
    (now() - pg_stat_activity.query_start)::text AS query_duration,
    query,
    usename AS username,
    application_name AS application,
    client_addr AS ip,
    pg_blocking_pids(pid) AS blocked_by_process_id
FROM pg_stat_activity
ORDER BY now() - pg_stat_activity.query_start DESC;`

const DefaultActiveVacuums = `SELECT 
    p.pid                                                                          AS pid,
    (now() - a.xact_start)::text                                                   AS duration,
    wait_event_type || '.' || wait_event                         				   AS wait,
    CASE WHEN a.query ~ 'to prevent wraparound' THEN 'freeze' ELSE 'regular' END   AS mode,
    (SELECT datname FROM pg_database WHERE oid = p.datid)                          AS dat,
    p.relid::regclass                                                              AS tab,
    p.phase,
    round((p.heap_blks_total * current_setting('block_size')::int)/1024.0/1024)    AS tab_mb,
    round(pg_total_relation_size(relid)/1024.0/1024)                               AS ttl_mb,
    round((p.heap_blks_scanned * current_setting('block_size')::int)/1024.0/1024)  AS scan_mb,
    round((p.heap_blks_vacuumed * current_setting('block_size')::int)/1024.0/1024) AS vac_mb,
    (100 * p.heap_blks_scanned / nullif(p.heap_blks_total, 0))                     AS scan_pct,
    (100 * p.heap_blks_vacuumed / nullif(p.heap_blks_total, 0))                    AS vac_pct,
    p.index_vacuum_count                                                           AS ind_vac_cnt,
    round(p.num_dead_tuples * 100.0 / nullif(p.max_dead_tuples, 0),1)              AS dead_pct
FROM pg_stat_progress_vacuum p 
    JOIN pg_stat_activity a using (pid) 
ORDER BY duration DESC;`

const DefaultAllQueriesByState = `SELECT 
    datname AS db,
    state,
    count(*) 
FROM pg_stat_activity
GROUP BY db, state ORDER BY db, state;`

const DefaultAllLocks = `SELECT
    loc.pid,
	loc.mode AS lock,
    (now() - sa.state_change)::text as lock_duration,
    loc.locktype AS lock_type,
    sa.datname AS db,
    sa.usename AS username,
    sa.client_addr AS ip,
    sa.application_name AS application,
    sa.query AS query
FROM pg_locks loc
    LEFT JOIN pg_catalog.pg_database db ON db.oid = loc.database
    LEFT JOIN pg_stat_activity sa ON loc.pid = sa.pid
ORDER BY now() - sa.state_change DESC;`

const DefaultAllLocksByLock = `SELECT 
    mode AS lock,
    count(*) 
FROM pg_locks
GROUP BY lock ORDER BY lock;`

const DefaultRatioOfDeadTuples = `SELECT
    schemaname || '.' || relname AS table_name,
    n_dead_tup AS dead_tuples,
    n_live_tup AS live_tuples,
    (n_dead_tup + n_live_tup) AS total_tuples,
    (n_dead_tup::numeric / (n_dead_tup + n_live_tup))::numeric(30,2) AS ratio
FROM pg_stat_user_tables
WHERE n_live_tup <> 0 OR n_dead_tup <> 0
ORDER BY n_dead_tup DESC
LIMIT 100;`

const DefaultPureNumberOfDeadTuples = `SELECT 
    schemaname || '.' || relname AS table_name,
    n_dead_tup AS dead_tuples, 
    n_live_tup AS live_tuples,
    last_vacuum, 
    last_autovacuum, 
    last_analyze, 
    last_autoanalyze
FROM pg_stat_user_tables
ORDER BY n_dead_tup, n_live_tup DESC
LIMIT 100;`

const DefaultTableBloat = `SELECT 
    table_name,
	pg_size_pretty(relation_size + toast_relation_size) AS total_size,
	pg_size_pretty(toast_relation_size) AS toast_size,
	round(((relation_size - (relation_size - free_space) * 100 / fillfactor) * 100 / greatest(relation_size, 1))::numeric, 1) AS table_waste_percent,
	pg_size_pretty((relation_size - (relation_size - free_space) * 100 / fillfactor)::bigint) AS table_waste,
	round(((toast_free_space + relation_size - (relation_size - free_space) * 100 / fillfactor) * 100 / greatest(relation_size + toast_relation_size, 1))::numeric, 1) AS total_waste_percent,
	pg_size_pretty((toast_free_space + relation_size - (relation_size - free_space) * 100 / fillfactor)::bigint) AS total_waste
FROM (
    SELECT
		n.nspname || '.' || c.relname AS table_name,
		(SELECT free_space FROM pgstattuple(c.oid)) AS free_space,
		pg_relation_size(c.oid) AS relation_size,
		(CASE WHEN c.reltoastrelid = 0 THEN 0 ELSE (SELECT free_space FROM pgstattuple(c.reltoastrelid)) END) AS toast_free_space,
		coalesce(pg_relation_size(c.reltoastrelid), 0) AS toast_relation_size,
		coalesce((SELECT (regexp_matches(c.reloptions::text, E'.*fillfactor=(\\d+).*'))[1]),'100')::real AS fillfactor
    FROM pg_class c
    	LEFT JOIN pg_namespace n ON (n.oid = c.relnamespace)
    WHERE n.nspname NOT IN ('pg_catalog', 'information_schema')
      AND n.nspname !~ '^pg_toast' AND n.nspname !~ '^pg_temp' 
      AND c.relkind IN ('r', 'm') AND (c.relpersistence = 'p' OR NOT pg_is_in_recovery())
      AND n.nspname ~ $1 -- schema name
      AND c.relname ~ $2 -- table name
) t
ORDER BY (toast_free_space + relation_size - (relation_size - free_space) * 100 / fillfactor) DESC 
LIMIT 20;`

const DefaultTableBloatApproximate = `SELECT 
    table_name,
	pg_size_pretty(relation_size + toast_relation_size) AS total_size,
	pg_size_pretty(toast_relation_size) AS toast_size,
	round(((relation_size - (relation_size - free_space) * 100 / fillfactor) * 100 / greatest(relation_size, 1))::numeric, 1) AS table_waste_percent,
	pg_size_pretty((relation_size - (relation_size - free_space) * 100 / fillfactor)::bigint) AS table_waste,
	round(((toast_free_space + relation_size - (relation_size - free_space) * 100 / fillfactor) * 100 / greatest(relation_size + toast_relation_size, 1))::numeric, 1) AS total_waste_percent,
	pg_size_pretty((toast_free_space + relation_size - (relation_size - free_space) * 100 / fillfactor)::bigint) AS total_waste
FROM (
    SELECT
		n.nspname || '.' || c.relname AS table_name,
		(SELECT approx_free_space FROM pgstattuple_approx(c.oid)) AS free_space,
		pg_relation_size(c.oid) AS relation_size,
		(CASE WHEN c.reltoastrelid = 0 THEN 0 ELSE (SELECT free_space FROM pgstattuple(c.reltoastrelid)) END) AS toast_free_space,
		coalesce(pg_relation_size(c.reltoastrelid), 0) AS toast_relation_size,
		coalesce((SELECT (regexp_matches(c.reloptions::text, E'.*fillfactor=(\\d+).*'))[1]),'100')::real AS fillfactor
    FROM pg_class c
        LEFT JOIN pg_namespace n ON (n.oid = c.relnamespace)
    WHERE n.nspname NOT IN ('pg_catalog', 'information_schema') 
      AND n.nspname !~ '^pg_toast' AND n.nspname !~ '^pg_temp' 
      AND c.relkind IN ('r', 'm') 
      AND (c.relpersistence = 'p' OR NOT pg_is_in_recovery())
      AND n.nspname ~ $1 -- schema name
      AND c.relname ~ $2 -- table name
) t
ORDER BY (toast_free_space + relation_size - (relation_size - free_space) * 100 / fillfactor) DESC
LIMIT 20;`

const DefaultIndexBloat = `SELECT 
	table_name,
	index_name,
	pg_size_pretty(index_size) AS index_size,
	index_scans,
	round((free_space*100/index_size)::numeric, 1) AS waste_percent,
	pg_size_pretty(free_space) AS waste
FROM (
    SELECT 
		p.schemaname || '.' || p.relname AS table_name,
		p.indexrelname AS index_name,
		(SELECT (
			CASE WHEN avg_leaf_density = 'NaN' THEN 0
			ELSE greatest(ceil(index_size * (1 - avg_leaf_density / (coalesce((SELECT (regexp_matches(reloptions::text, E'.*fillfactor=(\\d+).*'))[1]),'90')::real)))::bigint, 0) END
		) FROM pgstatindex(p.indexrelid::regclass::text)) AS free_space,
		pg_relation_size(p.indexrelid) AS index_size,
		p.idx_scan AS index_scans
    FROM pg_stat_user_indexes p
		JOIN pg_class c ON p.indexrelid = c.oid
		JOIN pg_index i ON p.indexrelid = i.indexrelid
    WHERE pg_get_indexdef(p.indexrelid) LIKE '%USING btree%' 
	  AND i.indisvalid 
	  AND (c.relpersistence = 'p' OR NOT pg_is_in_recovery())
	  AND p.schemaname ~ $1  -- schema name
	  AND p.relname ~ $2     -- table name
      AND p.indexrelname~ $3 -- index name
) t
ORDER BY free_space DESC
LIMIT 100;`

const DefaultCheckTableBloat = `SELECT * FROM pgstattuple($1)`
const DefaultCheckIndexBloat = `SELECT * FROM pgstatindex($1)`

const DefaultPostgresConfig = `SELECT 
    name, context, vartype, source,
    min_val, max_val, enumvals, boot_val, 
    reset_val, unit, pending_restart
FROM pg_settings;`

const DefaultPostgresConfigDescription = `SELECT 
    name, short_desc
FROM pg_settings;`

const DefaultPostgresUsers = `SELECT * FROM pg_user;`

const DefaultSimpleReplication = `SELECT
    application_name AS name,
    client_addr AS ip,
    sent_lsn,
    write_lsn,
    flush_lsn,
    replay_lsn
FROM pg_stat_replication;`

const DefaultPrettyReplication = `SELECT 
    application_name  AS name,
    client_addr       AS ip,
    usename           AS username, 
    state,
    sync_state        AS mode, 
    backend_xmin,
    (pg_wal_lsn_diff(CASE WHEN pg_is_in_recovery() THEN pg_last_wal_replay_lsn() ELSE pg_current_wal_lsn() END,sent_lsn)/1024.0/1024)::numeric(10,1) AS pending_mb,
    (pg_wal_lsn_diff(sent_lsn,write_lsn)/1024.0/1024)::numeric(10,1)                                                                                 AS write_mb,
    (pg_wal_lsn_diff(write_lsn,flush_lsn)/1024.0/1024)::numeric(10,1)                                                                                AS flush_mb,
    (pg_wal_lsn_diff(flush_lsn,replay_lsn)/1024.0/1024)::numeric(10,1)                                                                               AS replay_mb,
    ((pg_wal_lsn_diff(CASE WHEN pg_is_in_recovery() THEN sent_lsn ELSE pg_current_wal_lsn() END,replay_lsn))::bigint/1024.0/1024)::numeric(10,1)     AS total_mb,
    replay_lag::interval(0) AS replay_lag
FROM pg_stat_replication;`

const DefaultPureReplication = `SELECT * FROM pg_stat_replication;`

const DefaultDatabaseSize = `SELECT 
    datname AS db, 
    pg_size_pretty(pg_database_size(datname)) AS size
FROM pg_database
ORDER BY pg_database_size(datname) DESC;`

const DefaultTableSize = `SELECT
	u.schemaname as schema_name,
	s.relname as table_name,
	pg_size_pretty(pg_total_relation_size(s.relid) - pg_indexes_size(s.relid)) as table_size,
	pg_size_pretty(pg_indexes_size(s.relid)) As index_size,
	pg_size_pretty(pg_total_relation_size(s.relid)) As total_size
FROM pg_catalog.pg_statio_user_tables AS s 
    JOIN pg_stat_user_tables AS u ON s.relid = u.relid
ORDER BY pg_total_relation_size(s.relid) DESC;`

const DefaultIndexInCache = `SELECT 
    sum(idx_blks_read) AS idx_read, 
    sum(idx_blks_hit)  AS idx_hit, 
    (sum(idx_blks_hit) - sum(idx_blks_read)) / sum(idx_blks_hit) AS ratio
FROM pg_statio_user_indexes;`

const DefaultIndexUnused = `SELECT 
    s.schemaname AS schema_name, 
    s.relname AS table_name,
    s.indexrelname AS index_name,
    pg_size_pretty(pg_relation_size(s.indexrelid)) AS index_size
FROM pg_catalog.pg_stat_user_indexes s
   JOIN pg_catalog.pg_index i ON s.indexrelid = i.indexrelid
WHERE s.idx_scan = 0     
  AND 0 <> ALL (i.indkey) 
  AND NOT i.indisunique  
  AND NOT EXISTS (SELECT 1 FROM pg_catalog.pg_constraint c WHERE c.conindid = s.indexrelid)
  AND NOT EXISTS (SELECT 1 FROM pg_catalog.pg_inherits inh WHERE inh.inhrelid = s.indexrelid)
ORDER BY pg_relation_size(s.indexrelid) DESC;`

const DefaultIndexInvalid = `SELECT
    oid,
    indrelid::regclass AS table_name,
    relname AS index,
    indisvalid AS valid
FROM pg_class, pg_index
WHERE pg_index.indisvalid = false
  AND pg_index.indexrelid = pg_class.oid;`

const DefaultReplicationSlots = `SELECT
    slot_name,
    slot_type,
    database,
    active,
    active_pid,
    restart_lsn,
    pg_size_pretty(pg_wal_lsn_diff(pg_current_wal_lsn(), restart_lsn)) AS retained_wal,
    wal_status,
    pg_size_pretty(safe_wal_size) AS safe_wal_size
FROM pg_replication_slots
ORDER BY pg_wal_lsn_diff(pg_current_wal_lsn(), restart_lsn) DESC NULLS LAST;`

const DefaultBlockingSessions = `SELECT
    blocked.pid AS blocked_pid,
    blocked.usename AS blocked_user,
    (now() - blocked.query_start)::text AS blocked_for,
    blocked.wait_event_type || '.' || blocked.wait_event AS waiting_on,
    blocked.query AS blocked_query,
    blocking.pid AS blocking_pid,
    blocking.usename AS blocking_user,
    blocking.state AS blocking_state,
    (now() - blocking.state_change)::text AS blocking_state_age,
    blocking.query AS blocking_query
FROM pg_stat_activity blocked
    JOIN LATERAL unnest(pg_blocking_pids(blocked.pid)) AS blocking_pid ON true
    JOIN pg_stat_activity blocking ON blocking.pid = blocking_pid
ORDER BY blocked.query_start;`

const DefaultLongTransactions = `SELECT
    pid,
    state,
    usename AS username,
    application_name AS application,
    (now() - xact_start)::text AS transaction_age,
    (now() - state_change)::text AS state_age,
    wait_event_type || '.' || wait_event AS wait,
    backend_xmin::text AS holding_snapshot,
    query
FROM pg_stat_activity
WHERE xact_start IS NOT NULL
  AND backend_type = 'client backend'
  AND now() - xact_start > interval '1 minute'
ORDER BY xact_start;`

const DefaultCacheHitRatio = `SELECT
    sum(heap_blks_read) AS disk_reads,
    sum(heap_blks_hit) AS cache_hits,
    round(100.0 * sum(heap_blks_hit) / nullif(sum(heap_blks_hit) + sum(heap_blks_read), 0), 2) AS cache_hit_percent
FROM pg_statio_user_tables;`

const DefaultVacuumProgress = `SELECT
    p.pid,
    p.datname AS database,
    p.relid::regclass AS table_name,
    p.phase,
    pg_size_pretty(p.heap_blks_total * current_setting('block_size')::bigint) AS table_size,
    round(100 * p.heap_blks_scanned / nullif(p.heap_blks_total, 0), 2) AS scanned_percent,
    p.index_vacuum_count,
    a.query
FROM pg_stat_progress_vacuum p
    LEFT JOIN pg_stat_activity a ON a.pid = p.pid;`

const DefaultSettingsPendingRestart = `SELECT
    name,
    setting,
    unit,
    source,
    sourcefile,
    sourceline
FROM pg_settings
WHERE pending_restart
ORDER BY name;`

const DefaultCheckpoints = `SELECT * FROM pg_stat_bgwriter;`

const DefaultArchiver = `SELECT
    archived_count,
    last_archived_wal,
    last_archived_time,
    failed_count,
    last_failed_wal,
    last_failed_time,
    stats_reset
FROM pg_stat_archiver;`

const DefaultSequentialScans = `SELECT
    schemaname AS schema_name,
    relname AS table_name,
    seq_scan,
    seq_tup_read,
    idx_scan,
    n_live_tup AS live_tuples,
    pg_size_pretty(pg_relation_size(relid)) AS table_size
FROM pg_stat_user_tables
WHERE seq_scan > 0
ORDER BY seq_tup_read DESC
LIMIT 50;`

const DefaultReplicaStatus = `SELECT
    pg_is_in_recovery() AS in_recovery,
    pg_last_wal_receive_lsn()::text AS received_lsn,
    pg_last_wal_replay_lsn()::text AS replayed_lsn,
    pg_size_pretty(pg_wal_lsn_diff(pg_last_wal_receive_lsn(), pg_last_wal_replay_lsn())) AS replay_backlog,
    pg_last_xact_replay_timestamp() AS last_replayed_at,
    (now() - pg_last_xact_replay_timestamp())::text AS replay_delay,
    pg_is_wal_replay_paused() AS replay_paused;`

const DefaultStatementsByTime = `SELECT
    calls,
    round(total_exec_time::numeric, 2) AS total_ms,
    round(mean_exec_time::numeric, 2) AS mean_ms,
    rows,
    query
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 50;`
