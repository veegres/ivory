package management

import (
	"ivory/src/features/auth"
	"ivory/src/features/permission"
	"ivory/src/features/secret"
	"ivory/src/storage/env"
)

// COMMON (WEB AND SERVER)

type Response struct {
	Response any `json:"response"`
	Error    any `json:"error"`
}

type AppInfo struct {
	Auth    AuthInfo            `json:"auth"`
	Config  ConfigInfo          `json:"config"`
	Version env.Version         `json:"version"`
	Secret  secret.SecretStatus `json:"secret"`
}

type ConfigInfo struct {
	Configured bool   `json:"configured"`
	Error      string `json:"error,omitempty"`
	Company    string `json:"company"`
}

type AuthInfo struct {
	Supported  []auth.AuthType `json:"supported"`
	Authorised bool            `json:"authorised"`
	User       *UserInfo       `json:"user"`
	Error      string          `json:"error,omitempty"`
}

type UserInfo struct {
	Username    string                   `json:"username"`
	Permissions permission.PermissionMap `json:"permissions"`
}

type SecretUpdateRequest struct {
	PreviousKey string `json:"previousKey"`
	NewKey      string `json:"newKey"`
}

// SPECIFIC (SERVER)

// Backup this struct shouldn't change its API over time it is important for backward compatability
// which Ivory and users rely on. It shouldn't import any model, you should always create its own
// struct even if it requires duplicating code
type Backup struct {
	Clusters    []backupCluster     `json:"clusters"`
	Queries     []backupQuery       `json:"queries"`
	Permissions []backupPermissions `json:"permissions"`
}

type backupCluster struct {
	Name     string          `json:"name"`
	Tags     []string        `json:"tags"`
	Sidecars []backupSidecar `json:"sidecars"`
}

type backupSidecar struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type backupQuery struct {
	Name        string               `json:"name"`
	Type        backupQueryType      `json:"type"`
	Varieties   []backupQueryVariety `json:"varieties"`
	Params      []string             `json:"params"`
	Description *string              `json:"description"`
	Default     string               `json:"default"`
	Custom      string               `json:"custom"`
}

type backupQueryType int8

const (
	BLOAT backupQueryType = iota
	ACTIVITY
	REPLICATION
	STATISTIC
	OTHER
)

type backupQueryVariety int8

const (
	DatabaseSensitive backupQueryVariety = iota
	MasterOnly
	ReplicaRecommended
)

type backupPermissions struct {
	Username    string                          `json:"username"`
	Permissions map[string]backupPermissionType `json:"permissions"`
}

type backupPermissionType int

const (
	NOT_PERMITTED backupPermissionType = iota
	PENDING
	GRANTED
)
