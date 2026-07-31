package keeper

import (
	"ivory/core/config"
	"regexp"
	"slices"
)

// Metadata is implemented by keeper plugins to describe themselves without
// touching a running keeper: which features they support and how their image
// should be deployed. It is intentionally separate from Adapter so a plugin
// without a management API yet (plain postgres) can still declare metadata.
type Metadata interface {
	// SupportedFeatures reports, for every keeper-related config.Feature this
	// plugin knows about, whether it supports it. A feature absent from the
	// map is not a keeper capability at all and is left unrestricted.
	SupportedFeatures() map[config.Feature]bool
	DeploymentSpec() DeploymentSpec
}

// DeploymentSpec declares, in platform-agnostic terms, what a keeper image
// needs at deploy time. Values may contain {{placeholder}} templates that are
// interpolated right before deployment; quoting inside values is kept as-is
// because the rendered text stays user-editable.
type DeploymentSpec struct {
	DefaultImage string
	Env          []EnvVar
	Ports        []string
	Volumes      []VolumeSpec
	// Defaults declares which built-in variables the deployment uses, with
	// their default values:
	//   - VarDbPort: the database endpoint default port, always required
	//   - VarKeeperPort: present = the plugin has a separate keeper endpoint
	//     (the value is its default port); absent = the keeper endpoint is
	//     the database itself
	//   - VarDbUser: present = the deployment consumes database credentials
	//     ({{dbUser}}/{{dbPass}}, resolved from the vault); a non-empty value
	//     is the engine-required username (forms prefill and lock it)
	Defaults map[Var]string
	Fields   []FieldSpec
	// PostScript, if set, is executed once inside an already-deployed
	// container after the whole deployment succeeds (e.g. enabling
	// authentication). Unlike EntryScript, this doesn't replace the
	// container's startup command - it's exec'd into the already-running
	// container afterward - so a plugin needing several steps joins them
	// itself (e.g. with "&&") into one script, the same way EntryScript is
	// one script rather than a list of commands. It uses the regular
	// {{placeholder}} vocabulary, including credentials, interpolated the
	// same way as EntryScript.
	PostScript string
	// EntryScript, if set, replaces the image's own default startup command
	// for every node, including the deploy's first (always the leader; see
	// KeeperDeployPlan) - unless EntryScriptReplicasOnly is set, see its doc.
	// Unlike PostScript, this isn't exec'd into an already-running container
	// after the fact - it IS the container's own command, so it must itself
	// end by starting the plugin's actual server process (e.g. native
	// postgres rebases via pg_basebackup before postgres itself ever
	// starts, since streaming replication needs a real base backup - a
	// fresh initdb can't be turned into a replica by just pointing it at a
	// leader). It uses the regular {{placeholder}} vocabulary, including
	// VarLeaderHost, interpolated by KeeperDeployPlan the same way as
	// Options.
	EntryScript string
	// EntryScriptReplicasOnly, if true, skips EntryScript for the deploy's
	// first node (the leader) instead of KeeperDeployPlan's default of
	// applying it uniformly to every node. Set this when the first node
	// genuinely is special and must get the image's plain startup command -
	// e.g. postgres/redis, where only a replica needs to rebase/attach
	// itself to the leader, and the leader would break if it ran that same
	// script. Leave it false (the default) for keepers with no
	// single-leader distinction at all, where every node - the first one
	// too - needs the identical startup-time setup (e.g. ClickHouse's
	// EntryScript generates the same cluster config file on every node).
	EntryScriptReplicasOnly bool
}

// Var is a built-in {{placeholder}} variable that Ivory provides to every
// deployment, in its interpolated form; plugin-declared FieldSpec names
// extend this set. Templates are built from these values, so a misspelled
// variable cannot exist in a template. It is a distinct type from string so
// a Var constant can't be passed where any plain string would do (and vice
// versa) without an explicit conversion - the interpolation machinery
// (Interpolate, UnresolvedPlaceholders, and the values maps built around
// them) still works in plain strings, since a rendered template mixes Var
// values with arbitrary literal text.
type Var string

// NOTE: every variable lives in this single block — built-ins provided by
// Ivory first, then the plugin field variables (declared by plugins as
// FieldSpec names), so names are all visible at a glance and never collide.
const (
	VarCluster    Var = "{{cluster}}"    // cluster name
	VarHost       Var = "{{host}}"       // node host
	VarKeeperPort Var = "{{keeperPort}}" // keeper endpoint port (the database port when there is no separate keeper endpoint)
	VarDbPort     Var = "{{dbPort}}"     // database endpoint port
	VarDbUser     Var = "{{dbUser}}"     // database endpoint credentials username, resolved from the vault
	VarDbPass     Var = "{{dbPass}}"     // database endpoint credentials password, resolved from the vault
	VarLeaderHost Var = "{{leaderHost}}" // the deploy request's first node's host (see keeper.DeploymentSpec.EntryScript and KeeperDeployPlan) - named after Role's Leader/Replica vocabulary, not the underlying engine's own term (postgres says "primary", redis says "master")
	VarIndex      Var = "{{index}}"      // this node's 1-based position in the deploy request's node list (see KeeperDeployPlanNode.Index) - the only built-in that isn't derived from the node's own connection details, for engines whose ensemble config needs a genuinely unique small integer per member rather than a hostname (e.g. ZooKeeper's myid/server.N)

	VarDcs          Var = "{{dcs}}"          // external coordinator address a plugin points at instead of deploying its own: patroni's etcd/zookeeper DCS, clickhouse's zookeeper/clickhouse-keeper ensemble
	VarPeerPort     Var = "{{peerPort}}"     // a port only cluster members dial among themselves, never Ivory itself - only etcd declares it today (its raft peer listener, unique per node in single-host mode), but any plugin needing a cluster-internal-only port (as opposed to VarDbPort/VarKeeperPort, which Ivory's own Adapter calls) can reuse this same field rather than inventing its own
	VarClusterHosts Var = "{{clusterHosts}}" // every node in this deploy, built via a plugin's own FieldSpec.Template (etcd: member list name=http://host:peerPort,...; clickhouse: <replica> entries for its own <remote_servers> shard) - distinct from {{dcs}}, the external coordinator this deploy does NOT include
)

// Vars lists the built-in variables Ivory provides to every deployment;
// plugin field variables count as known only when the spec declares them.
var Vars = []Var{VarCluster, VarHost, VarKeeperPort, VarDbPort, VarDbUser, VarDbPass, VarLeaderHost, VarIndex}

var placeholderPattern = regexp.MustCompile(`{{\w+}}`)

// Interpolate substitutes {{placeholder}} variables with values; the map is
// keyed by the variable's interpolated form (Var values and field names).
// Missing or empty values leave the variable in place so it can be reported
// by UnresolvedPlaceholders.
func Interpolate(template string, values map[string]string) string {
	return placeholderPattern.ReplaceAllStringFunc(template, func(match string) string {
		if val, ok := values[match]; ok && val != "" {
			return val
		}
		return match
	})
}

// UnresolvedPlaceholders reports {{placeholder}} templates left after
// interpolation, which means the deploy request misses values the rendered
// options actually reference.
func UnresolvedPlaceholders(text string) []string {
	var unresolved []string
	seen := map[string]bool{}
	for _, match := range placeholderPattern.FindAllString(text, -1) {
		if !seen[match] {
			seen[match] = true
			unresolved = append(unresolved, match)
		}
	}
	return unresolved
}

// UnknownVariables reports every {{placeholder}} referenced by the spec's
// templates that is neither a built-in Var nor a declared field name; plugin
// metadata tests use it to catch misspelled variables.
func (s DeploymentSpec) UnknownVariables() []string {
	known := make(map[string]bool, len(Vars)+len(s.Fields))
	for _, v := range Vars {
		known[string(v)] = true
	}
	for _, f := range s.Fields {
		known[string(f.Name)] = true
	}

	texts := []string{s.DefaultImage, s.EntryScript, s.PostScript}
	texts = append(texts, s.Ports...)
	for _, e := range s.Env {
		texts = append(texts, e.Value)
	}
	for _, v := range s.Volumes {
		texts = append(texts, v.HostPath, v.ContainerPath)
	}
	for _, f := range s.Fields {
		texts = append(texts, f.Template, f.Default)
	}

	unknown := make([]string, 0)
	for _, text := range texts {
		for _, match := range placeholderPattern.FindAllString(text, -1) {
			if !known[match] && !slices.Contains(unknown, match) {
				unknown = append(unknown, match)
			}
		}
	}
	return unknown
}

// FieldType tells the forms how to render a field and the planner how to
// treat its value.
type FieldType string

const (
	FieldText FieldType = "text"
	// FieldPort values are per-node: the field value is the base port and
	// nodes get base+index when they all share one host (single-host mode),
	// so their listeners don't collide.
	FieldPort FieldType = "port"
)

// FieldSpec declares an image-level interpolation variable shown as an
// editable form field; Name is its interpolated form, e.g. "{{peerPort}}",
// declared as a plugin constant next to the built-in Vars. A user-provided
// value always wins; otherwise the value is derived from the node list
// (Template interpolated once per node, entries joined with Separator — e.g.
// VarClusterHosts' member list) or falls back to Default (e.g.
// patroni's external DCS address, empty until typed). Fields are resolved in
// declaration order, so a Template may reference earlier fields.
type FieldSpec struct {
	Name      Var
	Label     string
	Example   string
	Type      FieldType
	Default   string
	Template  string
	Separator string
}

// EnvVar is a single environment variable to set on the deployed container.
type EnvVar struct {
	Name  string
	Value string
}

// VolumeSpec is a single host-path-to-container-path volume mount required
// by the deployed container.
type VolumeSpec struct {
	HostPath      string
	ContainerPath string
}
