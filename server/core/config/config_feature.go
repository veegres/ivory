package env

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
	ViewNodeKeeperOverview       Feature = "view.node.keeper.overview"
	ViewNodeKeeperConfig         Feature = "view.node.keeper.config"
	ManageNodeKeeperConfigUpdate Feature = "manage.node.keeper.config.update"
	ManageNodeKeeperSwitchover   Feature = "manage.node.keeper.switchover"
	ManageNodeKeeperReinitialize Feature = "manage.node.keeper.reinitialize"
	ManageNodeKeeperRestart      Feature = "manage.node.keeper.restart"
	ManageNodeKeeperReload       Feature = "manage.node.keeper.reload"
	ManageNodeKeeperFailover     Feature = "manage.node.keeper.failover"
	ManageNodeKeeperActivation   Feature = "manage.node.keeper.activation"

	ViewNodePlatform            Feature = "view.node.platform"
	ManageNodePlatform          Feature = "manage.node.platform"
	ViewNodePlatformContainer   Feature = "view.node.platform.container"
	ManageNodePlatformContainer Feature = "manage.node.platform.container"

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
	ViewToolPgCompactTableList  Feature = "view.tool.pg_compacttable.list"
	ViewToolPgCompactTableItem  Feature = "view.tool.pg_compacttable.item"
	ViewToolPgCompactTableLogs  Feature = "view.tool.pg_compacttable.logs"
	ManageToolPgCompactTableJob Feature = "manage.tool.pg_compacttable.job"

	// Management features
	ViewManagementSecret   Feature = "view.management.secret"
	ManageManagementSecret Feature = "manage.management.secret"
	ManageManagementErase  Feature = "manage.management.erase"
	ManageManagementFree   Feature = "manage.management.free"
	ManageManagementBackup Feature = "manage.management.backup"
)

type Plugin interface {
	String() string
}

var All = []Feature{
	ViewClusterList,
	ViewClusterItem,
	ViewClusterOverview,
	ManageClusterCreate,
	ManageClusterUpdate,
	ManageClusterDelete,
	ViewTagList,
	ViewNodeKeeperOverview,
	ViewNodeKeeperConfig,
	ViewNodePlatform,
	ManageNodePlatform,
	ManageNodeKeeperConfigUpdate,
	ManageNodeKeeperSwitchover,
	ManageNodeKeeperReinitialize,
	ManageNodeKeeperRestart,
	ManageNodeKeeperReload,
	ManageNodeKeeperFailover,
	ManageNodeKeeperActivation,
	ManageNodePlatformContainer,
	ViewNodePlatformContainer,
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
	ViewToolPgCompactTableList,
	ViewToolPgCompactTableItem,
	ViewToolPgCompactTableLogs,
	ManageToolPgCompactTableJob,
	ViewManagementSecret,
	ManageManagementSecret,
	ManageManagementErase,
	ManageManagementFree,
	ManageManagementBackup,
}
