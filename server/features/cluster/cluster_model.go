package cluster

import (
	"fmt"
	"ivory/core/config"
	"ivory/core/service/cert"
	"ivory/core/utils"
	"ivory/features/node"
	"ivory/features/query"
	"maps"
	"slices"

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
	Name       string `json:"name" form:"name" binding:"required,notblank"`
	Host       string `json:"host" form:"host" binding:"required,notblank"`
	SshPort    *int   `json:"sshPort" form:"sshPort" binding:"omitempty,min=1,max=65535"`
	KeeperPort *int   `json:"keeperPort" form:"keeperPort" binding:"omitempty,min=1,max=65535"`
	DbPort     *int   `json:"dbPort" form:"dbPort" binding:"omitempty,min=1,max=65535"`
}

type Options struct {
	Plugins Plugins    `json:"plugins" binding:"required"`
	Tls     Tls        `json:"tls"`
	Certs   cert.Certs `json:"certs"`
	Vaults  Vaults     `json:"vaults"`
	Tags    []string   `json:"tags"`
}

type Request struct {
	Name  string       `json:"name" binding:"required,notblank"`
	Nodes []NodeConfig `json:"nodes" binding:"required,min=1,dive"`
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
	Keeper   node.KeeperPlugin `json:"keeper" binding:"required,valid"`
	Database query.DbPlugin    `json:"database" binding:"required,valid"`
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
	Name string `json:"name" binding:"required,notblank"`
	Host string `json:"host" binding:"required,notblank"`
	Port int    `json:"port" binding:"required,min=1,max=65535"`
	Options
}

// DeployRequest is one whole cluster's deployment: every node states its own
// name, host, ports and command, and nothing is resolved server-side from the
// keeper plugin.
type DeployRequest struct {
	Parallel       bool                `json:"parallel"`
	Nodes          []DeployNode        `json:"nodes" binding:"required,min=1,dive"`
	Platform       node.PlatformPlugin `json:"platform" binding:"required,valid"`
	CommonConfig   CommonConfig        `json:"commonConfig"`
	ClusterOptions Options             `json:"clusterOptions"`
}

// DeployNode pairs one node with the command that deploys it, for the length
// of one request only - the command is never persisted on the cluster.
type DeployNode struct {
	NodeConfig
	Command     string   `json:"command" binding:"required,notblank"`
	PostScripts []string `json:"postScripts"`
}

type CommonConfig struct {
	Cluster string `json:"cluster" binding:"required,notblank"`
	// Dcs is the address of the coordination store the whole deployment runs
	// against, resolved into {{dcs}}. It is one answer for the cluster rather
	// than a per-node field, and it is never stored on the cluster: like the
	// commands themselves it belongs to the deploy request only.
	Dcs        string `json:"dcs"`
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
	Tags     []string           `json:"tags,omitempty"`
	Keeper   *node.KeeperPlugin `json:"keeper,omitempty" binding:"omitempty,valid"`
	Database *query.DbPlugin    `json:"database,omitempty" binding:"omitempty,valid"`
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

// RouterNameRequest and RouterDeployRequest carry what the http route takes
// from the url path. They exist because a tool's arguments are one typed
// object and there is no path to read a name out of - not to restate a model.
type RouterNameRequest struct {
	Name string `json:"name" binding:"required,notblank" jsonschema:"name of the cluster as registered in Ivory"`
}

// RouterOverviewRequest narrows the overview to what one keeper endpoint
// reports when Host and Port are given, the way the http route's query does.
type RouterOverviewRequest struct {
	Name string `json:"name" binding:"required,notblank" jsonschema:"name of the cluster as registered in Ivory"`
	Host string `json:"host,omitempty" binding:"required_with=Port" jsonschema:"ask only this keeper host instead of every configured node"`
	Port int    `json:"port,omitempty" binding:"required_with=Host,omitempty,min=1,max=65535" jsonschema:"keeper port of that host"`
}

// NodeConnections is the input every node-level operation takes, resolved
// from one node's stored configuration. A half is nil when the node lacks
// what it needs: a keeper port, a database port, or an ssh port and ssh vault.
type NodeConnections struct {
	Keeper   *node.KeeperOneRequest        `json:"keeper"`
	Platform *node.PlatformVaultConnection `json:"platform"`
	Database *query.Connection             `json:"database"`
}

// NodeCheck is how far Ivory gets reaching one node, through each of its
// connections. A half is missing where NodeConnections has none.
type NodeCheck struct {
	Keeper   []utils.CheckStep `json:"keeper,omitempty"`
	Platform []utils.CheckStep `json:"platform,omitempty"`
	Database []utils.CheckStep `json:"database,omitempty"`
}

// RouterDeployRequest is DeployRequest with CommonConfig's inline credentials left
// out: a deploy names vault entries instead, so no password passes through the
// agent's context. Deploy refuses a vault id beside an inline pair anyway.
type RouterDeployRequest struct {
	Cluster  string              `json:"cluster" binding:"required,notblank" jsonschema:"name to register the new cluster under"`
	Dcs      string              `json:"dcs,omitempty" jsonschema:"address of the coordination store the whole deployment registers with, substituted for {{dcs}}"`
	Parallel bool                `json:"parallel,omitempty" jsonschema:"deploy every node at once instead of in order"`
	Platform node.PlatformPlugin `json:"platform" binding:"required,valid" jsonschema:"deployment target, e.g. docker"`
	Nodes    []DeployNode        `json:"nodes" binding:"required,min=1,dive" jsonschema:"one entry per node: its name, host, ssh/keeper/database ports, the command that deploys it and that command's post-scripts"`
	Options  Options             `json:"options" binding:"required" jsonschema:"keeper and database plugins, TLS flags, certificates, tags, and the vault ids holding the ssh key and the keeper/database credentials"`
}

// DeployResponse reports each node's steps, keyed by node name. Getting one
// back at all means the cluster was registered with every node, whatever
// their steps say.
type DeployResponse struct {
	Logs  []string                    `json:"logs"`
	Nodes map[string]node.DeploySteps `json:"nodes"`
}

func (in RouterDeployRequest) toDeployRequest() DeployRequest {
	return DeployRequest{
		Parallel:       in.Parallel,
		Nodes:          in.Nodes,
		Platform:       in.Platform,
		CommonConfig:   CommonConfig{Cluster: in.Cluster, Dcs: in.Dcs},
		ClusterOptions: in.Options,
	}
}

// SPECIFIC (SERVER)

func (r DeployResponse) done() bool {
	for _, steps := range r.Nodes {
		if !steps.Done() {
			return false
		}
	}
	return true
}

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

// mapUniqueKeeperResponse suffixes a discovered name that is already spoken
// for, so Detect and Fix cannot store a cluster whose nodes share a name: an
// engine that identifies its members by endpoint reports no name at all, so
// every node of a single-host cluster falls back to the same host, and
// validateNodeNames would then reject every later write to that cluster.
func mapUniqueKeeperResponse(r node.KeeperOneResponse, taken map[string]bool) NodeConfig {
	config := mapKeeperResponse(r)
	config.Name = uniqueNodeName(config.Name, taken)
	return config
}

func mapKeeperResponseList(keeperNodes []node.KeeperOneResponse) []NodeConfig {
	nodes := make([]NodeConfig, 0, len(keeperNodes))
	taken := make(map[string]bool, len(keeperNodes))
	for _, item := range keeperNodes {
		if !describesNode(item) {
			continue
		}
		nodes = append(nodes, mapUniqueKeeperResponse(item, taken))
	}
	return nodes
}

func mapKeeperResponseMap(keeperNodes map[string]node.KeeperOneResponse) []NodeConfig {
	nodes := make([]NodeConfig, 0, len(keeperNodes))
	taken := make(map[string]bool, len(keeperNodes))
	// NOTE: map iteration order would otherwise decide which node keeps the
	// bare name and which one is suffixed
	for _, key := range slices.Sorted(maps.Keys(keeperNodes)) {
		item := keeperNodes[key]
		if !describesNode(item) {
			continue
		}
		nodes = append(nodes, mapUniqueKeeperResponse(item, taken))
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

// uniqueNodeName returns base, or the first free base-2, base-3, ... when it
// is already spoken for, and records the result as taken.
func uniqueNodeName(base string, taken map[string]bool) string {
	name := base
	for i := 2; taken[name]; i++ {
		name = fmt.Sprintf("%s-%d", base, i)
	}
	taken[name] = true
	return name
}

// nodeConnections applies the same rule discovery does: certificates travel
// only on the side whose tls flag is on.
func (o Options) nodeConnections(cluster string, n NodeConfig) NodeConnections {
	var connections NodeConnections
	if n.KeeperPort != nil {
		var certs *cert.Certs
		if o.Tls.Keeper {
			certs = &o.Certs
		}
		connections.Keeper = &node.KeeperOneRequest{
			KeeperConnection: node.KeeperConnection{Host: n.Host, Port: *n.KeeperPort},
			KeeperOptions: node.KeeperOptions{
				Cluster: cluster,
				Plugin:  o.Plugins.Keeper,
				VaultId: o.Vaults.KeeperId,
				Certs:   certs,
			},
		}
	}
	if n.SshPort != nil && o.Vaults.SshKeyId != nil {
		connections.Platform = &node.PlatformVaultConnection{
			Host:     n.Host,
			Port:     *n.SshPort,
			VaultId:  *o.Vaults.SshKeyId,
			Platform: node.DefaultPlatform,
		}
	}
	if n.DbPort != nil {
		var certs *cert.Certs
		if o.Tls.Database {
			certs = &o.Certs
		}
		connections.Database = &query.Connection{
			Db:      query.DbConfig{Plugin: o.Plugins.Database, Host: n.Host, Port: *n.DbPort},
			Certs:   certs,
			VaultId: o.Vaults.DatabaseId,
		}
	}
	return connections
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
