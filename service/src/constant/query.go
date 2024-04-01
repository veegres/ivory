package constant

import . "ivory/src/model"

const GetAllDatabases = `SELECT datname AS name FROM pg_database WHERE datistemplate = false AND datname LIKE $1 LIMIT 100;`
const GetAllSchemas = `SELECT nspname AS NAME FROM pg_namespace WHERE nspname LIKE $1 LIMIT 100;`
const GetAllTables = `SELECT relname AS name FROM pg_stat_all_tables WHERE schemaname = $1 AND relname LIKE $2 LIMIT 100;`

const DefaultActiveRunningQueries = `SELECT
    pid,
    state,
    wait_event AS event,
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
    wait_event AS event,
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
    p.pid,
    (now() - a.xact_start)::text                                                   AS duration,
    coalesce(wait_event_type || '.' || wait_event, 'f')                            AS wait,
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

const DefaultIndexBloat = `WITH indexes AS (SELECT * from pg_stat_user_indexes)
SELECT 
	table_name,
	pg_size_pretty(table_size) AS table_size,
	index_name,
	pg_size_pretty(index_size) AS index_size,
	index_scans,
	round((free_space*100/index_size)::numeric, 1) AS waste_percent,
	pg_size_pretty(free_space) AS waste
FROM (
    SELECT 
		p.schemaname || '.' || c.relname AS table_name,
		p.indexrelname AS index_name,
		(SELECT (
			CASE WHEN avg_leaf_density = 'NaN' THEN 0
			ELSE greatest(ceil(index_size * (1 - avg_leaf_density / (coalesce((SELECT (regexp_matches(reloptions::text, E'.*fillfactor=(\\d+).*'))[1]),'90')::real)))::bigint, 0) END
		) FROM pgstatindex(p.indexrelid::regclass::text)) AS free_space,
		pg_relation_size(p.indexrelid) AS index_size,
		pg_relation_size(p.relid) AS table_size,
		p.idx_scan AS index_scans
    FROM indexes p
		JOIN pg_class c ON p.indexrelid = c.oid
		JOIN pg_index i ON i.indexrelid = p.indexrelid
    WHERE pg_get_indexdef(p.indexrelid) LIKE '%USING btree%' 
	  AND i.indisvalid 
	  AND (c.relpersistence = 'p' OR NOT pg_is_in_recovery())
	  AND p.indexrelname ~ $1 -- index name
	  AND p.schemaname ~ $2   -- schema name
      AND c.relname ~ $3      -- table name
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

func CreateChartsMap() map[QueryChartType]QueryRequest {
	return map[QueryChartType]QueryRequest{
		Databases:      {Name: "Databases", Query: "SELECT count(*) FROM pg_database;"},
		Connections:    {Name: "Connections", Query: "SELECT count(*) FROM pg_stat_activity;"},
		DatabaseSize:   {Name: "Database Size", Query: "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_database_size(datname) AS size FROM pg_database) AS sizes;"},
		DatabaseUptime: {Name: "Database Uptime", Query: "SELECT date_trunc('seconds', now() - pg_postmaster_start_time())::text;"},
		Schemas:        {Name: "Schemas", Query: "SELECT count(*) FROM pg_namespace;"},
		TablesSize:     {Name: "Tables Size", Query: "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_table_size(relid) AS size FROM pg_stat_all_tables) AS sizes;"},
		IndexesSize:    {Name: "Indexes Size", Query: "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_indexes_size(relid) AS size FROM pg_stat_all_tables) AS sizes;"},
		TotalSize:      {Name: "Total Size", Query: "SELECT pg_size_pretty(sum(size)) FROM (SELECT pg_total_relation_size(relid) AS size FROM pg_stat_all_tables) AS sizes;"},
	}
}
