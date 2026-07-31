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
	}
}
