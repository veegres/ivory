package docker

import (
	"ivory/core/config"
	"ivory/plugins/platform"
)

// NOTE: validate that is matches interface in compile-time
var _ platform.Metadata = (*Plugin)(nil)

// SupportedFeatures reports the system view as supported: docker runs
// containers directly on a host Ivory can reach over ssh, so the node's own
// metrics, processes, logs and authorized_keys are all reachable. A platform
// that only ever addresses a scheduler (kubernetes) reports false instead.
func (p *Plugin) SupportedFeatures() map[config.Feature]bool {
	return map[config.Feature]bool{
		config.ViewNodeSystem:   true,
		config.ManageNodeSystem: true,
	}
}
