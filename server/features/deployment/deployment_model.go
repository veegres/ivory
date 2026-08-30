package deployment

import (
	"ivory/features/node"
	"ivory/plugins/keeper"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type KeeperPlugin = node.KeeperPlugin
type PlatformPlugin = node.PlatformPlugin

// defaultIdPrefix marks the synthetic ids of the shipped templates. They are
// computed rather than stored, so they need an id shape that can never collide
// with a stored template's uuid and that a write can recognise and refuse.
const defaultIdPrefix = "default:"

// defaultId builds a shipped template's id. It is deliberately a slug rather
// than the display name: an id ends up in a URL path, and a name like
// "Etcd (Multi Host)" would make a malformed one.
func defaultId(plugin KeeperPlugin, platform PlatformPlugin, name string) string {
	return defaultIdPrefix + string(platform) + ":" + string(plugin) + ":" + slugify(name)
}

func slugify(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case b.Len() > 0 && !strings.HasSuffix(b.String(), "-"):
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

// CreationType distinguishes the templates Ivory ships from the ones a user
// owns, using the same vocabulary as query.CreationType.
type CreationType string

const (
	// Manual templates are the user's own: editable, deletable, deployable.
	Manual CreationType = "manual"
	// System templates come from the keeper plugins. They are computed rather
	// than stored, so they are read-only by construction and are copied into a
	// manual template rather than deployed directly.
	System CreationType = "system"
)

// Template is a saved deployment: an ordered list of commands, one per node,
// and nothing about where they run. Host, ports and node names belong to a
// deploy request, which is what makes a template reusable.
type Template struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Keeper      KeeperPlugin      `json:"keeper"`
	Platform    PlatformPlugin    `json:"platform"`
	Defaults    TemplateDefaults  `json:"defaults"`
	Commands    []TemplateCommand `json:"commands"`
	Creation    CreationType      `json:"creation"`
	CreatedAt   int64             `json:"createdAt"`
}

// TemplateCommand is one node's deployment: it has no identity beyond its
// position, the node it lands on being chosen at deploy time. It is an alias
// rather than a copy so a template shipped by a keeper plugin needs no
// per-command mapping.
type TemplateCommand = keeper.DeploymentCommand

// CommandDefaults is what one command fills its node card in with at deploy
// time. Aliased for the same reason TemplateCommand is: a template shipped by a
// keeper plugin needs no per-command mapping.
type CommandDefaults = keeper.DeploymentCommandDefaults

// TemplateDefaults is what the whole template fills the deploy screen's
// credential fields in with. It replaced keeper.Requirements, which asked the
// engine what a deployment consumes: what a deployment ends up with is decided
// by the command that creates it, so the template is the only thing that can
// answer, and a template the user wrote answers for itself.
type TemplateDefaults = keeper.DeploymentTemplateDefaults

type TemplateRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Keeper      KeeperPlugin      `json:"keeper"`
	Platform    PlatformPlugin    `json:"platform"`
	Defaults    TemplateDefaults  `json:"defaults"`
	Commands    []TemplateCommand `json:"commands"`
}

type ListRequest struct {
	Keeper   *KeeperPlugin   `json:"keeper" form:"keeper"`
	Platform *PlatformPlugin `json:"platform" form:"platform"`
}

func isDefaultId(id string) bool {
	return strings.HasPrefix(id, defaultIdPrefix)
}

// unknownPlaceholders reports every variable a template references that is
// outside the closed vocabulary, across all of its commands.
func (r TemplateRequest) unknownPlaceholders() []string {
	unknown := make([]string, 0)
	for _, c := range r.Commands {
		for _, text := range append([]string{c.Command}, c.PostScripts...) {
			for _, v := range keeper.UnknownPlaceholders(text) {
				if !slices.Contains(unknown, v) {
					unknown = append(unknown, v)
				}
			}
		}
	}
	return unknown
}

// invalidPort reports the first default port outside the range a port can take
// at all. Zero is allowed and means the command states none, leaving the deploy
// form's port empty for the user to type.
func (r TemplateRequest) invalidPort() (int, bool) {
	for _, c := range r.Commands {
		for _, port := range []int{c.Defaults.SshPort, c.Defaults.KeeperPort, c.Defaults.DbPort} {
			if port < 0 || port > 65535 {
				return port, true
			}
		}
	}
	return 0, false
}

// trimmed is what actually gets stored and compared: a name is identified by
// what it reads as, so " etcd" and "etcd" are the same name and neither is
// allowed to sit next to the other in the list. A command's default node name
// is trimmed for the same reason - it is written into --name and matched
// against a keeper's own member names later. The default usernames are trimmed
// for the same reason: one is written into a vault entry and authenticates with
// it later.
func (r TemplateRequest) trimmed() TemplateRequest {
	r.Name = strings.TrimSpace(r.Name)
	r.Description = strings.TrimSpace(r.Description)
	r.Defaults.KeeperUser = strings.TrimSpace(r.Defaults.KeeperUser)
	r.Defaults.DbUser = strings.TrimSpace(r.Defaults.DbUser)
	commands := make([]TemplateCommand, len(r.Commands))
	for i, c := range r.Commands {
		c.Defaults.Name = strings.TrimSpace(c.Defaults.Name)
		commands[i] = c
	}
	r.Commands = commands
	return r
}

func (r TemplateRequest) toTemplate(id uuid.UUID, createdAt int64) Template {
	r = r.trimmed()
	return Template{
		Id:          id.String(),
		Name:        r.Name,
		Description: r.Description,
		Keeper:      r.Keeper,
		Platform:    r.Platform,
		Defaults:    r.Defaults,
		Commands:    r.Commands,
		Creation:    Manual,
		CreatedAt:   createdAt,
	}
}
