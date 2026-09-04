package basic

type Login struct {
	Username string `form:"username" json:"username,omitempty"`
	Password string `form:"password" json:"password,omitempty"`
}

// Config carries no credentials on purpose: basic auth authenticates against
// the Ivory users, so the only thing the app config states is that it is on.
type Config struct{}
