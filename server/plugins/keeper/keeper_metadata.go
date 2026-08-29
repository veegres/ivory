package keeper

import (
	"ivory/core/config"
	"ivory/plugins/platform"
	"regexp"
	"slices"
)

// Metadata is implemented by keeper plugins to describe themselves without
// touching a running keeper. It deliberately says nothing about how the engine
// is deployed: deployment is a command the user writes, stored as a template.
type Metadata interface {
	// SupportedFeatures reports, for every keeper-related config.Feature this
	// plugin knows about, whether it supports it. A feature absent from the
	// map is not a keeper capability at all and is left unrestricted.
	SupportedFeatures() map[config.Feature]bool
	Requirements() Requirements
	// DefaultTemplates returns the ready-to-copy deployments Ivory ships for
	// this engine, one set per platform it supports. It lives on the plugin
	// rather than in a central catalog for the same reason
	// database.Adapter.SystemRequests does: a new plugin cannot be registered
	// without answering it, so its deployments can never be forgotten.
	DefaultTemplates() []DeploymentTemplate
}

// DeploymentTemplate is one shipped deployment: an ordered list of commands,
// one per node, written for a specific platform. It carries nothing about
// where the nodes run - host, ports and names are supplied at deploy time.
type DeploymentTemplate struct {
	Platform    platform.PluginType
	Name        string
	Description string
	Commands    []DeploymentCommand
}

// DeploymentCommand is one node's deployment. It has no identity beyond its
// position: which node it lands on is chosen at deploy time.
type DeploymentCommand struct {
	Command    string `json:"command"`
	PostScript string `json:"postScript"`
}

// Requirements is what Ivory must know to talk to the engine, not how to
// deploy it.
type Requirements struct {
	// DbPort is the database endpoint's default port.
	DbPort int
	// KeeperPort is the keeper endpoint's default port, declared even when it
	// is the database port.
	KeeperPort int
	// KeeperCredentials reports whether the keeper endpoint consumes credentials.
	KeeperCredentials bool
	// KeeperUser is the username the keeper endpoint locks itself to; empty
	// leaves the choice to the user.
	KeeperUser string
	// DbCredentials reports whether the database endpoint consumes credentials.
	DbCredentials bool
	// DbUser is the username the engine locks itself to; empty leaves the
	// choice to the user.
	DbUser string
}

// Var is a {{placeholder}} variable, in its interpolated form. The set is
// closed and hard-coded: a template referencing anything outside Vars is
// rejected rather than treated as a new variable, so every command in the
// system draws from one reviewable vocabulary. Adding one means adding a Var
// here and a field to Values (and, when the user supplies it, to Defaults) -
// keeper_metadata_test.go asserts both directions stay exhaustive.
type Var string

const (
	VarCluster    Var = "{{cluster}}"    // cluster name
	VarName       Var = "{{name}}"       // node name, unique within the cluster; the deployment's own name
	VarHost       Var = "{{host}}"       // node host
	VarSshPort    Var = "{{sshPort}}"    // node ssh port
	VarKeeperPort Var = "{{keeperPort}}" // keeper endpoint port
	VarDbPort     Var = "{{dbPort}}"     // database endpoint port
	VarKeeperUser Var = "{{keeperUser}}" // keeper credentials username, resolved from the keeper vault at execution time
	VarKeeperPass Var = "{{keeperPass}}" // keeper credentials password, resolved from the keeper vault at execution time
	VarDbUser     Var = "{{dbUser}}"     // database credentials username, resolved from the database vault at execution time
	VarDbPass     Var = "{{dbPass}}"     // database credentials password, resolved from the database vault at execution time
)

// Vars is the complete closed list of variables a command may reference. It is
// exactly the node's own identity and endpoints - what Ivory needs to register
// and reach the cluster afterwards - plus the credentials it resolves from the
// vault. Anything else an engine needs (a peer port, a member list, a
// coordinator address, which node is the leader) is written literally into the
// command, because only the user knows it and it is plain text they can read.
var Vars = []Var{
	VarCluster, VarName, VarHost, VarSshPort, VarKeeperPort, VarDbPort,
	VarKeeperUser, VarKeeperPass, VarDbUser, VarDbPass,
}

// Values is one command's complete interpolation scope. There is no shared
// map: everything a command can reference belongs to that command's own
// deployment, so one node's values can never reach another node's command.
type Values struct {
	Cluster    string
	Name       string
	Host       string
	SshPort    string
	KeeperPort string
	DbPort     string
	KeeperUser string
	KeeperPass string
	DbUser     string
	DbPass     string
}

func (v Values) lookup() map[Var]string {
	return map[Var]string{
		VarCluster:    v.Cluster,
		VarName:       v.Name,
		VarHost:       v.Host,
		VarSshPort:    v.SshPort,
		VarKeeperPort: v.KeeperPort,
		VarDbPort:     v.DbPort,
		VarKeeperUser: v.KeeperUser,
		VarKeeperPass: v.KeeperPass,
		VarDbUser:     v.DbUser,
		VarDbPass:     v.DbPass,
	}
}

var placeholderPattern = regexp.MustCompile(`{{\w+}}`)

// Interpolate substitutes {{placeholder}} variables with the command's own
// values. Missing or empty values leave the variable in place so it can be
// reported by UnresolvedPlaceholders.
func Interpolate(template string, values Values) string {
	lookup := values.lookup()
	return placeholderPattern.ReplaceAllStringFunc(template, func(match string) string {
		if val, ok := lookup[Var(match)]; ok && val != "" {
			return val
		}
		return match
	})
}

// UnresolvedPlaceholders reports {{placeholder}} templates left after
// interpolation, which means the deployment misses values its command
// actually references.
func UnresolvedPlaceholders(text string) []string {
	return placeholders(text, func(string) bool { return true })
}

// UnknownPlaceholders reports every {{placeholder}} outside Vars. Templates
// are rejected on save and commands on deploy when this is non-empty, which is
// what keeps the vocabulary closed - a typo is an error, never a new variable.
func UnknownPlaceholders(text string) []string {
	return placeholders(text, func(match string) bool {
		return !slices.Contains(Vars, Var(match))
	})
}

func placeholders(text string, keep func(match string) bool) []string {
	var found []string
	seen := map[string]bool{}
	for _, match := range placeholderPattern.FindAllString(text, -1) {
		if !seen[match] && keep(match) {
			seen[match] = true
			found = append(found, match)
		}
	}
	return found
}
