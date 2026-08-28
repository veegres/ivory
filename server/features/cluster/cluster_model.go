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
	// Name is the node's own name, unique within the cluster and independent
	// of its host: it is the deployment's identity ({{name}}, --name) and the
	// name the platform addresses the deployment by, so it is required on
	// every write (see Service.validateNodeNames). Host stays the connection
	// identity that keeper-reported nodes are matched against, so the two are
	// deliberately separate.
	Name       string `json:"name" form:"name"`
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
	Parallel       bool         `json:"parallel"`
	Nodes          []DeployNode `json:"nodes"`
	CommonConfig   CommonConfig `json:"commonConfig"`
	ClusterOptions Options      `json:"clusterOptions"`
}

// DeployNode pairs one node with the command that deploys it, for the length
// of one request only - the command is never persisted on the cluster.
type DeployNode struct {
	NodeConfig
	Command    string `json:"command"`
	PostScript string `json:"postScript"`
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
	// NOTE: a keeper that names its members by endpoint rather than by name
	// reports no name at all, and the node falls back to its host - the same
	// default withNodeNames gives clusters stored before names existed
	name := *r.DiscoveredHost
	if r.DiscoveredName != nil && *r.DiscoveredName != "" {
		name = *r.DiscoveredName
	}
	return NodeConfig{
		Name:       name,
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

// withNodeNames defaults every node's name to its host, so clusters stored
// before names existed keep rendering and deploying: host used to serve as
// the deployment's name, and that is exactly what the default preserves.
func (r Response) withNodeNames() Response {
	for i := range r.Nodes {
		if r.Nodes[i].Name == "" {
			r.Nodes[i].Name = r.Nodes[i].Host
		}
	}
	return r
}

// mapDeployNodeConfigs keeps only the connection details on the cluster: the
// commands themselves live in templates, never on the stored cluster.
func mapDeployNodeConfigs(nodes []DeployNode) []NodeConfig {
	configs := make([]NodeConfig, 0, len(nodes))
	for _, n := range nodes {
		configs = append(configs, n.NodeConfig)
	}
	return configs
}
