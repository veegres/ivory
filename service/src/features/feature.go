package features

// COMMON (WEB AND SERVER)

type Feature string

const (
	// Cluster features
	ViewClusterList     Feature = "view.cluster.list"
	ViewClusterItem     Feature = "view.cluster.item"
	ViewClusterOverview Feature = "view.cluster.overview"
	ManageClusterCreate Feature = "manage.cluster.create"
	ManageClusterUpdate Feature = "manage.cluster.update"
	ManageClusterDelete Feature = "manage.cluster.delete"

	// Tags features
	ViewTagList Feature = "view.tag.list"

	// Node features
	ViewNodeDbOverview       Feature = "view.node.db.overview"
	ViewNodeDbConfig         Feature = "view.node.db.config"
	ViewNodeSshMetrics       Feature = "view.node.ssh.metrics"
	ManageNodeDbConfigUpdate Feature = "manage.node.db.config.update"
	ManageNodeDbSwitchover   Feature = "manage.node.db.switchover"
	ManageNodeDbReinitialize Feature = "manage.node.db.reinitialize"
	ManageNodeDbRestart      Feature = "manage.node.db.restart"
	ManageNodeDbReload       Feature = "manage.node.db.reload"
	ManageNodeDbFailover     Feature = "manage.node.db.failover"
	ManageNodeDbActivation   Feature = "manage.node.db.activation"
	ManageNodeSshDocker      Feature = "manage.node.ssh.docker"
	ViewNodeSshDocker        Feature = "view.node.ssh.docker"

	// Query features
	ViewQueryList               Feature = "view.query.list"
	ManageQueryCreate           Feature = "manage.query.create"
	ManageQueryUpdate           Feature = "manage.query.update"
	ManageQueryDelete           Feature = "manage.query.delete"
	ViewQueryExecuteInfo        Feature = "view.query.execute.info"
	ViewQueryExecuteChart       Feature = "view.query.execute.chart"
	ManageQueryExecuteTemplate  Feature = "manage.query.execute.template"
	ManageQueryExecuteConsole   Feature = "manage.query.execute.console"
	ManageQueryExecuteCancel    Feature = "manage.query.execute.cancel"
	ManageQueryExecuteTerminate Feature = "manage.query.execute.terminate"
	ViewQueryLogList            Feature = "view.query.log.list"
	ManageQueryLogDelete        Feature = "manage.query.log.delete"

	// Vault features
	ViewVaultList     Feature = "view.vault.list"
	ManageVaultCreate Feature = "manage.vault.create"
	ManageVaultUpdate Feature = "manage.vault.update"
	ManageVaultDelete Feature = "manage.vault.delete"

	// Cert features
	ViewCertList     Feature = "view.cert.list"
	ManageCertCreate Feature = "manage.cert.create"
	ManageCertDelete Feature = "manage.cert.delete"

	// Permission management features
	ViewPermissionList     Feature = "view.permission.list"
	ManagePermissionUpdate Feature = "manage.permission.update"
	ManagePermissionDelete Feature = "manage.permission.delete"

	// Bloat features
	ViewBloatList  Feature = "view.bloat.list"
	ViewBloatItem  Feature = "view.bloat.item"
	ViewBloatLogs  Feature = "view.bloat.logs"
	ManageBloatJob Feature = "manage.bloat.job"

	// Management features
	ViewManagementSecret   Feature = "view.management.secret"
	ManageManagementSecret Feature = "manage.management.secret"
	ManageManagementErase  Feature = "manage.management.erase"
	ManageManagementFree   Feature = "manage.management.free"
	ManageManagementBackup Feature = "manage.management.backup"
)

var All = []Feature{
	ViewClusterList,
	ViewClusterItem,
	ViewClusterOverview,
	ManageClusterCreate,
	ManageClusterUpdate,
	ManageClusterDelete,
	ViewTagList,
	ViewNodeDbOverview,
	ViewNodeDbConfig,
	ViewNodeSshMetrics,
	ManageNodeDbConfigUpdate,
	ManageNodeDbSwitchover,
	ManageNodeDbReinitialize,
	ManageNodeDbRestart,
	ManageNodeDbReload,
	ManageNodeDbFailover,
	ManageNodeDbActivation,
	ManageNodeSshDocker,
	ViewNodeSshDocker,
	ViewQueryList,
	ManageQueryCreate,
	ManageQueryUpdate,
	ManageQueryDelete,
	ViewQueryExecuteInfo,
	ViewQueryExecuteChart,
	ManageQueryExecuteTemplate,
	ManageQueryExecuteConsole,
	ManageQueryExecuteCancel,
	ManageQueryExecuteTerminate,
	ViewQueryLogList,
	ManageQueryLogDelete,
	ViewVaultList,
	ManageVaultCreate,
	ManageVaultUpdate,
	ManageVaultDelete,
	ViewCertList,
	ManageCertCreate,
	ManageCertDelete,
	ViewPermissionList,
	ManagePermissionUpdate,
	ManagePermissionDelete,
	ViewBloatList,
	ViewBloatItem,
	ViewBloatLogs,
	ManageBloatJob,
	ViewManagementSecret,
	ManageManagementSecret,
	ManageManagementErase,
	ManageManagementFree,
	ManageManagementBackup,
}
