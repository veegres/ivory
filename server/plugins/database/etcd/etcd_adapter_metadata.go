package etcd

import (
	"ivory/core/config"
	"ivory/plugins/database"
)

const DefaultAllKeys = `get "" --prefix --keys-only --limit 500`
const DefaultMembers = `member list`
const DefaultEndpointStatus = `endpoint status`
const DefaultAlarms = `alarm list`
const DefaultPatroniKeys = `get /service/ --prefix --keys-only`
const DefaultPatroniState = `get /service/ --prefix`

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
			Description: "Shows active cluster alarms such as NOSPACE or CORRUPT. A NOSPACE alarm makes the whole cluster reject writes until the backend is defragmented and the alarm is disarmed, so it explains a patroni cluster that suddenly cannot elect anyone.",
			Query:       DefaultAlarms,
		},
		{
			Name: "Patroni cluster keys", Type: database.REPLICATION,
			Description: "Lists the keys patroni stores under its default /service namespace, without their values. An empty result on a running cluster means patroni is coordinating through a different namespace or a different store entirely.",
			Query:       DefaultPatroniKeys,
		},
		{
			Name: "Patroni cluster state", Type: database.REPLICATION,
			Description: "Shows the same keys with their values: the leader key and its lease, each member's state, and the dynamic config the whole cluster shares. This is the ground truth patroni itself elects on, as opposed to what any one node reports about itself.",
			Query:       DefaultPatroniState,
		},
	}
}
