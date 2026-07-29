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
	Plugins    Plugins    `json:"plugins"`
	Tls        Tls        `json:"tls"`
	Certs      cert.Certs `json:"certs"`
	Vaults     Vaults     `json:"vaults"`
	Tags       []string   `json:"tags"`
	SingleHost bool       `json:"singleHost"`
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
	Nodes    map[string]Node         `json:"nodes"`
	Features map[config.Feature]bool `json:"features"`
}

type CreateAutoRequest struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Options
}

// DeployRequest describes a deployment intent: node ports, the image, aux
// ports, the DCS value and the per-node options are resolved server-side from
// the keeper plugin's DeploymentSpec unless explicitly provided. Values
// carries plugin-required inputs (e.g. "dcs" for patroni).
type DeployRequest struct {
	Parallel       bool              `json:"parallel"`
	SingleHost     bool              `json:"singleHost"`
	Image          string            `json:"image"`
	Nodes          []DeployNode      `json:"nodes"`
	Values         map[string]string `json:"values"`
	CommonConfig   CommonConfig      `json:"commonConfig"`
	ClusterOptions Options           `json:"clusterOptions"`
}

type DeployNode struct {
	NodeConfig
	// Options overrides the rendered options template for this node.
	Options string `json:"options"`
}

type CommonConfig struct {
	Cluster string `json:"cluster"`
	SshUser string `json:"sshUser"`
	SshPass string `json:"sshPass"`
	DbUser  string `json:"dbUser"`
	DbPass  string `json:"dbPass"`
}

// SearchRequest narrows Search down to matching clusters. Tags is resolved
// to cluster names before hitting the repository; Keeper and Database are
// passed through as-is. A nil/empty field is skipped (no restriction).
type SearchRequest struct {
	Tags     []string
	Keeper   *node.KeeperPlugin
	Database *query.DbPlugin
}

// SearchCriteria narrows List down to matching clusters. A nil field is
// skipped; Names is nil-vs-empty sensitive so a resolved-but-empty tag
// search (no cluster has the requested tags) still returns no results
// instead of falling back to an unfiltered list.
type SearchCriteria struct {
	Names    []string
	Keeper   *node.KeeperPlugin
	Database *query.DbPlugin
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

func mapPlanNodeConfigs(planNodes []node.KeeperDeployPlanNode) []NodeConfig {
	nodes := make([]NodeConfig, 0, len(planNodes))
	for _, pn := range planNodes {
		sshPort, keeperPort, dbPort := pn.SshPort, pn.KeeperPort, pn.DbPort
		nodes = append(nodes, NodeConfig{
			Host:       pn.Host,
			SshPort:    &sshPort,
			KeeperPort: &keeperPort,
			DbPort:     &dbPort,
		})
	}
	return nodes
}
