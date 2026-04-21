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
	Name     string            `json:"name"`
	Sidecars []sidecar.Sidecar `json:"sidecars"`
}

type ClusterAuto struct {
	ClusterOptions
	Name string          `json:"name"`
	Node sidecar.Sidecar `json:"node"`
}

type ClusterTls struct {
	Sidecar  bool `json:"sidecar"`
	Database bool `json:"database"`
}

type Credentials struct {
	PatroniId  *uuid.UUID `json:"patroniId"`
	PostgresId *uuid.UUID `json:"postgresId"`
}

type ClusterOverview struct {
	Nodes      NodeOverview     `json:"nodes"`
	DetectedBy *sidecar.Sidecar `json:"detectedBy"`
	MainNode   *Node            `json:"mainNode"`
}

type NodeOverview map[string]*Node

type Node struct {
	sidecar.Instance
	InCluster bool `json:"inCluster"`
	InSidecar bool `json:"inSidecar"`
}

// SPECIFIC (SERVER)
