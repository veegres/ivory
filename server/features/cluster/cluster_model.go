package cluster

import (
	"fmt"
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

// DeployRequest is one whole cluster's deployment: every node states its own
// name, host, ports and command, and nothing is resolved server-side from the
// keeper plugin.
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
	Command     string   `json:"command"`
	PostScripts []string `json:"postScripts"`
}

type CommonConfig struct {
	Cluster    string `json:"cluster"`
	SshUser    string `json:"sshUser"`
	SshPass    string `json:"sshPass"`
	KeeperUser string `json:"keeperUser"`
	KeeperPass string `json:"keeperPass"`
	DbUser     string `json:"dbUser"`
	DbPass     string `json:"dbPass"`
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

// describesNode reports whether a keeper response describes a node at all. One
// that states no host describes something else - native postgres reporting a
// standby's sync state from the primary, an etcd member with no client url yet -
// and a node config with no host is a node nothing can reach, so Detect and Fix
// leave those out of the cluster they write rather than storing one.
func describesNode(r node.KeeperOneResponse) bool {
	return r.DiscoveredHost != nil
}

func mapKeeperResponseList(keeperNodes []node.KeeperOneResponse) []NodeConfig {
	nodes := make([]NodeConfig, 0, len(keeperNodes))
	for _, item := range keeperNodes {
		if !describesNode(item) {
			continue
		}
		nodes = append(nodes, mapKeeperResponse(item))
	}
	return nodes
}

func mapKeeperResponseMap(keeperNodes map[string]node.KeeperOneResponse) []NodeConfig {
	nodes := make([]NodeConfig, 0, len(keeperNodes))
	for _, item := range keeperNodes {
		if !describesNode(item) {
			continue
		}
		nodes = append(nodes, mapKeeperResponse(item))
	}
	return nodes
}

// withNodeNames defaults every node's name to its host, so clusters stored
// before names existed keep rendering and deploying: host used to serve as
// the deployment's name, and that is exactly what the default preserves.
//
// Repeats are suffixed because several nodes on one VM share a host, and a
// cluster whose names collide is rejected by validateNodeNames - defaulting
// them all to the same host is what made an existing single-host cluster
// impossible to update or edit in the list at all.
func (r Response) withNodeNames() Response {
	taken := make(map[string]bool, len(r.Nodes))
	for _, n := range r.Nodes {
		if n.Name != "" {
			taken[n.Name] = true
		}
	}
	for i := range r.Nodes {
		if r.Nodes[i].Name == "" {
			r.Nodes[i].Name = uniqueNodeName(r.Nodes[i].Host, taken)
		}
	}
	return r
}

// uniqueNodeName returns host, or the first free host-2, host-3, ... when it
// is already spoken for, and records the result as taken.
func uniqueNodeName(host string, taken map[string]bool) string {
	name := host
	for i := 2; taken[name]; i++ {
		name = fmt.Sprintf("%s-%d", host, i)
	}
	taken[name] = true
	return name
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
