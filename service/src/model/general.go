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
	Host         string    `json:"host"`
	Port         int       `json:"port"`
	CredentialId uuid.UUID `json:"credentialId"`
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
	Label  string `json:"label"`
}

type FileUsageType int8

const (
	UPLOAD FileUsageType = iota
	PATH
)

type Database struct {
	Host string  `json:"host"`
	Port int     `json:"port"`
	Name *string `json:"name"`
}

type SidecarStatus string

const (
	Paused SidecarStatus = "PAUSED"
	Active               = "ACTIVE"
)

type Sidecar struct {
	Host   string         `json:"host" form:"host"`
	Port   int            `json:"port" form:"port"`
	Name   *string        `json:"name" form:"name"`
	Status *SidecarStatus `json:"status,omitempty" form:"status"`
}

// SPECIFIC (SERVER)
