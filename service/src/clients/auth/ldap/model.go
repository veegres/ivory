package ldap

type Login struct {
	Username string `form:"username" json:"username,omitempty"`
	Password string `form:"password" json:"password,omitempty"`
}

type Config struct {
	Url      string     `json:"url"`
	BindDN   string     `json:"bindDN"`
	BindPass string     `json:"bindPass"`
	BaseDN   string     `json:"baseDN"`
	Filter   string     `json:"filter"`
	Tls      *TlsConfig `json:"tls,omitempty"`
}

type TlsConfig struct {
	CaCert string `json:"caCert"`
}
