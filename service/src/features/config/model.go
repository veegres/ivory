package config

import (
	"ivory/src/clients/auth/basic"
	"ivory/src/clients/auth/ldap"
	"ivory/src/clients/auth/oidc"
)

type NewAppConfig struct {
	Secret    string    `json:"secret,omitempty"`
	AppConfig AppConfig `json:"appConfig"`
}

type AppConfig struct {
	Company      string             `json:"company"`
	Availability AvailabilityConfig `json:"availability"`
	Auth         AuthConfig         `json:"auth"`
}

type AuthConfig struct {
	Basic *basic.Config `json:"basic"`
	Ldap  *ldap.Config  `json:"ldap"`
	Oidc  *oidc.Config  `json:"oidc"`
}

type AvailabilityConfig struct {
	ManualQuery bool `json:"manualQuery"`
}
