package model

import "github.com/google/uuid"

// COMMON (WEB AND SERVER)

type Login struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type DbConnection struct {
	Host   string    `json:"host"`
	Port   int       `json:"port"`
	CredId uuid.UUID `json:"credId"`
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
