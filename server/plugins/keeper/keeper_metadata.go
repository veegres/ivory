package keeper

import "ivory/core/config"

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
	DefaultImage  string
	Env           []EnvVar
	Ports         []string
	Volumes       []VolumeSpec
	DefaultValues map[string]string
}

type EnvVar struct {
	Name  string
	Value string
}

type VolumeSpec struct {
	HostPath      string
	ContainerPath string
}
