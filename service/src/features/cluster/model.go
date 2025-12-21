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

type ClusterOverview struct {
	Instances    map[string]Instance `json:"instances"`
	DetectedBy   *sidecar.Sidecar    `json:"detectedBy"`
	MainInstance *Instance           `json:"mainInstance"`
}

type Instance struct {
	sidecar.Instance
	InCluster bool `json:"inCluster"`
	InSidecar bool `json:"inSidecar"`
}

// SPECIFIC (SERVER)
