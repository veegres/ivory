package auth

import "ivory/features/user"

type AuthType int8

const (
	BASIC AuthType = iota
	LDAP
	OIDC
)

// User is the way of signing in this auth type stands for. The vocabulary
// belongs to features/user, because a user is registered for it, and spelling
// the string form out with it keeps the two from drifting apart.
func (s AuthType) User() user.AuthType {
	return []user.AuthType{user.AuthBasic, user.AuthLdap, user.AuthOidc}[s]
}

func (s AuthType) String() string {
	return string(s.User())
}
