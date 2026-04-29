package database

import (
	"crypto/tls"
)

// COMMON (WEB AND SERVER)

type Type int8

const (
	POSTGRES Type = iota
	ETCD
)

type Config struct {
	Type   Type    `json:"type"`
	Host   string  `json:"host"`
	Port   int     `json:"port"`
	Name   *string `json:"name"`
	Schema *string `json:"schema"`
}

// SPECIFIC (SERVER)

type Connection struct {
	Config      Config       `json:"config" form:"config"`
	Credentials *Credentials `json:"credentials" form:"credentials"`
	TlsConfig   *tls.Config  `json:"tlsConfig" form:"tlsConfig"`
}

type Context struct {
	Connection *Connection `json:"connection"`
	Session    string      `json:"session"`
}

type Credentials struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type QueryAnalysis struct {
	LIMIT     int
	UPDATE    int
	DELETE    int
	INSERT    int
	SELECT    int
	FROM      int
	EXPLAIN   int
	Semicolon bool
}

type QueryOptions struct {
	Params []any   `json:"params"`
	Limit  *string `json:"limit"`
	Trim   *bool   `json:"trim"`
}

type QueryField struct {
	Name        string `json:"name"`
	DataType    string `json:"dataType"`
	DataTypeOID uint32 `json:"dataTypeOID"`
}

type QueryFields struct {
	Fields    []QueryField  `json:"fields"`
	Rows      [][]any       `json:"rows"`
	StartTime int64         `json:"startTime"`
	EndTime   int64         `json:"endTime"`
	Url       string        `json:"url"`
	Options   *QueryOptions `json:"options"`
}
