package permission

import (
	"ivory/core/config"
	"slices"
)

// COMMON (WEB AND SERVER)

type Status int

const (
	NOT_PERMITTED Status = iota
	PENDING
	GRANTED
)

type PermissionMap map[config.Feature]Status

// without drops the features nobody may hold, so an answer never offers a
// permission the caller has already said is not on the table.
func (p PermissionMap) without(features []config.Feature) PermissionMap {
	if len(features) == 0 {
		return p
	}
	kept := make(PermissionMap, len(p))
	for feature, status := range p {
		if !slices.Contains(features, feature) {
			kept[feature] = status
		}
	}
	return kept
}

type UserPermissions struct {
	Username    string        `json:"username"`
	Permissions PermissionMap `json:"permissions"`
}

// SPECIFIC (SERVER)

type RouterPermissionRequest struct {
	Permissions []config.Feature `json:"permissions" binding:"required"`
}

// UserPermissionsMap is a map where the key is username/email and value is the permissions map
type UserPermissionsMap map[string]PermissionMap
