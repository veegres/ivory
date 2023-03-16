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

const DefaultNumberOfDeadTuples = `SELECT 
    schemaname || '.' || relname AS table_name,
    last_vacuum, 
    last_autovacuum, 
    last_analyze, 
    last_autoanalyze, 
    n_dead_tup, 
    n_live_tup
FROM pg_stat_user_tables
WHERE n_dead_tup > 50000
ORDER BY n_dead_tup DESC;`

const DefaultPostgresConfig = `SELECT 
    name, context, vartype, source,
    min_val, max_val, enumvals, boot_val, 
    reset_val, unit, pending_restart
FROM pg_settings;`

const DefaultPostgresConfigDescription = `SELECT 
    name, short_desc, extra_desc
FROM pg_settings;`

const DefaultPostgresUsers = `SELECT * FROM pg_user;`

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
