package zookeeper

import (
	"ivory/core/config"
	"ivory/plugins/database"
)

const DefaultRootList = `ls /`
const DefaultConfig = `get /zookeeper/config`
const DefaultInternals = `ls /zookeeper`
const DefaultPatroniNamespace = `ls /service`
const DefaultClickhouseTree = `ls /clickhouse`
const DefaultMonitor = `mntr`
const DefaultServer = `srvr`
const DefaultStat = `stat`
const DefaultRuok = `ruok`
const DefaultServerConfig = `conf`

func (a *Adapter) SupportedFeatures() map[config.Feature]bool {
	return map[config.Feature]bool{
		config.ViewQueryDbInfo:        false,
		config.ViewQueryDbChart:       false,
		config.ManageQueryDbTemplate:  true,
		config.ManageQueryDbConsole:   true,
		config.ManageQueryDbCancel:    false,
		config.ManageQueryDbTerminate: false,
	}
}

func (a *Adapter) SystemCharts() map[database.SystemChartType]string {
	return map[database.SystemChartType]string{}
}

func (a *Adapter) SystemRequests() []database.SystemRequest {
	return []database.SystemRequest{
		{
			Name: "Root children", Type: database.OTHER,
			Description: "Lists the znodes directly under the root path.",
			Query:       DefaultRootList,
		},
		{
			Name: "Ensemble config", Type: database.REPLICATION,
			Description: "Shows the dynamic ensemble membership config znode.",
			Query:       DefaultConfig,
		},
		{
			Name: "Server internals", Type: database.OTHER,
			Description: "Lists zookeeper's own reserved subtree, where the ensemble config and quota znodes live.",
			Query:       DefaultInternals,
		},
		{
			Name: "Patroni namespace", Type: database.REPLICATION,
			Description: "Lists the clusters registered under patroni's default /service namespace, where the leader key and member entries patroni elects on are kept. An empty result on a running cluster means patroni is coordinating somewhere else.",
			Query:       DefaultPatroniNamespace,
		},
		{
			Name: "ClickHouse coordination", Type: database.REPLICATION,
			Description: "Lists the subtree ClickHouse registers its replicated tables under. It is where a replica whose session dropped stops appearing, which is what turns its tables read-only.",
			Query:       DefaultClickhouseTree,
		},
		{
			Name: "Server health", Type: database.STATISTIC,
			Description: "The server's own metrics (mntr): its state, outstanding requests, latency, and on the leader how many followers are synced. A server that has lost its quorum answers that it is not currently serving requests.",
			Query:       DefaultMonitor,
		},
		{
			Name: "Server status", Type: database.REPLICATION,
			Description: "The server's mode (leader, follower, observer, standalone), its zxid and latency (srvr). It is what this node itself thinks its role is.",
			Query:       DefaultServer,
		},
		{
			Name: "Server connections", Type: database.ACTIVITY,
			Description: "Like Server status, plus every client connected to this server (stat), to find who holds the sessions.",
			Query:       DefaultStat,
		},
		{
			Name: "Server alive", Type: database.OTHER,
			Description: "Asks the server whether it is running (ruok); it answers imok. It says nothing about the quorum - a server without one still answers.",
			Query:       DefaultRuok,
		},
		{
			Name: "Server config", Type: database.REPLICATION,
			Description: "The server's effective configuration (conf): its ports, data directories and the ensemble's server.N lines.",
			Query:       DefaultServerConfig,
		},
	}
}
