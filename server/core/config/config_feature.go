package config

import "slices"

// COMMON (WEB AND SERVER)

type Feature string

func (f Feature) Valid() bool {
	return slices.Contains(All, f)
}

const (
	// Cluster features
	ViewClusterList     Feature = "view.cluster.list"
	ViewClusterItem     Feature = "view.cluster.item"
	ViewClusterOverview Feature = "view.cluster.overview"
	ViewClusterCheck    Feature = "view.cluster.check"
	ManageClusterCreate Feature = "manage.cluster.create"
	ManageClusterUpdate Feature = "manage.cluster.update"
	ManageClusterDelete Feature = "manage.cluster.delete"

	// Tags features
	ViewTagList Feature = "view.tag.list"

	// Node features
	ViewNodeKeeperOverview       Feature = "view.node.keeper.overview"
	ViewNodeKeeperConfig         Feature = "view.node.keeper.config"
	ViewNodeKeeperCheck          Feature = "view.node.keeper.check"
	ManageNodeKeeperConfigUpdate Feature = "manage.node.keeper.config.update"
	ManageNodeKeeperSwitchover   Feature = "manage.node.keeper.switchover"
	ManageNodeKeeperReinitialize Feature = "manage.node.keeper.reinitialize"
	ManageNodeKeeperRestart      Feature = "manage.node.keeper.restart"
	// Scheduling is its own feature because only patroni takes a time for these.
	ManageNodeKeeperSwitchoverSchedule Feature = "manage.node.keeper.switchover.schedule"
	ManageNodeKeeperRestartSchedule    Feature = "manage.node.keeper.restart.schedule"
	ManageNodeKeeperReload             Feature = "manage.node.keeper.reload"
	ManageNodeKeeperFailover           Feature = "manage.node.keeper.failover"
	ManageNodeKeeperActivation         Feature = "manage.node.keeper.activation"

	ViewNodePlatformCheck       Feature = "view.node.platform.check"
	ViewNodeSystem              Feature = "view.node.system"
	ManageNodeSystem            Feature = "manage.node.system"
	ViewNodePlatformContainer   Feature = "view.node.platform.container"
	ManageNodePlatformContainer Feature = "manage.node.platform.container"

	ViewDeploymentTemplateList     Feature = "view.deployment.template.list"
	ManageDeploymentTemplateCreate Feature = "manage.deployment.template.create"
	ManageDeploymentTemplateUpdate Feature = "manage.deployment.template.update"
	ManageDeploymentTemplateDelete Feature = "manage.deployment.template.delete"

	// Query features
	ViewQueryCrudList     Feature = "view.query.crud.list"
	ManageQueryCrudCreate Feature = "manage.query.crud.create"
	ManageQueryCrudUpdate Feature = "manage.query.crud.update"
	ManageQueryCrudDelete Feature = "manage.query.crud.delete"

	ViewQueryDbInfo        Feature = "view.query.db.info"
	ViewQueryDbChart       Feature = "view.query.db.chart"
	ViewQueryDbCheck       Feature = "view.query.db.check"
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

	// User management features
	ViewUserList            Feature = "view.user.list"
	ManageUserCreate        Feature = "manage.user.create"
	ManageUserUpdate        Feature = "manage.user.update"
	ManageUserPasswordReset Feature = "manage.user.password.reset"
	ManageUserDelete        Feature = "manage.user.delete"

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

// authOnly are the features that only mean anything when Ivory knows who is
// asking. Without authentication there is nobody to register and nobody to
// grant anything to, so a loginless session is not given them at all - the
// alternative is a user manager that can only ever create users who cannot
// sign in.
var authOnly = []Feature{
	ViewUserList,
	ManageUserCreate,
	ManageUserUpdate,
	ManageUserPasswordReset,
	ManageUserDelete,
	ViewPermissionList,
	ManagePermissionUpdate,
	ManagePermissionDelete,
}

// FeaturesWithheldWithoutAuth names the features nobody holds in this Ivory
// when authEnabled is false. It is stated here rather than by whoever answers
// a request, so the permission middleware and the app info the UI reads
// cannot disagree about what a loginless session may do.
func FeaturesWithheldWithoutAuth(authEnabled bool) []Feature {
	if authEnabled {
		return nil
	}
	return authOnly
}

var All = []Feature{
	ViewClusterList,
	ViewClusterItem,
	ViewClusterOverview,
	ViewClusterCheck,
	ManageClusterCreate,
	ManageClusterUpdate,
	ManageClusterDelete,
	ViewTagList,
	ViewNodeKeeperOverview,
	ViewNodeKeeperConfig,
	ViewNodeKeeperCheck,
	ViewNodePlatformCheck,
	ViewNodeSystem,
	ManageNodeSystem,
	ManageNodeKeeperConfigUpdate,
	ManageNodeKeeperSwitchover,
	ManageNodeKeeperSwitchoverSchedule,
	ManageNodeKeeperReinitialize,
	ManageNodeKeeperRestart,
	ManageNodeKeeperRestartSchedule,
	ManageNodeKeeperReload,
	ManageNodeKeeperFailover,
	ManageNodeKeeperActivation,
	ManageNodePlatformContainer,
	ViewNodePlatformContainer,
	ViewDeploymentTemplateList,
	ManageDeploymentTemplateCreate,
	ManageDeploymentTemplateUpdate,
	ManageDeploymentTemplateDelete,
	ViewQueryCrudList,
	ManageQueryCrudCreate,
	ManageQueryCrudUpdate,
	ManageQueryCrudDelete,
	ViewQueryDbInfo,
	ViewQueryDbChart,
	ViewQueryDbCheck,
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
	ViewUserList,
	ManageUserCreate,
	ManageUserUpdate,
	ManageUserPasswordReset,
	ManageUserDelete,
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
