package auth

import (
	"ivory/features/user"
	"time"
)

// AuthType is user.AuthType under another name: this is where a person's way
// of signing in gets verified, but the vocabulary of what that way even is
// belongs to features/user, since a user is who gets registered for it.
type AuthType = user.AuthType

const (
	BASIC = user.AuthBasic
	LDAP  = user.AuthLdap
	OIDC  = user.AuthOidc

	SESSION = user.AuthSession
)

// Audience is what a minted session token was issued for, carried in its aud claim.
type Audience string

const (
	AudienceMcp Audience = "mcp"
)

// RouterTokenResponse is a freshly minted session token and the moment it dies.
// It mirrors user.RegistrationResponse because it is the same kind of answer: a
// token is shown once, when it is issued, and Ivory never stores it to show
// again.
type RouterTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
