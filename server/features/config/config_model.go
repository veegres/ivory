package config

import (
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/features/user"
)

type NewAppConfig struct {
	Secret    string    `json:"secret,omitempty"`
	AppConfig AppConfig `json:"appConfig"`
	// User is the first user Ivory gets, and it is registered a superuser. Setup
	// is the one place a password is typed on somebody else's behalf, because
	// there is nobody yet who could be sent a registration; every user after
	// this one sets their own.
	User *UserSetupRequest `json:"user,omitempty"`
}

type UserSetupRequest struct {
	Username  string          `json:"username"`
	Password  string          `json:"password"`
	AuthTypes []user.AuthType `json:"authTypes"`
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
