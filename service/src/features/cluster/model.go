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

type Options struct {
	Plugins Plugins    `json:"plugins"`
	Tls     Tls        `json:"tls"`
	Certs   cert.Certs `json:"certs"`
	Vaults  Vaults     `json:"vaults"`
	Tags    []string   `json:"tags"`
}

type Request struct {
	Name  string            `json:"name"`
	Nodes []node.Connection `json:"nodes"`
	Options
}

type Response struct {
	Name  string            `json:"name"`
	Nodes []node.Connection `json:"nodes"`
	Options
}

type AutoRequest struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Options
}

type Tls struct {
	Keeper   bool `json:"keeper"`
	Database bool `json:"database"`
}

type Plugins struct {
	Keeper   keeper.Plugin   `json:"keeper"`
	Database database.Plugin `json:"database"`
}

type Vaults struct {
	KeeperId   *uuid.UUID `json:"keeperId"`
	DatabaseId *uuid.UUID `json:"databaseId"`
	SshKeyId   *uuid.UUID `json:"sshKeyId"`
}

type Overview struct {
	Nodes          NodeMap            `json:"nodes"`
	DetectedDomain string             `json:"detectedDomain"`
	Features       []features.Feature `json:"features"`
}

type NodeMap = map[string]Node

type Node struct {
	node.KeeperResponse
	Warnings []string `json:"warnings"`
}

// SPECIFIC (SERVER)
