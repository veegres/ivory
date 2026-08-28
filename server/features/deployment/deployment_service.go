package deployment

import (
	"errors"
	"ivory/core/utils"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
)

var ErrTemplateNameRequired = errors.New("template name is required")
var ErrTemplateNameTaken = errors.New("template name is already taken")
var ErrTemplateCommandsRequired = errors.New("template needs at least one command")
var ErrTemplateReadOnly = errors.New("shipped templates cannot be changed, copy one instead")
var ErrTemplatePluginImmutable = errors.New("template keeper and platform cannot be changed")
var ErrUnknownPlatform = errors.New("unknown platform")

type Service struct {
	repository             *Repository
	keeperMetadataRegistry *utils.Registry[keeper.Plugin, keeper.Metadata]
	platformRegistry       *utils.Registry[platform.Plugin, platform.Adapter]
}

func NewService(
	repository *Repository,
	keeperMetadataRegistry *utils.Registry[keeper.Plugin, keeper.Metadata],
	platformRegistry *utils.Registry[platform.Plugin, platform.Adapter],
) *Service {
	return &Service{
		repository:             repository,
		keeperMetadataRegistry: keeperMetadataRegistry,
		platformRegistry:       platformRegistry,
	}
}
