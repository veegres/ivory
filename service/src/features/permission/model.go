package permission

// COMMON (WEB AND SERVER)

type PermissionStatus int

const (
	NOT_PERMITTED PermissionStatus = iota
	PENDING
	GRANTED
)

type UserPermissions struct {
	Username    string                      `json:"username"`
	Permissions map[string]PermissionStatus `json:"permissions"`
}

// SPECIFIC (SERVER)

// UserPermissionsMap is a map where the key is username/email and value is the permissions map
type UserPermissionsMap map[string]map[string]PermissionStatus

var PermissionList = []string{
	// Cluster permissions
	"view.cluster.list",
	"view.cluster.item",
	"manage.cluster.create",
	"manage.cluster.update",
	"manage.cluster.delete",

	// Tags permissions
	"view.tag.list",

	// Instance permissions
	"view.instance.overview",
	"view.instance.config",
	"manage.instance.config.update",
	"manage.instance.switchover",
	"manage.instance.reinitialize",
	"manage.instance.restart",
	"manage.instance.reload",
	"manage.instance.failover",
	"manage.instance.activation",

	// Query permissions
	"view.query.list",
	"manage.query.create",
	"manage.query.update",
	"manage.query.delete",
	"view.query.execute.info",
	"view.query.execute.chart",
	"manage.query.execute.template",
	"manage.query.execute.console",
	"manage.query.execute.cancel",
	"manage.query.execute.terminate",
	"view.query.log.list",
	"manage.query.log.delete",

	// Password permissions
	"view.password.list",
	"manage.password.create",
	"manage.password.update",
	"manage.password.delete",

	// Cert permissions
	"view.cert.list",
	"manage.cert.create",
	"manage.cert.delete",

	// Permission management
	"view.permission.list",
	"manage.permission.update",
	"manage.permission.delete",

	// Bloat permissions
	"view.bloat.list",
	"view.bloat.item",
	"view.bloat.logs",
	"manage.bloat.update",
	"manage.bloat.job",

	// Management permissions
	"view.management.secret",
	"manage.management.secret",
	"manage.management.erase",
}
