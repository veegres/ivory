package zookeeper

import (
	"ivory/core/config"
	"ivory/plugins/database"
)

const DefaultRootList = `ls /`
const DefaultConfig = `get /zookeeper/config`

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
	}
}
