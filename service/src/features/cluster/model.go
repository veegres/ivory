package cluster

import (
	"ivory/src/features"
	"ivory/src/features/cert"
	"ivory/src/features/node"
	"ivory/src/plugins/database"
	"ivory/src/plugins/keeper"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type ClusterOptions struct {
	DbType     database.Type `json:"dbType"`
	KeeperType keeper.Type   `json:"keeperType"`
	Tls        ClusterTls    `json:"tls"`
	Certs      cert.Certs    `json:"certs"`
	Vaults     Vaults        `json:"vaults"`
	Tags       []string      `json:"tags"`
}

type Cluster struct {
	ClusterOptions
	Name  string            `json:"name"`
	Nodes []node.Connection `json:"nodes"`
}

type ClusterAuto struct {
	ClusterOptions
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
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
	Nodes          NodeOverview       `json:"nodes"`
	DetectedDomain string             `json:"detectedDomain"`
	Features       []features.Feature `json:"features"`
}

type NodeOverview map[string]Node

type Node struct {
	node.KeeperResponse
	Warnings []string `json:"warnings"`
}

// SPECIFIC (SERVER)
