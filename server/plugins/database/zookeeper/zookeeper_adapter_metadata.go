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
	}
}
