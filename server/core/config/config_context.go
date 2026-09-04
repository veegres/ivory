package config

// AuthContextKey are the gin.Context keys the auth middlewares populate and other features read.
var AuthContextKey = struct {
	Enabled    string
	Authorised string
	Username   string
	Type       string
	Superuser  string
	Error      string
	Session    string
}{
	Enabled:    "auth",
	Authorised: "authorised",
	Username:   "username",
	Type:       "authType",
	Superuser:  "superuser",
	Error:      "authError",
	Session:    "session",
}
