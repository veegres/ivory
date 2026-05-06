package permission

import "ivory/src/features"

// COMMON (WEB AND SERVER)

type Status int

const (
	NOT_PERMITTED Status = iota
	PENDING
	GRANTED
)

type PermissionMap map[features.Feature]Status

type UserPermissions struct {
	Username    string        `json:"username"`
	Permissions PermissionMap `json:"permissions"`
}

// SPECIFIC (SERVER)

type PermissionRequest struct {
	Permissions []features.Feature `json:"permissions" binding:"required"`
}

// UserPermissionsMap is a map where the key is username/email and value is the permissions map
type UserPermissionsMap map[string]PermissionMap
