package config

type AuthType int8

const (
	NONE AuthType = iota
	BASIC
	LDAP
	OIDC
)

type AppConfig struct {
	Company      string       `json:"company"`
	Availability Availability `json:"availability"`
	Auth         AuthConfig   `json:"auth"`
}

type AuthConfig struct {
	Type  AuthType     `json:"type"`
	Basic *BasicConfig `json:"basic"`
	Ldap  *LdapConfig  `json:"ldap"`
	Oidc  *OidcConfig  `json:"oidc"`
}

type LdapConfig struct {
	Url      string `json:"url"`
	BindDN   string `json:"bindDN"`
	BindPass string `json:"bindPass"`
	BaseDN   string `json:"baseDN"`
	Filter   string `json:"filter,omitempty"`
}

type OidcConfig struct {
	IssuerURL    string `json:"issuerUrl"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURL  string `json:"redirectUrl"`
}

type BasicConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Availability struct {
	ManualQuery bool `json:"manualQuery"`
}
