package config

import (
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
)

type NewAppConfig struct {
	Secret    string    `json:"secret,omitempty"`
	AppConfig AppConfig `json:"appConfig"`
	// Superuser is the first user Ivory gets. Setup is the one place a password
	// is typed on somebody else's behalf, because there is nobody yet who could
	// be sent a link; every user after this one sets their own.
	Superuser *SuperuserRequest `json:"superuser,omitempty"`
}

type SuperuserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AppConfig struct {
	Company string     `json:"company"`
	Auth    AuthConfig `json:"auth"`
}

type AuthConfig struct {
	Basic *basic.Config `json:"basic"`
	Ldap  *ldap.Config  `json:"ldap"`
	Oidc  *oidc.Config  `json:"oidc"`
}

func (c AuthConfig) enabled() bool {
	return c.Basic != nil || c.Ldap != nil || c.Oidc != nil
}
