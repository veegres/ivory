package cluster

import (
	"ivory/core/config"
	"ivory/core/service/cert"
	"ivory/features/node"
	"ivory/features/query"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type NodeConfig struct {
	Host       string `json:"host" form:"host"`
	SshPort    *int   `json:"sshPort" form:"sshPort"`
	KeeperPort *int   `json:"keeperPort" form:"keeperPort"`
	DbPort     *int   `json:"dbPort" form:"dbPort"`
}

type Options struct {
	Plugins Plugins    `json:"plugins"`
	Tls     Tls        `json:"tls"`
	Certs   cert.Certs `json:"certs"`
	Vaults  Vaults     `json:"vaults"`
	Tags    []string   `json:"tags"`
}

type Request struct {
	Name  string       `json:"name"`
	Nodes []NodeConfig `json:"nodes"`
	Options
}

type Response struct {
	Name  string       `json:"name"`
	Nodes []NodeConfig `json:"nodes"`
	Options
}

type Tls struct {
	Keeper   bool `json:"keeper"`
	Database bool `json:"database"`
}

type Plugins struct {
	Keeper   node.KeeperPlugin `json:"keeper"`
	Database query.DbPlugin    `json:"database"`
}

type Vaults struct {
	KeeperId   *uuid.UUID `json:"keeperId"`
	DatabaseId *uuid.UUID `json:"databaseId"`
	SshKeyId   *uuid.UUID `json:"sshKeyId"`
}

type Node struct {
	Config   NodeConfig             `json:"config"`
	Keeper   node.KeeperOneResponse `json:"keeper"`
	Warnings []string               `json:"warnings"`
}

type Overview struct {
	Nodes    map[string]Node `json:"nodes"`
	Features []env.Feature   `json:"features"`
}

type CreateAutoRequest struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Options
}

type DeployRequest struct {
	Uri                 string            `json:"uri"`
	Parallel            bool              `json:"parallel"`
	NodeConfig          []NodeConfig      `json:"nodeConfig"`
	CommonConfig        CommonConfig      `json:"commonConfig"`
	NodeRawImageOptions map[string]string `json:"nodeRawImageOptions"`
	ClusterOptions      Options           `json:"clusterOptions"`
}

type CommonConfig struct {
	Cluster string `json:"cluster"`
	Dcs     string `json:"dcs"`
	SshUser string `json:"sshUser"`
	SshPass string `json:"sshPass"`
	DbUser  string `json:"dbUser"`
	DbPass  string `json:"dbPass"`
}

// SPECIFIC (SERVER)

func mapKeeperResponse(r node.KeeperOneResponse) NodeConfig {
	return NodeConfig{
		Host:       *r.DiscoveredHost,
		KeeperPort: r.DiscoveredKeeperPort,
		DbPort:     r.DiscoveredDbPort,
	}
}

func mapKeeperResponseList(keeperNodes []node.KeeperOneResponse) []NodeConfig {
	nodes := make([]NodeConfig, 0, len(keeperNodes))
	for _, item := range keeperNodes {
		nodes = append(nodes, mapKeeperResponse(item))
	}
	return nodes
}

func mapKeeperResponseMap(keeperNodes map[string]node.KeeperOneResponse) []NodeConfig {
	nodes := make([]NodeConfig, 0, len(keeperNodes))
	for _, item := range keeperNodes {
		nodes = append(nodes, mapKeeperResponse(item))
	}
	return nodes
}
