package mongo

import (
	"ivory/core/config"
	"ivory/plugins/database"
)

const DefaultServerStatus = `db.runCommand({"serverStatus": 1})`
const DefaultDbStats = `db.runCommand({"dbStats": 1})`
const DefaultListCollections = `db.runCommand({"listCollections": 1})`
const DefaultCurrentOp = `db.runCommand({"currentOp": 1})`

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
	}
}
