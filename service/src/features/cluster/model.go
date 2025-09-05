package cluster

import (
	"ivory/src/clients/sidecar"
	"ivory/src/features/cert"

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
	Name      string            `json:"name"`
	Instances []sidecar.Sidecar `json:"instances"`
}

type ClusterAuto struct {
	ClusterOptions
	Name     string          `json:"name"`
	Instance sidecar.Sidecar `json:"instance"`
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
