package management

import (
	"ivory/src/clients/env"
	"ivory/src/features/auth"
	"ivory/src/features/config"
	"ivory/src/features/secret"
)

// COMMON (WEB AND SERVER)

type Response struct {
	Response any `json:"response"`
	Error    any `json:"error"`
}

type AppInfo struct {
	Company      string              `json:"company"`
	Configured   bool                `json:"configured"`
	Availability config.Availability `json:"availability"`
	Version      env.Version         `json:"version"`
	Secret       secret.SecretStatus `json:"secret"`
	Auth         auth.AuthInfo       `json:"auth"`
}
