package auth

import "ivory/src/features/config"

type Info struct {
	Type       config.AuthType `json:"type"`
	Authorised bool            `json:"authorised"`
	Error      string          `json:"error"`
}

type Login struct {
	Username string `form:"username" json:"username,omitempty"`
	Password string `form:"password" json:"password,omitempty"`
}
