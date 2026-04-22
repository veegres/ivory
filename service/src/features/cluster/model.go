package cluster

import (
	"ivory/src/clients/keeper"
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
	Name    string          `json:"name"`
	Keepers []keeper.Keeper `json:"sidecars"`
}

type ClusterAuto struct {
	ClusterOptions
	Name string        `json:"name"`
	Node keeper.Keeper `json:"node"`
}

type ClusterTls struct {
	Keeper   bool `json:"sidecar"`
	Database bool `json:"database"`
}

type Credentials struct {
	PatroniId  *uuid.UUID `json:"patroniId"`
	PostgresId *uuid.UUID `json:"postgresId"`
}

type ClusterOverview struct {
	Nodes      NodeOverview   `json:"nodes"`
	DetectedBy *keeper.Keeper `json:"detectedBy"`
	MainNode   *Node          `json:"mainNode"`
}

type NodeOverview map[string]*Node

type Node struct {
	keeper.Node
	InCluster bool `json:"inCluster"`
	InKeeper  bool `json:"inSidecar"`
}

// SPECIFIC (SERVER)
