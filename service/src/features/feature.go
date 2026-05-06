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
	ManageNodeDbConfigUpdate Feature = "manage.node.db.config.update"
	ManageNodeDbSwitchover   Feature = "manage.node.db.switchover"
	ManageNodeDbReinitialize Feature = "manage.node.db.reinitialize"
	ManageNodeDbRestart      Feature = "manage.node.db.restart"
	ManageNodeDbReload       Feature = "manage.node.db.reload"
	ManageNodeDbFailover     Feature = "manage.node.db.failover"
	ManageNodeDbActivation   Feature = "manage.node.db.activation"

	ViewNodeSshMetrics  Feature = "view.node.ssh.metrics"
	ViewNodeSshDocker   Feature = "view.node.ssh.docker"
	ManageNodeSshDocker Feature = "manage.node.ssh.docker"

	// Query features
	ViewQueryCrudList     Feature = "view.query.crud.list"
	ManageQueryCrudCreate Feature = "manage.query.crud.create"
	ManageQueryCrudUpdate Feature = "manage.query.crud.update"
	ManageQueryCrudDelete Feature = "manage.query.crud.delete"

	ViewQueryDbInfo        Feature = "view.query.db.info"
	ViewQueryDbChart       Feature = "view.query.db.chart"
	ManageQueryDbTemplate  Feature = "manage.query.db.template"
	ManageQueryDbConsole   Feature = "manage.query.db.console"
	ManageQueryDbCancel    Feature = "manage.query.db.cancel"
	ManageQueryDbTerminate Feature = "manage.query.db.terminate"

	ViewQueryLogList     Feature = "view.query.log.list"
	ManageQueryLogDelete Feature = "manage.query.log.delete"

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

	// Tool features
	ViewToolBloatList  Feature = "view.tool.bloat.list"
	ViewToolBloatItem  Feature = "view.tool.bloat.item"
	ViewToolBloatLogs  Feature = "view.tool.bloat.logs"
	ManageToolBloatJob Feature = "manage.tool.bloat.job"

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
	ViewQueryCrudList,
	ManageQueryCrudCreate,
	ManageQueryCrudUpdate,
	ManageQueryCrudDelete,
	ViewQueryDbInfo,
	ViewQueryDbChart,
	ManageQueryDbTemplate,
	ManageQueryDbConsole,
	ManageQueryDbCancel,
	ManageQueryDbTerminate,
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
	ViewToolBloatList,
	ViewToolBloatItem,
	ViewToolBloatLogs,
	ManageToolBloatJob,
	ViewManagementSecret,
	ManageManagementSecret,
	ManageManagementErase,
	ManageManagementFree,
	ManageManagementBackup,
}
