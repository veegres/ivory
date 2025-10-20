package auth

type AuthType int8

const (
	NONE AuthType = iota
	BASIC
	LDAP
	OIDC
)

type Info struct {
	Type       AuthType `json:"type"`
	Authorised bool     `json:"authorised"`
	Error      string   `json:"error"`
}
