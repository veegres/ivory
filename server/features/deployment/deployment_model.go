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
	Commands    []TemplateCommand `json:"commands"`
	Creation    CreationType      `json:"creation"`
	CreatedAt   int64             `json:"createdAt"`
}

// TemplateCommand is one node's deployment: it has no identity beyond its
// position, the node it lands on being chosen at deploy time. It is an alias
// rather than a copy so a template shipped by a keeper plugin needs no
// per-command mapping.
type TemplateCommand = keeper.DeploymentCommand

type TemplateRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Keeper      KeeperPlugin      `json:"keeper"`
	Platform    PlatformPlugin    `json:"platform"`
	Commands    []TemplateCommand `json:"commands"`
}

type ListRequest struct {
	Keeper   *KeeperPlugin   `json:"keeper" form:"keeper"`
	Platform *PlatformPlugin `json:"platform" form:"platform"`
}

// withCreation defaults a stored template to Manual, so one saved before the
// field existed reads back as what it actually is - only the shipped templates
// are ever System, and those are computed rather than stored.
func (t Template) withCreation() Template {
	if t.Creation == "" {
		t.Creation = Manual
	}
	return t
}

// normalized is what every read path returns: a stored template carries the
// field names and platform key it was written with, and the rest of the code
// only ever sees the current ones.
func (t Template) normalized() Template {
	t = t.withCreation()
	t.Platform = t.Platform.Current()
	return t
}

// normalized resolves a platform sent under a name that has since been renamed,
// so a backup written before the rename imports as the platform it means.
func (r TemplateRequest) normalized() TemplateRequest {
	r.Platform = r.Platform.Current()
	return r
}

// normalized resolves a platform filter sent under a name that has since been
// renamed, so it matches the templates it was meant to select.
func (r ListRequest) normalized() ListRequest {
	if r.Platform != nil {
		current := r.Platform.Current()
		r.Platform = &current
	}
	return r
}

func isDefaultId(id string) bool {
	return strings.HasPrefix(id, defaultIdPrefix)
}

// unknownPlaceholders reports every variable a template references that is
// outside the closed vocabulary, across all of its commands.
func (r TemplateRequest) unknownPlaceholders() []string {
	unknown := make([]string, 0)
	for _, c := range r.Commands {
		for _, text := range []string{c.Command, c.PostScript} {
			for _, v := range keeper.UnknownPlaceholders(text) {
				if !slices.Contains(unknown, v) {
					unknown = append(unknown, v)
				}
			}
		}
	}
	return unknown
}

// trimmed is what actually gets stored and compared: a name is identified by
// what it reads as, so " etcd" and "etcd" are the same name and neither is
// allowed to sit next to the other in the list.
func (r TemplateRequest) trimmed() TemplateRequest {
	r.Name = strings.TrimSpace(r.Name)
	r.Description = strings.TrimSpace(r.Description)
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
		Commands:    r.Commands,
		Creation:    Manual,
		CreatedAt:   createdAt,
	}
}
