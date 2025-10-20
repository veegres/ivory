package config

import (
	"ivory/src/clients/auth/basic"
	"ivory/src/clients/auth/ldap"
	"ivory/src/clients/auth/oidc"
	"ivory/src/features/auth"
)

type Info struct {
	Configured   bool         `json:"configured"`
	Error        string       `json:"error"`
	Company      string       `json:"company"`
	Availability Availability `json:"availability"`
}

type AppConfig struct {
	Company      string       `json:"company"`
	Availability Availability `json:"availability"`
	Auth         AuthConfig   `json:"auth"`
}

type AuthConfig struct {
	Type  auth.AuthType `json:"type"`
	Basic *basic.Config `json:"basic"`
	Ldap  *ldap.Config  `json:"ldap"`
	Oidc  *oidc.Config  `json:"oidc"`
}

type Availability struct {
	ManualQuery bool `json:"manualQuery"`
}
