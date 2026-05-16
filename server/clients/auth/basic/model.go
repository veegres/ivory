package basic

type Login struct {
	Username string `form:"username" json:"username,omitempty"`
	Password string `form:"password" json:"password,omitempty"`
}

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
