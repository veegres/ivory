package model

// COMMON (WEB AND SERVER)

type Password struct {
	Username string       `json:"username"`
	Password string       `json:"password"`
	Type     PasswordType `json:"type"`
}

type PasswordType int

const (
	POSTGRES PasswordType = iota
	PATRONI
)

type PasswordMap map[string]Password

// SPECIFIC (SERVER)
