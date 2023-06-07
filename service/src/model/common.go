package model

import (
	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type Response struct {
	Response any `json:"response"`
	Error    any `json:"error"`
}

type Login struct {
	Username string `form:"username" json:"username,omitempty"`
	Password string `form:"password" json:"password,omitempty"`
}

type DbConnection struct {
	Host   string    `json:"host"`
	Port   int       `json:"port"`
	CredId uuid.UUID `json:"credId"`
}

type AuthType int8

const (
	NONE AuthType = iota
	BASIC
)

type AppInfo struct {
	Company      string       `json:"company"`
	Configured   bool         `json:"configured"`
	Availability Availability `json:"availability"`
	Secret       SecretStatus `json:"secret"`
	Version      Version      `json:"version"`
	Auth         AuthInfo     `json:"auth"`
}

type AuthInfo struct {
	Type       AuthType `json:"type"`
	Authorised bool     `json:"authorised"`
	Error      string   `json:"error"`
}

type AppConfig struct {
	Company      string       `json:"company"`
	Availability Availability `json:"availability"`
	Auth         AuthConfig   `json:"auth"`
}

type AuthConfig struct {
	Type AuthType          `json:"type"`
	Body map[string]string `json:"body"`
}

type Availability struct {
	ManualQuery bool `json:"manualQuery"`
}

type Version struct {
	Tag    string `json:"tag"`
	Commit string `json:"commit"`
}

type FileUsageType int8

const (
	UPLOAD PasswordType = iota
	PATH
)

type Database struct {
	Host     string  `json:"host"`
	Port     int     `json:"port"`
	Database *string `json:"database"`
}

type Sidecar struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// SPECIFIC (SERVER)
