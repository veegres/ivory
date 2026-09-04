package auth

import "ivory/features/permission"

type AuthType int8

const (
	BASIC AuthType = iota
	LDAP
	OIDC
)

// Prefix is the authority a permission record is filed under when this auth
// type vouches for somebody, which is what the string form of an auth type is
// for in the first place.
func (s AuthType) Prefix() permission.Prefix {
	return []permission.Prefix{permission.PrefixBasic, permission.PrefixLdap, permission.PrefixOidc}[s]
}

func (s AuthType) String() string {
	return string(s.Prefix())
}
