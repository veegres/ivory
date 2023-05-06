package model

import "github.com/google/uuid"

// COMMON (WEB AND SERVER)

type ClusterOptions struct {
	Certs       Certs       `json:"certs"`
	Credentials Credentials `json:"credentials"`
	Tags        []string    `json:"tags"`
}

type Cluster struct {
	ClusterOptions
	Name      string    `json:"name"`
	Instances []Sidecar `json:"instances"`
}

type ClusterAuto struct {
	ClusterOptions
	Name     string  `json:"name"`
	Instance Sidecar `json:"instance"`
}

type Certs struct {
	ClientCAId   *uuid.UUID `json:"clientCAId"`
	ClientKeyId  *uuid.UUID `json:"clientKeyId"`
	ClientCertId *uuid.UUID `json:"clientCertId"`
}

type Credentials struct {
	PatroniId  *uuid.UUID `json:"patroniId"`
	PostgresId *uuid.UUID `json:"postgresId"`
}

// SPECIFIC (SERVER)
