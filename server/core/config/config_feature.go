package config

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

// renamedFeatures maps a feature key that used to be stored under a different
// name to the one it goes by now. Permissions persist the key they were written
// with, so a rename has to stay readable or every stored grant silently resets
// to the default the next time the database is normalized.
var renamedFeatures = map[Feature]Feature{
	"view.node.platform":   ViewNodeSystem,
	"manage.node.platform": ManageNodeSystem,

	// v1 called a node an instance, and its keeper features carried no keeper
	// in the name because there was only one
	"view.instance.overview":        ViewNodeKeeperOverview,
	"view.instance.config":          ViewNodeKeeperConfig,
	"manage.instance.config.update": ManageNodeKeeperConfigUpdate,
	"manage.instance.switchover":    ManageNodeKeeperSwitchover,
	"manage.instance.reinitialize":  ManageNodeKeeperReinitialize,
	"manage.instance.restart":       ManageNodeKeeperRestart,
	"manage.instance.reload":        ManageNodeKeeperReload,
	"manage.instance.failover":      ManageNodeKeeperFailover,
	"manage.instance.activation":    ManageNodeKeeperActivation,

	// v1's saved queries became query.crud, and running one became query.db
	"view.query.list":                ViewQueryCrudList,
	"manage.query.create":            ManageQueryCrudCreate,
	"manage.query.update":            ManageQueryCrudUpdate,
	"manage.query.delete":            ManageQueryCrudDelete,
	"view.query.execute.info":        ViewQueryDbInfo,
	"view.query.execute.chart":       ViewQueryDbChart,
	"manage.query.execute.template":  ManageQueryDbTemplate,
	"manage.query.execute.console":   ManageQueryDbConsole,
	"manage.query.execute.cancel":    ManageQueryDbCancel,
	"manage.query.execute.terminate": ManageQueryDbTerminate,

	// v1's password store became the vault
	"view.password.list":     ViewVaultList,
	"manage.password.create": ManageVaultCreate,
	"manage.password.update": ManageVaultUpdate,
	"manage.password.delete": ManageVaultDelete,

	// v1's bloat feature became one tool among several
	"view.bloat.list":  ViewToolPgCompactTableList,
	"view.bloat.item":  ViewToolPgCompactTableItem,
	"view.bloat.logs":  ViewToolPgCompactTableLogs,
	"manage.bloat.job": ManageToolPgCompactTableJob,
}

// Current resolves a stored feature key to the name it goes by now, leaving
// anything already current untouched.
func (f Feature) Current() Feature {
	if current, ok := renamedFeatures[f]; ok {
		return current
	}
	return f
}

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

// Withheld names the features nobody holds in this Ivory. It is stated here
// rather than by whoever answers a request, so the permission middleware and
// the app info the UI reads cannot disagree about what a loginless session may
// do.
func Withheld(authEnabled bool) []Feature {
	if authEnabled {
		return nil
	}
	return authOnly
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
	ViewNodeSystem,
	ManageNodeSystem,
	ManageNodeKeeperConfigUpdate,
	ManageNodeKeeperSwitchover,
	ManageNodeKeeperReinitialize,
	ManageNodeKeeperRestart,
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
