package platform

import "ivory/core/config"

// Metadata is implemented by platform plugins to describe themselves without
// touching a node, the way keeper.Metadata does for keepers. It is a separate
// interface from Adapter on purpose: Adapter covers operations against a
// reachable node, Metadata covers what the platform can be asked for at all.
type Metadata interface {
	// SupportedFeatures reports, for every platform-related config.Feature this
	// plugin knows about, whether it supports it. A feature absent from the
	// map is not a platform capability at all and is left unrestricted.
	SupportedFeatures() map[config.Feature]bool
}
