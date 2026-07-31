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
			Description: "Lists keys matching a pattern (default: all keys in the selected database).",
			Query:       DefaultKeys,
		},
	}
}
