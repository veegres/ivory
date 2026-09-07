package mongo

import (
	"ivory/core/config"
	"ivory/plugins/database"
)

const DefaultServerStatus = `db.runCommand({"serverStatus": 1})`
const DefaultDbStats = `db.runCommand({"dbStats": 1})`
const DefaultListCollections = `db.runCommand({"listCollections": 1})`
const DefaultCurrentOp = `db.runCommand({"currentOp": 1})`
const DefaultReplSetStatus = `db.runCommand({"replSetGetStatus": 1})`
const DefaultReplSetConfig = `db.runCommand({"replSetGetConfig": 1})`
const DefaultLongOperations = `db.runCommand({"currentOp": 1, "secs_running": {"$gt": 5}})`
const DefaultHostInfo = `db.runCommand({"hostInfo": 1})`
const DefaultConnPoolStats = `db.runCommand({"connPoolStats": 1})`
const DefaultTopCollections = `db.runCommand({"top": 1})`
const DefaultCollectionStats = `db.runCommand({"collStats": "collection"})`
const DefaultIndexStats = `collection.aggregate([{"$indexStats": {}}])`
const DefaultProfilerSlowQueries = `system.profile.find({"millis": {"$gt": 100}})`

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
			Name: "Server status", Type: database.STATISTIC,
			Description: "Shows server info: version, uptime, memory and connections.",
			Query:       DefaultServerStatus,
		},
		{
			Name: "Database stats", Type: database.STATISTIC,
			Description: "Shows the selected database's size, collection and index counts.",
			Query:       DefaultDbStats,
		},
		{
			Name: "Collections", Type: database.OTHER,
			Description: "Lists every collection in the selected database.",
			Query:       DefaultListCollections,
		},
		{
			Name: "Current operations", Type: database.ACTIVITY,
			Description: "Lists every in-progress operation, the same data behind the query table's terminate action.",
			Query:       DefaultCurrentOp,
		},
		{
			Name: "Replica set status", Type: database.REPLICATION,
			Description: "Shows every member's state, health and last applied optime, the same view the cluster overview is built from. Run it against the admin database.",
			Query:       DefaultReplSetStatus,
		},
		{
			Name: "Replica set config", Type: database.REPLICATION,
			Description: "Shows the member list with each member's priority, votes, hidden and arbiter flags - the settings that decide who can be elected. Run it against the admin database.",
			Query:       DefaultReplSetConfig,
		},
		{
			Name: "Long running operations", Type: database.ACTIVITY,
			Description: "Lists operations running longer than five seconds, with the lock and namespace each is waiting on. Run it against the admin database.",
			Query:       DefaultLongOperations,
		},
		{
			Name: "Host info", Type: database.STATISTIC,
			Description: "Shows the machine mongod runs on: cpu count, memory, and whether transparent huge pages or NUMA are set in a way mongo warns about.",
			Query:       DefaultHostInfo,
		},
		{
			Name: "Connection pool stats", Type: database.ACTIVITY,
			Description: "Shows the outgoing connections this node holds to its replica set peers, which is where a member that is reachable but slow to answer shows up.",
			Query:       DefaultConnPoolStats,
		},
		{
			Name: "Collection usage", Type: database.STATISTIC,
			Description: "Shows time spent per collection split by operation, so the collection actually carrying the load is visible rather than guessed. Run it against the admin database.",
			Query:       DefaultTopCollections,
		},
		{
			Name: "Collection stats", Type: database.STATISTIC,
			Description: "Shows one collection's document count, storage size and per-index sizes. Replace 'collection' with the collection name.",
			Params:      []string{"collection"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultCollectionStats,
		},
		{
			Name: "Index usage", Type: database.STATISTIC,
			Description: "Shows how many times each of a collection's indexes has actually been used since the server started - an index with no accesses is paid for on every write and never read. Replace 'collection' with the collection name.",
			Params:      []string{"collection"},
			Varieties:   []database.SystemRequestVariety{database.DatabaseSensitive},
			Query:       DefaultIndexStats,
		},
		{
			Name: "Profiler slow queries", Type: database.ACTIVITY,
			Description: "Lists profiled operations slower than 100ms, with the plan they used. Requires profiling to have been enabled on this database first.",
			Query:       DefaultProfilerSlowQueries,
		},
	}
}
