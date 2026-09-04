package auth

import "ivory/features/user"

// AuthType is user.AuthType under another name: this is where a person's way
// of signing in gets verified, but the vocabulary of what that way even is
// belongs to features/user, since a user is who gets registered for it.
type AuthType = user.AuthType

const (
	BASIC = user.AuthBasic
	LDAP  = user.AuthLdap
	OIDC  = user.AuthOidc
)
