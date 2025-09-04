package cluster

import (
	"ivory/src/features/cert"
	"ivory/src/features/instance"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type ClusterOptions struct {
	Tls         ClusterTls  `json:"tls"`
	Certs       cert.Certs  `json:"certs"`
	Credentials Credentials `json:"credentials"`
	Tags        []string    `json:"tags"`
}

type Cluster struct {
	ClusterOptions
	Name      string             `json:"name"`
	Instances []instance.Sidecar `json:"instances"`
}

type ClusterAuto struct {
	ClusterOptions
	Name     string           `json:"name"`
	Instance instance.Sidecar `json:"instance"`
}

type ClusterTls struct {
	Sidecar  bool `json:"sidecar"`
	Database bool `json:"database"`
}

type Credentials struct {
	PatroniId  *uuid.UUID `json:"patroniId"`
	PostgresId *uuid.UUID `json:"postgresId"`
}

// SPECIFIC (SERVER)
