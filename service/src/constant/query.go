package constant

const DefaultRunningQuery = `SELECT
	pid AS process_id,
	pg_blocking_pids(pid) AS blocked_by_process_id,
	now() - pg_stat_activity.query_start AS duration,
	query,
	state,
	usename,
	application_name,
	client_addr
FROM pg_stat_activity
WHERE now() - pg_stat_activity.query_start IS NOT NULL
AND state <> 'idle'
AND query NOT LIKE 'START_REPLICATION SLOT%'
ORDER BY now() - pg_stat_activity.query_start DESC;`

const DefaultConnectionQuery = `SELECT datname, state, count(*)
FROM pg_stat_activity
GROUP BY 1, 2 ORDER BY 1, 2;`
