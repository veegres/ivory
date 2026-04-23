package cluster

import (
	"ivory/src/clients/keeper"
	"ivory/src/features/cert"
	"ivory/src/features/node"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type ClusterOptions struct {
	Tls    ClusterTls `json:"tls"`
	Certs  cert.Certs `json:"certs"`
	Vaults Vaults     `json:"vaults"`
	Tags   []string   `json:"tags"`
}

type Cluster struct {
	ClusterOptions
	Name  string            `json:"name"`
	Nodes []node.Connection `json:"nodes"`
}

type ClusterAuto struct {
	ClusterOptions
	Name string          `json:"name"`
	Node node.Connection `json:"node"`
}

type ClusterTls struct {
	Keeper   bool `json:"sidecar"`
	Database bool `json:"database"`
}

type Vaults struct {
	KeeperId   *uuid.UUID `json:"keeperId"`
	DatabaseId *uuid.UUID `json:"databaseId"`
}

type ClusterOverview struct {
	Nodes      NodeOverview   `json:"nodes"`
	DetectedBy *keeper.Keeper `json:"detectedBy"`
	MainNode   *Node          `json:"mainNode"`
}

type NodeOverview map[string]*Node

type Node struct {
	node.Node
	InCluster bool `json:"inCluster"`
	InKeeper  bool `json:"inSidecar"`
}

// SPECIFIC (SERVER)
