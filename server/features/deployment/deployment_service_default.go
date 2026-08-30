package deployment

import (
	"ivory/plugins/keeper"
	"sort"
)

// Defaults returns the templates the registered keeper plugins ship for the
// requested platform. They are computed from the plugins on every request,
// never stored: that is what makes them read-only by construction, keeps them
// in step with whatever the current version ships, and - because
// keeper.Metadata requires DefaultTemplates - makes it impossible to add a
// keeper plugin that has no deployments to copy from.
func (s *Service) Defaults(criteria ListRequest) []Template {
	defaults := make([]Template, 0)
	for pluginType, plugin := range s.keeperRegistry.All() {
		if plugin == nil {
			continue
		}
		if criteria.Keeper != nil && pluginType != *criteria.Keeper {
			continue
		}
		for _, t := range plugin.DefaultTemplates() {
			if criteria.Platform != nil && t.Platform != *criteria.Platform {
				continue
			}
			if _, err := s.platformRegistry.Get(t.Platform); err != nil {
				continue
			}
			defaults = append(defaults, mapDefaultTemplate(pluginType, t))
		}
	}

	// NOTE: Registry.All() is a map, so without this the shipped list would
	// shuffle between requests
	sort.Slice(defaults, func(i, j int) bool {
		if defaults[i].Platform != defaults[j].Platform {
			return defaults[i].Platform < defaults[j].Platform
		}
		return defaults[i].Name < defaults[j].Name
	})
	return defaults
}

func mapDefaultTemplate(plugin keeper.PluginType, t keeper.DeploymentTemplate) Template {
	return Template{
		Id:          defaultId(plugin, t.Platform, t.Name),
		Name:        t.Name,
		Description: t.Description,
		Keeper:      plugin,
		Platform:    t.Platform,
		Defaults:    t.Defaults,
		Commands:    t.Commands,
		Creation:    System,
	}
}
