package permission

// COMMON (WEB AND SERVER)

type PermissionStatus int

const (
	NOT_PERMITTED PermissionStatus = iota
	PENDING
	GRANTED
)

type PermissionMap map[Permission]PermissionStatus

type UserPermissions struct {
	Username    string        `json:"username"`
	Permissions PermissionMap `json:"permissions"`
}

// SPECIFIC (SERVER)

type PermissionRequest struct {
	Permissions []Permission `json:"permissions" binding:"required"`
}

// UserPermissionsMap is a map where the key is username/email and value is the permissions map
type UserPermissionsMap map[string]PermissionMap

type Permission string

const (
	// Cluster permissions
	ViewClusterList     Permission = "view.cluster.list"
	ViewClusterItem     Permission = "view.cluster.item"
	ViewClusterOverview Permission = "view.cluster.overview"
	ManageClusterCreate Permission = "manage.cluster.create"
	ManageClusterUpdate Permission = "manage.cluster.update"
	ManageClusterDelete Permission = "manage.cluster.delete"

	// Tags permissions
	ViewTagList Permission = "view.tag.list"

	// Instance permissions
	ViewInstanceOverview       Permission = "view.instance.overview"
	ViewInstanceConfig         Permission = "view.instance.config"
	ManageInstanceConfigUpdate Permission = "manage.instance.config.update"
	ManageInstanceSwitchover   Permission = "manage.instance.switchover"
	ManageInstanceReinitialize Permission = "manage.instance.reinitialize"
	ManageInstanceRestart      Permission = "manage.instance.restart"
	ManageInstanceReload       Permission = "manage.instance.reload"
	ManageInstanceFailover     Permission = "manage.instance.failover"
	ManageInstanceActivation   Permission = "manage.instance.activation"

	// Query permissions
	ViewQueryList               Permission = "view.query.list"
	ManageQueryCreate           Permission = "manage.query.create"
	ManageQueryUpdate           Permission = "manage.query.update"
	ManageQueryDelete           Permission = "manage.query.delete"
	ViewQueryExecuteInfo        Permission = "view.query.execute.info"
	ViewQueryExecuteChart       Permission = "view.query.execute.chart"
	ManageQueryExecuteTemplate  Permission = "manage.query.execute.template"
	ManageQueryExecuteConsole   Permission = "manage.query.execute.console"
	ManageQueryExecuteCancel    Permission = "manage.query.execute.cancel"
	ManageQueryExecuteTerminate Permission = "manage.query.execute.terminate"
	ViewQueryLogList            Permission = "view.query.log.list"
	ManageQueryLogDelete        Permission = "manage.query.log.delete"

	// Password permissions
	ViewPasswordList     Permission = "view.password.list"
	ManagePasswordCreate Permission = "manage.password.create"
	ManagePasswordUpdate Permission = "manage.password.update"
	ManagePasswordDelete Permission = "manage.password.delete"

	// Cert permissions
	ViewCertList     Permission = "view.cert.list"
	ManageCertCreate Permission = "manage.cert.create"
	ManageCertDelete Permission = "manage.cert.delete"

	// Permission management
	ViewPermissionList     Permission = "view.permission.list"
	ManagePermissionUpdate Permission = "manage.permission.update"
	ManagePermissionDelete Permission = "manage.permission.delete"

	// Bloat permissions
	ViewBloatList  Permission = "view.bloat.list"
	ViewBloatItem  Permission = "view.bloat.item"
	ViewBloatLogs  Permission = "view.bloat.logs"
	ManageBloatJob Permission = "manage.bloat.job"

	// Management permissions
	ViewManagementSecret   Permission = "view.management.secret"
	ManageManagementSecret Permission = "manage.management.secret"
	ManageManagementErase  Permission = "manage.management.erase"
	ManageManagementFree   Permission = "manage.management.free"
	ManageManagementExport Permission = "manage.management.export"
	ManageManagementImport Permission = "manage.management.import"
)

var Permissions = []Permission{
	ViewClusterList,
	ViewClusterItem,
	ViewClusterOverview,
	ManageClusterCreate,
	ManageClusterUpdate,
	ManageClusterDelete,
	ViewTagList,
	ViewInstanceOverview,
	ViewInstanceConfig,
	ManageInstanceConfigUpdate,
	ManageInstanceSwitchover,
	ManageInstanceReinitialize,
	ManageInstanceRestart,
	ManageInstanceReload,
	ManageInstanceFailover,
	ManageInstanceActivation,
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
	ViewPasswordList,
	ManagePasswordCreate,
	ManagePasswordUpdate,
	ManagePasswordDelete,
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
	ManageManagementExport,
	ManageManagementImport,
}
