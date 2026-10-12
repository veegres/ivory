package ldap

type Login struct {
	Username string `form:"username" json:"username,omitempty" binding:"required"`
	Password string `form:"password" json:"password,omitempty" binding:"required"`
}

type Config struct {
	Url      string     `json:"url" binding:"required,notblank"`
	BindDN   string     `json:"bindDN" binding:"required,notblank"`
	BindPass string     `json:"bindPass" binding:"required"`
	BaseDN   string     `json:"baseDN" binding:"required,notblank"`
	Filter   string     `json:"filter" binding:"required,notblank"`
	Tls      *TlsConfig `json:"tls,omitempty"`
}

type TlsConfig struct {
	CaCert string `json:"caCert"`
}
