package cluster

import (
	"ivory/features"
	"ivory/features/cert"
	"ivory/features/node"
	"ivory/plugins/database"
	"ivory/plugins/keeper"

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
	Keeper   keeper.Plugin   `json:"keeper"`
	Database database.Plugin `json:"database"`
}

type Vaults struct {
	KeeperId   *uuid.UUID `json:"keeperId"`
	DatabaseId *uuid.UUID `json:"databaseId"`
	SshKeyId   *uuid.UUID `json:"sshKeyId"`
}

type Node struct {
	Config   NodeConfig          `json:"config"`
	Keeper   node.KeeperResponse `json:"keeper"`
	Warnings []string            `json:"warnings"`
}

type Overview struct {
	Nodes    map[string]Node    `json:"nodes"`
	Features []features.Feature `json:"features"`
}

type CreateAutoRequest struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Options
}

type ImageConfig struct {
	Cluster string `json:"cluster"`
	Dcs     string `json:"dcs"`
	DbPass  string `json:"dbPass"`
	DbUser  string `json:"dbUser"`
	SshPass string `json:"sshPass"`
	SshUser string `json:"sshUser"`
}

type DeployRequest struct {
	Uri            string            `json:"uri"`
	NodeRawOptions map[string]string `json:"nodeRawOptions"`
	NodeConfig     []NodeConfig      `json:"nodeConfig"`
	CommonConfig   ImageConfig       `json:"commonConfig"`
	ClusterOptions Options           `json:"clusterOptions"`
}

// SPECIFIC (SERVER)
