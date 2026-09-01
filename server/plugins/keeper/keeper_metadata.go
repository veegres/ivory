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
	// HasLeader reports whether the engine elects a single primary at all, so
	// a consumer can tell a missing leader from an engine that has none.
	HasLeader() bool
	// DefaultTemplates returns the ready-to-copy deployments Ivory ships for
	// this engine, one set per platform it supports. It lives on the plugin
	// rather than in a central catalog for the same reason
	// database.Adapter.SystemRequests does: a new plugin cannot be registered
	// without answering it, so its deployments can never be forgotten.
	DefaultTemplates() []DeploymentTemplate
}

// DeploymentTemplate is one shipped deployment: an ordered list of commands,
// one per node, written for a specific platform. It says nothing about where
// the nodes run - the host is supplied at deploy time.
type DeploymentTemplate struct {
	Platform    platform.PluginType
	Name        string
	Description string
	Defaults    DeploymentTemplateDefaults
	Commands    []DeploymentCommand
}

// DeploymentCommand is one node's deployment. It has no identity beyond its
// position: which node it lands on is chosen at deploy time.
type DeploymentCommand struct {
	Command string `json:"command"`
	// PostScripts run inside the deployment once the batch is up, in order,
	// each as its own execution. It is a list rather than one script because
	// nothing may assume a shell: the images these run in are increasingly
	// distroless (etcd's holds only etcd, etcdctl and etcdutl), so "&&" is not
	// available to chain with. One command per step also means each argument is
	// interpolated on its own, with no second parse to escape for.
	PostScripts []string                  `json:"postScripts"`
	Defaults    DeploymentCommandDefaults `json:"defaults"`
}

// DeploymentCommandDefaults is what the deploy form fills this command's node
// card in with. A command and the endpoints it answers on are one fact rather
// than two: a single-host template writes a distinct peer port into each of its
// commands, and only that command knows which client port has to match it. The
// name is here for the same reason - a member list naming etcd2 is only correct
// if the node running that command is called etcd2. SshPort is here for the
// same reason as the other two ports: a single-host template's nodes are
// reached through the same VM but can still be forwarded on distinct ports, and
// only the command that pairs with a given peer port knows which one.
//
// Host is empty in a multi-host template - each of its nodes is a distinct
// real machine a template can never know ahead of time. A single-host template
// names one, because all of its commands land on the same machine by
// definition: it defaults to "localhost", the deploy form's stand-in for
// wherever the operator runs it, and stays a plain string the user can
// overwrite like any other default.
type DeploymentCommandDefaults struct {
	Name       string `json:"name"`
	Host       string `json:"host"`
	SshPort    int    `json:"sshPort"`
	KeeperPort int    `json:"keeperPort"`
	DbPort     int    `json:"dbPort"`
}

// DeploymentTemplateDefaults is what a template fills the deploy screen's
// credential fields in with. It sits on the template rather than on a command
// because credentials are one answer for the whole cluster: every node of a
// deployment authenticates the same way, so three commands each naming their
// own username could only ever disagree.
//
// It states usernames and nothing else. Whether a deployment has credentials at
// all is the user's answer on the deploy screen, not the template's: a username
// here means the deployment ends up with that account - spilo names its
// superuser postgres, etcd can only enable auth through root - so the screen
// opens on it, already filled in. Where it names none, the screen opens with
// that credential switched off and the user says what they want instead.
// Passwords are never here: a template is stored, read and copied, and a secret
// in one would be a secret in all three.
type DeploymentTemplateDefaults struct {
	KeeperUser string `json:"keeperUser"`
	DbUser     string `json:"dbUser"`
	// Dcs is the coordination store {{dcs}} resolves to. It is a default rather
	// than literal text in the command for the same reason the usernames are:
	// one cluster coordinates through one store, so three commands each naming
	// their own address could only ever disagree, and the address is the one
	// thing about an external DCS that changes between deployments of the same
	// template. A template whose commands never name {{dcs}} leaves it empty.
	Dcs string `json:"dcs"`
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
	VarDcs        Var = "{{dcs}}"        // address of the coordination store the cluster runs against, one value for the whole deployment
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
// vault and the coordination store the cluster runs against. Anything else an
// engine needs (a peer port, a member list, which node is the leader) is
// written literally into the command, because only the user knows it and it is
// plain text they can read.
var Vars = []Var{
	VarCluster, VarDcs, VarName, VarHost, VarSshPort, VarKeeperPort, VarDbPort,
	VarKeeperUser, VarKeeperPass, VarDbUser, VarDbPass,
}

// Values is one command's complete interpolation scope. There is no shared
// map: everything a command can reference belongs to that command's own
// deployment, so one node's values can never reach another node's command.
type Values struct {
	Cluster    string
	Dcs        string
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
		VarDcs:        v.Dcs,
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

// placeholderPattern matches a variable and the near misses of one: a bare
// token in braces, with or without surrounding space. {{ host }} and {{db-port}}
// used to slip past a stricter pattern and reach the shell verbatim, which is
// the one thing the closed vocabulary exists to prevent - they match here and
// are reported against Vars like any other typo. Go template syntax
// ({{json .}}, {{.Name}}) still does not match: docker interprets it itself and
// a command may legitimately carry it.
var placeholderPattern = regexp.MustCompile(`{{\s*[\w-]+\s*}}`)

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
