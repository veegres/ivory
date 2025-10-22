package auth

type AuthType int8

const (
	BASIC AuthType = iota
	LDAP
	OIDC
)
