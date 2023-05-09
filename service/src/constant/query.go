package constant

const GetAllDatabases = `SELECT datname AS name FROM pg_database WHERE datistemplate = false AND datname LIKE $1 LIMIT 100;`
const GetAllSchemas = `SELECT nspname AS NAME FROM pg_namespace WHERE nspname LIKE $1 LIMIT 100;`
const GetAllTables = `SELECT relname AS name FROM pg_stat_all_tables WHERE schemaname = $1 AND relname LIKE $2 LIMIT 100;`

const DefaultActiveRunningQueries = `SELECT
    pid,
    state,
    (now() - pg_stat_activity.query_start)::text AS duration,
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

const DefaultActiveVacuums = `SELECT 
    p.pid,
    date_trunc('second',now() - a.xact_start)                                      AS duration,
    coalesce(wait_event_type ||'.'|| wait_event, 'f')                              AS wait,
    CASE WHEN a.query ~ 'to prevent wraparound' THEN 'freeze' ELSE 'regular' END   AS mode,
    (SELECT datname FROM pg_database WHERE oid = p.datid)                          AS dat,
    p.relid::regclass                                                              AS tab,
    p.phase,
    round((p.heap_blks_total * current_setting('block_size')::int)/1024.0/1024)    AS tab_mb,
    round(pg_total_relation_size(relid)/1024.0/1024)                               AS ttl_mb,
    round((p.heap_blks_scanned * current_setting('block_size')::int)/1024.0/1024)  AS scan_mb,
    round((p.heap_blks_vacuumed * current_setting('block_size')::int)/1024.0/1024) AS vac_mb,
    (100 * p.heap_blks_scanned / nullif(p.heap_blks_total,0))                      AS scan_pct,
    (100 * p.heap_blks_vacuumed / nullif(p.heap_blks_total,0))                     AS vac_pct,
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
GROUP BY 1, 2 ORDER BY 1, 2;`

const DefaultAllLocks = `SELECT
    LOC.pid,
    transactionid AS tid,
    virtualtransaction AS vtid,
    extract(epoch from (NOW() - state_change)) as locks_duration,
    locktype,
    mode,
    relation::regclass,
    SA.datname AS db,
    usename AS username,
    client_addr AS ip,
    application_name AS application,
    SA.query
FROM pg_locks LOC
    LEFT JOIN pg_catalog.pg_database db ON db.oid = LOC.database
    LEFT JOIN pg_stat_activity AS SA ON LOC.pid = SA.pid;`

const DefaultSimpleNumberOfDeadTuples = `SELECT 
    schemaname || '.' || relname AS table_name,
    n_dead_tup AS dead_tuples, 
    n_live_tup AS live_tuples
FROM pg_stat_user_tables
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
ORDER BY n_dead_tup DESC
LIMIT 100;`

const DefaultTableBloat = `SELECT 
    table_name,
	pg_size_pretty(relation_size + toast_relation_size) AS total_size,
	pg_size_pretty(toast_relation_size) AS toast_size,
	round(((relation_size - (relation_size - free_space)*100/fillfactor)*100/greatest(relation_size, 1))::numeric, 1) table_waste_percent,
	pg_size_pretty((relation_size - (relation_size - free_space)*100/fillfactor)::bigint) table_waste,
	round(((toast_free_space + relation_size - (relation_size - free_space)*100/fillfactor)*100/greatest(relation_size + toast_relation_size, 1))::numeric, 1) total_waste_percent,
	pg_size_pretty((toast_free_space + relation_size - (relation_size - free_space)*100/fillfactor)::bigint) total_waste
FROM (
    SELECT
		(CASE WHEN n.nspname = 'public' THEN format('%I', c.relname) ELSE format('%I.%I', n.nspname, c.relname) END) AS table_name,
		(SELECT free_space FROM pgstattuple(c.oid)) AS free_space,
		pg_relation_size(c.oid) AS relation_size,
		(CASE WHEN reltoastrelid = 0 THEN 0 ELSE (SELECT free_space FROM pgstattuple(c.reltoastrelid)) END) AS toast_free_space,
		coalesce(pg_relation_size(c.reltoastrelid), 0) AS toast_relation_size,
		coalesce((SELECT (regexp_matches(reloptions::text, E'.*fillfactor=(\\d+).*'))[1]),'100')::real AS fillfactor
    FROM pg_class c
    	LEFT JOIN pg_namespace n ON (n.oid = c.relnamespace)
    WHERE nspname NOT IN ('pg_catalog', 'information_schema')
      AND nspname !~ '^pg_toast' AND nspname !~ '^pg_temp' AND relkind IN ('r', 'm') AND (relpersistence = 'p' OR NOT pg_is_in_recovery())
      -- AND relname ~ :'tablename'
) t
ORDER BY (toast_free_space + relation_size - (relation_size - free_space)*100/fillfactor) DESC 
LIMIT 20;`

const DefaultTableBloatApproximate = `SELECT 
    table_name,
	pg_size_pretty(relation_size + toast_relation_size) AS total_size,
	pg_size_pretty(toast_relation_size) AS toast_size,
	round(((relation_size - (relation_size - free_space)*100/fillfactor)*100/greatest(relation_size, 1))::numeric, 1) table_waste_percent,
	pg_size_pretty((relation_size - (relation_size - free_space)*100/fillfactor)::bigint) table_waste,
	round(((toast_free_space + relation_size - (relation_size - free_space)*100/fillfactor)*100/greatest(relation_size + toast_relation_size, 1))::numeric, 1) total_waste_percent,
	pg_size_pretty((toast_free_space + relation_size - (relation_size - free_space)*100/fillfactor)::bigint) total_waste
FROM (
    SELECT
		(CASE WHEN n.nspname = 'public' THEN format('%I', c.relname) ELSE format('%I.%I', n.nspname, c.relname) END) AS table_name,
		(SELECT approx_free_space FROM pgstattuple_approx(c.oid)) AS free_space,
		pg_relation_size(c.oid) AS relation_size,
		(CASE WHEN reltoastrelid = 0 THEN 0 ELSE (SELECT free_space FROM pgstattuple(c.reltoastrelid)) END) AS toast_free_space,
		coalesce(pg_relation_size(c.reltoastrelid), 0) AS toast_relation_size,
		coalesce((SELECT (regexp_matches(reloptions::text, E'.*fillfactor=(\\d+).*'))[1]),'100')::real AS fillfactor
    FROM pg_class c
        LEFT JOIN pg_namespace n ON (n.oid = c.relnamespace)
    WHERE nspname NOT IN ('pg_catalog', 'information_schema') 
      AND nspname !~ '^pg_toast' AND nspname !~ '^pg_temp' AND relkind IN ('r', 'm') AND (relpersistence = 'p' OR NOT pg_is_in_recovery())
      -- AND relname ~ :'tablename'
) t
ORDER BY (toast_free_space + relation_size - (relation_size - free_space)*100/fillfactor) DESC
LIMIT 20;`

const DefaultIndexBloat = `WITH indexes AS (
    SELECT * from pg_stat_user_indexes
)
SELECT 
	table_name,
	pg_size_pretty(table_size) AS table_size,
	index_name,
	pg_size_pretty(index_size) AS index_size,
	idx_scan as index_scans,
	round((free_space*100/index_size)::numeric, 1) AS waste_percent,
	pg_size_pretty(free_space) AS waste
FROM (
    SELECT 
		(CASE WHEN schemaname = 'public' THEN format('%I', p.relname) ELSE format('%I.%I', schemaname, p.relname) END) AS table_name,
		indexrelname AS index_name,
		(SELECT (
			CASE WHEN avg_leaf_density = 'NaN' THEN 0
			ELSE greatest(ceil(index_size * (1 - avg_leaf_density / (coalesce((SELECT (regexp_matches(reloptions::text, E'.*fillfactor=(\\d+).*'))[1]),'90')::real)))::bigint, 0) END
		) FROM pgstatindex(p.indexrelid::regclass::text)) AS free_space,
		pg_relation_size(p.indexrelid) AS index_size,
		pg_relation_size(p.relid) AS table_size,
		idx_scan
    FROM indexes p
    JOIN pg_class c ON p.indexrelid = c.oid
    JOIN pg_index i ON i.indexrelid = p.indexrelid
    WHERE pg_get_indexdef(p.indexrelid) LIKE '%USING btree%' AND i.indisvalid 
      AND (c.relpersistence = 'p' OR NOT pg_is_in_recovery()) 
      -- AND indexrelname ~ :'indexname'
) t
ORDER BY free_space DESC
LIMIT 100;`

const DefaultCheckTableBloat = ` SELECT * FROM pgstattuple('TABLE_NAME')`
const DefaultCheckIndexBloat = ` SELECT * FROM pgstatindex('INDEX_NAME')`

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
    usename           AS user, 
    state,
    sync_state AS mode, 
    backend_xmin,
    (pg_wal_lsn_diff(CASE WHEN pg_is_in_recovery() THEN pg_last_wal_replay_lsn() ELSE pg_current_wal_lsn() END,sent_lsn)/1024.0/1024)::numeric(10,1) AS pending_mb,
    (pg_wal_lsn_diff(sent_lsn,write_lsn)/1024.0/1024)::numeric(10,1)                                                                                 AS write_mb,
    (pg_wal_lsn_diff(write_lsn,flush_lsn)/1024.0/1024)::numeric(10,1)                                                                                AS flush_mb,
    (pg_wal_lsn_diff(flush_lsn,replay_lsn)/1024.0/1024)::numeric(10,1)                                                                               AS replay_mb,
    ((pg_wal_lsn_diff(CASE WHEN pg_is_in_recovery() THEN sent_lsn ELSE pg_current_wal_lsn() END,replay_lsn))::bigint/1024.0/1024)::numeric(10,1)     AS total_mb,
    replay_lag::interval(0) replay_lag
  FROM pg_stat_replication;`

const DefaultPureReplication = `SELECT * FROM pg_stat_replication;`

const DefaultDatabaseSize = `SELECT 
    datname AS db, 
    pg_size_pretty(pg_database_size(datname)) AS size
FROM pg_database
ORDER BY pg_database_size(datname) DESC;`

const DefaultTableSize = `SELECT
   UT.schemaname as schema,
   SUT.relname as table_name,
   pg_size_pretty(pg_total_relation_size(SUT.relid) - pg_indexes_size(SUT.relid)) as table_size,
   pg_size_pretty(pg_indexes_size(SUT.relid)) As index_size,
   pg_size_pretty(pg_total_relation_size(SUT.relid)) As total_size
FROM pg_catalog.pg_statio_user_tables AS SUT JOIN pg_stat_user_tables AS UT ON SUT.relid = UT.relid
ORDER BY pg_total_relation_size(SUT.relid) DESC;`

const DefaultIndexInCache = `SELECT 
    sum(idx_blks_read) as idx_read, 
    sum(idx_blks_hit)  as idx_hit, 
    (sum(idx_blks_hit) - sum(idx_blks_read)) / sum(idx_blks_hit) as ratio
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
  AND NOT EXISTS (SELECT 1 FROM pg_catalog.pg_inherits AS inh WHERE inh.inhrelid = s.indexrelid)
ORDER BY pg_relation_size(s.indexrelid) DESC;`
