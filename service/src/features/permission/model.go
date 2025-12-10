package permission

// COMMON (WEB AND SERVER)

type PermissionStatus int

const (
	NOT_PERMITTED PermissionStatus = iota
	PENDING
	GRANTED
)

type Permission struct {
	Name   string           `json:"name"`
	Status PermissionStatus `json:"status"`
}

type UserPermissions struct {
	Username    string                      `json:"username"`
	Permissions map[string]PermissionStatus `json:"permissions"`
}

// UserPermissionsMap is a map where key is username/email and value is the permissions map
type UserPermissionsMap map[string]map[string]PermissionStatus

// SPECIFIC (SERVER)

var PermissionList = []string{
	// Cluster permissions
	"view.cluster.list",
	"view.cluster.item",
	"manage.cluster.create",
	"manage.cluster.update",
	"manage.cluster.delete",

	// Instance permissions
	"view.instance.list",
	"view.instance.item",
	"manage.instance.create",
	"manage.instance.update",
	"manage.instance.delete",

	// Query permissions
	"view.query.list",
	"view.query.item",
	"manage.query.execute",
	"manage.query.cancel",

	// Password/Credential permissions
	"view.password.list",
	"view.password.item",
	"manage.password.create",
	"manage.password.update",
	"manage.password.delete",

	// User permissions
	"view.user.list",
	"view.user.item",
	"manage.user.create",
	"manage.user.update",
	"manage.user.delete",

	// Permission management
	"view.permission.list",
	"manage.permission.update",
}
