package clickhouse

import (
	"ivory/core/config"
	"ivory/plugins/database"
)

const GetActiveQueries = `SELECT query_id, user, query, elapsed, read_rows, memory_usage FROM system.processes ORDER BY elapsed DESC`
const GetTableSizes = `SELECT database, table, formatReadableSize(sum(bytes_on_disk)) AS size, formatReadableSize(sum(data_uncompressed_bytes)) AS uncompressed FROM system.parts WHERE active GROUP BY database, table ORDER BY sum(bytes_on_disk) DESC`
const GetReplicationStatus = `SELECT database, table, is_readonly, absolute_delay, queue_size FROM system.replicas`
const GetSlowQueries = `SELECT query_start_time, query_duration_ms, query FROM system.query_log WHERE type = 'QueryFinish' ORDER BY query_start_time DESC LIMIT 100`
const GetMerges = `SELECT database, table, elapsed, progress, is_mutation FROM system.merges`
const GetClusters = `SELECT cluster, shard_num, replica_num, host_name, host_address, port, is_local, errors_count, slowdowns_count
FROM system.clusters
ORDER BY cluster, shard_num, replica_num`
const GetReplicationQueue = `SELECT database, table, node_name, type, create_time, num_tries, num_postponed, postpone_reason, last_exception
FROM system.replication_queue
ORDER BY create_time`
const GetMutations = `SELECT database, table, mutation_id, command, create_time, parts_to_do, latest_fail_reason
FROM system.mutations
WHERE NOT is_done
ORDER BY create_time`
const GetPartsPerPartition = `SELECT database, table, partition, count() AS parts, sum(rows) AS total_rows, formatReadableSize(sum(bytes_on_disk)) AS size
FROM system.parts
WHERE active
GROUP BY database, table, partition
ORDER BY parts DESC
LIMIT 100`
const GetDetachedParts = `SELECT database, table, partition_id, name, disk, reason FROM system.detached_parts`
const GetFailedQueries = `SELECT event_time, query_duration_ms, user, query, exception
FROM system.query_log
WHERE type IN ('ExceptionBeforeStart', 'ExceptionWhileProcessing')
ORDER BY event_time DESC
LIMIT 100`
const GetQueriesByMemory = `SELECT event_time, formatReadableSize(memory_usage) AS memory, query_duration_ms, read_rows, user, query
FROM system.query_log
WHERE type = 'QueryFinish'
ORDER BY memory_usage DESC
LIMIT 50`
const GetErrors = `SELECT name, code, value, last_error_time, last_error_message FROM system.errors WHERE value > 0 ORDER BY value DESC`
const GetDisks = `SELECT name, path, type, formatReadableSize(free_space) AS free, formatReadableSize(total_space) AS total FROM system.disks`
const GetChangedSettings = `SELECT name, value, "default" AS default_value, description FROM system.settings WHERE changed ORDER BY name`
const GetDistributionQueue = `SELECT database, table, is_blocked, error_count, data_files, formatReadableSize(data_compressed_bytes) AS size, last_exception
FROM system.distribution_queue`
const GetCoordinationRoot = `SELECT name, value, ctime, mtime, numChildren FROM system.zookeeper WHERE path = '/'`

func (a *Adapter) SupportedFeatures() map[config.Feature]bool {
	return map[config.Feature]bool{
		config.ViewQueryDbInfo:        true,
		config.ViewQueryDbChart:       true,
		config.ManageQueryDbTemplate:  true,
		config.ManageQueryDbConsole:   true,
		config.ManageQueryDbCancel:    false,
		config.ManageQueryDbTerminate: false,
	}
}

// SystemCharts omits Schemas/IndexesSize/TotalSize: clickhouse has no schema
// layer (see ListSchemas) and no separate index storage distinct from a
// table's own parts, so those charts have nothing accurate to show.
func (a *Adapter) SystemCharts() map[database.SystemChartType]string {
	return map[database.SystemChartType]string{
		database.Databases:      "SELECT count(*) FROM system.databases",
		database.Connections:    "SELECT value FROM system.metrics WHERE metric = 'TCPConnection'",
		database.DatabaseSize:   "SELECT formatReadableSize(sum(bytes_on_disk)) FROM system.parts WHERE active",
		database.DatabaseUptime: "SELECT formatReadableTimeDelta(uptime())",
		database.TablesSize:     "SELECT formatReadableSize(sum(bytes_on_disk)) FROM system.parts WHERE active",
	}
}

func (a *Adapter) SystemRequests() []database.SystemRequest {
	return []database.SystemRequest{
		{
			Name: "Active queries", Type: database.ACTIVITY,
			Description: "Shows currently running queries.",
			Query:       GetActiveQueries,
		},
		{
			Name: "Table sizes", Type: database.BLOAT,
			Description: "Shows compressed and uncompressed size per table.",
			Query:       GetTableSizes,
		},
		{
			Name: "Replication status", Type: database.REPLICATION,
			Description: "Shows replication queue size and read-only state per replicated table.",
			Query:       GetReplicationStatus,
		},
		{
			Name: "Slow queries", Type: database.STATISTIC,
			Description: "Shows the most recent finished queries ordered by duration.",
			Query:       GetSlowQueries,
		},
		{
			Name: "Merges", Type: database.OTHER,
			Description: "Shows currently running background merges.",
			Query:       GetMerges,
		},
		{
			Name: "Clusters", Type: database.REPLICATION,
			Description: "Shows every cluster from remote_servers with its shards and replicas. errors_count is how often this node failed to reach that replica, so a host that is misconfigured or unreachable shows up here even while its own server is healthy.",
			Query:       GetClusters,
		},
		{
			Name: "Replication queue", Type: database.REPLICATION,
			Description: "Shows the pending fetches, merges and mutations each replicated table still owes. A growing num_tries with a last_exception is a replica that cannot pull data - usually an interserver port that is not reachable.",
			Query:       GetReplicationQueue,
		},
		{
			Name: "Unfinished mutations", Type: database.OTHER,
			Description: "Shows ALTER UPDATE/DELETE mutations that have not finished. A latest_fail_reason means the mutation is stuck and will block every mutation queued behind it.",
			Query:       GetMutations,
		},
		{
			Name: "Parts per partition", Type: database.BLOAT,
			Description: "Shows how many active parts each partition is made of. Too many parts is what slows SELECTs down and eventually raises 'too many parts' on INSERT, and it usually means inserts are too small and too frequent.",
			Query:       GetPartsPerPartition,
		},
		{
			Name: "Detached parts", Type: database.OTHER,
			Description: "Shows parts clickhouse set aside instead of using, with the reason it did. They occupy disk without being queryable and have to be reattached or dropped deliberately.",
			Query:       GetDetachedParts,
		},
		{
			Name: "Failed queries", Type: database.STATISTIC,
			Description: "Shows the most recent queries that ended in an exception, with the exception text. Requires the query_log table to be enabled.",
			Query:       GetFailedQueries,
		},
		{
			Name: "Queries by memory", Type: database.STATISTIC,
			Description: "Shows the finished queries that used the most memory, which is what a memory limit exception is usually traced back to. Requires the query_log table to be enabled.",
			Query:       GetQueriesByMemory,
		},
		{
			Name: "Errors", Type: database.STATISTIC,
			Description: "Shows every error this server has counted since it started, with the last message and when it last happened.",
			Query:       GetErrors,
		},
		{
			Name: "Disks", Type: database.STATISTIC,
			Description: "Shows free and total space per configured disk. A full disk turns replicated tables read-only, which the cluster overview reports as a stopping node.",
			Query:       GetDisks,
		},
		{
			Name: "Changed settings", Type: database.OTHER,
			Description: "Shows only the settings that differ from their defaults, next to the default they replaced.",
			Query:       GetChangedSettings,
		},
		{
			Name: "Distribution queue", Type: database.REPLICATION,
			Description: "Shows the data Distributed tables have not yet delivered to their shards. A blocked queue with an exception means inserts are piling up locally instead of reaching the shard.",
			Query:       GetDistributionQueue,
		},
		{
			Name: "Coordination root", Type: database.OTHER,
			Description: "Lists the top-level znodes of the coordination store (ClickHouse Keeper or ZooKeeper) this server is configured with. It failing at all is itself the answer: this node has no working coordination connection, and replicated tables cannot work without one.",
			Query:       GetCoordinationRoot,
		},
	}
}
