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

// renamed resolves stored feature keys to the names they go by now, so a
// permission saved before a feature was renamed keeps its status instead of
// being dropped and reset to the default.
func (p PermissionMap) renamed() PermissionMap {
	renamed := make(PermissionMap, len(p))
	for feature, status := range p {
		current := feature.Current()
		// NOTE: a map written mid-rename can carry both keys - the current one
		// wins whichever order the range happens to visit them in
		if _, ok := p[current]; ok && current != feature {
			continue
		}
		renamed[current] = status
	}
	return renamed
}

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

type PermissionRequest struct {
	Permissions []config.Feature `json:"permissions" binding:"required"`
}

// UserPermissionsMap is a map where the key is username/email and value is the permissions map
type UserPermissionsMap map[string]PermissionMap
