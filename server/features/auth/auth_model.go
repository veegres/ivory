package auth

type AuthType int8

const (
	BASIC AuthType = iota
	LDAP
	OIDC
)

func (s AuthType) String() string {
	return []string{"basic", "ldap", "oidc"}[s]
}
