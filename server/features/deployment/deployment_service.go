package deployment

import (
	"errors"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
)

var ErrTemplateNameRequired = errors.New("template name is required")
var ErrTemplateNameTaken = errors.New("template name is already taken")
var ErrTemplateCommandsRequired = errors.New("template needs at least one command")
var ErrTemplateReadOnly = errors.New("shipped templates cannot be changed, copy one instead")
var ErrTemplatePluginImmutable = errors.New("template keeper and platform cannot be changed")
var ErrUnknownPlatform = errors.New("unknown platform")

// keeperRegistry is the narrow view deployment needs of the keeper plugin
// registry: it reads the deployments every registered plugin ships and never
// operates a running keeper.
type keeperRegistry interface {
	All() map[keeper.PluginType]keeper.Plugin
}

// platformRegistry is narrower still: deployment only asks whether a platform
// is registered at all.
type platformRegistry interface {
	Get(platform.PluginType) (platform.Plugin, error)
}

type Service struct {
	repository       *Repository
	keeperRegistry   keeperRegistry
	platformRegistry platformRegistry
}

func NewService(
	repository *Repository,
	keeperRegistry keeperRegistry,
	platformRegistry platformRegistry,
) *Service {
	return &Service{
		repository:       repository,
		keeperRegistry:   keeperRegistry,
		platformRegistry: platformRegistry,
	}
}
