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
