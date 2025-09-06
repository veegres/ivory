package management

import (
	"ivory/src/features/auth"
	"ivory/src/features/config"
	"ivory/src/features/secret"
	"ivory/src/storage/env"
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
	Auth         auth.Info           `json:"auth"`
}

type SecretUpdateRequest struct {
	PreviousKey string `json:"previousKey"`
	NewKey      string `json:"newKey"`
}

// SPECIFIC (SERVER)
