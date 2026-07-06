package etcd

import (
	"ivory/core/config"
	"ivory/plugins/database"
)

const DefaultAllKeys = `get "" --prefix --keys-only --limit 500`
const DefaultMembers = `member list`
const DefaultEndpointStatus = `endpoint status`
const DefaultAlarms = `alarm list`

func (a *Adapter) SupportedFeatures() map[env.Feature]bool {
	return map[env.Feature]bool{
		env.ViewQueryDbInfo:        false,
		env.ViewQueryDbChart:       false,
		env.ManageQueryDbTemplate:  true,
		env.ManageQueryDbConsole:   true,
		env.ManageQueryDbCancel:    false,
		env.ManageQueryDbTerminate: false,
	}
}

func (a *Adapter) SystemCharts() map[database.SystemChartType]string {
	return map[database.SystemChartType]string{}
}

func (a *Adapter) SystemRequests() []database.SystemRequest {
	return []database.SystemRequest{
		{
			Name: "All keys", Type: database.OTHER,
			Description: "Shows all keys in the cluster without values (limited to 500). Add a prefix to narrow the search or remove --keys-only to see values.",
			Query:       DefaultAllKeys,
		},
		{
			Name: "Members", Type: database.REPLICATION,
			Description: "Shows all cluster members with their peer and client urls.",
			Query:       DefaultMembers,
		},
		{
			Name: "Endpoint status", Type: database.STATISTIC,
			Description: "Shows per-member status: version, database size, leader flag and raft progress.",
			Query:       DefaultEndpointStatus,
		},
		{
			Name: "Alarms", Type: database.ACTIVITY,
			Description: "Shows active cluster alarms such as NOSPACE or CORRUPT.",
			Query:       DefaultAlarms,
		},
	}
}
