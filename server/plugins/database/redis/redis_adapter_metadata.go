package redis

import (
	"ivory/core/config"
	"ivory/plugins/database"
)

const DefaultInfo = `INFO`
const DefaultKeyspace = `INFO keyspace`
const DefaultReplication = `INFO replication`
const DefaultClientList = `CLIENT LIST`
const DefaultKeys = `KEYS *`
const DefaultSlowLog = `SLOWLOG GET 25`
const DefaultLatency = `LATENCY LATEST`
const DefaultMemoryDoctor = `MEMORY DOCTOR`
const DefaultMemory = `INFO memory`
const DefaultPersistence = `INFO persistence`
const DefaultCommandStats = `INFO commandstats`
const DefaultClients = `INFO clients`
const DefaultConfig = `CONFIG GET *`
const DefaultClusterInfo = `CLUSTER INFO`
const DefaultKeyCount = `DBSIZE`

func (a *Adapter) SupportedFeatures() map[config.Feature]bool {
	return map[config.Feature]bool{
		config.ViewQueryDbInfo:        false,
		config.ViewQueryDbChart:       false,
		config.ManageQueryDbTemplate:  true,
		config.ManageQueryDbConsole:   true,
		config.ManageQueryDbCancel:    false,
		config.ManageQueryDbTerminate: true,
	}
}

func (a *Adapter) SystemCharts() map[database.SystemChartType]string {
	return map[database.SystemChartType]string{}
}

func (a *Adapter) SystemRequests() []database.SystemRequest {
	return []database.SystemRequest{
		{
			Name: "Server info", Type: database.STATISTIC,
			Description: "Shows server info: version, memory, clients and stats.",
			Query:       DefaultInfo,
		},
		{
			Name: "Keyspace", Type: database.STATISTIC,
			Description: "Shows key counts per numbered database.",
			Query:       DefaultKeyspace,
		},
		{
			Name: "Replication", Type: database.REPLICATION,
			Description: "Shows this node's replication role and connected replicas.",
			Query:       DefaultReplication,
		},
		{
			Name: "Connected clients", Type: database.ACTIVITY,
			Description: "Lists every connected client and its last command.",
			Query:       DefaultClientList,
		},
		{
			Name: "All keys", Type: database.OTHER,
			Description: "Lists keys matching a pattern (default: all keys in the selected database). KEYS scans the whole keyspace and blocks the server while it runs - prefer SCAN on a busy instance.",
			Query:       DefaultKeys,
		},
		{
			Name: "Slow log", Type: database.ACTIVITY,
			Description: "Shows the commands that took longer than slowlog-log-slower-than, newest first. Redis runs commands one at a time, so anything here also delayed every other client.",
			Query:       DefaultSlowLog,
		},
		{
			Name: "Latency events", Type: database.ACTIVITY,
			Description: "Shows the latest latency spike per event type (fork, expire cycle, AOF write, ...), which is what explains a stall the slow log has no command for. Reports nothing unless latency-monitor-threshold is set.",
			Query:       DefaultLatency,
		},
		{
			Name: "Memory advice", Type: database.STATISTIC,
			Description: "Asks redis itself to report what looks wrong with its memory use, in plain sentences.",
			Query:       DefaultMemoryDoctor,
		},
		{
			Name: "Memory", Type: database.STATISTIC,
			Description: "Shows used memory against maxmemory, the fragmentation ratio and the eviction policy. Reaching maxmemory is what makes a healthy-looking redis start evicting keys or refusing writes.",
			Query:       DefaultMemory,
		},
		{
			Name: "Persistence", Type: database.STATISTIC,
			Description: "Shows RDB and AOF state: when the last save succeeded, whether one is running, and whether the last one failed. A failing background save is silent until a restart loses data.",
			Query:       DefaultPersistence,
		},
		{
			Name: "Command stats", Type: database.STATISTIC,
			Description: "Shows call counts and average time per command since the last reset, which is where an expensive command running far more often than expected shows up.",
			Query:       DefaultCommandStats,
		},
		{
			Name: "Clients", Type: database.ACTIVITY,
			Description: "Shows connected client counts, the longest running client and how many are blocked waiting on a list or stream.",
			Query:       DefaultClients,
		},
		{
			Name: "Config", Type: database.OTHER,
			Description: "Shows every runtime configuration parameter and its current value.",
			Query:       DefaultConfig,
		},
		{
			Name: "Cluster info", Type: database.REPLICATION,
			Description: "Shows redis cluster state, slot coverage and known nodes. Only meaningful on an instance started in cluster mode.",
			Query:       DefaultClusterInfo,
		},
		{
			Name: "Key count", Type: database.STATISTIC,
			Description: "Counts the keys in the selected database without scanning them, unlike KEYS.",
			Query:       DefaultKeyCount,
		},
	}
}
