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
	// SupportedFeatures reports, for every keeper-related env.Feature this
	// plugin knows about, whether it supports it. A feature absent from the
	// map is not a keeper capability at all and is left unrestricted.
	SupportedFeatures() map[env.Feature]bool
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
	// PostDeploy declares commands executed inside one deployed container
	// after the whole deployment succeeded (e.g. enabling authentication).
	// Commands are interpolated with the same {{placeholder}} vocabulary,
	// including credentials.
	PostDeploy []string
}

// Var is a built-in {{placeholder}} variable that Ivory provides to every
// deployment, in its interpolated form; plugin-declared FieldSpec names
// extend this set. Templates are built from these values, so a misspelled
// variable cannot exist in a template.
type Var = string

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

	VarDcs            Var = "{{dcs}}"            // patroni: address of the external DCS it coordinates through
	VarPeerPort       Var = "{{peerPort}}"       // etcd: peer listener port, unique per node in single-host mode
	VarInitialCluster Var = "{{initialCluster}}" // etcd: member list (name=http://host:peerPort,...) derived from the node list
)

// Vars lists the built-in variables Ivory provides to every deployment;
// plugin field variables count as known only when the spec declares them.
var Vars = []Var{VarCluster, VarHost, VarKeeperPort, VarDbPort, VarDbUser, VarDbPass}

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
		known[v] = true
	}
	for _, f := range s.Fields {
		known[f.Name] = true
	}

	texts := []string{s.DefaultImage}
	texts = append(texts, s.Ports...)
	texts = append(texts, s.PostDeploy...)
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
// etcd's initial cluster member list) or falls back to Default (e.g.
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
